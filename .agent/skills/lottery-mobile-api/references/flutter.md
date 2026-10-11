# Flutter / Dart client template

Uses `dio` for HTTP. Swap in `http` if the project already uses it; keep the same layering.

```yaml
# pubspec.yaml
dependencies:
  dio: ^5.4.0
```

## Config

```dart
// lib/data/api/api_config.dart
class ApiConfig {
  // flutter run --dart-define=API_BASE_URL=http://10.0.2.2:8080/api/v1 --dart-define=API_KEY=xxx
  static const baseUrl = String.fromEnvironment('API_BASE_URL',
      defaultValue: 'http://10.0.2.2:8080/api/v1');
  static const apiKey = String.fromEnvironment('API_KEY');
}

/// Path code of each game. Rules (ranges, prizes) come from GET /games, not from here.
enum GameCode {
  ftn, // 今彩539
  bl,  // 大樂透
  pw;  // 威力彩
}
```

## Envelope and errors

```dart
// lib/data/api/api_exception.dart
class ApiException implements Exception {
  final int code;          // envelope code, -1 for network/parse errors
  final String message;
  final int? httpStatus;
  ApiException(this.code, this.message, [this.httpStatus]);

  bool get isRetryable => code == -1 || code == 50301;
  @override
  String toString() => 'ApiException($code, $httpStatus): $message';
}
```

## Client

```dart
// lib/data/api/lottery_api.dart
import 'package:dio/dio.dart';
import 'api_config.dart';
import 'api_exception.dart';
import 'models.dart';

class LotteryApi {
  LotteryApi({Dio? dio})
      : _dio = dio ??
            Dio(BaseOptions(
              baseUrl: ApiConfig.baseUrl,
              connectTimeout: const Duration(seconds: 10),
              receiveTimeout: const Duration(seconds: 15),
              headers: {
                if (ApiConfig.apiKey.isNotEmpty) 'X-API-Key': ApiConfig.apiKey,
              },
              // read the error envelope instead of throwing on 4xx/5xx
              validateStatus: (_) => true,
            ));

  final Dio _dio;

  Future<T> _call<T>(Future<Response<dynamic>> req, T Function(Object? data) parse) async {
    final Response<dynamic> res;
    try {
      res = await req;
    } on DioException catch (e) {
      throw ApiException(-1, e.message ?? 'network error');
    }
    final body = res.data;
    if (body is! Map<String, dynamic>) {
      throw ApiException(-1, 'unexpected response', res.statusCode);
    }
    final code = body['code'] as int? ?? -1;
    if (code != 0) {
      throw ApiException(code, body['message'] as String? ?? '', res.statusCode);
    }
    return parse(body['data']);
  }

  Map<String, dynamic> _map(Object? d) => d as Map<String, dynamic>;

  Future<Health> health() => _call(_dio.get('/health'), (d) => Health.fromJson(_map(d)));

  Future<List<GameInfo>> games() => _call(_dio.get('/games'),
      (d) => (d as List).map((e) => GameInfo.fromJson(_map(e))).toList());

  Future<DrawPage> draws(GameCode game, {int limit = 20, int offset = 0, String? from, String? to}) =>
      _call(
        _dio.get('/${game.name}/draws', queryParameters: {
          'limit': limit,
          'offset': offset,
          if (from != null) 'from': from,
          if (to != null) 'to': to,
        }),
        (d) => DrawPage.fromJson(_map(d)),
      );

  Future<Draw> latestDraw(GameCode game) =>
      _call(_dio.get('/${game.name}/draws/latest'), (d) => Draw.fromJson(_map(d)));

  Future<Draw> draw(GameCode game, String period) =>
      _call(_dio.get('/${game.name}/draws/$period'), (d) => Draw.fromJson(_map(d)));

  Future<NumberStats> numberStats(GameCode game, {int window = 100}) => _call(
      _dio.get('/${game.name}/numbers/stats', queryParameters: {'window': window}),
      (d) => NumberStats.fromJson(_map(d)));

  /// [rules] comes from games(); validates before sending.
  Future<AnalyzeResult> analyze(GameInfo rules, List<int> numbers, {int? special}) {
    validatePick(rules, numbers, special);
    return _call(
        _dio.post('/${rules.code}/analyze', data: {
          'numbers': numbers,
          if (special != null) 'special': special,
        }),
        (d) => AnalyzeResult.fromJson(_map(d)));
  }
}

void validatePick(GameInfo rules, List<int> numbers, int? special) {
  if (numbers.length != rules.pickCount ||
      numbers.toSet().length != rules.pickCount ||
      numbers.any((n) => n < 1 || n > rules.maxNumber)) {
    throw ArgumentError('需選 ${rules.pickCount} 個不重複的 1~${rules.maxNumber} 號碼');
  }
  final sp = rules.special;
  if (sp != null && sp.userPick) {
    if (special == null || special < 1 || special > sp.maxNumber) {
      throw ArgumentError('需選 1 個 1~${sp.maxNumber} 的${sp.name}號碼');
    }
  } else if (special != null) {
    throw ArgumentError('${rules.name}不需另外選特別號');
  }
}
```

## Models

Hand-written below for clarity; `json_serializable` / `freezed` is fine if the project already uses it. Keep `period` a `String`; `special` fields are nullable.

```dart
// lib/data/api/models.dart
List<int> _ints(Object? v) => (v as List).map((e) => e as int).toList();

class SpecialSpec {
  final String name;
  final int maxNumber;
  final bool samePool, userPick;
  SpecialSpec.fromJson(Map<String, dynamic> j)
      : name = j['name'], maxNumber = j['maxNumber'], samePool = j['samePool'], userPick = j['userPick'];
}

class Prize {
  final int main;
  final bool special;
  final String name;
  Prize.fromJson(Map<String, dynamic> j) : main = j['main'], special = j['special'], name = j['name'];
}

class GameInfo {
  final String code, name;
  final int maxNumber, pickCount, bigFrom;
  final SpecialSpec? special;
  final List<Prize> prizes;
  GameInfo.fromJson(Map<String, dynamic> j)
      : code = j['code'], name = j['name'],
        maxNumber = j['maxNumber'], pickCount = j['pickCount'], bigFrom = j['bigFrom'],
        special = j['special'] == null ? null : SpecialSpec.fromJson(j['special']),
        prizes = (j['prizes'] as List).map((e) => Prize.fromJson(e)).toList();
}

class Features {
  final int sum, odd, even, big, small, prime, ac, span, tailSum, sameTailGroups, maxGap, minGap;
  final List<int> tenGroup; // length varies by game
  final List<int> mod3;
  final bool fullTenGroup;
  final String consecutive; // open set: none, 2, 3, 22, 32, 222, ...

  Features.fromJson(Map<String, dynamic> j)
      : sum = j['sum'], odd = j['odd'], even = j['even'],
        big = j['big'], small = j['small'], prime = j['prime'],
        ac = j['ac'], span = j['span'], tailSum = j['tailSum'],
        sameTailGroups = j['sameTailGroups'], maxGap = j['maxGap'], minGap = j['minGap'],
        tenGroup = _ints(j['tenGroup']), mod3 = _ints(j['mod3']),
        fullTenGroup = j['fullTenGroup'], consecutive = j['consecutive'];
}

class Draw {
  final String period;
  final DateTime date; // local date, no time
  final int yearIndex;
  final List<int> numbers;
  final int? special; // bl 特別號 / pw 第二區; null for ftn
  final Features features;

  Draw.fromJson(Map<String, dynamic> j)
      : period = j['period'],
        date = DateTime.parse(j['date']),
        yearIndex = j['yearIndex'],
        numbers = _ints(j['numbers']),
        special = j['special'],
        features = Features.fromJson(j['features']);
}

class DrawPage {
  final List<Draw> items;
  final int total, limit, offset;
  bool get hasMore => offset + items.length < total;

  DrawPage.fromJson(Map<String, dynamic> j)
      : items = (j['items'] as List).map((e) => Draw.fromJson(e)).toList(),
        total = j['total'], limit = j['limit'], offset = j['offset'];
}

class NumberStat {
  final int number, windowCount, totalCount, omission, maxOmission;
  NumberStat.fromJson(Map<String, dynamic> j)
      : number = j['number'], windowCount = j['windowCount'], totalCount = j['totalCount'],
        omission = j['omission'], maxOmission = j['maxOmission'];
}

List<NumberStat>? _stats(Object? v) =>
    v == null ? null : (v as List).map((e) => NumberStat.fromJson(e)).toList();

class NumberStats {
  final int window, totalDraws;
  final String latestPeriod;
  final double expectedWindow;
  final List<NumberStat> numbers;
  final double? specialExpectedWindow;
  final List<NumberStat>? special; // bl / pw only

  NumberStats.fromJson(Map<String, dynamic> j)
      : window = j['window'], totalDraws = j['totalDraws'], latestPeriod = j['latestPeriod'],
        expectedWindow = (j['expectedWindow'] as num).toDouble(),
        numbers = _stats(j['numbers'])!,
        specialExpectedWindow = (j['specialExpectedWindow'] as num?)?.toDouble(),
        special = _stats(j['special']);
}

class MatchCount {
  final int main, count;
  final bool special;
  final String? prize; // null when the combination wins nothing
  MatchCount.fromJson(Map<String, dynamic> j)
      : main = j['main'], special = j['special'], count = j['count'], prize = j['prize'];
}

class AnalyzeResult {
  final List<int> numbers;
  final int? special;
  final Features features;
  final List<MatchCount> history;
  final int totalDraws;
  final Draw? lastMatch3;
  AnalyzeResult.fromJson(Map<String, dynamic> j)
      : numbers = _ints(j['numbers']),
        special = j['special'],
        features = Features.fromJson(j['features']),
        history = (j['history'] as List).map((e) => MatchCount.fromJson(e)).toList(),
        totalDraws = j['totalDraws'],
        lastMatch3 = j['lastMatch3'] == null ? null : Draw.fromJson(j['lastMatch3']);
}

class GameHealth {
  final int draws;
  final String? latestPeriod, latestDate;
  GameHealth.fromJson(Map<String, dynamic> j)
      : draws = j['draws'], latestPeriod = j['latestPeriod'], latestDate = j['latestDate'];
}

class Health {
  final String status; // ok, partial, data_not_ready
  final Map<String, GameHealth> games;
  Health.fromJson(Map<String, dynamic> j)
      : status = j['status'],
        games = (j['games'] as Map<String, dynamic>)
            .map((k, v) => MapEntry(k, GameHealth.fromJson(v)));
}
```

## Platform setup for local http

- Android: in `android/app/src/debug/AndroidManifest.xml` add `android:usesCleartextTraffic="true"` to `<application>` (debug only).
- iOS: add `NSAppTransportSecurity` → `NSAllowsLocalNetworking = YES` in the debug Info.plist.

## Test

```dart
test('decodes ftn draw without special', () {
  final d = Draw.fromJson(jsonDecode(ftnLatestSample)['data']);
  expect(d.special, isNull);
  expect(d.features.tenGroup.length, 4);
});

test('decodes error envelope', () async {
  final dio = Dio()..httpClientAdapter = /* mock adapter returning
      400 {"code":40001,"message":"special is required for pw"} */;
  final api = LotteryApi(dio: dio);
  expect(() => api.draws(GameCode.pw),
      throwsA(isA<ApiException>().having((e) => e.code, 'code', 40001)));
});
```

## Game menu

All three games must appear. Entries come from `games()`, with a built-in fallback.

```dart
// lib/ui/game_menu.dart
const fallbackGames = [
  (code: 'ftn', name: '今彩539'),
  (code: 'bl', name: '大樂透'),
  (code: 'pw', name: '威力彩'),
];

class GameMenu extends StatelessWidget {
  const GameMenu({super.key, required this.games, required this.selected, required this.onSelect});

  final List<GameInfo>? games; // null while loading or when /games failed
  final String selected;       // game code
  final ValueChanged<String> onSelect;

  @override
  Widget build(BuildContext context) {
    final entries = games?.map((g) => (code: g.code, name: g.name)).toList() ?? fallbackGames;
    return NavigationBar(
      selectedIndex: entries.indexWhere((e) => e.code == selected).clamp(0, entries.length - 1),
      onDestinationSelected: (i) => onSelect(entries[i].code),
      destinations: [
        for (final e in entries) NavigationDestination(icon: const Icon(Icons.casino_outlined), label: e.name),
      ],
    );
  }
}
```

Persist `selected` with `shared_preferences`, and route the three per-game screens by code (e.g. `go_router` paths `/:game/draws`, `/:game/stats`, `/:game/analyze`).
