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

See [SPECIFICATION.md](SPECIFICATION.md) for technical details.
