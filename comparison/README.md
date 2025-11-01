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

BenchmarkString_Mixed/clipperhouse/displaywidth-8         10929 ns/op	    154.36 MB/s	      0 B/op     0 allocs/op
BenchmarkString_Mixed/mattn/go-runewidth-8                14540 ns/op	    116.02 MB/s	      0 B/op     0 allocs/op
BenchmarkString_Mixed/rivo/uniseg-8                       19751 ns/op	     85.41 MB/s	      0 B/op     0 allocs/op

BenchmarkString_EastAsian/clipperhouse/displaywidth-8     10885 ns/op	    154.98 MB/s	      0 B/op     0 allocs/op
BenchmarkString_EastAsian/mattn/go-runewidth-8            23969 ns/op	     70.38 MB/s	      0 B/op     0 allocs/op
BenchmarkString_EastAsian/rivo/uniseg-8                   19852 ns/op	     84.98 MB/s	      0 B/op     0 allocs/op

BenchmarkString_ASCII/clipperhouse/displaywidth-8          1103 ns/op	    116.01 MB/s	      0 B/op     0 allocs/op
BenchmarkString_ASCII/mattn/go-runewidth-8                 1166 ns/op	    109.79 MB/s	      0 B/op     0 allocs/op
BenchmarkString_ASCII/rivo/uniseg-8                        1584 ns/op	     80.83 MB/s	      0 B/op     0 allocs/op

BenchmarkString_Emoji/clipperhouse/displaywidth-8          3108 ns/op	    232.93 MB/s	      0 B/op     0 allocs/op
BenchmarkString_Emoji/mattn/go-runewidth-8                 4802 ns/op	    150.76 MB/s	      0 B/op     0 allocs/op
BenchmarkString_Emoji/rivo/uniseg-8                        6607 ns/op	    109.58 MB/s	      0 B/op     0 allocs/op

BenchmarkRune_Mixed/clipperhouse/displaywidth-8            3456 ns/op	    488.20 MB/s	      0 B/op     0 allocs/op
BenchmarkRune_Mixed/mattn/go-runewidth-8                   5400 ns/op	    312.39 MB/s	      0 B/op     0 allocs/op

BenchmarkRune_EastAsian/clipperhouse/displaywidth-8        3475 ns/op	    485.41 MB/s	      0 B/op     0 allocs/op
BenchmarkRune_EastAsian/mattn/go-runewidth-8              15701 ns/op	    107.44 MB/s	      0 B/op     0 allocs/op

BenchmarkRune_ASCII/clipperhouse/displaywidth-8             257.0 ns/op	    498.13 MB/s	      0 B/op     0 allocs/op
BenchmarkRune_ASCII/mattn/go-runewidth-8                    266.4 ns/op	    480.50 MB/s	      0 B/op     0 allocs/op

BenchmarkRune_Emoji/clipperhouse/displaywidth-8            1384 ns/op	    523.02 MB/s	      0 B/op     0 allocs/op
BenchmarkRune_Emoji/mattn/go-runewidth-8                   2273 ns/op	    318.45 MB/s	      0 B/op     0 allocs/op
```
