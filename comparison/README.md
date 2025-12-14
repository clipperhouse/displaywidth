## Compatibility

In real-world text, you mostly should see the same outputs from
`clipperhouse/displaywidth`,`mattn/go-runewidth`, and `rivo/uniseg`. It's
mostly the same data and logic.

The tests in this package exercise the behaviors of the three libraries.
Extensive details are available in the
[compatibility analysis](COMPATIBILITY_ANALYSIS.md).

## Benchmarks

```bash
go test -bench=. -benchmem
```

```
goos: darwin
goarch: arm64
pkg: github.com/clipperhouse/displaywidth/comparison
cpu: Apple M2

BenchmarkString_Mixed/clipperhouse/displaywidth-8             10400 ns/op	  162.21 MB/s	      0 B/op	     0 allocs/op
BenchmarkString_Mixed/mattn/go-runewidth-8                    14296 ns/op	  118.00 MB/s	      0 B/op	     0 allocs/op
BenchmarkString_Mixed/rivo/uniseg-8                           19770 ns/op	   85.33 MB/s	      0 B/op	     0 allocs/op

BenchmarkString_EastAsian/clipperhouse/displaywidth-8         10593 ns/op	  159.26 MB/s	      0 B/op	     0 allocs/op
BenchmarkString_EastAsian/mattn/go-runewidth-8                23980 ns/op	   70.35 MB/s	      0 B/op	     0 allocs/op
BenchmarkString_EastAsian/rivo/uniseg-8                       19777 ns/op	   85.30 MB/s	      0 B/op	     0 allocs/op

BenchmarkString_ASCII/clipperhouse/displaywidth-8              1032 ns/op	  124.09 MB/s	      0 B/op	     0 allocs/op
BenchmarkString_ASCII/mattn/go-runewidth-8                     1162 ns/op	  110.16 MB/s	      0 B/op	     0 allocs/op
BenchmarkString_ASCII/rivo/uniseg-8                            1586 ns/op	   80.69 MB/s	      0 B/op	     0 allocs/op

BenchmarkString_Emoji/clipperhouse/displaywidth-8              3017 ns/op	  240.01 MB/s	      0 B/op	     0 allocs/op
BenchmarkString_Emoji/mattn/go-runewidth-8                     4745 ns/op	  152.58 MB/s	      0 B/op	     0 allocs/op
BenchmarkString_Emoji/rivo/uniseg-8                            6745 ns/op	  107.34 MB/s	      0 B/op	     0 allocs/op

BenchmarkRune_Mixed/clipperhouse/displaywidth-8                3381 ns/op	  498.90 MB/s	      0 B/op	     0 allocs/op
BenchmarkRune_Mixed/mattn/go-runewidth-8                       5383 ns/op	  313.41 MB/s	      0 B/op	     0 allocs/op

BenchmarkRune_EastAsian/clipperhouse/displaywidth-8            3395 ns/op	  496.96 MB/s	      0 B/op	     0 allocs/op
BenchmarkRune_EastAsian/mattn/go-runewidth-8                  15645 ns/op	  107.83 MB/s	      0 B/op	     0 allocs/op

BenchmarkRune_ASCII/clipperhouse/displaywidth-8                 257.8 ns/op	  496.57 MB/s	      0 B/op	     0 allocs/op
BenchmarkRune_ASCII/mattn/go-runewidth-8                        267.3 ns/op	  478.89 MB/s	      0 B/op	     0 allocs/op

BenchmarkRune_Emoji/clipperhouse/displaywidth-8                1338 ns/op	  541.24 MB/s	      0 B/op	     0 allocs/op
BenchmarkRune_Emoji/mattn/go-runewidth-8                       2287 ns/op	  316.58 MB/s	      0 B/op	     0 allocs/op

BenchmarkTruncateWithTail/clipperhouse/displaywidth-8          3689 ns/op	   47.98 MB/s	    192 B/op	    14 allocs/op
BenchmarkTruncateWithTail/mattn/go-runewidth-8                 8069 ns/op	   21.93 MB/s	    192 B/op	    14 allocs/op

BenchmarkTruncateWithoutTail/clipperhouse/displaywidth-8       3457 ns/op	   66.24 MB/s	      0 B/op	     0 allocs/op
BenchmarkTruncateWithoutTail/mattn/go-runewidth-8             10441 ns/op	   21.93 MB/s	      0 B/op	     0 allocs/op
```
