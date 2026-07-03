<!-- SPDX-License-Identifier: BSD-3-Clause -->
# `go-ruby-time` library-level benchmark harness

Reproducible, cross-runtime benchmark of the **pure-Go `go-ruby-time` library**
against the reference Ruby runtimes (MRI, MRI + YJIT, JRuby, TruffleRuby). It
measures the **library primitive** through its Go API, isolated from any
interpreter, so the numbers answer: *is the pure-Go implementation as fast as the
reference runtime's own `Time`?*

## Layout

- `go/`           — self-contained Go driver; `go.mod` pins the published library
  by pseudo-version (no `replace`). The built `go/bench` binary is git-ignored.
- `ruby/time.rb`  — the equivalent workload; `ruby/_harness.rb` is the shared timer.
- `run.sh`        — runs every available runtime and prints one Markdown table per
  sub-benchmark (ns/op + ratio vs MRI).

## Run

```sh
bash benchmarks/run.sh
```

Environment knobs: `OUTER` (timed passes, default 25), `WARM` (untimed warm-up
passes, default 3), and `RUBY`/`JRUBY`/`TRUFFLERUBY` to select runtime binaries.

## Verify (outputs identical to MRI)

Both drivers accept a `verify` argument that prints each operation's canonical
output for a **fixed** input instant (never the wall clock), so correctness is
checked before any timing is trusted:

```sh
diff <(cd go && go run . verify) <(ruby ruby/time.rb verify)   # must be empty
```

## Method

Each process runs `WARM` untimed passes (to let the JVM/GraalVM JITs warm up),
then `OUTER` timed passes of a fixed inner loop, timed with a monotonic clock;
the **best** pass is reported as **ns/op**. Interpreter start-up is outside the
timed region. The Go driver and the Ruby script build **identical fixed inputs**
(an ISO-8601 and an RFC-2822 string, plus one fixed UTC instant) and their
outputs are checked identical to MRI before timing. Results are published, dated,
in `../docs/performance.md`.
