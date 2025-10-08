package stringwidth

import (
	"github.com/clipperhouse/stringwidth/internal/stringish"
	"github.com/clipperhouse/uax29/v2/graphemes"
)

// Has returns true if the property flag is set
func (p property) Has(flag property) bool {
	return p&flag != 0
}

// IsEastAsianWide returns true if the character is East Asian Wide
func (p property) IsEastAsianWide() bool {
	return p.Has(EAW_Fullwidth) || p.Has(EAW_Wide)
}

// IsEastAsianAmbiguous returns true if the character is East Asian Ambiguous
func (p property) IsEastAsianAmbiguous() bool {
	return p.Has(EAW_Ambiguous)
}

// IsCombining returns true if the character is a combining mark
func (p property) IsCombining() bool {
	return p.Has(IsCombiningMark)
}

// IsControl returns true if the character is a control character
func (p property) IsControl() bool {
	return p.Has(IsControlChar)
}

// IsZeroWidth returns true if the character has zero width
func (p property) IsZeroWidth() bool {
	return p.Has(IsZeroWidth)
}

// IsEmoji returns true if the character is an emoji
func (p property) IsEmoji() bool {
	return p.Has(IsEmoji)
}

// LookupCharPropertiesString returns the properties for the first character in a string
func LookupCharProperties[T stringish.Interface](s T) (property, int) {
	if len(s) == 0 {
		return 0, 0
	}

	// Use the generated trie for lookup
	props, size := lookup(s)
	return props, size
}

// calculateWidth determines the display width of a character based on its properties
// and configuration options
func calculateWidth(props property, eastAsianWidth bool, strictEmojiNeutral bool) int {
	// Handle control characters (width 0)
	if props.IsControl() {
		return 0
	}

	// Handle combining marks (width 0)
	if props.IsCombining() {
		return 0
	}

	// Handle zero-width characters (width 0)
	if props.IsZeroWidth() {
		return 0
	}

	// Handle East Asian Ambiguous characters (before emoji check)
	if props.IsEastAsianAmbiguous() {
		if eastAsianWidth {
			return 2
		}
		return 1
	}

	// Handle emoji - match go-runewidth logic exactly
	if props.IsEmoji() {
		// go-runewidth logic: emoji get width 2 by default
		// Only ambiguous emoji get width 1 in strict mode
		if strictEmojiNeutral && props.IsEastAsianAmbiguous() {
			return 1
		}
		return 2
	}

	// Handle East Asian Width properties
	if props.IsEastAsianWide() {
		return 2
	}

	// Default width for all other characters
	return 1
}

// getDefaultWidth returns the default width for a character not in the trie
// This handles characters that aren't in our trie (props == 0)
func getDefaultWidth() int {
	// Default width for unmapped characters
	// Most characters default to width 1
	return 1
}

// processStringWidth calculates the total width of a string using grapheme clusters
func processStringWidth(s string, eastAsianWidth bool, strictEmojiNeutral bool) int {
	if len(s) == 0 {
		return 0
	}

	totalWidth := 0
	g := graphemes.FromString(s)

	for g.Next() {
		// Look up character properties from trie for the first character in the grapheme cluster
		props, _ := LookupCharProperties(g.Value())

		// Calculate width based on properties
		var chWidth int
		if props == 0 {
			// Character not in trie, use default behavior
			chWidth = getDefaultWidth()
		} else {
			// Use trie properties to calculate width
			chWidth = calculateWidth(props, eastAsianWidth, strictEmojiNeutral)
		}

		totalWidth += chWidth
	}

	return totalWidth
}

// StringWidth calculates the display width of a string
// eastAsianWidth: when true, treat ambiguous width characters as wide (width 2)
// strictEmojiNeutral: when true, use strict emoji width calculation (some emoji become width 1)
func StringWidth(s string, eastAsianWidth bool, strictEmojiNeutral bool) int {
	return processStringWidth(s, eastAsianWidth, strictEmojiNeutral)
}
