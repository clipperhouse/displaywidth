## Compatibility

In real-world text, you should see the same outputs from
`clipperhouse/displaywidth`,`mattn/go-runewidth`, and `rivo/uniseg`.

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

BenchmarkString_Mixed/clipperhouse/displaywidth-8         	   5460 ns/op	      308.96 MB/s	     0 B/op	     0 allocs/op
BenchmarkString_Mixed/mattn/go-runewidth-8                	  14301 ns/op	      117.96 MB/s	     0 B/op	     0 allocs/op
BenchmarkString_Mixed/rivo/uniseg-8                       	  19562 ns/op	       86.24 MB/s	     0 B/op	     0 allocs/op

BenchmarkString_EastAsian/clipperhouse/displaywidth-8     	   5546 ns/op	      304.20 MB/s	     0 B/op	     0 allocs/op
BenchmarkString_EastAsian/mattn/go-runewidth-8            	  23801 ns/op	       70.88 MB/s	     0 B/op	     0 allocs/op
BenchmarkString_EastAsian/rivo/uniseg-8                   	  19768 ns/op	       85.34 MB/s	     0 B/op	     0 allocs/op

BenchmarkString_ASCII/clipperhouse/displaywidth-8         	     54.58 ns/op	 2345.21 MB/s	     0 B/op	     0 allocs/op
BenchmarkString_ASCII/mattn/go-runewidth-8                	   1167 ns/op	      109.73 MB/s	     0 B/op	     0 allocs/op
BenchmarkString_ASCII/rivo/uniseg-8                       	   1577 ns/op	       81.17 MB/s	     0 B/op	     0 allocs/op

BenchmarkString_Emoji/clipperhouse/displaywidth-8         	   3127 ns/op	      231.51 MB/s	     0 B/op	     0 allocs/op
BenchmarkString_Emoji/mattn/go-runewidth-8                	   4722 ns/op	      153.31 MB/s	     0 B/op	     0 allocs/op
BenchmarkString_Emoji/rivo/uniseg-8                       	   6562 ns/op	      110.34 MB/s	     0 B/op	     0 allocs/op

BenchmarkRune_Mixed/clipperhouse/displaywidth-8           	   3452 ns/op	      488.68 MB/s	     0 B/op	     0 allocs/op
BenchmarkRune_Mixed/mattn/go-runewidth-8                  	   5367 ns/op	      314.33 MB/s	     0 B/op	     0 allocs/op

BenchmarkRune_EastAsian/clipperhouse/displaywidth-8       	   3757 ns/op	      449.06 MB/s	     0 B/op	     0 allocs/op
BenchmarkRune_EastAsian/mattn/go-runewidth-8              	  15390 ns/op	      109.62 MB/s	     0 B/op	     0 allocs/op

BenchmarkRune_ASCII/clipperhouse/displaywidth-8           	    256.3 ns/op 	  499.40 MB/s	     0 B/op	     0 allocs/op
BenchmarkRune_ASCII/mattn/go-runewidth-8                  	    262.3 ns/op 	  487.91 MB/s	     0 B/op	     0 allocs/op

BenchmarkRune_Emoji/clipperhouse/displaywidth-8           	   1436 ns/op	      504.16 MB/s	     0 B/op	     0 allocs/op
BenchmarkRune_Emoji/mattn/go-runewidth-8                  	   2267 ns/op	      319.32 MB/s	     0 B/op	     0 allocs/op

BenchmarkTruncateWithTail/clipperhouse/displaywidth-8     	   3120 ns/op	       56.73 MB/s	   192 B/op	    14 allocs/op
BenchmarkTruncateWithTail/mattn/go-runewidth-8            	   8134 ns/op	       21.76 MB/s	   192 B/op	    14 allocs/op

BenchmarkTruncateWithoutTail/clipperhouse/displaywidth-8  	   3427 ns/op	       66.82 MB/s	     0 B/op	     0 allocs/op
BenchmarkTruncateWithoutTail/mattn/go-runewidth-8         	  10410 ns/op	       22.00 MB/s	     0 B/op	     0 allocs/op
```
