package displaywidth

import (
	"unicode/utf8"

	"github.com/clipperhouse/stringish"
	"github.com/clipperhouse/uax29/v2/graphemes"
)

// String calculates the display width of a string
// using the [DefaultOptions]
func String(s string) int {
	return DefaultOptions.String(s)
}

// Bytes calculates the display width of a []byte
// using the [DefaultOptions]
func Bytes(s []byte) int {
	return DefaultOptions.Bytes(s)
}

func Rune(r rune) int {
	return DefaultOptions.Rune(r)
}

type Options struct {
	EastAsianWidth bool
}

var DefaultOptions = Options{
	EastAsianWidth: false,
}

// String calculates the display width of a string
// for the given options
func (options Options) String(s string) int {
	if len(s) == 0 {
		return 0
	}

	total := 0
	g := graphemes.FromString(s)
	for g.Next() {
		// The first character in the grapheme cluster determines the width;
		// we use lookupProperties which can consider immediate VS15/VS16.
		props := lookupProperties(g.Value())
		total += props.width(options)
	}
	return total
}

// BytesOptions calculates the display width of a []byte
// for the given options
func (options Options) Bytes(s []byte) int {
	if len(s) == 0 {
		return 0
	}

	total := 0
	g := graphemes.FromBytes(s)
	for g.Next() {
		// The first character in the grapheme cluster determines the width;
		// we use lookupProperties which can consider immediate VS15/VS16.
		props := lookupProperties(g.Value())
		total += props.width(options)
	}
	return total
}

func (options Options) Rune(r rune) int {
	// Fast path for ASCII
	if r < utf8.RuneSelf {
		if isASCIIControl(byte(r)) {
			// Control (0x00-0x1F) and DEL (0x7F)
			return 0
		}
		// ASCII printable (0x20-0x7E)
		return 1
	}

	// Surrogates (U+D800-U+DFFF) are invalid UTF-8 and have zero width
	// Other packages might turn them into the replacement character (U+FFFD)
	// in which case, we won't see it.
	if r >= 0xD800 && r <= 0xDFFF {
		return 0
	}

	// Stack-allocated to avoid heap allocation
	var buf [4]byte // UTF-8 is at most 4 bytes
	n := utf8.EncodeRune(buf[:], r)
	// Skip the grapheme iterator and directly lookup properties
	props := lookupProperties(buf[:n])
	return props.width(options)
}

func isASCIIControl(b byte) bool {
	return b < 0x20 || b == 0x7F
}

// isRIPrefix checks if the slice matches the Regional Indicator prefix
// (F0 9F 87). It assumes len(s) >= 3.
func isRIPrefix[T stringish.Interface](s T) bool {
	return s[0] == 0xF0 && s[1] == 0x9F && s[2] == 0x87
}

// isVSPrefix checks if the slice matches the Variation Selector prefix
// (EF B8). It assumes len(s) >= 2.
func isVSPrefix[T stringish.Interface](s T) bool {
	return s[0] == 0xEF && s[1] == 0xB8
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
			if isVSPrefix(vs) {
				switch vs[2] {
				case 0x8E:
					// VS15 requests text presentation (width 1)
					return _Default
				case 0x8F:
					// VS16 requests emoji presentation (width 2)
					return _East_Asian_Wide
				}
			}
		}
		return _Default
	}

	// Regional indicator pair (flag)
	if l >= 8 {
		// Subslice may help eliminate bounds checks
		ri := s[:8]
		if isRIPrefix(ri[0:3]) {
			b3 := ri[3]
			if b3 >= 0xA6 && b3 <= 0xBF && isRIPrefix(ri[4:7]) {
				b7 := ri[7]
				if b7 >= 0xA6 && b7 <= 0xBF {
					return _East_Asian_Wide
				}
			}
		}
	}

	props, size := lookup(s)
	p := property(props)

	// Variation Selectors
	if size > 0 && l >= size+3 {
		// Subslice may help eliminate bounds checks
		vs := s[size : size+3]
		if isVSPrefix(vs) {
			switch vs[2] {
			case 0x8E:
				// VS15 requests text presentation.
				// For emoji-only characters (Extended_Pictographic +
				// Emoji_Presentation, but not East Asian Wide), text
				// presentation means width 1
				if p == _Emoji_Only {
					return _Default
				}
				// Otherwise, preserve the base character's width:
				// East Asian Wide remains wide, East Asian Ambiguous remains
				// ambiguous.
				return p
			case 0x8F:
				// VS16 requests emoji presentation (width 2)
				return _East_Asian_Wide
			}
		}
	}

	return p
}

const _Default property = 0

// a jump table of sorts, instead of a switch
var widthTable = [6]int{
	_Default:              1,
	_Zero_Width:           0,
	_East_Asian_Wide:      2,
	_East_Asian_Ambiguous: 1,
	_Emoji_Only:           2,
}

// width determines the display width of a character based on its properties
// and configuration options
func (p property) width(options Options) int {
	if options.EastAsianWidth && p == _East_Asian_Ambiguous {
		return 2
	}

	return widthTable[p]
}
