# displaywidth

A high-performance Go package for measuring the display width of strings in monospace terminals.

```bash
go get github.com/clipperhouse/displaywidth
```

## API

```go
// StringWidth calculates the display width of a string
// eastAsianWidth: when true, treat ambiguous width characters as wide (width 2)
// strictEmojiNeutral: when true, use strict emoji width calculation (some emoji become width 1)
func String(s string, eastAsianWidth bool, strictEmojiNeutral bool) int

// BytesWidth calculates the display width of a []byte
// eastAsianWidth: when true, treat ambiguous width characters as wide (width 2)
// strictEmojiNeutral: when true, use strict emoji width calculation (some emoji become width 1)
func Bytes(s []byte, eastAsianWidth bool, strictEmojiNeutral bool) int
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/clipperhouse/displaywidth"
)

func main() {
    // Basic usage
    width := displaywidth.String("Hello, 世界!", false, false)
    fmt.Println(width) // Output: 13

    // With East Asian width mode
    width = displaywidth.String("café", true, false)
    fmt.Println(width) // Output: 4

    // Using BytesWidth
    width = displaywidth.Bytes([]byte("🌍"), false, false)
    fmt.Println(width) // Output: 2
}
```

## Details

This package implements the Unicode East Asian Width standard (UAX #11) and is intended to be compatible with `go-runewidth` and `wcwidth`. It operates on bytes without decoding runes for better performance.

## Prior Art

[mattn/go-runewidth](https://github.com/mattn/go-runewidth), which is excellent and popular.
I've made an effort so that `displaywidth` returns identical results.

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
BenchmarkString/displaywidth-8   	  108600	     11077 ns/op	 154.82 MB/s	       0 B/op	       0 allocs/op
BenchmarkString/go-runewidth-8   	   82090	     14666 ns/op	 116.94 MB/s	       0 B/op	       0 allocs/op
BenchmarkString_EAW/displaywidth-8         	  108531	     11104 ns/op	 154.45 MB/s	       0 B/op	       0 allocs/op
BenchmarkString_EAW/go-runewidth-8         	   49471	     24450 ns/op	  70.14 MB/s	       0 B/op	       0 allocs/op
BenchmarkString_StrictEmoji/displaywidth-8 	  108102	     11100 ns/op	 154.50 MB/s	       0 B/op	       0 allocs/op
BenchmarkString_StrictEmoji/go-runewidth-8 	   81535	     14779 ns/op	 116.04 MB/s	       0 B/op	       0 allocs/op
BenchmarkString_ASCII/displaywidth-8       	 1000000	      1087 ns/op	 117.79 MB/s	       0 B/op	       0 allocs/op
BenchmarkString_ASCII/go-runewidth-8       	 1000000	      1169 ns/op	 109.47 MB/s	       0 B/op	       0 allocs/op
BenchmarkString_Unicode/displaywidth-8     	 1311325	       914.3 ns/op	 145.46 MB/s	       0 B/op	       0 allocs/op
BenchmarkString_Unicode/go-runewidth-8     	  820699	      1437 ns/op	  92.56 MB/s	       0 B/op	       0 allocs/op
BenchmarkStringWidth_Emoji/displaywidth-8  	  385572	      3138 ns/op	 230.70 MB/s	       0 B/op	       0 allocs/op
BenchmarkStringWidth_Emoji/go-runewidth-8  	  246482	      4839 ns/op	 149.63 MB/s	       0 B/op	       0 allocs/op
BenchmarkString_Mixed/displaywidth-8       	  300392	      3995 ns/op	 126.92 MB/s	       0 B/op	       0 allocs/op
BenchmarkString_Mixed/go-runewidth-8       	  257757	      4636 ns/op	 109.36 MB/s	       0 B/op	       0 allocs/op
BenchmarkString_ControlChars/displaywidth-8         	 3498014	       343.5 ns/op	  96.08 MB/s	       0 B/op	       0 allocs/op
BenchmarkString_ControlChars/go-runewidth-8         	 3278515	       367.3 ns/op	  89.84 MB/s	       0 B/op	       0 allocs/op
```

I use a similar technique in [this grapheme cluster library](https://github.com/clipperhouse/uax29).
