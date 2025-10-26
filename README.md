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
BenchmarkStringDefault/clipperhouse/displaywidth-8         	     11040 ns/op	 152.81 MB/s     0 B/op	     0 allocs/op
BenchmarkStringDefault/mattn/go-runewidth-8                	     14468 ns/op	 116.60 MB/s     0 B/op	     0 allocs/op
BenchmarkStringDefault/rivo/uniseg-8                       	     19274 ns/op	  87.53 MB/s     0 B/op	     0 allocs/op
BenchmarkString_EAW/clipperhouse/displaywidth-8            	     11537 ns/op	 146.22 MB/s     0 B/op	     0 allocs/op
BenchmarkString_EAW/mattn/go-runewidth-8                   	     23753 ns/op	  71.02 MB/s     0 B/op	     0 allocs/op
BenchmarkString_EAW/rivo/uniseg-8                          	     19739 ns/op	  85.47 MB/s     0 B/op	     0 allocs/op
BenchmarkString_StrictEmoji/clipperhouse/displaywidth-8    	     11641 ns/op	 144.92 MB/s     0 B/op	     0 allocs/op
BenchmarkString_StrictEmoji/mattn/go-runewidth-8           	     14337 ns/op	 117.67 MB/s     0 B/op	     0 allocs/op
BenchmarkString_StrictEmoji/rivo/uniseg-8                  	     19890 ns/op	  84.82 MB/s     0 B/op	     0 allocs/op
BenchmarkString_ASCII/clipperhouse/displaywidth-8          	      1108 ns/op	 115.51 MB/s     0 B/op	     0 allocs/op
BenchmarkString_ASCII/mattn/go-runewidth-8                 	      1166 ns/op	 109.73 MB/s     0 B/op	     0 allocs/op
BenchmarkString_ASCII/rivo/uniseg-8                        	      1582 ns/op	  80.92 MB/s     0 B/op	     0 allocs/op
BenchmarkString_Unicode/clipperhouse/displaywidth-8        	       981.9 ns/op	 135.45 MB/s     0 B/op	     0 allocs/op
BenchmarkString_Unicode/mattn/go-runewidth-8               	      1428 ns/op	  93.15 MB/s     0 B/op	     0 allocs/op
BenchmarkString_Unicode/rivo/uniseg-8                      	      2023 ns/op	  65.74 MB/s     0 B/op	     0 allocs/op
BenchmarkStringWidth_Emoji/clipperhouse/displaywidth-8     	      3243 ns/op	 223.28 MB/s     0 B/op	     0 allocs/op
BenchmarkStringWidth_Emoji/mattn/go-runewidth-8            	      4743 ns/op	 152.63 MB/s     0 B/op	     0 allocs/op
BenchmarkStringWidth_Emoji/rivo/uniseg-8                   	      6574 ns/op	 110.13 MB/s     0 B/op	     0 allocs/op
BenchmarkString_Mixed/clipperhouse/displaywidth-8          	      4105 ns/op	 123.50 MB/s     0 B/op	     0 allocs/op
BenchmarkString_Mixed/mattn/go-runewidth-8                 	      4616 ns/op	 109.83 MB/s     0 B/op	     0 allocs/op
BenchmarkString_Mixed/rivo/uniseg-8                        	      6311 ns/op	  80.34 MB/s     0 B/op	     0 allocs/op
BenchmarkString_ControlChars/clipperhouse/displaywidth-8   	       347.5 ns/op	  94.95 MB/s     0 B/op	     0 allocs/op
BenchmarkString_ControlChars/mattn/go-runewidth-8          	       365.2 ns/op	  90.37 MB/s     0 B/op	     0 allocs/op
BenchmarkString_ControlChars/rivo/uniseg-8                 	       409.3 ns/op	  80.62 MB/s     0 B/op	     0 allocs/op
```

I use a similar technique in [this grapheme cluster library](https://github.com/clipperhouse/uax29).

## Compatibility

`clipperhouse/displaywidth`, `mattn/go-runewidth`, and `rivo/uniseg` should give the
same outputs for real-world text. See [comparison/README.md](comparison/README.md).
