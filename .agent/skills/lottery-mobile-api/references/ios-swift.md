# iOS / Swift client template

`URLSession` + `async/await` + `Codable`, no third-party dependencies.

## Config

Put values in an `.xcconfig` per configuration and expose them through Info.plist keys `API_BASE_URL` and `API_KEY`.

```swift
enum ApiConfig {
    static let baseURL = URL(string: Bundle.main.object(forInfoDictionaryKey: "API_BASE_URL") as? String
                             ?? "http://localhost:8080/api/v1")!
    static let apiKey = Bundle.main.object(forInfoDictionaryKey: "API_KEY") as? String ?? ""
}

/// Path code of each game. Rules come from GET /games.
enum GameCode: String, CaseIterable {
    case ftn   // 今彩539
    case bl    // 大樂透
    case pw    // 威力彩
}
```

## Models

```swift
struct Envelope<T: Decodable>: Decodable {
    let code: Int
    let message: String
    let data: T?
}

struct SpecialSpec: Codable, Hashable {
    let name: String
    let maxNumber: Int
    let samePool: Bool
    let userPick: Bool
}

struct Prize: Codable, Hashable {
    let main: Int
    let special: Bool
    let name: String
}

struct GameInfo: Codable, Hashable, Identifiable {
    let code, name: String
    let maxNumber, pickCount, bigFrom: Int
    let special: SpecialSpec?
    let prizes: [Prize]
    var id: String { code }
}

struct Features: Codable, Hashable {
    let sum, odd, even, big, small, prime: Int
    let tenGroup: [Int]       // length varies by game
    let fullTenGroup: Bool
    let consecutive: String   // open set: none, 2, 3, 22, 32, 222, ...
    let ac, span: Int
    let mod3: [Int]
    let tailSum, sameTailGroups, maxGap, minGap: Int
}

struct Draw: Codable, Hashable, Identifiable {
    let period: String
    let date: String          // YYYY-MM-DD local date
    let yearIndex: Int
    let numbers: [Int]
    let special: Int?         // bl 特別號 / pw 第二區; nil for ftn
    let features: Features
    var id: String { period }
}

struct DrawPage: Codable {
    let items: [Draw]
    let total, limit, offset: Int
    var hasMore: Bool { offset + items.count < total }
}

struct NumberStat: Codable, Hashable, Identifiable {
    let number, windowCount, totalCount, omission, maxOmission: Int
    var id: Int { number }
}

struct NumberStats: Codable {
    let window, totalDraws: Int
    let latestPeriod: String
    let expectedWindow: Double
    let numbers: [NumberStat]
    let specialExpectedWindow: Double?
    let special: [NumberStat]?   // bl / pw only
}

struct MatchCount: Codable, Hashable {
    let main: Int
    let special: Bool
    let count: Int
    let prize: String?        // nil when the combination wins nothing
}

struct AnalyzeResult: Codable {
    let numbers: [Int]
    let special: Int?
    let features: Features
    let history: [MatchCount]
    let totalDraws: Int
    let lastMatch3: Draw?
}

struct GameHealth: Codable {
    let draws: Int
    let latestPeriod: String?
    let latestDate: String?
}

struct Health: Codable {
    let status: String        // ok, partial, data_not_ready
    let games: [String: GameHealth]
}
```

## Client

```swift
struct ApiError: Error, LocalizedError {
    let code: Int          // envelope code, -1 for network/parse errors
    let message: String
    let httpStatus: Int?
    var retryable: Bool { code == -1 || code == 50301 }
    var errorDescription: String? { message }
}

private struct Empty: Decodable {}

final class LotteryAPI {
    static let shared = LotteryAPI()

    private let session: URLSession = {
        let c = URLSessionConfiguration.default
        c.timeoutIntervalForRequest = 15
        return URLSession(configuration: c)
    }()
    private let decoder = JSONDecoder()

    private func request<T: Decodable>(_ path: String, query: [URLQueryItem] = [],
                                       method: String = "GET", body: Encodable? = nil) async throws -> T {
        var comps = URLComponents(url: ApiConfig.baseURL.appendingPathComponent(path), resolvingAgainstBaseURL: false)!
        if !query.isEmpty { comps.queryItems = query }
        var req = URLRequest(url: comps.url!)
        req.httpMethod = method
        if !ApiConfig.apiKey.isEmpty { req.setValue(ApiConfig.apiKey, forHTTPHeaderField: "X-API-Key") }
        if let body {
            req.setValue("application/json", forHTTPHeaderField: "Content-Type")
            req.httpBody = try JSONEncoder().encode(body)
        }

        let data: Data, response: URLResponse
        do { (data, response) = try await session.data(for: req) }
        catch { throw ApiError(code: -1, message: error.localizedDescription, httpStatus: nil) }
        let status = (response as? HTTPURLResponse)?.statusCode

        guard let env = try? decoder.decode(Envelope<T>.self, from: data) else {
            if let e = try? decoder.decode(Envelope<Empty>.self, from: data) {
                throw ApiError(code: e.code, message: e.message, httpStatus: status)
            }
            throw ApiError(code: -1, message: "unexpected response", httpStatus: status)
        }
        guard env.code == 0, let payload = env.data else {
            throw ApiError(code: env.code, message: env.message, httpStatus: status)
        }
        return payload
    }

    func health() async throws -> Health { try await request("health") }

    func games() async throws -> [GameInfo] { try await request("games") }

    func draws(_ game: GameCode, limit: Int = 20, offset: Int = 0, from: String? = nil, to: String? = nil) async throws -> DrawPage {
        var q = [URLQueryItem(name: "limit", value: "\(limit)"), URLQueryItem(name: "offset", value: "\(offset)")]
        if let from { q.append(.init(name: "from", value: from)) }
        if let to { q.append(.init(name: "to", value: to)) }
        return try await request("\(game.rawValue)/draws", query: q)
    }

    func latestDraw(_ game: GameCode) async throws -> Draw { try await request("\(game.rawValue)/draws/latest") }

    func draw(_ game: GameCode, period: String) async throws -> Draw {
        try await request("\(game.rawValue)/draws/\(period)")
    }

    func numberStats(_ game: GameCode, window: Int = 100) async throws -> NumberStats {
        try await request("\(game.rawValue)/numbers/stats", query: [.init(name: "window", value: "\(window)")])
    }

    /// `rules` comes from games(); validates before sending.
    func analyze(_ rules: GameInfo, numbers: [Int], special: Int? = nil) async throws -> AnalyzeResult {
        guard numbers.count == rules.pickCount, Set(numbers).count == rules.pickCount,
              numbers.allSatisfy({ (1...rules.maxNumber).contains($0) }) else {
            throw ApiError(code: 40001, message: "需選 \(rules.pickCount) 個不重複的 1~\(rules.maxNumber) 號碼", httpStatus: nil)
        }
        if let sp = rules.special, sp.userPick {
            guard let special, (1...sp.maxNumber).contains(special) else {
                throw ApiError(code: 40001, message: "需選 1 個 1~\(sp.maxNumber) 的\(sp.name)號碼", httpStatus: nil)
            }
        } else if special != nil {
            throw ApiError(code: 40001, message: "\(rules.name)不需另外選特別號", httpStatus: nil)
        }
        struct Body: Encodable { let numbers: [Int]; let special: Int? }  // nil is omitted by JSONEncoder
        return try await request("\(rules.code)/analyze", method: "POST", body: Body(numbers: numbers, special: special))
    }
}
```

## Local http (debug only)

In the Debug Info.plist:

```xml
<key>NSAppTransportSecurity</key>
<dict>
  <key>NSAllowsLocalNetworking</key><true/>
</dict>
```

On a real device reaching a LAN IP, iOS also shows the local-network permission prompt; add `NSLocalNetworkUsageDescription`.

## Game menu (SwiftUI)

All three games must appear. Entries come from `/games`, with a built-in fallback.

```swift
private let fallbackGames: [(code: String, name: String)] = [("ftn", "今彩539"), ("bl", "大樂透"), ("pw", "威力彩")]

struct GameMenu: View {
    let games: [GameInfo]?          // nil while loading or when /games failed
    @AppStorage("selectedGame") private var selected = "ftn"

    var body: some View {
        let entries = games?.map { (code: $0.code, name: $0.name) } ?? fallbackGames
        TabView(selection: $selected) {
            ForEach(entries, id: \.code) { entry in
                GameHomeView(gameCode: entry.code)   // hosts 開獎紀錄 / 號碼統計 / 選號分析
                    .tabItem { Label(entry.name, systemImage: "die.face.5") }
                    .tag(entry.code)
            }
        }
    }
}
```
