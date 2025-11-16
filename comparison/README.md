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

BenchmarkString_Mixed/clipperhouse/displaywidth-8     	     10469 ns/op	   161.15 MB/s      0 B/op      0 allocs/op
BenchmarkString_Mixed/mattn/go-runewidth-8            	     14250 ns/op	   118.39 MB/s      0 B/op      0 allocs/op
BenchmarkString_Mixed/rivo/uniseg-8                   	     19258 ns/op	    87.60 MB/s      0 B/op      0 allocs/op

BenchmarkString_EastAsian/clipperhouse/displaywidth-8 	     10518 ns/op	   160.39 MB/s      0 B/op      0 allocs/op
BenchmarkString_EastAsian/mattn/go-runewidth-8        	     23827 ns/op	    70.80 MB/s      0 B/op      0 allocs/op
BenchmarkString_EastAsian/rivo/uniseg-8               	     19537 ns/op	    86.35 MB/s      0 B/op      0 allocs/op

BenchmarkString_ASCII/clipperhouse/displaywidth-8     	      1027 ns/op	   124.61 MB/s      0 B/op      0 allocs/op
BenchmarkString_ASCII/mattn/go-runewidth-8            	      1166 ns/op	   109.78 MB/s      0 B/op      0 allocs/op
BenchmarkString_ASCII/rivo/uniseg-8                   	      1551 ns/op	    82.52 MB/s      0 B/op      0 allocs/op

BenchmarkString_Emoji/clipperhouse/displaywidth-8     	      3164 ns/op	   228.84 MB/s      0 B/op      0 allocs/op
BenchmarkString_Emoji/mattn/go-runewidth-8            	      4728 ns/op	   153.13 MB/s      0 B/op      0 allocs/op
BenchmarkString_Emoji/rivo/uniseg-8                   	      6489 ns/op	   111.57 MB/s      0 B/op      0 allocs/op

BenchmarkRune_Mixed/clipperhouse/displaywidth-8       	      3429 ns/op	   491.96 MB/s      0 B/op      0 allocs/op
BenchmarkRune_Mixed/mattn/go-runewidth-8              	      5308 ns/op	   317.81 MB/s      0 B/op      0 allocs/op

BenchmarkRune_EastAsian/clipperhouse/displaywidth-8   	      3419 ns/op	   493.49 MB/s      0 B/op      0 allocs/op
BenchmarkRune_EastAsian/mattn/go-runewidth-8          	     15321 ns/op	   110.11 MB/s      0 B/op      0 allocs/op

BenchmarkRune_ASCII/clipperhouse/displaywidth-8       	       254.4 ns/op	   503.19 MB/s      0 B/op      0 allocs/op
BenchmarkRune_ASCII/mattn/go-runewidth-8              	       264.3 ns/op	   484.31 MB/s      0 B/op      0 allocs/op

BenchmarkRune_Emoji/clipperhouse/displaywidth-8       	      1374 ns/op	   527.02 MB/s      0 B/op      0 allocs/op
BenchmarkRune_Emoji/mattn/go-runewidth-8              	      2210 ns/op	   327.66 MB/s      0 B/op      0 allocs/op
```
