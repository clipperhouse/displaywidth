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
		{"NULL", 0x00, _ControlChar, "NULL character"},
		{"Backspace", 0x08, _ControlChar, "Backspace"},
		{"Tab", 0x09, _ControlChar, "Tab"},
		{"Line Feed", 0x0A, _ControlChar, "Line Feed"},
		{"Carriage Return", 0x0D, _ControlChar, "Carriage Return"},
		{"DEL", 0x7F, _ControlChar, "DEL character"},

		// ASCII printable characters (no properties, default to width 1)
		{"Space", 0x20, 0, "Space character"},
		{"A", 'A', 0, "Latin A"},
		{"!", '!', 0, "Exclamation mark"},
		{"0", '0', 0, "Digit 0"},

		// East Asian Wide characters
		{"Chinese 中", '中', _EAW_Wide, "Chinese character"},
		{"Japanese あ", 'あ', _EAW_Wide, "Hiragana character"},
		{"Korean 가", '가', _EAW_Wide, "Hangul character"},

		// East Asian Fullwidth characters
		{"Fullwidth A", 'Ａ', _EAW_Fullwidth, "Fullwidth Latin A"},
		{"Fullwidth !", '！', _EAW_Fullwidth, "Fullwidth exclamation"},

		// East Asian Ambiguous characters
		{"Inverted !", '¡', _EAW_Ambiguous, "Inverted exclamation mark"},
		{"Degree", '°', _EAW_Ambiguous, "Degree sign"},
		{"Plus-Minus", '±', _EAW_Ambiguous, "Plus-minus sign"},

		// Combining marks
		{"Combining Grave", 0x0300, _CombiningMark, "Combining grave accent"},
		{"Combining Acute", 0x0301, _CombiningMark, "Combining acute accent"},
		{"Combining Tilde", 0x0303, _CombiningMark, "Combining tilde"},

		// Zero-width characters
		{"Zero Width Space", 0x200B, _ZeroWidth, "Zero width space"},
		{"Zero Width Joiner", 0x200D, _ZeroWidth, "Zero width joiner"},
		{"Zero Width Non-Joiner", 0x200C, _ZeroWidth, "Zero width non-joiner"},
		{"Word Joiner", 0x2060, _ZeroWidth, "Word joiner"},
		{"Function Application", 0x2061, _ZeroWidth, "Function application"},
		{"Zero Width No-Break Space", 0xFEFF, _ZeroWidth, "Zero width no-break space"},

		// Emoji (basic detection)
		{"Grinning Face", 0x1F600, _Emoji, "Grinning face emoji"},
		{"Party Popper", 0x1F389, _Emoji, "Party popper emoji"},
		{"Rocket", 0x1F680, _Emoji, "Rocket emoji"},
		{"Sun", 0x2600, _Emoji, "Sun emoji"},
		{"Check Mark", 0x2705, _Emoji, "Check mark emoji"},

		// Invalid runes
		{"Replacement Char", 0xFFFD, 0, "Unicode replacement character"},
		{"Max Rune + 1", 0x110000, 0, "Invalid rune beyond Unicode range"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := lookupProperties(string(tt.r))

			// For zero expected value, check that no special properties are set
			if tt.expected == 0 {
				if result != 0 {
					t.Errorf("LookupCharPropertiesString(%q) = %v, want 0 (%s)",
						string(tt.r), result, tt.desc)
				}
			} else {
				// For non-zero expected value, check that the specific property is set
				if !result.is(tt.expected) {
					t.Errorf("LookupCharPropertiesString(%q) = %v, want property %v (%s)",
						string(tt.r), result, tt.expected, tt.desc)
				}
			}
		})
	}
}

func TestCharPropertiesHas(t *testing.T) {
	props := _EAW_Wide | _ControlChar | _Emoji

	if !props.is(_EAW_Wide) {
		t.Error("Has(EAW_Wide) should return true")
	}
	if !props.is(_ControlChar) {
		t.Error("Has(IsControlChar) should return true")
	}
	if !props.is(_Emoji) {
		t.Error("Has(IsEmoji) should return true")
	}
	if props.is(_EAW_Ambiguous) {
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
		{'中', []property{_EAW_Wide}, "Chinese character should be wide"},
		{'Ａ', []property{_EAW_Fullwidth}, "Fullwidth A should be fullwidth"},
		{'¡', []property{_EAW_Ambiguous}, "Inverted exclamation should be ambiguous"},
		{0x0300, []property{_CombiningMark}, "Combining grave should be combining"},
		{0x200B, []property{_ZeroWidth}, "Zero width space should be zero width"},
		{0x1F600, []property{_Emoji}, "Grinning face should be emoji"},
	}

	for _, tc := range testCases {
		props, _ := lookupProperties(string(tc.char))
		for _, expectedProp := range tc.hasProps {
			if !props.is(expectedProp) {
				t.Errorf("Character %q (%U) should have property %v (%s)",
					string(tc.char), tc.char, expectedProp, tc.desc)
			}
		}
	}
}
