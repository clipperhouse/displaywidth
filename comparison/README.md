This package is for comparison of `clipperhouse/displaywidth` with the excellent
 `mattn/go-runewidth` and `rivo/uniseg` packages.

## Compatibility

In real-world text, you mostly should see the same outputs from
`clipperhouse/displaywidth`,`mattn/go-runewidth`, and `rivo/uniseg`. It's
mostly the same data and logic. Run the various tests in this package to see
the minor differences.

All of these libraries (I believe) are based on the Unicode standards of
[UAX #11 (East Asian Width)](https://www.unicode.org/reports/tr11/).
`displaywidth` and `uniseg` support
[UAX #51 (Unicode Emoji)](https://unicode.org/reports/tr51/).

In each of those libraries, you almost certainly want to use the `String` or `StringWidth`
methods. They operate on graphemes, not runes, which is what appears on your screen.

`mattn/go-runewidth` has some automatic defaults based on a machine's locale and
environment variables. `clipperhouse/displaywidth` does not, but you can set them
manually with the `Options` struct.

Some libraries choose to import both `go-runewidth` and `displaywidth`, branching
to use `go-runewidth` where backwards compatibility is needed, and `displaywidth`
otherwise.

## Benchmarks

```bash
go test -bench=. -benchmem
```

```
goos: darwin
goarch: arm64
pkg: github.com/clipperhouse/displaywidth/comparison
cpu: Apple M2
BenchmarkString_Mixed/clipperhouse/displaywidth-8      	     11290 ns/op	   149.42 MB/s     0 B/op	   0 allocs/op
BenchmarkString_Mixed/mattn/go-runewidth-8             	     14439 ns/op	   116.84 MB/s     0 B/op	   0 allocs/op
BenchmarkString_Mixed/rivo/uniseg-8                    	     20076 ns/op	    84.03 MB/s     0 B/op	   0 allocs/op

BenchmarkString_EastAsian/clipperhouse/displaywidth-8  	     11248 ns/op	   149.98 MB/s     0 B/op	   0 allocs/op
BenchmarkString_EastAsian/mattn/go-runewidth-8         	     24063 ns/op	    70.11 MB/s     0 B/op	   0 allocs/op
BenchmarkString_EastAsian/rivo/uniseg-8                	     20051 ns/op	    84.14 MB/s     0 B/op	   0 allocs/op

BenchmarkString_ASCII/clipperhouse/displaywidth-8      	      1116 ns/op	   114.71 MB/s     0 B/op	   0 allocs/op
BenchmarkString_ASCII/mattn/go-runewidth-8             	      1182 ns/op	   108.27 MB/s     0 B/op	   0 allocs/op
BenchmarkString_ASCII/rivo/uniseg-8                    	      1620 ns/op	    79.04 MB/s     0 B/op	   0 allocs/op

BenchmarkString_Emoji/clipperhouse/displaywidth-8      	      3264 ns/op	   221.82 MB/s     0 B/op	   0 allocs/op
BenchmarkString_Emoji/mattn/go-runewidth-8             	      4804 ns/op	   150.71 MB/s     0 B/op	   0 allocs/op
BenchmarkString_Emoji/rivo/uniseg-8                    	      6783 ns/op	   106.74 MB/s     0 B/op	   0 allocs/op

BenchmarkRune_Mixed/clipperhouse/displaywidth-8        	      3759 ns/op	   448.83 MB/s     0 B/op	   0 allocs/op
BenchmarkRune_Mixed/mattn/go-runewidth-8               	      5417 ns/op	   311.40 MB/s     0 B/op	   0 allocs/op

BenchmarkRune_EastAsian/clipperhouse/displaywidth-8    	      3678 ns/op	   458.69 MB/s     0 B/op	   0 allocs/op
BenchmarkRune_EastAsian/mattn/go-runewidth-8           	     15908 ns/op	   106.05 MB/s     0 B/op	   0 allocs/op

BenchmarkRune_ASCII/clipperhouse/displaywidth-8        	       265.2 ns/op	   482.70 MB/s     0 B/op	   0 allocs/op
BenchmarkRune_ASCII/mattn/go-runewidth-8               	       265.2 ns/op	   482.67 MB/s     0 B/op	   0 allocs/op

BenchmarkRune_Emoji/clipperhouse/displaywidth-8        	      1522 ns/op	   475.65 MB/s     0 B/op	   0 allocs/op
BenchmarkRune_Emoji/mattn/go-runewidth-8               	      2295 ns/op	   315.53 MB/s     0 B/op	   0 allocs/op
```
