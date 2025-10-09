package displaywidth

import (
	"github.com/clipperhouse/displaywidth/internal/stringish"
	"github.com/clipperhouse/uax29/v2/graphemes"
)

// String calculates the display width of a string
// eastAsianWidth: when true, treat ambiguous width characters as wide (width 2)
// strictEmojiNeutral: when true, use strict emoji width calculation (some emoji become width 1)
func String(s string, eastAsianWidth bool, strictEmojiNeutral bool) int {
	if len(s) == 0 {
		return 0
	}

	total := 0
	g := graphemes.FromString(s)
	for g.Next() {
		// Look up character properties from trie for the first character in the grapheme cluster
		props, _ := lookupProperties(g.Value())
		total += props.width(eastAsianWidth, strictEmojiNeutral)
	}
	return total
}

// Bytes calculates the display width of a []byte
// eastAsianWidth: when true, treat ambiguous width characters as wide (width 2)
// strictEmojiNeutral: when true, use strict emoji width calculation (some emoji become width 1)
func Bytes(s []byte, eastAsianWidth bool, strictEmojiNeutral bool) int {
	if len(s) == 0 {
		return 0
	}

	total := 0
	g := graphemes.FromBytes(s)
	for g.Next() {
		// Look up character properties from trie for the first character in the grapheme cluster
		props, _ := lookupProperties(g.Value())
		total += props.width(eastAsianWidth, strictEmojiNeutral)
	}
	return total
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

const controlCombiningZero = _ControlChar | _CombiningMark | _ZeroWidth

// width determines the display width of a character based on its properties
// and configuration options
func (p property) width(eastAsianWidth bool, strictEmojiNeutral bool) int {
	if p == 0 {
		// Character not in trie, use default behavior
		return defaultWidth
	}

	if p.is(controlCombiningZero) {
		return 0
	}

	// Handle East Asian Ambiguous characters (before emoji check)
	if eastAsianWidth && p.is(_EAW_Ambiguous) {
		return 2
	}

	if eastAsianWidth && p.is(_Emoji) &&
		!strictEmojiNeutral && p.is(_EAW_Ambiguous) {
		return 2
	}

	// Handle East Asian Width properties
	if p.is(_EAW_Fullwidth) || p.is(_EAW_Wide) {
		return 2
	}

	// Handle East Asian Width properties
	if p.is(_EAW_Fullwidth) || p.is(_EAW_Wide) {
		return 2
	}

	// Default width for all other characters
	return defaultWidth
}
