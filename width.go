package stringwidth

import "unicode"

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

// Global trie instance
var stringWidthTrieInstance = newStringWidthTrie(0)

// LookupCharPropertiesBytes returns the properties for the first character in a byte slice
func LookupCharPropertiesBytes(s []byte) (property, int) {
	if len(s) == 0 {
		return 0, 0
	}

	// Use the generated trie for lookup
	props, size := stringWidthTrieInstance.lookup([]byte(s))
	return props, size
}

// LookupCharPropertiesString returns the properties for the first character in a string
func LookupCharPropertiesString(s string) (property, int) {
	if len(s) == 0 {
		return 0, 0
	}

	// Use the generated trie for lookup
	props, size := stringWidthTrieInstance.lookup([]byte(s))
	return props, size
}

// LookupCharProperties returns the properties for a given rune
func LookupCharProperties(r rune) property {
	// Check for invalid runes
	if r == unicode.ReplacementChar || r > unicode.MaxRune {
		return 0
	}

	// Convert rune to string and use stringish lookup
	props, _ := LookupCharPropertiesString(string(r))
	return props
}
