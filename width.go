package displaywidth

import (
	"unicode/utf8"
	"unsafe"

	"github.com/clipperhouse/stringish"
	"github.com/clipperhouse/uax29/v2/graphemes"
)

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

// String calculates the display width of a string,
// by iterating over grapheme clusters in the string
// and summing their widths.
func String(s string) int {
	return DefaultOptions.String(s)
}

// String calculates the display width of a string, for the given options, by
// iterating over grapheme clusters in the string and summing their widths.
func (options Options) String(s string) int {
	width := 0
	pos := 0

	for pos < len(s) {
		// Try ASCII optimization (need >= 8 bytes for it to be worth it)
		asciiLen := printableASCIILength(s[pos:])
		if asciiLen >= 0 {
			width += asciiLen
			pos += asciiLen
			continue
		}

		// Not ASCII (or < 8 bytes), use grapheme parsing
		g := graphemes.FromString(s[pos:])

		hitASCII := false
		for g.Next() {
			width += graphemeWidth(g.Value(), options)
			absEnd := pos + g.End()

			// Quick check: if remaining might be an ASCII run, break to outer loop
			if len(s)-absEnd >= 8 && s[absEnd] >= 0x20 && s[absEnd] <= 0x7E {
				pos = absEnd
				hitASCII = true
				break
			}
		}

		if !hitASCII {
			// Consumed all remaining via graphemes
			break
		}
	}

	return width
}

// Bytes calculates the display width of a []byte,
// by iterating over grapheme clusters in the byte slice
// and summing their widths.
func Bytes(s []byte) int {
	return DefaultOptions.Bytes(s)
}

// Bytes calculates the display width of a []byte, for the given options, by
// iterating over grapheme clusters in the slice and summing their widths.
func (options Options) Bytes(s []byte) int {
	width := 0
	pos := 0

	for pos < len(s) {
		// Try ASCII optimization (need >= 8 bytes for it to be worth it)
		asciiLen := printableASCIILengthBytes(s[pos:])
		if asciiLen >= 0 {
			width += asciiLen
			pos += asciiLen
			continue
		}

		// Not ASCII (or < 8 bytes), use grapheme parsing
		g := graphemes.FromBytes(s[pos:])

		hitASCII := false
		for g.Next() {
			width += graphemeWidth(g.Value(), options)
			absEnd := pos + g.End()

			// Quick check: if remaining might be an ASCII run, break to outer loop
			if len(s)-absEnd >= 8 && s[absEnd] >= 0x20 && s[absEnd] <= 0x7E {
				pos = absEnd
				hitASCII = true
				break
			}
		}

		if !hitASCII {
			// Consumed all remaining via graphemes
			break
		}
	}

	return width
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

// Rune calculates the display width of a rune, for the given options.
//
// You should almost certainly use [String] or [Bytes] for most purposes.
//
// The smallest unit of display width is a grapheme cluster, not a rune.
// Iterating over runes to measure width is incorrect in many cases.
func (options Options) Rune(r rune) int {
	if r < utf8.RuneSelf {
		return asciiWidth(byte(r))
	}

	// Surrogates (U+D800-U+DFFF) are invalid UTF-8.
	if r >= 0xD800 && r <= 0xDFFF {
		return 0
	}

	var buf [4]byte
	n := utf8.EncodeRune(buf[:], r)

	// Skip the grapheme iterator
	return graphemeWidth(buf[:n], options)
}

const _Default property = 0

// TruncateString truncates a string to the given maxWidth, and appends the
// given tail if the string is truncated.
//
// It ensures the total width, including the width of the tail, is less than or
// equal to maxWidth.
func (options Options) TruncateString(s string, maxWidth int, tail string) string {
	maxWidthWithoutTail := maxWidth - options.String(tail)

	var pos, total int
	g := graphemes.FromString(s)
	for g.Next() {
		gw := graphemeWidth(g.Value(), options)
		if total+gw <= maxWidthWithoutTail {
			pos = g.End()
		}
		total += gw
		if total > maxWidth {
			return s[:pos] + tail
		}
	}
	// No truncation
	return s
}

// TruncateString truncates a string to the given maxWidth, and appends the
// given tail if the string is truncated.
//
// It ensures the total width, including the width of the tail, is less than or
// equal to maxWidth.
func TruncateString(s string, maxWidth int, tail string) string {
	return DefaultOptions.TruncateString(s, maxWidth, tail)
}

// TruncateBytes truncates a []byte to the given maxWidth, and appends the
// given tail if the []byte is truncated.
//
// It ensures the total width, including the width of the tail, is less than or
// equal to maxWidth.
func (options Options) TruncateBytes(s []byte, maxWidth int, tail []byte) []byte {
	maxWidthWithoutTail := maxWidth - options.Bytes(tail)

	var pos, total int
	g := graphemes.FromBytes(s)
	for g.Next() {
		gw := graphemeWidth(g.Value(), options)
		if total+gw <= maxWidthWithoutTail {
			pos = g.End()
		}
		total += gw
		if total > maxWidth {
			result := make([]byte, 0, pos+len(tail))
			result = append(result, s[:pos]...)
			result = append(result, tail...)
			return result
		}
	}
	// No truncation
	return s
}

// TruncateBytes truncates a []byte to the given maxWidth, and appends the
// given tail if the []byte is truncated.
//
// It ensures the total width, including the width of the tail, is less than or
// equal to maxWidth.
func TruncateBytes(s []byte, maxWidth int, tail []byte) []byte {
	return DefaultOptions.TruncateBytes(s, maxWidth, tail)
}

// graphemeWidth returns the display width of a grapheme cluster.
// The passed string must be a single grapheme cluster.
func graphemeWidth[T stringish.Interface](s T, options Options) int {
	// Optimization: no need to look up properties
	switch len(s) {
	case 0:
		return 0
	case 1:
		return asciiWidth(s[0])
	}

	p, sz := lookup(s)
	prop := property(p)

	// Variation Selector 16 (VS16) requests emoji presentation
	if prop != _Wide && sz > 0 && len(s) >= sz+3 {
		vs := s[sz : sz+3]
		if isVS16(vs) {
			prop = _Wide
		}
		// VS15 (0x8E) requests text presentation but does not affect width,
		// in my reading of Unicode TR51. Falls through to return the base
		// character's property.
	}

	if options.EastAsianWidth && prop == _East_Asian_Ambiguous {
		prop = _Wide
	}

	if prop > upperBound {
		prop = _Default
	}

	return propertyWidths[prop]
}

func asciiWidth(b byte) int {
	if b <= 0x1F || b == 0x7F {
		return 0
	}
	return 1
}

// printableASCIILength returns the length of consecutive printable ASCII bytes
// starting at the beginning of s. Returns -1 if fewer than 8 consecutive
// printable ASCII bytes are found (not worth optimizing). Uses SWAR to check
// 8 bytes at a time.
func printableASCIILength(s string) int {
	if len(s) < 8 {
		return -1
	}

	i := 0
	for ; i+8 <= len(s); i += 8 {
		x := *(*uint64)(unsafe.Add(unsafe.Pointer(unsafe.StringData(s)), i))
		// Check for non-ASCII (high bit set)
		if x&0x8080808080808080 != 0 {
			break
		}
		// Check for control chars (< 0x20): add 0x60, printable bytes overflow to set high bit
		if (x+0x6060606060606060)&0x8080808080808080 != 0x8080808080808080 {
			break
		}
		// Check for DEL (0x7F) using zero-byte detection
		xored := x ^ 0x7F7F7F7F7F7F7F7F
		if ((xored - 0x0101010101010101) & ^xored & 0x8080808080808080) != 0 {
			break
		}
	}

	// If we didn't get at least 8 bytes, not worth optimizing
	if i == 0 {
		return -1
	}

	// Check remaining bytes individually to extend the run
	for ; i < len(s); i++ {
		if b := s[i]; b < 0x20 || b > 0x7E {
			break
		}
	}

	return i
}

// printableASCIILengthBytes returns the length of consecutive printable ASCII bytes
// starting at the beginning of s. Returns -1 if fewer than 8 consecutive
// printable ASCII bytes are found (not worth optimizing). Uses SWAR to check
// 8 bytes at a time.
func printableASCIILengthBytes(s []byte) int {
	if len(s) < 8 {
		return -1
	}

	i := 0
	for ; i+8 <= len(s); i += 8 {
		x := *(*uint64)(unsafe.Pointer(&s[i]))
		// Check for non-ASCII (high bit set)
		if x&0x8080808080808080 != 0 {
			break
		}
		// Check for control chars (< 0x20): add 0x60, printable bytes overflow to set high bit
		if (x+0x6060606060606060)&0x8080808080808080 != 0x8080808080808080 {
			break
		}
		// Check for DEL (0x7F) using zero-byte detection
		xored := x ^ 0x7F7F7F7F7F7F7F7F
		if ((xored - 0x0101010101010101) & ^xored & 0x8080808080808080) != 0 {
			break
		}
	}

	// If we didn't get at least 8 bytes, not worth optimizing
	if i == 0 {
		return -1
	}

	// Check remaining bytes individually to extend the run
	for ; i < len(s); i++ {
		if b := s[i]; b < 0x20 || b > 0x7E {
			break
		}
	}

	return i
}

// isVS16 checks if the slice matches VS16 (U+FE0F) UTF-8 encoding
// (EF B8 8F). It assumes len(s) >= 3.
func isVS16[T stringish.Interface](s T) bool {
	return s[0] == 0xEF && s[1] == 0xB8 && s[2] == 0x8F
}

// propertyWidths is a jump table of sorts, instead of a switch
var propertyWidths = [4]int{
	_Default:              1,
	_Zero_Width:           0,
	_Wide:                 2,
	_East_Asian_Ambiguous: 1,
}

const upperBound = property(len(propertyWidths) - 1)
