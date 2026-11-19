---
name: lottery-backend
description: >-
  Develop, maintain, test, and debug the Go (Gin + Cobra) Taiwan Lottery backend service (`module lottery`),
  including data scraping/CSV ingestion, in-memory store & deduplication, number feature algorithms,
  backtesting & filtering models (今彩539 ftn, 大樂透 bl, 威力彩 pw, 四星彩 star4), REST API endpoints (/api/v1),
  API Key & HMAC Premium Token authentication, and CLI commands. Use whenever working on Go source code,
  tests, configuration, or data pipelines in this repository.
---

# Lottery Backend Skill (`module lottery`)

This skill provides the architecture, domain rules, data flow, and development runbook for the Taiwan Lottery Go backend project.

For deep-dive details on legacy statistical models, Monte Carlo backtesting, filter chains, and CSV column layouts, read [references/models-and-algorithms.md](./references/models-and-algorithms.md).
For mobile client integration contracts, see the companion skill [lottery-mobile-api](../lottery-mobile-api/SKILL.md).

---

## 1. Project Overview & Tech Stack

- **Module**: `lottery` (Go 1.24)
- **Web Framework**: `github.com/gin-gonic/gin`
- **CLI Framework**: `github.com/spf13/cobra`
- **HTML Scraping**: `github.com/PuerkitoBio/goquery`
- **Math / Combinations / PRNG**: `gonum.org/v1/gonum/stat/combin`, `github.com/goark/mt/v2` (Mersenne Twister `mt19937`)
- **Storage**: Local CSV files (`./data/csv/`), JSON backtest snapshots (`./gendata/`), optional SQLite cache (`cache/sqllite.go`)
- **Configuration**: `config.yaml` loaded via `config.LoadConfig("./config.yaml")` into global `config.Config`

---

## 2. Directory & Package Architecture

| Package / Path | Responsibility |
|---|---|
| `main.go` | Entry point: runs `cmd.RootCmd`, initializes `lotteryapi.Store` for each game (`ftn`, `bl`, `pw`), registers daily reload hook, starts Gin HTTP server on `config.Config.Server.Addr` (default `:8080`). |
| `cmd/` | Cobra CLI commands: `root.go` (starts server + daily 10:00 AM scraper goroutine), `ftn.go` (`lot ftn` Monte Carlo backtest/pick CLI), `premium.go` (`lot premium` token issuer CLI). |
| `config/` | Parses `config.yaml` (`server`, `http`, `local`, `aokurls`, `combinations`). Note: `LoadConfig` calls `os.Exit(1)` on error. |
| `api/` | Shared Gin helpers (`OK`, `Fail`, error codes), `APIKeyAuth` middleware (`X-API-Key`), and HMAC-SHA256 `PremiumAuth` middleware (`X-Premium-Token`) + `/api/v1/me` handler. |
| `api/lottery/` | REST API v1 implementation (`game.go`, `store.go`, `handler.go`, `dto.go`, `handler_test.go`). |
| `model/df/` | Core domain features (`Feature` in `df.go`, extended metrics in `extra.go`), ball metadata (`IBall`), prime tables, ten-group & consecutive pattern logic. |
| `model/common/` | HTML table scraper (`GetAll()`), `.doc` downloader (`GetAllFromURL()`), CSV reader (`ReadCSV()`), and PRNG helpers (`SetRandomGenerator`, `RandomNuber`). |
| `model/ftn/` | 今彩539 (5 of 1–39) model, `LoadHistory()`, `FTNsManager`, filter chains (`arrfilter.go`), ten-group analysis, and backtesting (`backtest.go`). |
| `model/bl/` | 大樂透 (6 of 1–49 + 1 特別號) model, `LoadHistory()`, `BLManager`, statistics. |
| `model/pw/` | 威力彩 (6 of 1–38 + 第二區 1 of 1–8) model, `LoadHistory()`, `PowerManager`, group & backtesting. |
| `model/star4/` | 四星彩 (4 digits 0000–9999) model, permutations, and frequency statistics (CLI/model only; not exposed in `/api/v1/games`). |
| `csv/` | Reads/writes yearly CSV files under `./data/csv/<game-path>/<label>.csv`. |
| `algorithm/` | Generic bitmask power-set and $C(n, k)$ combination generators (`All[T]`, `Combinations[T]`). |
| `cache/` | Standalone SQLite wrapper (`NewLotteryDB`) for `record` caching. |

---

## 3. Data Pipeline & In-Memory Store (`api/lottery/store.go`)

### Scraping & CSV Storage
1. **Source**: `common.GetAll()` scrapes yearly HTML tables from `http://www.nfd.com.tw/lottery/<path>/<label>.htm` configured in `config.yaml`:
   - `39-year` (`39-%d.htm`, baseyear 2007) → 今彩539 (`df.InfoFTN`)
   - `49-year` (`49-%d.htm`, baseyear 2004) → 大樂透 (`df.Info49`)
   - `power-38` (`%d.htm`, baseyear 2008) → 威力彩 (`df.InfoPOWER`)
   - `4-star` (`%d.htm`, baseyear 2003) → 四星彩 (`df.Info4STAR`)
2. **Caching rule**: Past years (`year < time.Now().Year()`) are skipped if the local CSV file already exists in `./data/csv/`. The current year's CSV is overwritten on every `GetAll()` run.
3. **Daily Cron**: `cmd/root.go` runs `checkRawData()` in a background goroutine that triggers every day at **10:00 AM local time**, calls `common.GetAll()`, and invokes `cmd.OnRawDataUpdated(reloadAll)` to hot-reload all `lotteryapi.Store` instances.

### Store Cleaning & Deduplication (`Store.clean`)
When `Store.Reload()` loads records from `ftn.LoadHistory()`, `bl.LoadHistory()`, or `pw.LoadHistory()`:
- **Concurrency**: `s.loadMu` serializes `Reload()` calls because `model/ftn` uses package-level maps (`ballPeriotStatic`, `ballContinueStatic`) during CSV parsing. Read requests use `s.mu.RLock()` via `s.snapshot()`.
- **Validation**: `g.validRecord(&r)` filters out malformed dates, out-of-range numbers, duplicate numbers within a draw, or invalid special numbers.
- **Period normalization**: Strips leading zeros from `Period` (`"037"` → `"37"`).
- **Sorting**: Main numbers are sorted ascending; records are sorted by integer `Period` **descending (newest first)**.
- **Deduplication (`dedupPeriods`)**: Raw upstream data occasionally has typos or lunar new year bonus draws sharing a period number. `dedupPeriods` checks adjacent unique periods' dates (`lo <= r.Date() <= hi`):
  - If all records matching the date window are identical, keeps one.
  - If one record fits the date window between neighboring periods, keeps it and drops the typo.
  - If ambiguous, drops the conflicting period records.

---

## 4. REST API (`/api/v1`) & Auth Rules

### Authentication & Entitlements
Configured under `server:` in `config.yaml`:
1. **`X-API-Key` (`api.APIKeyAuth`)**:
   - Checked on all `/api/v1/*` routes **except** `GET /api/v1/health`.
   - If `server.api_key` is `""`, check is bypassed. Otherwise uses constant-time comparison; returns `401` (`code: 40101`) on mismatch.
2. **`X-Premium-Token` (`api.PremiumAuth`)**:
   - Format: `v1.<userID>.<expiresUnixSec>.<HMAC-SHA256-hex>` signed with `server.premium_secret`.
   - Missing header → request proceeds as **free user** (`api.IsPremium(c) == false`).
   - Present but invalid/expired (or `premium_secret` is empty) → returns `401` (`code: 40102`).
3. **Free Tier Truncation (`server.free_draw_limit`)**:
   - When `free_draw_limit > 0` and user is not premium:
     - `h.records(c)` slices history to `records[:freeDrawLimit]` (newest `N` draws) and sets `"limited": true` in responses.
     - `GET /{game}/draws/:period` for a period older than `freeDrawLimit` returns `403` (`code: 40301`).
     - `GET /{game}/numbers/stats?window=W` with `W > freeDrawLimit` returns `403` (`code: 40301`).
     - `POST /{game}/analyze` only matches against the newest `freeDrawLimit` draws.

### Error Codes (`api/common.go`)
| Constant | Code | HTTP Status | Meaning |
|---|---|---|---|
| `CodeOK` | `0` | 200 | Success |
| `CodeInvalidParam` | `40001` | 400 | Invalid query/body parameter |
| `CodeUnauthorized` | `40101` | 401 | Missing or invalid `X-API-Key` |
| `CodePremiumTokenInvalid` | `40102` | 401 | Invalid or expired `X-Premium-Token` |
| `CodePremiumRequired` | `40301` | 403 | Requested data exceeds free tier limit |
| `CodeNotFound` | `40401` | 404 | Draw period not found |
| `CodeDataNotReady` | `50301` | 503 | Store has 0 loaded records |

---

## 5. Common Development Workflows

### Run the Server Locally
```bash
go run .
```
*(Requires `./config.yaml` in the current working directory and CSV data in `./data/csv/`).*

### Fetch / Update Lottery CSV Data
Run the scraper test or trigger `common.GetAll()`:
```bash
go test -v -run Test_getAll .
```
*(Note: Requires network access to `www.nfd.com.tw` and `www.aok.url.tw`).*

### Issue a Premium Token via CLI
Make sure `server.premium_secret` is set in `config.yaml`, then run:
```bash
go run . premium -u <user_id> -d <days>
```

### Run Unit Tests
Run the API and package unit tests (sandbox-safe, uses in-memory test stores):
```bash
go test ./api/... ./config/... ./csv/... ./model/df/...
```
> **Important Note on Tests**:
> - `api/lottery/handler_test.go`, `api/premium_test.go`, and `model/df/extra_test.go` are self-contained unit tests that run fast without network or external files.
> - `getall_test.go` makes live HTTP requests to scrape lottery sites.
> - Some legacy tests in `model/ftn/`, `model/bl/`, `model/pw/`, and `model/common/` expect specific local paths, `config.yaml` in the working directory, or long-running Monte Carlo simulations.

### Adding a New Game or Modifying Features
1. **New Feature Metric**:
   - Implement the calculation method on `*df.Feature` in `model/df/extra.go` (assuming `f.IBalls` is sorted ascending).
   - Add the field to `lotteryapi.Features` and `g.newFeatures(numbers)` in `api/lottery/dto.go`.
   - Add unit tests in `model/df/extra_test.go` and `api/lottery/handler_test.go`.
   - Update `.agent/skills/lottery-mobile-api/references/openapi.yaml` if mobile clients consume it.
2. **New Game in API**:
   - Define the `Game` constructor in `api/lottery/game.go` (specify `Code`, `Name`, `MaxNumber`, `PickCount`, `BigFrom`, `Special`, `Prizes`, and `load` function returning `[]Record`).
   - Add it to `Games()` in `api/lottery/game.go` — routes and `/health` / `/games` endpoints register automatically.
