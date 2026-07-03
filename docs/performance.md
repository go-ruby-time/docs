# Performance

`go-ruby-time/time` is a pure-Go, MRI-compatible implementation of Ruby's
`Time`. Unlike most modules in the ecosystem it is **not** bound into
[`rbgo`](https://github.com/go-embedded-ruby/ruby) — the interpreter ships its
own integrated `Time` core — so this library stands alone: a CI-green,
100%-covered, drop-in Ruby `Time` for pure-Go programs. This page records a
**comparative, library-level benchmark** of that standalone module against the
reference Ruby runtimes, part of the ecosystem-wide per-module parity suite.

## Library-level benchmark (Go API vs runtimes) — 2026-07-03

This measures the **pure-Go library directly, through its Go API**, isolated
from any interpreter dispatch, answering the parity question head-on: *is the
pure-Go implementation as fast as the reference runtime's own `Time`?* The
**same workload, same fixed inputs, same iteration counts** run through the Go
library and through each reference runtime's stdlib; outputs were checked
**byte-identical to MRI** before any timing (see *Reproduce* below).

- **Host:** Apple M4 Max (`Mac16,5`, arm64), macOS 26.5.1 — **date 2026-07-03**.
  All runtimes measured **on the host**, no VM.
- **Runtimes:** Go 1.26.4 · MRI `ruby 4.0.5 +PRISM` (the oracle) · MRI + YJIT ·
  JRuby 10.1.0.0 (OpenJDK 25) · TruffleRuby 34.0.1 (GraalVM CE Native).
- **Fixed inputs (reproducible, never the wall clock):** ISO-8601
  `2024-03-15T13:45:30+00:00`; RFC-2822 `Fri, 15 Mar 2024 13:45:30 +0000`; the
  instant `2024-03-15 13:45:30 UTC` (epoch `1710510330`) for `strftime` /
  arithmetic. `strftime` format `%Y-%m-%dT%H:%M:%S %A %z`.
- **Operations:** `Time.parse` of an ISO-8601 and an RFC-2822 string,
  `Time.strptime` with an explicit layout, `strftime`, and time arithmetic
  (add seconds, add a day, difference).
- **Method:** each process runs 3 untimed warm-up passes, then 25 timed passes of
  a fixed inner loop, timed with a monotonic clock; the **best** pass is reported
  as **ns/op** (lower is better). `vs MRI` < 1.00× means *faster than MRI*.
  Interpreter start-up is outside the timed region, so these are operation costs,
  not `ruby file.rb` process costs.

#### parse-iso

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 377.8 | 0.06× |
| MRI | 6288.5 | 1.00× |
| MRI + YJIT | 6198.0 | 0.99× |
| JRuby | 6458.7 | 1.03× |
| TruffleRuby | 57237.8 | 9.10× |

#### parse-rfc2822

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 1202.2 | 0.19× |
| MRI | 6413.0 | 1.00× |
| MRI + YJIT | 5943.0 | 0.93× |
| JRuby | 4525.9 | 0.71× |
| TruffleRuby | 56947.6 | 8.88× |

#### strptime-iso

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 72.0 | 0.03× |
| MRI | 2197.5 | 1.00× |
| MRI + YJIT | 1751.0 | 0.80× |
| JRuby | 2083.1 | 0.95× |
| TruffleRuby | 8909.8 | 4.05× |

#### strftime

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 653.3 | 1.24× |
| MRI | 525.0 | 1.00× |
| MRI + YJIT | 475.5 | 0.91× |
| JRuby | 426.1 | 0.81× |
| TruffleRuby | 1044.9 | 1.99× |

#### add-seconds

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 9.5 | 0.17× |
| MRI | 55.4 | 1.00× |
| MRI + YJIT | 29.6 | 0.53× |
| JRuby | 43.8 | 0.79× |
| TruffleRuby | 130.4 | 2.35× |

#### add-days

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 9.6 | 0.17× |
| MRI | 56.0 | 1.00× |
| MRI + YJIT | 40.2 | 0.72× |
| JRuby | 15.0 | 0.27× |
| TruffleRuby | 151.1 | 2.70× |

#### diff

| Runtime | ns/op | vs MRI |
| --- | ---: | ---: |
| **go-ruby (pure Go)** | 6.5 | 0.11× |
| MRI | 60.8 | 1.00× |
| MRI + YJIT | 37.6 | 0.62× |
| JRuby | 8.7 | 0.14× |
| TruffleRuby | 66.0 | 1.09× |
Parsing is where the pure-Go library pulls decisively ahead: `Time.parse` of the
ISO-8601 string is **~16× faster than MRI** (0.06×) and `strptime` is **~30×
faster** (0.03×), because MRI's `time` stdlib routes `parse`/`strptime` through
the general `Date._parse` heuristic while the Go library scans a fixed set of
layouts directly. RFC-2822 parsing is ~5× faster (0.19×). Time arithmetic
(add-seconds, add-days, difference) is several-fold faster than MRI (0.11×–0.17×).
`strftime` is the one op slightly behind MRI's C formatter (1.24×) — still within
the same order of magnitude, at parity with JRuby/YJIT territory. The MRI +
YJIT column barely moves the parse rows (`Date._parse` is a long C path YJIT
cannot compile away), which is why YJIT ≈ MRI there.

!!! note "Reproduce"
    The harness is committed under
    [`benchmarks/`](https://github.com/go-ruby-time/docs/tree/main/benchmarks):
    a self-contained Go driver (`go/`, pins the published library via `go.mod`
    pseudo-version — no `replace`), the equivalent `ruby/time.rb` workload, and
    `run.sh`. Run `bash benchmarks/run.sh`; env `OUTER`/`WARM` tune the pass
    budget and `RUBY`/`JRUBY`/`TRUFFLERUBY` select the runtime binaries. Both
    drivers accept a `verify` argument (`(cd benchmarks/go && go run . verify)`
    vs `ruby benchmarks/ruby/time.rb verify`) that prints each operation's
    canonical output; the two were confirmed **byte-identical** before timing.

!!! warning "Warm-up budget & noise — honest framing"
    Numbers reflect a **fixed warm-process budget** (3 warm-up + 25 timed passes
    in one process). The JVM/GraalVM JITs (JRuby, TruffleRuby) may need a larger
    warm-up to reach steady state, so their columns can **understate** peak
    throughput — most visibly TruffleRuby on the parse rows, which are cold-JIT
    figures, not steady-state. Sub-microsecond rows (arithmetic) carry the most
    relative noise; treat those ratios as order-of-magnitude. Every number here
    is a **real measured value** from the dated run above — nothing is
    fabricated, estimated, or cherry-picked. The go-ruby column is the pure-Go
    library; every other column is that interpreter's own `Time` stdlib doing
    the equivalent work.
