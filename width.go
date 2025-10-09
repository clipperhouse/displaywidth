package displaywidth

import (
	"github.com/clipperhouse/displaywidth/internal/stringish"
	"github.com/clipperhouse/uax29/v2/graphemes"
)

func String(s string) int {
	return StringOptions(s, DefaultOptions)
}

func Bytes(s []byte) int {
	return BytesOptions(s, DefaultOptions)
}

// StringOptions calculates the display width of a string
func StringOptions(s string, options Options) int {
	if len(s) == 0 {
		return 0
	}

	total := 0
	g := graphemes.FromString(s)
	for g.Next() {
		// Look up character properties from trie for the first character in the grapheme cluster
		props, _ := lookupProperties(g.Value())
		total += props.width(options)
	}
	return total
}

// BytesOptions calculates the display width of a []byte
func BytesOptions(s []byte, options Options) int {
	if len(s) == 0 {
		return 0
	}

	total := 0
	g := graphemes.FromBytes(s)
	for g.Next() {
		// Look up character properties from trie for the first character in the grapheme cluster
		props, _ := lookupProperties(g.Value())
		total += props.width(options)
	}
	return total
}

type Options struct {
	EastAsianWidth     bool
	StrictEmojiNeutral bool
}

var DefaultOptions = Options{
	EastAsianWidth:     false,
	StrictEmojiNeutral: true,
}

const defaultWidth = 1

// is returns true if the property flag is set
func (p property) is(flag property) bool {
	return p&flag != 0
}

// lookupProperties returns the properties for the first character in a string
func lookupProperties[T stringish.Interface](s T) (property, int) {
	if len(s) == 0 {
		return 0, 0
	}

	// Fast path for ASCII characters (single byte)
	b := s[0]
	if b < 0x80 { // Single-byte ASCII
		if b < 0x20 || b == 0x7F {
			// Control characters (0x00-0x1F) and DEL (0x7F) - width 0
			return _ControlChar, 1
		}
		// ASCII printable characters (0x20-0x7E) - width 1
		// Return 0 properties, width calculation will default to 1
		return 0, 1
	}

	// Use the generated trie for lookup
	props, size := lookup(s)
	return props, size
}

// width determines the display width of a character based on its properties
// and configuration options
func (p property) width(options Options) int {
	if p == 0 {
		// Character not in trie, use default behavior
		return defaultWidth
	}

	if p.is(_ControlChar | _CombiningMark | _ZeroWidth) {
		return 0
	}

	if options.EastAsianWidth {
		if p.is(_EAW_Ambiguous) {
			return 2
		}
		if p.is(_EAW_Ambiguous|_Emoji) && !options.StrictEmojiNeutral {
			return 2
		}
	}

	if p.is(_EAW_Fullwidth | _EAW_Wide) {
		return 2
	}

	// Default width for all other characters
	return defaultWidth
}
