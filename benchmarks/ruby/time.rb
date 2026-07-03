# frozen_string_literal: true
# SPDX-License-Identifier: BSD-3-Clause
require "time"
require_relative "_harness"

# Fixed, reproducible inputs — never Time.now — mirroring benchmarks/go/main.go.
ISO       = "2024-03-15T13:45:30+00:00"      # ISO-8601, explicit +00:00
RFC2822   = "Fri, 15 Mar 2024 13:45:30 +0000" # RFC-2822
FMT       = "%Y-%m-%dT%H:%M:%S %A %z"
T         = Time.at(Time.utc(2024, 3, 15, 13, 45, 30).to_i).utc
T2        = T + 86_400

if ARGV[0] == "verify"
  pi = Time.parse(ISO)
  pr = Time.parse(RFC2822)
  sp = Time.strptime(ISO, "%Y-%m-%dT%H:%M:%S%z")
  printf("parse-iso\t%d\t%d\n",     pi.to_i, pi.utc_offset)
  printf("parse-rfc2822\t%d\t%d\n", pr.to_i, pr.utc_offset)
  printf("strptime-iso\t%d\t%d\n",  sp.to_i, sp.utc_offset)
  printf("strftime\t%s\n",          T.strftime(FMT))
  printf("add-seconds\t%d\n",       (T + 3600).to_i)
  printf("add-days\t%d\n",          (T + 86_400).to_i)
  printf("diff\t%d\n",              (T2 - T).to_i)
  exit
end

bench("parse-iso",     2000) { Time.parse(ISO) }
bench("parse-rfc2822", 2000) { Time.parse(RFC2822) }
bench("strptime-iso",  2000) { Time.strptime(ISO, "%Y-%m-%dT%H:%M:%S%z") }
bench("strftime",      2000) { T.strftime(FMT) }
bench("add-seconds",   5000) { T + 3600 }
bench("add-days",      5000) { T + 86_400 }
bench("diff",          5000) { T2 - T }
