# Roadmap

`go-ruby-time/time` is grown **test-first**, each capability differential-tested against
MRI rather than built in isolation. The deterministic, interpreter-independent
slice of Ruby's `Time` extracted from rbgo's internals is **complete**.

| Stage | What | Status |
| --- | --- | --- |
| Constructors | `Now`, `At` / `AtNsec` / `AtTime`, `UTC`, `Local`, `New`, `Parse`, `Strptime`. | **Done** |
| Components & conversions | All civil accessors, sub-second (`USec`/`NSec`/`Subsec`), zone fields, weekday predicates, `to_i`/`to_f`/`to_r`/`to_a`. | **Done** |
| Zone shifts | `UTCTime`, `GetLocal`, `ToLocal` — same instant, different observed offset. | **Done** |
| Arithmetic & ordering | `Add`/`Sub`, `Diff`, `Round`/`Floor`/`Ceil`, `Cmp`, `Before`/`After`, `Equal`, `Hash`. | **Done** |
| Formatting | Full `Strftime` directive set, `Inspect`, `ToS`, `CTime`, `ISO8601`, `RFC2822`, `HTTPDate` — MRI's exact bytes. | **Done** |
| Differential oracle & coverage | Formatters, `strftime` corpus, accessors, arithmetic, `Parse`/`Strptime` checked against `ruby` (TZ=UTC, ≥ 4.0); 100% coverage, green across 6 arches and 3 OSes. | **Done** |

## Documented out-of-scope boundaries

These are **deliberate**, recorded so the module's surface is unambiguous:

- **No interpreter.** The library implements the deterministic algorithm; it
  never runs arbitrary Ruby. Anything that needs a live binding or evaluation is
  the consumer's job — that is why `rbgo` binds this module rather than the
  reverse.
- **Reference is reference Ruby (MRI).** Byte-for-byte conformance targets MRI's
  behaviour; differences across MRI releases are matched to the reference used by
  the differential oracle.
- **Standalone & reusable.** The module has no dependency on the Ruby runtime;
  the dependency runs the other way.

See [Usage & API](api.md) for the surface and [Why pure Go](why.md) for the
deterministic/interpreter split.
