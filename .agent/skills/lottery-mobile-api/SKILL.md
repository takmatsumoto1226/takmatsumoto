---
name: lottery-mobile-api
description: Integrate a mobile app (Flutter, Android/Kotlin, iOS/Swift) with the Lottery REST API served by the Go/gin backend, covering 今彩539 (ftn), 大樂透 (bl) and 威力彩 (pw). Use when building or changing app screens, models, repositories or network code that fetch draws, number statistics, game rules or pick analysis from /api/v1, or when debugging those API calls.
---

# Lottery Mobile API

This skill describes the backend contract and the rules for calling it from a mobile app. The backend lives in the Go repo `lottery` (gin): entry point `main.go`, handlers in `api/lottery/`.

The full machine-readable contract is `references/openapi.yaml`. Read it before generating models or network code. Platform templates:

- Flutter / Dart: `references/flutter.md`
- Android / Kotlin: `references/android-kotlin.md`
- iOS / Swift: `references/ios-swift.md`

## Games

All three games share the same endpoints and models; only the path prefix `{game}` and the rules differ. Fetch the rules from `GET /games` instead of hardcoding them, so new games or rule changes need no app release.

| `{game}` | Name | Main numbers | Special | User picks |
|---|---|---|---|---|
| `ftn` | 今彩539 | 5 of 1–39 | none | 5 numbers |
| `bl` | 大樂透 | 6 of 1–49 | 特別號: 1 more number drawn from the **same pool** | 6 numbers (no special pick) |
| `pw` | 威力彩 | 6 of 1–38 (第一區) | 第二區: 1 of 1–8, **separate pool** | 6 numbers + 1 第二區 number |

`special.samePool == true` (大樂透) means the special number is matched against the user's 6 numbers. `special.userPick == true` (威力彩) means the user must choose a separate special number.

## Base URL and environments

| Environment | Base URL |
|---|---|
| Android emulator → server on dev machine | `http://10.0.2.2:8080/api/v1` |
| iOS simulator → server on dev machine | `http://localhost:8080/api/v1` |
| Real device on same Wi-Fi | `http://<dev-machine-LAN-IP>:8080/api/v1` |
| Production | ask the user; must be `https://` |

Never hardcode the base URL inside feature code. Put it in one config object (flavor / build config / `--dart-define`) so dev and prod differ only by configuration.

Plain `http://` is only for local development. It needs a cleartext exception on Android and an ATS exception on iOS. Do not ship those exceptions in release builds.

## Authentication

- If the server sets `server.api_key` in `config.yaml`, every request except `/health` must send header `X-API-Key: <key>`.
- Inject the key from build configuration, not source code. A key inside an app binary can be extracted, so treat it as a casual-abuse guard, not as user authentication.

## Response envelope

Every response, success or error, is JSON:

```json
{ "code": 0, "message": "ok", "data": { ... } }
```

- `code == 0` means success; `data` holds the payload.
- On error, `data` is absent and the HTTP status is non-2xx.

| HTTP | code | Meaning | App behaviour |
|---|---|---|---|
| 400 | 40001 | Invalid parameter | `message` is developer-readable English; map it to localized UI text |
| 401 | 40101 | Missing / wrong API key | Configuration bug; log it, show generic error |
| 404 | 40401 | Period not found (or unknown `{game}`: plain 404 without envelope) | Show "找不到該期" |
| 503 | 50301 | Server has not loaded this game's data yet | Show retry UI; retry with backoff |

Write one shared decode function that turns non-zero `code` into a typed error (`ApiException(code, message, httpStatus)`); never re-implement this per call site.

## Endpoints

Paths are relative to the base URL. Main numbers are always sorted ascending.

| Method | Path | Purpose |
|---|---|---|
| GET | `/health` | Overall status and per-game draw count / latest period (no API key) |
| GET | `/games` | Rules of every game: ranges, pick count, special spec, prize table |
| GET | `/{game}` | Rules of one game |
| GET | `/{game}/draws?limit=20&offset=0&from=YYYY-MM-DD&to=YYYY-MM-DD` | Paged draw history, newest first. `limit` 1–200 (default 20) |
| GET | `/{game}/draws/latest` | Latest draw |
| GET | `/{game}/draws/{period}` | One draw by global period number, e.g. `/bl/draws/2489` |
| GET | `/{game}/numbers/stats?window=100` | Per-number counts in the last `window` draws, total count, current and max omission; plus `special` stats for bl/pw |
| POST | `/{game}/analyze` | Features of a pick plus how often history matched it, grouped by prize |

### Core models

`GameInfo`

```json
{
  "code": "pw", "name": "威力彩", "maxNumber": 38, "pickCount": 6, "bigFrom": 20,
  "special": { "name": "第二區", "maxNumber": 8, "samePool": false, "userPick": true },
  "prizes": [ { "main": 6, "special": true, "name": "頭獎" }, ... ]
}
```

`special` is absent for ftn. Use `maxNumber`, `pickCount` and `special` to build the number picker and validation for any game.

`Draw`

```json
{
  "period": "1953",
  "date": "2026-10-08",
  "yearIndex": 81,
  "numbers": [4, 23, 27, 28, 29, 38],
  "special": 5,
  "features": { ...Features }
}
```

- `period` is a **string** identifier; keep it a string. It never has leading zeros.
- `date` is `YYYY-MM-DD`, Taiwan local date, no time. Parse as a local date, not a UTC instant.
- `special` is absent for ftn; for bl it is the 特別號, for pw the 第二區 number. Display it visually separate from `numbers`.

`Features` (computed on the main numbers only)

| Field | Type | Meaning |
|---|---|---|
| `sum` | int | Sum of main numbers |
| `odd`, `even` | int | Odd / even count |
| `big`, `small` | int | Count of numbers ≥ / < the game's `bigFrom` (ftn 20, bl 25, pw 20) |
| `prime` | int | Count of primes |
| `tenGroup` | int[] | Counts in 1–10, 11–20, …; length 4 for ftn/pw, 5 for bl |
| `fullTenGroup` | bool | Every group present |
| `consecutive` | string | Run lengths of consecutive numbers, longest first: `none`, `2`, `3`, `22`, `32`, `222`, `33`, … |
| `ac` | int | AC value (distinct pairwise differences − (count − 1)) |
| `span` | int | Max − min |
| `mod3` | int[3] | Count with remainder 0, 1, 2 when divided by 3 |
| `tailSum` | int | Sum of last digits |
| `sameTailGroups` | int | Number of last digits that appear 2+ times |
| `maxGap`, `minGap` | int | Largest / smallest gap between adjacent numbers |

Do not assume a fixed `tenGroup` length or a closed set of `consecutive` values. Ignore unknown JSON fields so new backend fields do not break old app versions.

`DrawPage`: `{ items: Draw[], total, limit, offset }`. More pages exist while `offset + items.length < total`.

`NumberStats`: `{ window, totalDraws, latestPeriod, expectedWindow, numbers: NumberStat[maxNumber], specialExpectedWindow?, special?: NumberStat[special.maxNumber] }`, where `NumberStat` is `{ number, windowCount, totalCount, omission, maxOmission }`. `expectedWindow` / `specialExpectedWindow` are the average counts under pure chance; use them as the baseline line when drawing hot/cold bars.

`AnalyzeRequest`

| Game | Body |
|---|---|
| ftn | `{"numbers":[3,7,15,17,30]}` |
| bl | `{"numbers":[3,9,11,25,34,41]}` (sending `special` is a 400) |
| pw | `{"numbers":[4,23,27,28,29,38],"special":5}` (`special` required) |

`AnalyzeResult`: `{ numbers, special?, features, history: MatchCount[], totalDraws, lastMatch3?: Draw }`.

`MatchCount` is `{ main, special, count, prize? }`: how many past draws matched `main` main numbers, with (`special: true`) or without the special number. `prize` is the prize name for that combination (e.g. `貳獎`) and is omitted when it wins nothing. Entries are ordered by `main` descending, `special: true` first; impossible combinations (bl 6 + 特別號) are omitted. Render the list as-is rather than rebuilding prize logic in the app.

## Screens and navigation

The app supports **all three games**. A menu or screen that only offers 今彩539 is a bug.

**Game menu**

- The main menu (drawer, bottom navigation, tab bar or segmented control, whichever the app uses) has one entry per game, in the order `GET /games` returns them: 今彩539, 大樂透, 威力彩.
- Build the entries from `/games` (`code` for routing, `name` for the label). Do not hardcode a list containing only `ftn`.
- If `/games` fails and there is no cached copy, fall back to the built-in list `ftn 今彩539`, `bl 大樂透`, `pw 威力彩` so the menu never shows a single game or an empty state.
- Remember the last selected game across launches (local preference) and open it on start.

**Per-game screens.** Every game gets the same three screens, routed by game code (e.g. `/{game}/draws`), sharing one implementation:

| Screen | Data | 大樂透 / 威力彩 differences |
|---|---|---|
| 開獎紀錄 (draw history, newest first, infinite scroll) | `/{game}/draws` | Show `special` as an extra ball in a different colour, labelled with `special.name` (特別號 / 第二區) |
| 號碼統計 (hot/cold, omission) | `/{game}/numbers/stats` | Grid sized by `maxNumber` (39 / 49 / 38); add a second section or tab for `special` stats (49 / 8 numbers) |
| 選號分析 (pick numbers, show features and history matches) | `POST /{game}/analyze` | Picker sized by `maxNumber`, select exactly `pickCount`; for 威力彩 (`special.userPick`) add a separate 第二區 picker 1–8; for 大樂透 no special picker. Show `history` rows with their `prize` labels |

Switching games from the menu keeps the user on the same kind of screen (e.g. 號碼統計 for 539 → 號碼統計 for 大樂透) and clears any half-made pick, since ranges differ.

## Implementation rules

1. **One API client layer.** A single client/repository owns the base URL, headers, timeouts (connect 10 s, read 15 s), envelope decoding and error mapping. UI code calls repository methods that return domain models, never raw JSON. Pass the game code as a parameter or enum; do not write three copies of each call.
2. **Rules come from `/games`.** Load once at startup (cache it; it changes rarely) and drive pickers, validation and labels from it.
3. **Validate before POST /analyze.** Exactly `pickCount` distinct integers in 1–`maxNumber`; for `userPick` games also a special number in 1–`special.maxNumber`; for others no special. The server rejects anything else with 40001.
4. **Paging.** Use `limit`/`offset` with infinite scroll; default page size 20. Stop when `offset >= total`.
5. **Caching.** The server refreshes its data once a day, so cache draws and stats for at least several minutes and serve cached data while refreshing. Use `/health.games[code].latestPeriod` as a cheap "has anything changed" check.
6. **Retries.** Retry only idempotent GETs, only on network errors / 503, with exponential backoff (1 s, 2 s, 4 s, max 3 tries). Do not retry 4xx.
7. **Wording.** Statistics are descriptive history. The app must not present them as predictions or as raising the chance of winning; every pick has the same odds. Use labels like "近 100 期出現次數" or "目前遺漏期數", not "推薦號碼" or "必中".
8. **Tests.** Unit-test JSON decoding with the sample payloads in `references/openapi.yaml` for all three games (including absent `special` for ftn), plus one error envelope (40001) and one 503.

## Verifying against a running backend

Start the server from the Go repo root (`go run .`), then:

```bash
curl -s http://localhost:8080/api/v1/health
curl -s http://localhost:8080/api/v1/games
curl -s 'http://localhost:8080/api/v1/bl/draws?limit=2'
curl -s 'http://localhost:8080/api/v1/pw/numbers/stats?window=100'
curl -s -X POST http://localhost:8080/api/v1/pw/analyze \
  -H 'Content-Type: application/json' -d '{"numbers":[4,23,27,28,29,38],"special":5}'
# add  -H 'X-API-Key: <key>'  when server.api_key is set
```

If the app and curl disagree, compare the request URL, headers and body first; the backend logs every request line (gin logger) in its console.
