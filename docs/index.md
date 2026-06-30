# go-ruby-time documentation

**Ruby's Time class and the time stdlib's parsing extensions in pure Go — MRI byte-exact, no cgo.**

`go-ruby-time/time` is a faithful, pure-Go (zero cgo) reimplementation of Ruby's `Time`,
matching reference Ruby (MRI) byte-for-byte. The module path is
`github.com/go-ruby-time/time`.

It is a **standalone, reusable** library importable by any Go program, and the
backend bound into [go-embedded-ruby](https://github.com/go-embedded-ruby/ruby)
by `rbgo` as a native module — the same pattern as
[go-ruby-yaml](https://github.com/go-ruby-yaml/yaml). The dependency runs the
other way: this library has **no dependency on the Ruby runtime**.

!!! success "Status: complete — MRI byte-exact"
    A faithful pure-Go port of Ruby's `Time`, validated by a **differential oracle**
    against the system `ruby` — results compared byte-for-byte — at 100%
    coverage, `gofmt` + `go vet` clean, CI green across the six 64-bit Go targets
    and three OSes.

## Quick taste

```go
t := rtime.UTC(2026, 6, 29, 5, 18, 32)

t.Inspect()              // 2026-06-29 05:18:32 UTC
t.Strftime("%A, %d %B %Y") // Monday, 29 June 2026
t.ISO8601()              // 2026-06-29T05:18:32Z
t.RFC2822()              // Mon, 29 Jun 2026 05:18:32 -0000
t.Add(90).Inspect()     // 2026-06-29 05:20:02 UTC
```

## Repositories

| Repo | What it is |
| --- | --- |
| [`time`](https://github.com/go-ruby-time/time) | the library — Time in pure Go |
| [`docs`](https://github.com/go-ruby-time/docs) | this documentation site (MkDocs Material, versioned with mike) |
| [`go-ruby-time.github.io`](https://github.com/go-ruby-time/go-ruby-time.github.io) | the organization landing page (Hugo) |
| [`brand`](https://github.com/go-ruby-time/brand) | logo and brand assets |

## Principles

- **Pure Go, `CGO_ENABLED=0`** — trivial cross-compilation, a single static
  binary, no C toolchain.
- **MRI byte-exact.** Output matches reference Ruby exactly, not approximately,
  validated by a differential oracle against the `ruby` binary.
- **Standalone & reusable.** Extracted from rbgo's internals; no dependency on
  the Ruby runtime — the dependency runs the other way.
- **100% test coverage** is the target, enforced as a CI gate, across 6 arches
  and 3 OSes.

## Where to go next

- [Why pure Go](why.md) — why this slice of Ruby is deterministic enough to live
  as a standalone, interpreter-independent Go library.
- [Usage & API](api.md) — the public surface and worked examples.
- [Roadmap](roadmap.md) — what is done and what is downstream by design.

Source lives at [github.com/go-ruby-time/time](https://github.com/go-ruby-time/time).
