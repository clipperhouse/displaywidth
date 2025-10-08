package stringwidth

import (
	"testing"
)

func TestLookupCharProperties(t *testing.T) {
	tests := []struct {
		name     string
		r        rune
		expected property
		desc     string
	}{
		// Control characters
		{"NULL", 0x00, IsControlChar, "NULL character"},
		{"Backspace", 0x08, IsControlChar, "Backspace"},
		{"Tab", 0x09, IsControlChar, "Tab"},
		{"Line Feed", 0x0A, IsControlChar, "Line Feed"},
		{"Carriage Return", 0x0D, IsControlChar, "Carriage Return"},
		{"DEL", 0x7F, IsControlChar, "DEL character"},

		// ASCII printable characters (have EAW_Narrow property)
		{"Space", 0x20, EAW_Narrow, "Space character"},
		{"A", 'A', EAW_Narrow, "Latin A"},
		{"!", '!', EAW_Narrow, "Exclamation mark"},
		{"0", '0', EAW_Narrow, "Digit 0"},

		// East Asian Wide characters
		{"Chinese 中", '中', EAW_Wide, "Chinese character"},
		{"Japanese あ", 'あ', EAW_Wide, "Hiragana character"},
		{"Korean 가", '가', EAW_Wide, "Hangul character"},

		// East Asian Fullwidth characters
		{"Fullwidth A", 'Ａ', EAW_Fullwidth, "Fullwidth Latin A"},
		{"Fullwidth !", '！', EAW_Fullwidth, "Fullwidth exclamation"},

		// East Asian Ambiguous characters
		{"Inverted !", '¡', EAW_Ambiguous, "Inverted exclamation mark"},
		{"Degree", '°', EAW_Ambiguous, "Degree sign"},
		{"Plus-Minus", '±', EAW_Ambiguous, "Plus-minus sign"},

		// Combining marks
		{"Combining Grave", 0x0300, IsCombiningMark, "Combining grave accent"},
		{"Combining Acute", 0x0301, IsCombiningMark, "Combining acute accent"},
		{"Combining Tilde", 0x0303, IsCombiningMark, "Combining tilde"},

		// Zero-width characters
		{"Zero Width Space", 0x200B, IsZeroWidth, "Zero width space"},
		{"Zero Width Joiner", 0x200D, IsZeroWidth, "Zero width joiner"},
		{"Zero Width Non-Joiner", 0x200C, IsZeroWidth, "Zero width non-joiner"},
		{"Word Joiner", 0x2060, IsZeroWidth, "Word joiner"},
		{"Function Application", 0x2061, IsZeroWidth, "Function application"},
		{"Zero Width No-Break Space", 0xFEFF, IsZeroWidth, "Zero width no-break space"},

		// Emoji (basic detection)
		{"Grinning Face", 0x1F600, IsEmoji, "Grinning face emoji"},
		{"Party Popper", 0x1F389, IsEmoji, "Party popper emoji"},
		{"Rocket", 0x1F680, IsEmoji, "Rocket emoji"},
		{"Sun", 0x2600, IsEmoji, "Sun emoji"},
		{"Check Mark", 0x2705, IsEmoji, "Check mark emoji"},

		// Invalid runes
		{"Replacement Char", 0xFFFD, 0, "Unicode replacement character"},
		{"Max Rune + 1", 0x110000, 0, "Invalid rune beyond Unicode range"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LookupCharProperties(tt.r)

			// For zero expected value, check that no special properties are set
			if tt.expected == 0 {
				if result != 0 {
					t.Errorf("LookupCharProperties(%q) = %v, want 0 (%s)",
						string(tt.r), result, tt.desc)
				}
			} else {
				// For non-zero expected value, check that the specific property is set
				if !result.Has(tt.expected) {
					t.Errorf("LookupCharProperties(%q) = %v, want property %v (%s)",
						string(tt.r), result, tt.expected, tt.desc)
				}
			}
		})
	}
}

func TestCharPropertiesMethods(t *testing.T) {
	// Test IsEastAsianWide
	wideProps := EAW_Wide | EAW_Fullwidth
	if !wideProps.IsEastAsianWide() {
		t.Error("IsEastAsianWide() should return true for wide characters")
	}

	neutralProps := EAW_Neutral
	if neutralProps.IsEastAsianWide() {
		t.Error("IsEastAsianWide() should return false for neutral characters")
	}

	// Test IsEastAsianAmbiguous
	ambiguousProps := EAW_Ambiguous
	if !ambiguousProps.IsEastAsianAmbiguous() {
		t.Error("IsEastAsianAmbiguous() should return true for ambiguous characters")
	}

	// Test IsCombining
	combiningProps := IsCombiningMark
	if !combiningProps.IsCombining() {
		t.Error("IsCombining() should return true for combining marks")
	}

	// Test IsControl
	controlProps := IsControlChar
	if !controlProps.IsControl() {
		t.Error("IsControl() should return true for control characters")
	}

	// Test IsZeroWidth
	zeroWidthProps := IsZeroWidth
	if !zeroWidthProps.IsZeroWidth() {
		t.Error("IsZeroWidth() should return true for zero-width characters")
	}

	// Test IsEmoji
	emojiProps := IsEmoji
	if !emojiProps.IsEmoji() {
		t.Error("IsEmoji() should return true for emoji characters")
	}
}

func TestCharPropertiesHas(t *testing.T) {
	props := EAW_Wide | IsControlChar | IsEmoji

	if !props.Has(EAW_Wide) {
		t.Error("Has(EAW_Wide) should return true")
	}
	if !props.Has(IsControlChar) {
		t.Error("Has(IsControlChar) should return true")
	}
	if !props.Has(IsEmoji) {
		t.Error("Has(IsEmoji) should return true")
	}
	if props.Has(EAW_Ambiguous) {
		t.Error("Has(EAW_Ambiguous) should return false")
	}
}

func TestSpecificCharacters(t *testing.T) {
	// Test some specific characters that should have known properties
	testCases := []struct {
		char     rune
		hasProps []property
		desc     string
	}{
		{'中', []property{EAW_Wide}, "Chinese character should be wide"},
		{'Ａ', []property{EAW_Fullwidth}, "Fullwidth A should be fullwidth"},
		{'¡', []property{EAW_Ambiguous}, "Inverted exclamation should be ambiguous"},
		{0x0300, []property{IsCombiningMark}, "Combining grave should be combining"},
		{0x200B, []property{IsZeroWidth}, "Zero width space should be zero width"},
		{0x1F600, []property{IsEmoji}, "Grinning face should be emoji"},
	}

	for _, tc := range testCases {
		props := LookupCharProperties(tc.char)
		for _, expectedProp := range tc.hasProps {
			if !props.Has(expectedProp) {
				t.Errorf("Character %q (%U) should have property %v (%s)",
					string(tc.char), tc.char, expectedProp, tc.desc)
			}
		}
	}
}
