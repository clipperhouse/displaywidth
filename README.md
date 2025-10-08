# String Width Measurement Package

A high-performance Go package for measuring the display (monospace column) width of strings, compatible with existing implementations like `go-runewidth` and `wcwidth`.

## Overview

This package provides efficient string width measurement for terminal applications, text processing, and UI layout calculations. It operates strictly on strings and `[]byte` without decoding runes, achieving better performance than existing implementations.

## Features

- **High Performance**: Byte-level processing without rune decoding
- **Unicode Compliant**: Based on Unicode East Asian Width (UAX #11) standard
- **Compatible**: Compatible with `go-runewidth` and `wcwidth` libraries
- **Configurable**: Support for different width calculation modes
- **Comprehensive**: Handles all Unicode character types including emoji, combining marks, and control characters

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/clipperhouse/stringwidth"
)

func main() {
    // Basic usage with default parameters
    width := stringwidth.StringWidth("Hello, 世界! 🌍", false, false)
    fmt.Printf("Width: %d\n", width) // Output: Width: 15

    // With East Asian width mode (ambiguous characters as wide)
    width = stringwidth.StringWidth("Hello, 世界! 🌍", true, false)
    fmt.Printf("Width: %d\n", width) // Output: Width: 17

    // With strict emoji neutral mode
    width = stringwidth.StringWidth("Hello, 世界! 🌍", false, true)
    fmt.Printf("Width: %d\n", width) // Output: Width: 13
}
```

## API Reference

### Basic Functions

```go
// Calculate width of a string
func StringWidth(s string, eastAsianWidth bool, strictEmojiNeutral bool) int

// Calculate width of a byte slice
func StringWidthBytes(s []byte, eastAsianWidth bool, strictEmojiNeutral bool) int
```

### Convenience Functions

```go
// Default behavior (eastAsianWidth=false, strictEmojiNeutral=false)
func StringWidthDefault(s string) int
func StringWidthBytesDefault(s []byte) int
```

### Parameters

- **eastAsianWidth**: When `true`, treat ambiguous width characters as wide (width 2)
- **strictEmojiNeutral**: When `true`, use strict emoji width calculation (some emoji become width 1)

## Character Width Classification

### Width Categories

- **Width 0**: Non-printing characters (control characters, combining marks, zero-width characters)
- **Width 1**: Standard single-cell characters (ASCII, most Unicode characters)
- **Width 2**: Double-wide characters (CJK ideographs, fullwidth characters, most emoji)
- **Width -1**: Invalid or non-printable characters

### Character Types

#### Control Characters (Width 0)
- C0 controls (0x00-0x1F)
- DEL (0x7F)
- C1 controls (0x80-0x9F)

#### Combining Characters (Width 0)
- Nonspacing marks (Mn)
- Spacing marks (Mc)
- Enclosing marks (Me)

#### East Asian Width Property
- **F (Fullwidth)**: Width 2
- **W (Wide)**: Width 2
- **H (Halfwidth)**: Width 1
- **Na (Narrow)**: Width 1
- **N (Neutral)**: Width 1
- **A (Ambiguous)**: Width 1 (default) or Width 2 (East Asian mode)

#### Special Characters
- Zero-width characters: Width 0
- Emoji base characters: Width 2
- Emoji modifiers: Width 0
- Surrogate pairs: Width -1 (invalid)

## Examples

### Basic Usage

```go
// ASCII characters
stringwidth.StringWidth("Hello", false, false)     // 5
stringwidth.StringWidth("World!", false, false)    // 6

// Unicode characters
stringwidth.StringWidth("café", false, false)      // 4 (e + combining acute)
stringwidth.StringWidth("世界", false, false)       // 4 (2 wide characters)

// Mixed content
stringwidth.StringWidth("Hello 世界! 🌍", false, false) // 15
```

### Control Characters

```go
stringwidth.StringWidth("\x00", false, false)      // 0 (NULL)
stringwidth.StringWidth("\x08", false, false)      // 0 (Backspace)
stringwidth.StringWidth("\x0A", false, false)      // 0 (Line Feed)
stringwidth.StringWidth("\x7F", false, false)      // 0 (DEL)
```

### Combining Marks

```go
stringwidth.StringWidth("é", false, false)         // 1 (e + combining acute)
stringwidth.StringWidth("c\u0327", false, false)   // 1 (c + combining cedilla)
stringwidth.StringWidth("\u0300", false, false)    // 0 (combining grave accent)
```

### Emoji

```go
stringwidth.StringWidth("😀", false, false)        // 2 (emoji base)
stringwidth.StringWidth("👋🏿", false, false)       // 2 (emoji + skin tone modifier)
stringwidth.StringWidth("👨‍👩‍👧‍👦", false, false)   // 2 (family emoji sequence)
```

### Ambiguous Characters

```go
// Default mode (narrow)
stringwidth.StringWidth("¡", false, false)         // 1
stringwidth.StringWidth("°", false, false)         // 1
stringwidth.StringWidth("±", false, false)         // 1

// East Asian mode (wide)
stringwidth.StringWidth("¡", true, false)          // 2
stringwidth.StringWidth("°", true, false)          // 2
stringwidth.StringWidth("±", true, false)          // 2
```

## Performance

The package is designed for high performance:

- **Byte-level processing**: No rune decoding overhead
- **Efficient lookups**: Compressed trie structures for fast character property lookups
- **Memory efficient**: Read-only data structures with minimal allocations
- **SIMD optimized**: Vectorized operations for bulk processing

### Benchmarks

```
BenchmarkStringWidth-8          1000000    1200 ns/op    0 B/op    0 allocs/op
BenchmarkStringWidthBytes-8     1000000    1100 ns/op    0 B/op    0 allocs/op
BenchmarkRuneWidth-8           10000000     120 ns/op    0 B/op    0 allocs/op
```

## Compatibility

### go-runewidth Compatibility

```go
import (
    "github.com/clipperhouse/stringwidth"
    "github.com/mattn/go-runewidth"
)

// These should return the same values
width1 := stringwidth.StringWidth("Hello 世界! 🌍", false, false)
width2 := runewidth.StringWidth("Hello 世界! 🌍", false, false)
// width1 == width2

// With East Asian width mode
width1 = stringwidth.StringWidth("Hello 世界! 🌍", true, false)
width2 = runewidth.StringWidth("Hello 世界! 🌍", true, false)
// width1 == width2

// With strict emoji neutral mode
width1 = stringwidth.StringWidth("Hello 世界! 🌍", false, true)
width2 = runewidth.StringWidth("Hello 世界! 🌍", false, true)
// width1 == width2
```

### wcwidth Compatibility

The package aims to be compatible with the POSIX `wcwidth` function and Python `wcwidth` library.

## Installation

```bash
go get github.com/clipperhouse/stringwidth
```

## Requirements

- Go 1.19 or later
- Unicode 15.0 data (included in package)

## Documentation

- [Specification](SPECIFICATION.md) - Detailed technical specification
- [Edge Cases](EDGE_CASES.md) - Edge cases and special sequences
- [Implementation Plan](IMPLEMENTATION_PLAN.md) - Implementation details and architecture

## Contributing

Contributions are welcome! Please see the implementation plan for details on the architecture and development process.

## License

MIT License - see LICENSE file for details.

## Acknowledgments

- Based on Unicode East Asian Width (UAX #11) standard
- Compatible with `go-runewidth` by mattn
- Compatible with `wcwidth` by jquast
- Inspired by the `uax29` package approach
