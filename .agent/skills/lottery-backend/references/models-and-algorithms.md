# Models, Feature Algorithms, and Legacy Backtesting Reference

## 1. CSV Data Formats (`./data/csv/`)

Each game's yearly HTML table from `www.nfd.com.tw` is stripped of non-digits per cell (`[^0-9]+`) and written to CSV:

| Game | Config Index | Directory | Columns (`[]string`) |
|---|---|---|---|
| 今彩539 (`ftn`) | `df.InfoFTN` (0) | `data/csv/39-year/39-YYYY.csv` | `[0] Year`, `[1] MonthDay (MMDD)`, `[2] LIdx (YearIndex)`, `[3..7] B1..B5`, `[8] TIdx (TotalPeriod)` (9 cols) |
| 大樂透 (`bl`) | `df.Info49` (1) | `data/csv/49-year/49-YYYY.csv` | `[0] Year`, `[1] MonthDay`, `[2] LIdx`, `[3..8] B1..B6`, `[9] B7 (特別號)`, `[10] TIdx` (11 cols) |
| 威力彩 (`pw`) | `df.InfoPOWER` (2) | `data/csv/power-38/YYYY.csv` | `[0] Year`, `[1] MonthDay`, `[2] LIdx`, `[3..8] B1..B6`, `[9] S1 (第二區)`, `[10] TIdx` (11 cols) |
| 四星彩 (`star4`) | `df.Info4STAR` (3) | `data/csv/4-star/YYYY.csv` | Two draws per row: `[0] Date1`, `[1] YearIdx1`, `[2] Balls1`, `[3] Space`, `[4] Date2`, `[5] YearIdx2`, `[6] Balls2` (7 cols) |

---

## 2. Number Feature Calculations (`model/df/`)

`df.NewFeature(numbers []int, ballsCount int)` and `model/df/extra.go` compute statistical properties for a sorted slice of main numbers:

- **`SUM`**: Sum of main numbers.
- **`OddNumberCount` / `EvenNumberCount`**: Count of odd and even numbers.
- **`Big` / `Small`**: Numbers `>= game.BigFrom` vs `< game.BigFrom` (`BigFrom` is `20` for FTN/PW, `25` for BL).
- **`PrimeCount`**: Count of prime numbers (`2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47`).
- **`TenGroupCount`**: Bucket counts by decade `(n-1)/10`:
  - Index 0: `1–10`, Index 1: `11–20`, Index 2: `21–30`, Index 3: `31–40` (or `31–39`/`31–38`), Index 4: `41–49`.
- **`TailDigit`**: Counts of `(n-1)%10` in `df.Feature` (legacy internal), whereas `extra.go` provides:
  - `TailSum()`: Sum of `b % 10` for each ball.
  - `SameTailGroups()`: Number of last digits (`0–9`) that appear at least twice.
- **`ConsecutivePattern()` (`extra.go`)**:
  - Finds all contiguous runs of adjacent numbers (`b[i] - b[i-1] == 1`) of length $\ge 2$.
  - Sorts run lengths descending and concatenates them as a string (e.g., `"none"`, `"2"`, `"3"`, `"22"`, `"32"`, `"222"`).
- **`AC()` (Arithmetic Complexity)**:
  - Number of distinct positive pairwise differences between numbers minus `(pickCount - 1)`.
- **`Span()`**: `max - min` (`IBalls[len-1] - IBalls[0]`).
- **`Mod3Counts()`**: Array `[3]int` counting numbers with `b % 3 == 0, 1, 2` (012路).
- **`MaxGap()` / `MinGap()`**: Maximum and minimum difference between adjacent sorted numbers `b[i+1] - b[i]`.

### Inter-Draw Transition Features (`model/df/transition.go`)

`df.NewTransition(current []int, history [][]int, maxNumber int, coldThreshold int) *df.Transition` computes inter-draw relationships between `current` and preceding draws (`history`, ordered newest-to-oldest, `history[0]` = $t-1$). Exposed as `transition` on `Draw` and `AnalyzeResult` in `/api/v1`, and usable offline via `(*FTN).Transition(history)` and `FTNArray.FilterByTransition(history, TransitionFilter)` in `model/ftn/transition.go`:

- **Repeat / Neighbor / Isolated**: `repeatCount`, `repeatNumbers`, `skip1RepeatCount` ($t-2$ repeat not in $t-1$), `neighborCount` (non-repeat adjacent $\pm 1$, satisfying `repeat + neighbor + isolated == pickCount`), `neighborPlus`, `neighborMinus`, `isolatedCount`, `maxRepeatStreak`, `maxDiagonalRun`.
- **Geometry & Position**: `sandwichCount` (prev had $x-1, x+1$, curr hits $x$), `splitCount` (prev had $x$, curr hits $x-1, x+1$), `sameTailPrevCount`, `positionDiffs`.
- **Structure Deltas**: `sumDiff`, `spanDiff`, `gapZoneRefill` (count of current balls landing in previous draw's empty ten-groups).
- **Multi-Draw Hot/Cold & Follow-up**: `recent3Count`, `recent5Count`, `omissionSum`, `maxOmissionHit`, `coldHitCount` (omission $\ge 15$ by default), `followUpScore` (average historical `prevBall -> currBall` co-occurrence count).

---

## 3. Monte Carlo Simulation & Filtering (`model/ftn/` and `model/pw/`)

In addition to the REST API, `model/ftn` and `model/pw` contain an offline analysis and backtesting engine invoked via `go run . ftn -a pre` or test suites:

1. **Full Combinations**:
   - Uses `gonum.org/v1/gonum/stat/combin.Combinations(N, K)` to generate all $C(39, 5) = 575,757$ (FTN) or $C(38, 6) = 2,760,681$ (PW) combinations.
   - Note: `NewFTNWithInts` / `NewPowerWithInts` converts 0-based combination indices from `combin.Combinations` into 1-based ball numbers (`arr[i] + 1`).
2. **Monte Carlo Sampling (`DoGenJSON`)**:
   - Seeds either standard `math/rand` (`Randomer: 0`) or Mersenne Twister `mt19937` (`Randomer: 1`) with `crypto/rand`.
   - Draws `Sample * SampleTime` random combinations and selects combinations whose hit count exceeds `Threshold.Value` (`ThresholdNumbers`).
3. **Feature & Group Filtering (`arrfilter.go`, `group.go`)**:
   - `FilterFeatureIncludes` / `FilterFeatureExcludes`: Keeps or excludes combinations matching historical `df.Feature` signatures (configurable via `df.DisableFilters(...)`).
   - `FilterByGroupIndex`: Partitions combinations into fixed-size groups (e.g., 200 per group) and filters by historical hit density.
   - `FilterIncludes` / `FilterExcludes`: Filters combinations containing or excluding specific balls from recent draws.
   - `FilterCol` / `FilterNeighber`: Filters by repeat count (`haveCol`) or adjacent neighbor count (`haveNeighber`, `|b1 - b2| == 1`) against a target draw.
4. **Global State Caveat**:
   - `df.filters` (`model/df/df.go`) and `ballPeriotStatic` / `ballContinueStatic` (`model/ftn/ftn.go`) are package-level variables. `Store.Reload()` in `api/lottery/store.go` protects `load()` with `s.loadMu` to prevent data races. Avoid mutating `df.DisableFilters` concurrently in HTTP handlers.
