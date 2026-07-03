// SPDX-License-Identifier: BSD-3-Clause
package main

import (
	"fmt"
	"os"

	rbtime "github.com/go-ruby-time/time"
)

// Fixed, reproducible inputs — never the wall clock — so every run is identical
// and the Go outputs can be checked byte-for-byte against MRI.
const (
	iso       = "2024-03-15T13:45:30+00:00"    // ISO-8601, explicit +00:00
	rfc2822   = "Fri, 15 Mar 2024 13:45:30 +0000" // RFC-2822
	isoLayout = "%Y-%m-%dT%H:%M:%S%z"
	fmtLayout = "%Y-%m-%dT%H:%M:%S %A %z"
)

// fixed builds the one fixed instant (2024-03-15 13:45:30 UTC) all the
// formatting / arithmetic benchmarks operate on, with a second instant one day
// later for the difference benchmark.
func fixed() (*rbtime.Time, *rbtime.Time) {
	t := rbtime.UTC(2024, 3, 15, 13, 45, 30)
	return t, t.Add(86400)
}

// verify prints the canonical output of every benchmarked operation, one per
// line, so it can be diffed against the Ruby side before any timing is trusted.
func verify() {
	pi, _ := rbtime.Parse(iso)
	pr, _ := rbtime.Parse(rfc2822)
	sp, _ := rbtime.Strptime(iso, isoLayout)
	t, t2 := fixed()
	fmt.Printf("parse-iso\t%d\t%d\n", pi.ToI(), pi.UTCOffset())
	fmt.Printf("parse-rfc2822\t%d\t%d\n", pr.ToI(), pr.UTCOffset())
	fmt.Printf("strptime-iso\t%d\t%d\n", sp.ToI(), sp.UTCOffset())
	fmt.Printf("strftime\t%s\n", t.Strftime(fmtLayout))
	fmt.Printf("add-seconds\t%d\n", t.Add(3600).ToI())
	fmt.Printf("add-days\t%d\n", t.Add(86400).ToI())
	fmt.Printf("diff\t%d\n", int64(t2.Diff(t)))
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "verify" {
		verify()
		return
	}
	t, t2 := fixed()
	bench("parse-iso", 2000, func() { v, _ := rbtime.Parse(iso); sink = v })
	bench("parse-rfc2822", 2000, func() { v, _ := rbtime.Parse(rfc2822); sink = v })
	bench("strptime-iso", 2000, func() { v, _ := rbtime.Strptime(iso, isoLayout); sink = v })
	bench("strftime", 2000, func() { sink = t.Strftime(fmtLayout) })
	bench("add-seconds", 5000, func() { sink = t.Add(3600) })
	bench("add-days", 5000, func() { sink = t.Add(86400) })
	bench("diff", 5000, func() { sink = t2.Diff(t) })
}
