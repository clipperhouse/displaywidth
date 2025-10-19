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
	EastAsianWidth     bool
	StrictEmojiNeutral bool
}

var DefaultOptions = Options{
	EastAsianWidth:     false,
	StrictEmojiNeutral: true,
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

const defaultWidth = 1

// is returns true if the property flag is set
func (p property) is(flag property) bool {
	return p&flag != 0
}

const (
	// if the trie values are changed, be careful to update these

	// VARIATION SELECTOR-16 (U+FE0F) requests emoji presentation
	_VS16 property = 1 << 4 // force emoji presentation (width 2)

	// VARIATION SELECTOR-15 (U+FE0E) requests text presentation
	_VS15 property = 1 << 5 // force text presentation (width 1)
)

// lookupProperties returns the properties for the first character in a string
func lookupProperties[T stringish.Interface](s T) property {
	if len(s) == 0 {
		return 0
	}

	var p property
	var size int

	// Determine base properties for the first rune
	b := s[0]
	if b < utf8.RuneSelf { // Single-byte ASCII
		size = 1
		if isASCIIControl(b) {
			return _ZeroWidth
		}
		p = 0 // default width will be 1
	} else {
		// Use the generated trie for lookup
		props, n := lookup(s)
		p = property(props)
		size = n
	}

	// After the first code point, check for VS15 (U+FE0E) or VS16 (U+FE0F)
	// encoded as 0xEF 0xB8 0x8E/0x8F immediately following.
	if size > 0 && len(s) >= size+3 {
		if s[size] == 0xEF && s[size+1] == 0xB8 {
			switch s[size+2] {
			case 0x8F: // VS16: request emoji presentation
				p |= _VS16

				// Special-case keycap sequences: base + VS16 + U+20E3
				// U+20E3 encodes as 0xE2 0x83 0xA3
				if len(s) >= size+6 {
					if s[size+3] == 0xE2 && s[size+4] == 0x83 && s[size+5] == 0xA3 {
						p |= _VS16
					}
				}
			case 0x8E: // VS15: request text presentation
				p |= _VS15
			}
		}
	}

	return p
}

// width determines the display width of a character based on its properties
// and configuration options
func (p property) width(options Options) int {
	if p == 0 {
		// Character not in trie, use default behavior
		return defaultWidth
	}

	if p.is(_ZeroWidth) {
		return 0
	}

	// Explicit presentation overrides from VS come first.
	if p.is(_VS16) {
		return 2
	}
	if p.is(_VS15) {
		return 1
	}

	if options.EastAsianWidth {
		if p.is(_East_Asian_Ambiguous) {
			return 2
		}
		if p.is(_East_Asian_Ambiguous|_Emoji) && !options.StrictEmojiNeutral {
			return 2
		}
	}

	if p.is(_East_Asian_Full_Wide) {
		return 2
	}

	// Default width for all other characters
	return defaultWidth
}
