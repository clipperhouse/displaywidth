package displaywidth

import (
	"unicode/utf8"

	"github.com/clipperhouse/stringish"
	"github.com/clipperhouse/uax29/v2/graphemes"
)

// String calculates the display width of a string,
// by iterating over grapheme clusters in the string
// and summing their widths.
func String(s string) int {
	return DefaultOptions.String(s)
}

// Bytes calculates the display width of a []byte,
// by iterating over grapheme clusters in the byte slice
// and summing their widths.
func Bytes(s []byte) int {
	return DefaultOptions.Bytes(s)
}

// Rune calculates the display width of a rune. You
// should almost certainly use [String] or [Bytes] for
// most purposes.
//
// The smallest unit of display width is a grapheme
// cluster, not a rune. Iterating over runes to measure
// width is incorrect in many cases.
func Rune(r rune) int {
	return DefaultOptions.Rune(r)
}

// Options allows you to specify the treatment of ambiguous East Asian
// characters. When EastAsianWidth is false (default), ambiguous East Asian
// characters are treated as width 1. When EastAsianWidth is true, ambiguous
// East Asian characters are treated as width 2.
type Options struct {
	EastAsianWidth bool
}

// DefaultOptions is the default options for the display width
// calculation, which is EastAsianWidth: false.
var DefaultOptions = Options{EastAsianWidth: false}

// graphemeWidth returns the display width of a grapheme cluster.
// The passed string must be a single grapheme cluster.
func graphemeWidth[T stringish.Interface](s T, options Options) int {
	return lookupProperties(s).width(options)
}

// Graphemes is a iterator over grapheme clusters.
//
// Iterate using the Next method, and get the width of the current grapheme
// using the Width method.
type Graphemes[T stringish.Interface] struct {
	iter    graphemes.Iterator[T]
	options Options
}

// Next advances the iterator to the next grapheme cluster.
func (g *Graphemes[T]) Next() bool {
	return g.iter.Next()
}

// Value returns the current grapheme cluster.
func (g *Graphemes[T]) Value() T {
	return g.iter.Value()
}

// Width returns the display width of the current grapheme cluster.
func (g *Graphemes[T]) Width() int {
	return graphemeWidth(g.Value(), g.options)
}

// StringGraphemes returns an iterator over grapheme clusters for the given
// string.
//
// Iterate using the Next method, and get the width of the current grapheme
// using the Width method.
func StringGraphemes(s string) Graphemes[string] {
	return DefaultOptions.StringGraphemes(s)
}

// StringGraphemes returns an iterator over grapheme clusters for the given
// string, with the given options.
//
// Iterate using the Next method, and get the width of the current grapheme
// using the Width method.
func (options Options) StringGraphemes(s string) Graphemes[string] {
	return Graphemes[string]{
		iter:    graphemes.FromString(s),
		options: options,
	}
}

// BytesGraphemes returns an iterator over grapheme clusters for the given
// []byte.
//
// Iterate using the Next method, and get the width of the current grapheme
// using the Width method.
func BytesGraphemes(s []byte) Graphemes[[]byte] {
	return DefaultOptions.BytesGraphemes(s)
}

// BytesGraphemes returns an iterator over grapheme clusters for the given
// []byte, with the given options.
//
// Iterate using the Next method, and get the width of the current grapheme
// using the Width method.
func (options Options) BytesGraphemes(s []byte) Graphemes[[]byte] {
	return Graphemes[[]byte]{
		iter:    graphemes.FromBytes(s),
		options: options,
	}
}

// String calculates the display width of a string, for the given options, by
// iterating over grapheme clusters in the string and summing their widths.
func (options Options) String(s string) int {
	switch len(s) {
	case 0:
		return 0
	case 1:
		return graphemeWidth(s, options)
	}

	width := 0
	g := graphemes.FromString(s)
	for g.Next() {
		width += graphemeWidth(g.Value(), options)
	}
	return width
}

// Bytes calculates the display width of a []byte, for the given options, by
// iterating over grapheme clusters in the slice and summing their widths.
func (options Options) Bytes(s []byte) int {
	switch len(s) {
	case 0:
		return 0
	case 1:
		return graphemeWidth(s, options)
	}

	width := 0
	g := graphemes.FromBytes(s)
	for g.Next() {
		width += graphemeWidth(g.Value(), options)
	}
	return width
}

// Rune calculates the display width of a rune, for the given options.
//
// You should almost certainly use [String] or [Bytes] for most purposes.
//
// The smallest unit of display width is a grapheme cluster, not a rune.
// Iterating over runes to measure width is incorrect in many cases.
func (options Options) Rune(r rune) int {
	if r < utf8.RuneSelf {
		if isASCIIControl(byte(r)) {
			return 0
		}
		return 1
	}

	// Surrogates (U+D800-U+DFFF) are invalid UTF-8.
	if r >= 0xD800 && r <= 0xDFFF {
		return 0
	}

	var buf [4]byte
	n := utf8.EncodeRune(buf[:], r)

	// Skip the grapheme iterator
	return lookupProperties(buf[:n]).width(options)
}

func isASCIIControl(b byte) bool {
	return b < 0x20 || b == 0x7F
}

// isRIPrefix checks if the slice matches the Regional Indicator prefix
// (F0 9F 87). It assumes len(s) >= 3.
func isRIPrefix[T stringish.Interface](s T) bool {
	return s[0] == 0xF0 && s[1] == 0x9F && s[2] == 0x87
}

// isVS16 checks if the slice matches VS16 (U+FE0F) UTF-8 encoding
// (EF B8 8F). It assumes len(s) >= 3.
func isVS16[T stringish.Interface](s T) bool {
	return s[0] == 0xEF && s[1] == 0xB8 && s[2] == 0x8F
}

// lookupProperties returns the properties for the first character in a string
func lookupProperties[T stringish.Interface](s T) property {
	l := len(s)

	if l == 0 {
		return 0
	}

	b := s[0]
	if isASCIIControl(b) {
		return _Zero_Width
	}

	if b < utf8.RuneSelf {
		// Check for variation selector after ASCII (e.g., keycap sequences like 1️⃣)
		if l >= 4 {
			// Subslice may help eliminate bounds checks
			vs := s[1:4]
			if isVS16(vs) {
				// VS16 requests emoji presentation (width 2)
				return _Emoji
			}
			// VS15 (0x8E) requests text presentation but does not affect width,
			// in my reading of Unicode TR51. Falls through to _Default.
		}
		return _Default
	}

	// Regional indicator pair (flag)
	if l >= 8 {
		// Subslice may help eliminate bounds checks
		ri := s[:8]
		// First rune
		if isRIPrefix(ri[0:3]) {
			b3 := ri[3]
			if b3 >= 0xA6 && b3 <= 0xBF {
				// Second rune
				if isRIPrefix(ri[4:7]) {
					b7 := ri[7]
					if b7 >= 0xA6 && b7 <= 0xBF {
						return _Emoji
					}
				}
			}
		}
	}

	p, sz := lookup(s)

	// Variation Selectors
	if sz > 0 && l >= sz+3 {
		// Subslice may help eliminate bounds checks
		vs := s[sz : sz+3]
		if isVS16(vs) {
			// VS16 requests emoji presentation (width 2)
			return _Emoji
		}
		// VS15 (0x8E) requests text presentation but does not affect width,
		// in my reading of Unicode TR51. Falls through to return the base
		// character's property.
	}

	return property(p)
}

const _Default property = 0

// a jump table of sorts, instead of a switch
var widthTable = [5]int{
	_Default:              1,
	_Zero_Width:           0,
	_East_Asian_Wide:      2,
	_East_Asian_Ambiguous: 1,
	_Emoji:                2,
}

// width determines the display width of a character based on its properties
// and configuration options
func (p property) width(options Options) int {
	if options.EastAsianWidth && p == _East_Asian_Ambiguous {
		return 2
	}

	return widthTable[p]
}
