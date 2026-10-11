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
              // let us read the error envelope instead of throwing on 4xx/5xx
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

  Future<Health> health() =>
      _call(_dio.get('/health'), (d) => Health.fromJson(d as Map<String, dynamic>));

  Future<DrawPage> draws({int limit = 20, int offset = 0, String? from, String? to}) =>
      _call(
        _dio.get('/ftn/draws', queryParameters: {
          'limit': limit,
          'offset': offset,
          if (from != null) 'from': from,
          if (to != null) 'to': to,
        }),
        (d) => DrawPage.fromJson(d as Map<String, dynamic>),
      );

  Future<Draw> latestDraw() =>
      _call(_dio.get('/ftn/draws/latest'), (d) => Draw.fromJson(d as Map<String, dynamic>));

  Future<Draw> draw(String period) =>
      _call(_dio.get('/ftn/draws/$period'), (d) => Draw.fromJson(d as Map<String, dynamic>));

  Future<NumberStats> numberStats({int window = 100}) => _call(
      _dio.get('/ftn/numbers/stats', queryParameters: {'window': window}),
      (d) => NumberStats.fromJson(d as Map<String, dynamic>));

  Future<AnalyzeResult> analyze(List<int> numbers) {
    validatePick(numbers);
    return _call(_dio.post('/ftn/analyze', data: {'numbers': numbers}),
        (d) => AnalyzeResult.fromJson(d as Map<String, dynamic>));
  }
}

void validatePick(List<int> numbers) {
  if (numbers.length != 5 ||
      numbers.toSet().length != 5 ||
      numbers.any((n) => n < 1 || n > 39)) {
    throw ArgumentError('需選 5 個不重複的 1~39 號碼');
  }
}
```

## Models

Hand-written below for clarity; `json_serializable` / `freezed` is fine if the project already uses it. Keep `period` a `String`.

```dart
// lib/data/api/models.dart
List<int> _ints(Object? v) => (v as List).map((e) => e as int).toList();

class Features {
  final int sum, odd, even, big, small, prime, ac, span, tailSum, sameTailGroups, maxGap, minGap;
  final List<int> tenGroup, mod3;
  final bool fullTenGroup;
  final String consecutive; // none, 2, 3, 4, 5, 22, 32 (unknown -> keep string)

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
  final Features features;

  Draw.fromJson(Map<String, dynamic> j)
      : period = j['period'],
        date = DateTime.parse(j['date']),
        yearIndex = j['yearIndex'],
        numbers = _ints(j['numbers']),
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

class NumberStats {
  final int window, totalDraws;
  final double expectedWindow;
  final String latestPeriod;
  final List<NumberStat> numbers;
  NumberStats.fromJson(Map<String, dynamic> j)
      : window = j['window'], totalDraws = j['totalDraws'],
        expectedWindow = (j['expectedWindow'] as num).toDouble(),
        latestPeriod = j['latestPeriod'],
        numbers = (j['numbers'] as List).map((e) => NumberStat.fromJson(e)).toList();
}

class HistoryMatch {
  final int match5, match4, match3, match2;
  HistoryMatch.fromJson(Map<String, dynamic> j)
      : match5 = j['match5'], match4 = j['match4'], match3 = j['match3'], match2 = j['match2'];
}

class AnalyzeResult {
  final List<int> numbers;
  final Features features;
  final HistoryMatch history;
  final int totalDraws;
  final Draw? lastMatch3;
  AnalyzeResult.fromJson(Map<String, dynamic> j)
      : numbers = _ints(j['numbers']),
        features = Features.fromJson(j['features']),
        history = HistoryMatch.fromJson(j['history']),
        totalDraws = j['totalDraws'],
        lastMatch3 = j['lastMatch3'] == null ? null : Draw.fromJson(j['lastMatch3']);
}

class Health {
  final String status;
  final int draws;
  final String? latestPeriod, latestDate;
  Health.fromJson(Map<String, dynamic> j)
      : status = j['status'], draws = j['draws'],
        latestPeriod = j['latestPeriod'], latestDate = j['latestDate'];
}
```

## Platform setup for local http

- Android: in `android/app/src/debug/AndroidManifest.xml` add `android:usesCleartextTraffic="true"` to `<application>` (debug only).
- iOS: add an `NSAppTransportSecurity` → `NSAllowsLocalNetworking = YES` entry in the debug Info.plist.

## Test

```dart
test('decodes error envelope', () async {
  final dio = Dio()..httpClientAdapter = /* mock adapter returning
      400 {"code":40001,"message":"numbers must contain exactly 5 numbers"} */;
  final api = LotteryApi(dio: dio);
  expect(() => api.analyze([1, 2, 3, 4, 5]),
      throwsA(isA<ApiException>().having((e) => e.code, 'code', 40001)));
});
```
