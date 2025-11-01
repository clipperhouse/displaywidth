# displaywidth

A high-performance Go package for measuring the monospace display width of strings, UTF-8 bytes, and runes.

[![Documentation](https://pkg.go.dev/badge/github.com/clipperhouse/displaywidth.svg)](https://pkg.go.dev/github.com/clipperhouse/displaywidth)
[![Test](https://github.com/clipperhouse/displaywidth/actions/workflows/gotest.yml/badge.svg)](https://github.com/clipperhouse/displaywidth/actions/workflows/gotest.yml)
[![Fuzz](https://github.com/clipperhouse/displaywidth/actions/workflows/gofuzz.yml/badge.svg)](https://github.com/clipperhouse/displaywidth/actions/workflows/gofuzz.yml)
## Install
```bash
go get github.com/clipperhouse/displaywidth
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/clipperhouse/displaywidth"
)

func main() {
    width := displaywidth.String("Hello, 世界!")
    fmt.Println(width)

    width = displaywidth.Bytes([]byte("🌍"))
    fmt.Println(width)

    width = displaywidth.Rune('🌍')
    fmt.Println(width)
}
```

For most purposes, you should use the `String` or `Bytes` methods.

### Options

You can specify East Asian Width settings. If unspecified, the default is `EastAsianWidth: false`.

```go
options := displaywidth.Options{
    EastAsianWidth: true,
}

width := options.String("Hello, 世界!")
fmt.Println(width)
```

## Details

This package implements the Unicode East Asian Width standard
([UAX #11](https://www.unicode.org/reports/tr11/)), and handles
[version selectors](https://en.wikipedia.org/wiki/Variation_Selectors_(Unicode_block)),
and [regional indicator pairs](https://en.wikipedia.org/wiki/Regional_indicator_symbol)
(flags). We cover much of [Unicode TR51](https://unicode.org/reports/tr51/).

## Prior Art

[mattn/go-runewidth](https://github.com/mattn/go-runewidth)

[rivo/uniseg](https://github.com/rivo/uniseg)

[x/text/width](https://pkg.go.dev/golang.org/x/text/width)

[x/text/internal/triegen](https://pkg.go.dev/golang.org/x/text/internal/triegen)

## Benchmarks

```bash
cd comparison
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

## Compatibility

`clipperhouse/displaywidth`, `mattn/go-runewidth`, and `rivo/uniseg` should give the
same outputs for real-world text. See [comparison/README.md](comparison/README.md).

If you wish to investigate the core logic, see the `lookupProperties` and `width`
functions in [width.go](width.go#L112). The core of the trie generation logic is in
`BuildPropertyBitmap` in [unicode.go](internal/gen/unicode.go#L309).
