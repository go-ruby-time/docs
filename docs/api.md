# Usage & API

The public API lives at the module root (`github.com/go-ruby-time/time`). It is **Ruby-shaped but Go-idiomatic**: the constructors mirror `Time.at` / `Time.utc` / `Time.parse` and the methods mirror `inspect` / `strftime` / `iso8601`, while the surface follows Go conventions — value types, an explicit `error` on parsing, no global state beyond the pinnable clock seam.

!!! success "Status: implemented"
    The library is built and importable as `github.com/go-ruby-time/time`, bound
    into `rbgo` as a native module; see [Roadmap](roadmap.md).

## Install

```sh
go get github.com/go-ruby-time/time
```

## The clock seam

`Now()` is the only non-deterministic constructor. It reads the wall clock
through the package-level `NowFunc` seam, which a test pins to a fixed instant:

```go
old := rtime.NowFunc
rtime.NowFunc = func() (sec int64, nsec int32) { return 1782710312, 0 }
defer func() { rtime.NowFunc = old }()
// rtime.Now() is now exactly 2026-06-29 05:18:32 …
```

## Worked example

```go
t := rtime.UTC(2026, 6, 29, 5, 18, 32) // Time.utc(2026,6,29,5,18,32)

t.Inspect()               // 2026-06-29 05:18:32 UTC
t.Strftime("%A, %d %B %Y") // Monday, 29 June 2026
t.ISO8601()               // 2026-06-29T05:18:32Z
t.RFC2822()               // Mon, 29 Jun 2026 05:18:32 -0000

off := rtime.New(2026, 6, 29, 5, 18, 32, 2*3600) // Time.new(..., "+02:00")
off.Inspect()             // 2026-06-29 05:18:32 +0200
off.UTCTime().Inspect()   // 2026-06-29 03:18:32 UTC

p, _ := rtime.Parse("2026-06-29T05:18:32+05:30")
p.Inspect()               // 2026-06-29 05:18:32 +0530
```

## Surface

```go
type Time struct{ /* (seconds, nanoseconds, utc-offset, is-utc) */ }

// Clock seam — the only non-determinism.
var NowFunc func() (sec int64, nsec int32)

// Constructors
func Now() *Time
func At(sec float64) *Time; func AtNsec(sec, nsec int64) *Time; func AtTime(t *Time) *Time
func UTC(year int, rest ...int) *Time   // Time.utc / gm
func Local(year int, rest ...int) *Time // Time.local / mktime
func New(year, month, day, hour, min, sec, offset int) *Time
func Parse(input string) (*Time, error); func Strptime(input, layout string) (*Time, error)

// Components, conversions, zone shifts, arithmetic, ordering, formatting
func (t *Time) Year() int; func (t *Time) Month() int /* … */
func (t *Time) ToI() int64; func (t *Time) ToF() float64; func (t *Time) ToR() *big.Rat
func (t *Time) UTCTime() *Time; func (t *Time) GetLocal(offset int) *Time
func (t *Time) Add(secs float64) *Time; func (t *Time) Diff(other *Time) float64
func (t *Time) Cmp(other *Time) int; func (t *Time) Equal(other *Time) bool
func (t *Time) Strftime(format string) string
func (t *Time) Inspect() string; func (t *Time) ISO8601(fraction ...int) string
func (t *Time) RFC2822() string; func (t *Time) HTTPDate() string; func (t *Time) CTime() string
```

The byte-level subtleties are matched on purpose: `inspect` shows a trimmed
sub-second fraction where `to_s` never does, and a UTC `rfc2822` zone renders as
`-0000`.

## MRI conformance

A **differential oracle** runs every formatter, the full `strftime` directive
corpus, the numeric accessors, arithmetic, and `Parse` / `Strptime` through both
the system `ruby` and this library and compares them **byte-for-byte**. The
oracle runs `ruby` with `TZ=UTC`, is gated on `RUBY_VERSION >= "4.0"`, and skips
itself where `ruby` is absent; `TestMain` pins the process zone to UTC so the
deterministic tests are reproducible on any runner.
