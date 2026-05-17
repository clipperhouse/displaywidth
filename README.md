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

For most purposes, you should use the `String` or `Bytes` methods. They sum
the widths of grapheme clusters in the string or byte slice.

> Note: in your application, iterating over runes to measure width is likely incorrect;
the smallest unit of display is a grapheme, not a rune.

### Iterating over graphemes

If you need the individual graphemes:

```go
import (
    "fmt"
    "github.com/clipperhouse/displaywidth"
)

func main() {
    g := displaywidth.StringGraphemes("Hello, 世界!")
    for g.Next() {
        width := g.Width()
        value := g.Value()
        // do something with the width or value
    }
}
```

### Options

Create the options you need, and then use methods on the options struct.

```go
var myOptions = displaywidth.Options{
    EastAsianWidth: true,
    ControlSequences: true,
}

width := myOptions.String("Hello, 世界!")
```

#### ControlSequences

`ControlSequences` specifies whether to ignore ECMA-48 escape sequences
when calculating the display width. When `false` (default), ANSI escape
sequences are treated as just a series of characters. When `true`, they are
treated as a single zero-width unit.

#### ControlSequences8Bit

`ControlSequences8Bit` specifies whether to ignore 8-bit ECMA-48 escape sequences
when calculating the display width. When `false` (default), these are treated
as just a series of characters. When `true`, they are treated as a single
zero-width unit.

Note: this option is ignored by the `Truncate` methods, as the concatenation
can lead to unintended UTF-8 semantics.

#### EastAsianWidth

`EastAsianWidth` defines how
[East Asian Ambiguous characters](https://www.unicode.org/reports/tr11/#Ambiguous)
are treated.

When `false` (default), East Asian Ambiguous characters are treated as width 1.
When `true`, they are treated as width 2.

You may wish to configure this based on environment variables or locale.
 `go-runewidth`, for example, does so
 [during package initialization](https://github.com/mattn/go-runewidth/blob/master/runewidth.go#L26C1-L45C2). `displaywidth` does not do this automatically, we prefer to leave it to you.


## Technical standards and compatibility

This package implements the Unicode East Asian Width standard
([UAX #11](https://www.unicode.org/reports/tr11/tr11-43.html)), and handles
[version selectors](https://en.wikipedia.org/wiki/Variation_Selectors_(Unicode_block)),
and [regional indicator pairs](https://en.wikipedia.org/wiki/Regional_indicator_symbol)
(flags). We implement [Unicode TR51](https://www.unicode.org/reports/tr51/tr51-27.html)
for emojis. We are keeping an eye on
[emerging standards](https://www.jeffquast.com/post/state-of-terminal-emulation-2025/).

For control sequences, we implement the [ECMA-48](https://ecma-international.org/publications-and-standards/standards/ecma-48/) standard for 7-bit and 8-bit control sequences.

`clipperhouse/displaywidth`, `mattn/go-runewidth`, and `rivo/uniseg` will
give the same outputs for most real-world text. Extensive details are in the
[compatibility analysis](comparison/COMPATIBILITY_ANALYSIS.md).

## Invalid UTF-8

This package does not validate UTF-8. If you pass invalid UTF-8, the results
are undefined. We fuzz against invalid UTF-8 to ensure we don't panic or
loop indefinitely.

The `ControlSequences8Bit` option means that we will segment valid 8-bit
control sequences, which are typically _not_ valid UTF-8. 8-bit control bytes
happen to also be UTF-8 continuation bytes. Use with caution.

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

BenchmarkString_Mixed/clipperhouse/displaywidth-8                 6085 ns/op         277.23 MB/s           0 B/op          0 allocs/op
BenchmarkString_Mixed/mattn/go-runewidth-8                        9970 ns/op         169.21 MB/s           0 B/op          0 allocs/op
BenchmarkString_Mixed/rivo/uniseg-8                              19060 ns/op          88.51 MB/s           0 B/op          0 allocs/op

BenchmarkString_EastAsian/clipperhouse/displaywidth-8             6118 ns/op         275.76 MB/s           0 B/op          0 allocs/op
BenchmarkString_EastAsian/mattn/go-runewidth-8                   13917 ns/op         121.22 MB/s           0 B/op          0 allocs/op
BenchmarkString_EastAsian/rivo/uniseg-8                          19263 ns/op          87.58 MB/s           0 B/op          0 allocs/op

BenchmarkString_ASCII/clipperhouse/displaywidth-8                   54.54 ns/op     2347.10 MB/s           0 B/op          0 allocs/op
BenchmarkString_ASCII/mattn/go-runewidth-8                         125.5 ns/op      1020.32 MB/s           0 B/op          0 allocs/op
BenchmarkString_ASCII/rivo/uniseg-8                               1478 ns/op          86.62 MB/s           0 B/op          0 allocs/op

BenchmarkString_Emoji/clipperhouse/displaywidth-8                 3265 ns/op         221.74 MB/s           0 B/op          0 allocs/op
BenchmarkString_Emoji/mattn/go-runewidth-8                        5110 ns/op         141.69 MB/s           0 B/op          0 allocs/op
BenchmarkString_Emoji/rivo/uniseg-8                               7137 ns/op         101.44 MB/s           0 B/op          0 allocs/op

BenchmarkRune_Mixed/clipperhouse/displaywidth-8                   3517 ns/op         479.72 MB/s           0 B/op          0 allocs/op
BenchmarkRune_Mixed/mattn/go-runewidth-8                          4746 ns/op         355.48 MB/s           0 B/op          0 allocs/op

BenchmarkRune_EastAsian/clipperhouse/displaywidth-8               3454 ns/op         488.36 MB/s           0 B/op          0 allocs/op
BenchmarkRune_EastAsian/mattn/go-runewidth-8                     11432 ns/op         147.56 MB/s           0 B/op          0 allocs/op

BenchmarkRune_ASCII/clipperhouse/displaywidth-8                    255.5 ns/op       500.88 MB/s           0 B/op          0 allocs/op
BenchmarkRune_ASCII/mattn/go-runewidth-8                           264.7 ns/op       483.48 MB/s           0 B/op          0 allocs/op

BenchmarkRune_Emoji/clipperhouse/displaywidth-8                   1320 ns/op         548.44 MB/s           0 B/op          0 allocs/op
BenchmarkRune_Emoji/mattn/go-runewidth-8                          2286 ns/op         316.72 MB/s           0 B/op          0 allocs/op

BenchmarkTruncateWithTail/clipperhouse/displaywidth-8             2495 ns/op          70.94 MB/s         192 B/op         14 allocs/op
BenchmarkTruncateWithTail/mattn/go-runewidth-8                    4569 ns/op          38.74 MB/s         192 B/op         14 allocs/op

BenchmarkTruncateWithoutTail/clipperhouse/displaywidth-8          2456 ns/op          93.25 MB/s           0 B/op          0 allocs/op
BenchmarkTruncateWithoutTail/mattn/go-runewidth-8                 5182 ns/op          44.19 MB/s           0 B/op          0 allocs/op
```

Here are some notes on [how to make Unicode things fast](https://clipperhouse.com/go-unicode/).
