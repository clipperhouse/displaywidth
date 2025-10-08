package stringwidth

// Global trie instance
var trie = &stringWidthTrie{}

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

// LookupCharPropertiesBytes returns the properties for the first character in a byte slice
func LookupCharPropertiesBytes(s []byte) (property, int) {
	if len(s) == 0 {
		return 0, 0
	}

	// Use the generated trie for lookup
	props, size := trie.lookup([]byte(s))
	return props, size
}

// LookupCharPropertiesString returns the properties for the first character in a string
func LookupCharPropertiesString(s string) (property, int) {
	if len(s) == 0 {
		return 0, 0
	}

	// Use the generated trie for lookup
	props, size := trie.lookup([]byte(s))
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

	// Handle emoji
	if props.IsEmoji() {
		if strictEmojiNeutral {
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

// processStringWidth calculates the total width of a string
func processStringWidth(s string, eastAsianWidth bool, strictEmojiNeutral bool) int {
	if len(s) == 0 {
		return 0
	}

	totalWidth := 0
	pos := 0

	for pos < len(s) {
		// Look up character properties from trie
		props, trieSize := LookupCharPropertiesString(s[pos:])

		// If trie lookup failed (size 0), invalid UTF-8 sequence
		if trieSize == 0 {
			return -1
		}

		// Calculate width based on properties
		var width int
		if props == 0 {
			// Character not in trie, use default behavior
			width = getDefaultWidth()
		} else {
			// Use trie properties to calculate width
			width = calculateWidth(props, eastAsianWidth, strictEmojiNeutral)
		}

		totalWidth += width
		pos += trieSize
	}

	return totalWidth
}

// processBytesWidth calculates the total width of a byte slice
func processBytesWidth(b []byte, eastAsianWidth bool, strictEmojiNeutral bool) int {
	if len(b) == 0 {
		return 0
	}

	totalWidth := 0
	offset := 0

	for offset < len(b) {
		// Look up character properties from trie
		props, trieSize := LookupCharPropertiesBytes(b[offset:])

		// If trie lookup failed (size 0), invalid UTF-8 sequence
		if trieSize == 0 {
			return -1
		}

		// Calculate width based on properties
		var width int
		if props == 0 {
			// Character not in trie, use default behavior
			width = getDefaultWidth()
		} else {
			// Use trie properties to calculate width
			width = calculateWidth(props, eastAsianWidth, strictEmojiNeutral)
		}

		totalWidth += width
		offset += trieSize
	}

	return totalWidth
}

// StringWidth calculates the display width of a string
// eastAsianWidth: when true, treat ambiguous width characters as wide (width 2)
// strictEmojiNeutral: when true, use strict emoji width calculation (some emoji become width 1)
func StringWidth(s string, eastAsianWidth bool, strictEmojiNeutral bool) int {
	return processStringWidth(s, eastAsianWidth, strictEmojiNeutral)
}

// StringWidthBytes calculates the display width of a byte slice
// eastAsianWidth: when true, treat ambiguous width characters as wide (width 2)
// strictEmojiNeutral: when true, use strict emoji width calculation (some emoji become width 1)
func StringWidthBytes(b []byte, eastAsianWidth bool, strictEmojiNeutral bool) int {
	return processBytesWidth(b, eastAsianWidth, strictEmojiNeutral)
}
