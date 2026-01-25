## Compatibility

In real-world text, you should see the same outputs from
`clipperhouse/displaywidth`, `mattn/go-runewidth`, and `rivo/uniseg`.

The tests in this `comparison` package exercise the behaviors of the three
libraries. Extensive details are available in the
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

BenchmarkString_Mixed/clipperhouse/displaywidth-8             5784 ns/op	      291.69 MB/s	      0 B/op	   0 allocs/op
BenchmarkString_Mixed/mattn/go-runewidth-8                   14751 ns/op	      114.36 MB/s	      0 B/op	   0 allocs/op
BenchmarkString_Mixed/rivo/uniseg-8                          19360 ns/op	       87.14 MB/s	      0 B/op	   0 allocs/op

BenchmarkString_ASCII/clipperhouse/displaywidth-8               54.60 ns/op	     2344.32 MB/s	      0 B/op	   0 allocs/op
BenchmarkString_ASCII/mattn/go-runewidth-8                    1195 ns/op	      107.08 MB/s	      0 B/op	   0 allocs/op
BenchmarkString_ASCII/rivo/uniseg-8                           1578 ns/op	       81.13 MB/s	      0 B/op	   0 allocs/op

BenchmarkString_EastAsian/clipperhouse/displaywidth-8         5837 ns/op	      289.01 MB/s	      0 B/op	   0 allocs/op
BenchmarkString_EastAsian/mattn/go-runewidth-8               24418 ns/op	       69.09 MB/s	      0 B/op	   0 allocs/op
BenchmarkString_EastAsian/rivo/uniseg-8                      19339 ns/op	       87.23 MB/s	      0 B/op	   0 allocs/op

BenchmarkString_Emoji/clipperhouse/displaywidth-8             3225 ns/op	      224.51 MB/s	      0 B/op	   0 allocs/op
BenchmarkString_Emoji/mattn/go-runewidth-8                    4851 ns/op	      149.25 MB/s	      0 B/op	   0 allocs/op
BenchmarkString_Emoji/rivo/uniseg-8                           6591 ns/op	      109.85 MB/s	      0 B/op	   0 allocs/op

BenchmarkRune_Mixed/clipperhouse/displaywidth-8               3385 ns/op	      498.34 MB/s	      0 B/op	   0 allocs/op
BenchmarkRune_Mixed/mattn/go-runewidth-8                      5354 ns/op	      315.07 MB/s	      0 B/op	   0 allocs/op

BenchmarkRune_EastAsian/clipperhouse/displaywidth-8           3397 ns/op	      496.56 MB/s	      0 B/op	   0 allocs/op
BenchmarkRune_EastAsian/mattn/go-runewidth-8                 15673 ns/op	      107.64 MB/s	      0 B/op	   0 allocs/op

BenchmarkRune_ASCII/clipperhouse/displaywidth-8                255.7 ns/op	      500.53 MB/s	      0 B/op	   0 allocs/op
BenchmarkRune_ASCII/mattn/go-runewidth-8                       261.5 ns/op	      489.55 MB/s	      0 B/op	   0 allocs/op

BenchmarkRune_Emoji/clipperhouse/displaywidth-8               1371 ns/op	      528.22 MB/s	      0 B/op	   0 allocs/op
BenchmarkRune_Emoji/mattn/go-runewidth-8                      2267 ns/op	      319.43 MB/s	      0 B/op	   0 allocs/op

BenchmarkTruncateWithTail/clipperhouse/displaywidth-8         3229 ns/op	       54.82 MB/s	    192 B/op	  14 allocs/op
BenchmarkTruncateWithTail/mattn/go-runewidth-8                8408 ns/op	       21.05 MB/s	    192 B/op	  14 allocs/op

BenchmarkTruncateWithoutTail/clipperhouse/displaywidth-8      3554 ns/op	       64.43 MB/s	      0 B/op	   0 allocs/op
BenchmarkTruncateWithoutTail/mattn/go-runewidth-8            11189 ns/op	       20.47 MB/s	      0 B/op	   0 allocs/op
```
