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

BenchmarkString_Mixed/clipperhouse/displaywidth-8      	     10326 ns/op	   163.37 MB/s	      0 B/op	    0 allocs/op
BenchmarkString_Mixed/mattn/go-runewidth-8             	     14415 ns/op	   117.03 MB/s	      0 B/op	    0 allocs/op
BenchmarkString_Mixed/rivo/uniseg-8                    	     19343 ns/op	    87.21 MB/s	      0 B/op	    0 allocs/op

BenchmarkString_EastAsian/clipperhouse/displaywidth-8  	     10561 ns/op	   159.74 MB/s	      0 B/op	    0 allocs/op
BenchmarkString_EastAsian/mattn/go-runewidth-8         	     23790 ns/op	    70.91 MB/s	      0 B/op	    0 allocs/op
BenchmarkString_EastAsian/rivo/uniseg-8                	     19322 ns/op	    87.31 MB/s	      0 B/op	    0 allocs/op

BenchmarkString_ASCII/clipperhouse/displaywidth-8      	      1033 ns/op	   123.88 MB/s	      0 B/op	    0 allocs/op
BenchmarkString_ASCII/mattn/go-runewidth-8             	      1168 ns/op	   109.59 MB/s	      0 B/op	    0 allocs/op
BenchmarkString_ASCII/rivo/uniseg-8                    	      1585 ns/op	    80.74 MB/s	      0 B/op	    0 allocs/op

BenchmarkString_Emoji/clipperhouse/displaywidth-8      	      3034 ns/op	   238.61 MB/s	      0 B/op	    0 allocs/op
BenchmarkString_Emoji/mattn/go-runewidth-8             	      4797 ns/op	   150.94 MB/s	      0 B/op	    0 allocs/op
BenchmarkString_Emoji/rivo/uniseg-8                    	      6612 ns/op	   109.50 MB/s	      0 B/op	    0 allocs/op

BenchmarkRune_Mixed/clipperhouse/displaywidth-8        	      3343 ns/op	   504.67 MB/s	      0 B/op	    0 allocs/op
BenchmarkRune_Mixed/mattn/go-runewidth-8               	      5414 ns/op	   311.62 MB/s	      0 B/op	    0 allocs/op

BenchmarkRune_EastAsian/clipperhouse/displaywidth-8    	      3393 ns/op	   497.17 MB/s	      0 B/op	    0 allocs/op
BenchmarkRune_EastAsian/mattn/go-runewidth-8           	     15312 ns/op	   110.17 MB/s	      0 B/op	    0 allocs/op

BenchmarkRune_ASCII/clipperhouse/displaywidth-8        	       256.9 ns/op	   498.32 MB/s	      0 B/op	    0 allocs/op
BenchmarkRune_ASCII/mattn/go-runewidth-8               	       265.7 ns/op	   481.75 MB/s	      0 B/op	    0 allocs/op

BenchmarkRune_Emoji/clipperhouse/displaywidth-8        	      1336 ns/op	   541.96 MB/s	      0 B/op	    0 allocs/op
BenchmarkRune_Emoji/mattn/go-runewidth-8               	      2304 ns/op	   314.23 MB/s	      0 B/op	    0 allocs/op
```
