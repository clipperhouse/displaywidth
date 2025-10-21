This package is for comparision of `clipperhouse/displaywidth` with the excellent
 `mattn/go-runewidth` and `rivo/uniseg` packages.

## Compatibility

In real-world text, you should see the same outputs from `clipperhouse/displaywidth`,
`mattn/go-runewidth`, and `rivo/uniseg`. It's mostly the same data and logic.

In each of those libraries, you almost certainly want to use the `String` or `StringWidth`
methods. They operate on graphemes, not runes, which is what appears on your screen.

`mattn/go-runewidth` has some automatic defaults based on a machine's locale and
environment variables. `clipperhouse/displaywidth` does not, but you can set them
 manually with the `Options` struct.

If you are operating on individual runes (which you shouldn't!):

- Unicode category Mn (Nonspacing Mark): `displaywidth` will return width 0, `go-runewidth` may return width 1 for some runes.
- Unicode category Cf (Format): `displaywidth` will return width 0, `go-runewidth` may return width 1 for some runes.
- Unicode category Mc (Spacing Mark): `displaywidth` will return width 1, `go-runewidth` may return width 0 for some runes.
- Unicode category Cs (Surrogate): `displaywidth` will return width 0, `go-runewidth` may return width 1 for some runes. Surrogates are not valid UTF-8; some packages may turn them into the replacement character (U+FFFD).
- Unicode category Zl (Line separator): `displaywidth` will return width 0, `go-runewidth` may return width 1.
- Unicode category Zp (Paragraph separator): `displaywidth` will return width 0, `go-runewidth` may return width 1.
- Unicode Noncharacters (U+FFFE and U+FFFF): `displaywidth` will return width 0, `go-runewidth` may return width 1.

## Benchmarks

```bash
go test -bench=. -benchmem
```

```
goos: darwin
goarch: arm64
pkg: github.com/clipperhouse/displaywidth/comparison
cpu: Apple M2
BenchmarkStringDefault/clipperhouse/displaywidth-8        	     11225 ns/op	     150.28 MB/s	     0 B/op	       0 allocs/op
BenchmarkStringDefault/mattn/go-runewidth-8               	     14732 ns/op	     114.51 MB/s	     0 B/op	       0 allocs/op
BenchmarkStringDefault/rivo/uniseg-8                      	     19820 ns/op	      85.12 MB/s	     0 B/op	       0 allocs/op

BenchmarkString_EAW/clipperhouse/displaywidth-8           	     11224 ns/op	     150.30 MB/s	     0 B/op	       0 allocs/op
BenchmarkString_EAW/mattn/go-runewidth-8                  	     27861 ns/op	      60.55 MB/s	     0 B/op	       0 allocs/op
BenchmarkString_EAW/rivo/uniseg-8                         	     19929 ns/op	      84.65 MB/s	     0 B/op	       0 allocs/op

BenchmarkString_StrictEmoji/clipperhouse/displaywidth-8   	     11221 ns/op	     150.35 MB/s	     0 B/op	       0 allocs/op
BenchmarkString_StrictEmoji/mattn/go-runewidth-8          	     14724 ns/op	     114.58 MB/s	     0 B/op	       0 allocs/op
BenchmarkString_StrictEmoji/rivo/uniseg-8                 	     19812 ns/op	      85.15 MB/s	     0 B/op	       0 allocs/op

BenchmarkString_ASCII/clipperhouse/displaywidth-8         	      1143 ns/op	     111.95 MB/s	     0 B/op	       0 allocs/op
BenchmarkString_ASCII/mattn/go-runewidth-8                	      1166 ns/op	     109.74 MB/s	     0 B/op	       0 allocs/op
BenchmarkString_ASCII/rivo/uniseg-8                       	      1628 ns/op	      78.65 MB/s	     0 B/op	       0 allocs/op

BenchmarkString_Unicode/clipperhouse/displaywidth-8       	       942.0 ns/op	     141.19 MB/s	     0 B/op	       0 allocs/op
BenchmarkString_Unicode/mattn/go-runewidth-8              	      1579 ns/op	      84.25 MB/s	     0 B/op	       0 allocs/op
BenchmarkString_Unicode/rivo/uniseg-8                     	      2026 ns/op	      65.66 MB/s	     0 B/op	       0 allocs/op

BenchmarkStringWidth_Emoji/clipperhouse/displaywidth-8    	      3158 ns/op	     229.28 MB/s	     0 B/op	       0 allocs/op
BenchmarkStringWidth_Emoji/mattn/go-runewidth-8           	      4805 ns/op	     150.66 MB/s	     0 B/op	       0 allocs/op
BenchmarkStringWidth_Emoji/rivo/uniseg-8                  	      6606 ns/op	     109.60 MB/s	     0 B/op	       0 allocs/op

BenchmarkString_Mixed/clipperhouse/displaywidth-8         	      4201 ns/op	     120.69 MB/s	     0 B/op	       0 allocs/op
BenchmarkString_Mixed/mattn/go-runewidth-8                	      4710 ns/op	     107.65 MB/s	     0 B/op	       0 allocs/op
BenchmarkString_Mixed/rivo/uniseg-8                       	      6311 ns/op	      80.34 MB/s	     0 B/op	       0 allocs/op

BenchmarkString_ControlChars/clipperhouse/displaywidth-8  	       351.5 ns/op	      93.87 MB/s	     0 B/op	       0 allocs/op
BenchmarkString_ControlChars/mattn/go-runewidth-8         	       365.0 ns/op	      90.40 MB/s	     0 B/op	       0 allocs/op
BenchmarkString_ControlChars/rivo/uniseg-8                	       411.9 ns/op	      80.12 MB/s	     0 B/op	       0 allocs/op
```
