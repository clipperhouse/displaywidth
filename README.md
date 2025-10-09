# displaywidth

A high-performance Go package for measuring the display width of strings in monospace terminals.

[![Documentation](https://pkg.go.dev/badge/github.com/clipperhouse/displaywidth.svg)](https://pkg.go.dev/github.com/clipperhouse/displaywidth)
![Go](https://github.com/clipperhouse/displaywidth/actions/workflows/gotest.yml/badge.svg)

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
}
```

### Options

You can specify East Asian Width and Strict Emoji Neutral settings. If
unspecified, the default is `EastAsianWidth: false, StrictEmojiNeutral: true`.


```go
options := displaywidth.Options{
    EastAsianWidth:     true,
    StrictEmojiNeutral: false,
}

width := options.String("Hello, 世界!")
fmt.Println(width)
```

## Details

This package implements the Unicode East Asian Width standard (UAX #11) and is
intended to be compatible with `go-runewidth`. It operates on bytes without
decoding runes for better performance.


## Prior Art

[mattn/go-runewidth](https://github.com/mattn/go-runewidth), which is excellent and popular.
In my testing, `displaywidth` returns identical outputs.

## Benchmarks

Part of my motivation is the insight that we can avoid decoding runes for better performance.

```bash
go test -bench=. -benchmem
```

```
goos: darwin
goarch: arm64
pkg: github.com/clipperhouse/displaywidth
cpu: Apple M2
BenchmarkStringDefault/displaywidth-8         	   96490	     10552 ns/op	 159.88 MB/s	       0 B/op	       0 allocs/op
BenchmarkStringDefault/go-runewidth-8         	   83907	     14369 ns/op	 117.41 MB/s	       0 B/op	       0 allocs/op
BenchmarkString_EAW/displaywidth-8            	  112807	     10646 ns/op	 158.46 MB/s	       0 B/op	       0 allocs/op
BenchmarkString_EAW/go-runewidth-8            	   50692	     23977 ns/op	  70.36 MB/s	       0 B/op	       0 allocs/op
BenchmarkString_StrictEmoji/displaywidth-8    	  113710	     10601 ns/op	 159.14 MB/s	       0 B/op	       0 allocs/op
BenchmarkString_StrictEmoji/go-runewidth-8    	   83220	     14403 ns/op	 117.13 MB/s	       0 B/op	       0 allocs/op
BenchmarkString_ASCII/displaywidth-8          	 1000000	      1077 ns/op	 118.83 MB/s	       0 B/op	       0 allocs/op
BenchmarkString_ASCII/go-runewidth-8          	 1000000	      1173 ns/op	 109.13 MB/s	       0 B/op	       0 allocs/op
BenchmarkString_Unicode/displaywidth-8        	 1367460	       881.1 ns/op	 150.94 MB/s	       0 B/op	       0 allocs/op
BenchmarkString_Unicode/go-runewidth-8        	  840982	      1437 ns/op	  92.57 MB/s	       0 B/op	       0 allocs/op
BenchmarkStringWidth_Emoji/displaywidth-8     	  403082	      3022 ns/op	 239.56 MB/s	       0 B/op	       0 allocs/op
BenchmarkStringWidth_Emoji/go-runewidth-8     	  247605	      4821 ns/op	 150.18 MB/s	       0 B/op	       0 allocs/op
BenchmarkString_Mixed/displaywidth-8          	  303606	      3956 ns/op	 128.17 MB/s	       0 B/op	       0 allocs/op
BenchmarkString_Mixed/go-runewidth-8          	  256921	      4639 ns/op	 109.30 MB/s	       0 B/op	       0 allocs/op
BenchmarkString_ControlChars/displaywidth-8   	 3795948	       315.2 ns/op	 104.70 MB/s	       0 B/op	       0 allocs/op
BenchmarkString_ControlChars/go-runewidth-8   	 3273128	       364.7 ns/op	  90.48 MB/s	       0 B/op	       0 allocs/op
BenchmarkRuneDefault/displaywidth-8           	 3772311	       318.1 ns/op	 433.82 MB/s	       0 B/op	       0 allocs/op
BenchmarkRuneDefault/go-runewidth-8           	 1753222	       684.4 ns/op	 201.63 MB/s	       0 B/op	       0 allocs/op
BenchmarkRuneWidth_EAW/displaywidth-8         	 8469133	       142.6 ns/op	 385.75 MB/s	       0 B/op	       0 allocs/op
BenchmarkRuneWidth_EAW/go-runewidth-8         	 2383420	       502.9 ns/op	 109.37 MB/s	       0 B/op	       0 allocs/op
BenchmarkRuneWidth_ASCII/displaywidth-8       	19660138	        62.01 ns/op	 467.63 MB/s	       0 B/op	       0 allocs/op
BenchmarkRuneWidth_ASCII/go-runewidth-8       	17664040	        67.34 ns/op	 430.68 MB/s	       0 B/op	       0 allocs/op
```

I use a similar technique in [this grapheme cluster library](https://github.com/clipperhouse/uax29).
