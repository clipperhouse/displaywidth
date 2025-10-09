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
BenchmarkStringDefault/displaywidth-8         	  112932	     10611 ns/op	 158.98 MB/s	       0 B/op	       0 allocs/op
BenchmarkStringDefault/go-runewidth-8         	   83205	     14570 ns/op	 115.78 MB/s	       0 B/op	       0 allocs/op
BenchmarkString_EAW/displaywidth-8            	  108841	     10915 ns/op	 154.55 MB/s	       0 B/op	       0 allocs/op
BenchmarkString_EAW/go-runewidth-8            	   48816	     25474 ns/op	  66.22 MB/s	       0 B/op	       0 allocs/op
BenchmarkString_StrictEmoji/displaywidth-8    	  113624	     10896 ns/op	 154.82 MB/s	       0 B/op	       0 allocs/op
BenchmarkString_StrictEmoji/go-runewidth-8    	   82554	     14486 ns/op	 116.46 MB/s	       0 B/op	       0 allocs/op
BenchmarkString_ASCII/displaywidth-8          	 1000000	      1081 ns/op	 118.38 MB/s	       0 B/op	       0 allocs/op
BenchmarkString_ASCII/go-runewidth-8          	 1000000	      1168 ns/op	 109.58 MB/s	       0 B/op	       0 allocs/op
BenchmarkString_Unicode/displaywidth-8        	 1344934	       889.3 ns/op	 149.56 MB/s	       0 B/op	       0 allocs/op
BenchmarkString_Unicode/go-runewidth-8        	  792960	      1490 ns/op	  89.27 MB/s	       0 B/op	       0 allocs/op
BenchmarkStringWidth_Emoji/displaywidth-8     	  396948	      3045 ns/op	 237.75 MB/s	       0 B/op	       0 allocs/op
BenchmarkStringWidth_Emoji/go-runewidth-8     	  248937	      4860 ns/op	 148.98 MB/s	       0 B/op	       0 allocs/op
BenchmarkString_Mixed/displaywidth-8          	  301688	      3956 ns/op	 128.15 MB/s	       0 B/op	       0 allocs/op
BenchmarkString_Mixed/go-runewidth-8          	  258480	      4675 ns/op	 108.45 MB/s	       0 B/op	       0 allocs/op
BenchmarkString_ControlChars/displaywidth-8   	 3741540	       320.9 ns/op	 102.85 MB/s	       0 B/op	       0 allocs/op
BenchmarkString_ControlChars/go-runewidth-8   	 3251212	       367.5 ns/op	  89.80 MB/s	       0 B/op	       0 allocs/op
```

I use a similar technique in [this grapheme cluster library](https://github.com/clipperhouse/uax29).
