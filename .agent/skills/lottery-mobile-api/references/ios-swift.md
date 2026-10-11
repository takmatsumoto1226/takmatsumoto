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
```

## Models

```swift
struct Envelope<T: Decodable>: Decodable {
    let code: Int
    let message: String
    let data: T?
}

struct Features: Codable, Hashable {
    let sum, odd, even, big, small, prime: Int
    let tenGroup: [Int]
    let fullTenGroup: Bool
    let consecutive: String   // none, 2, 3, 4, 5, 22, 32; keep unknown values
    let ac, span: Int
    let mod3: [Int]
    let tailSum, sameTailGroups, maxGap, minGap: Int
}

struct Draw: Codable, Hashable, Identifiable {
    let period: String
    let date: String          // YYYY-MM-DD local date
    let yearIndex: Int
    let numbers: [Int]
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
    let expectedWindow: Double
    let latestPeriod: String
    let numbers: [NumberStat]
}

struct HistoryMatch: Codable { let match5, match4, match3, match2: Int }

struct AnalyzeResult: Codable {
    let numbers: [Int]
    let features: Features
    let history: HistoryMatch
    let totalDraws: Int
    let lastMatch3: Draw?
}

struct Health: Codable {
    let status: String
    let draws: Int
    let latestPeriod: String?
    let latestDate: String?
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
            // error envelopes have no data; decode code/message only
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

    func draws(limit: Int = 20, offset: Int = 0, from: String? = nil, to: String? = nil) async throws -> DrawPage {
        var q = [URLQueryItem(name: "limit", value: "\(limit)"), URLQueryItem(name: "offset", value: "\(offset)")]
        if let from { q.append(.init(name: "from", value: from)) }
        if let to { q.append(.init(name: "to", value: to)) }
        return try await request("ftn/draws", query: q)
    }

    func latestDraw() async throws -> Draw { try await request("ftn/draws/latest") }

    func draw(period: String) async throws -> Draw { try await request("ftn/draws/\(period)") }

    func numberStats(window: Int = 100) async throws -> NumberStats {
        try await request("ftn/numbers/stats", query: [.init(name: "window", value: "\(window)")])
    }

    func analyze(_ numbers: [Int]) async throws -> AnalyzeResult {
        guard numbers.count == 5, Set(numbers).count == 5, numbers.allSatisfy({ (1...39).contains($0) }) else {
            throw ApiError(code: 40001, message: "需選 5 個不重複的 1~39 號碼", httpStatus: nil)
        }
        struct Body: Encodable { let numbers: [Int] }
        return try await request("ftn/analyze", method: "POST", body: Body(numbers: numbers))
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
