---
name: lottery-mobile-api
description: Integrate a mobile app (Flutter, Android/Kotlin, iOS/Swift) with the Lottery 539 REST API served by the Go/gin backend. Use when building or changing app screens, models, repositories or network code that fetch 今彩539 draws, number statistics or combination analysis from /api/v1, or when debugging those API calls.
---

# Lottery 539 Mobile API

This skill describes the backend contract and the rules for calling it from a mobile app. The backend lives in the Go repo `lottery` (gin), entry point `main.go`, handlers in `api/ftn/`.

The full machine-readable contract is `references/openapi.yaml`. Read it before generating models or network code. Platform templates:

- Flutter / Dart: `references/flutter.md`
- Android / Kotlin: `references/android-kotlin.md`
- iOS / Swift: `references/ios-swift.md`

## Base URL and environments

| Environment | Base URL |
|---|---|
| Android emulator → server on dev machine | `http://10.0.2.2:8080/api/v1` |
| iOS simulator → server on dev machine | `http://localhost:8080/api/v1` |
| Real device on same Wi-Fi | `http://<dev-machine-LAN-IP>:8080/api/v1` |
| Production | ask the user; must be `https://` |

Never hardcode the base URL inside feature code. Put it in one config object (flavor / build config / `--dart-define`) so dev and prod differ only by configuration.

Plain `http://` is only for local development. It needs `android:usesCleartextTraffic` (or a network security config) on Android and an ATS exception on iOS. Do not ship those exceptions in release builds.

## Authentication

- If the server sets `server.api_key` in `config.yaml`, every `/ftn/*` request must send header `X-API-Key: <key>`.
- `GET /health` never needs the key.
- Inject the key from build configuration, not source code. A key inside an app binary can be extracted, so treat it as a rate-limit / casual-abuse guard, not as real user authentication.

## Response envelope

Every response, success or error, is JSON:

```json
{ "code": 0, "message": "ok", "data": { ... } }
```

- `code == 0` means success; `data` holds the payload.
- On error, `data` is absent and the HTTP status is non-2xx.

| HTTP | code | Meaning | App behaviour |
|---|---|---|---|
| 400 | 40001 | Invalid parameter | Show `message`; it is a developer-readable English string, map it to localized UI text |
| 401 | 40101 | Missing / wrong API key | Configuration bug; log it, show generic error |
| 404 | 40401 | Period not found | Show "找不到該期" |
| 503 | 50301 | Server has not loaded draw data yet | Show retry UI; retry with backoff |

Always parse the envelope first, then `data`. Write one shared decode function that turns non-zero `code` into a typed error (`ApiException(code, message, httpStatus)`); never let every call site re-implement this.

## Endpoints

All paths below are relative to the base URL. All numbers are 今彩539 numbers 1–39, five per draw, always sorted ascending.

| Method | Path | Purpose |
|---|---|---|
| GET | `/health` | Service status, latest period, draw count |
| GET | `/ftn/draws?limit=20&offset=0&from=YYYY-MM-DD&to=YYYY-MM-DD` | Paged draw history, newest first. `limit` 1–200 (default 20) |
| GET | `/ftn/draws/latest` | Latest draw |
| GET | `/ftn/draws/{period}` | One draw by global period number, e.g. `5990` |
| GET | `/ftn/numbers/stats?window=100` | Per-number counts in the last `window` draws, total count, current and max omission |
| POST | `/ftn/analyze` body `{"numbers":[3,7,15,17,30]}` | Features of a 5-number pick plus how often history matched 2/3/4/5 of them |

### Core models

`Draw`

```json
{
  "period": "5990",
  "date": "2026-10-10",
  "yearIndex": 246,
  "numbers": [14, 21, 23, 32, 34],
  "features": { ...Features }
}
```

- `period` is a **string**; keep it a string in the app (it is an identifier, not a quantity).
- `date` is `YYYY-MM-DD`, Taiwan local date, no time part. Parse as a local date, not a UTC instant.

`Features`

| Field | Type | Meaning |
|---|---|---|
| `sum` | int | Sum of the 5 numbers |
| `odd`, `even` | int | Odd / even count |
| `big`, `small` | int | Count of numbers ≥ 20 / < 20 |
| `prime` | int | Count of primes |
| `tenGroup` | int[4] | Counts in 1–10, 11–20, 21–30, 31–39 |
| `fullTenGroup` | bool | All four groups present |
| `consecutive` | string | `none`, `2`, `3`, `4`, `5`, `22` (two pairs), `32` (a triple and a pair) |
| `ac` | int | AC value (distinct pairwise differences − 4) |
| `span` | int | Max − min |
| `mod3` | int[3] | Count with remainder 0, 1, 2 when divided by 3 |
| `tailSum` | int | Sum of last digits |
| `sameTailGroups` | int | Number of last digits that appear 2+ times |
| `maxGap`, `minGap` | int | Largest / smallest gap between adjacent numbers |

Treat unknown `consecutive` values as "other" instead of crashing, and ignore unknown JSON fields, so new backend fields do not break old app versions.

`DrawPage`: `{ items: Draw[], total, limit, offset }`. More pages exist while `offset + items.length < total`.

`NumberStats`: `{ window, totalDraws, expectedWindow, latestPeriod, numbers: NumberStat[39] }`, where `NumberStat` is `{ number, windowCount, totalCount, omission, maxOmission }`. `expectedWindow` is the count every number would have on average under pure chance; use it as the baseline when drawing hot/cold bars.

`AnalyzeResult`: `{ numbers, features, history: {match5, match4, match3, match2}, totalDraws, lastMatch3?: Draw }`.

## Implementation rules

1. **One API client layer.** A single client/repository owns the base URL, headers, timeouts (connect 10 s, read 15 s), envelope decoding and error mapping. UI code calls repository methods that return domain models, never raw JSON.
2. **Validate before POST /analyze.** Exactly 5 distinct integers 1–39. The server rejects anything else with 40001; checking first gives instant feedback.
3. **Paging.** Use `limit`/`offset` with infinite scroll; default page size 20. Stop when `offset >= total`.
4. **Caching.** Draws happen Monday–Saturday evenings and the server refreshes its data once a day, so cache `draws` and `numbers/stats` for at least several minutes and serve cached data while refreshing. Use `/health.latestPeriod` as a cheap "has anything changed" check.
5. **Retries.** Retry only idempotent GETs, only on network errors / 503, with exponential backoff (e.g. 1 s, 2 s, 4 s, max 3 tries). Do not retry 4xx.
6. **Wording.** Statistics are descriptive history. The app must not present them as predictions or as raising the chance of winning; every 5-number pick has the same 1 / 575,757 chance. Use labels like "近 100 期出現次數" or "目前遺漏期數", not "推薦號碼" or "必中".
7. **Tests.** Unit-test JSON decoding with the sample payloads in `references/openapi.yaml`, plus at least one error envelope (e.g. 40001) and one 503.

## Verifying against a running backend

Start the server from the Go repo root (`go run .`), then:

```bash
curl -s http://localhost:8080/api/v1/health
curl -s 'http://localhost:8080/api/v1/ftn/draws?limit=2'
curl -s 'http://localhost:8080/api/v1/ftn/numbers/stats?window=100'
curl -s -X POST http://localhost:8080/api/v1/ftn/analyze \
  -H 'Content-Type: application/json' -d '{"numbers":[3,7,15,17,30]}'
# add  -H 'X-API-Key: <key>'  when server.api_key is set
```

If the app and curl disagree, compare the request URL, headers and body first; the backend logs every request line (gin logger) in its console.
