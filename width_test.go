package stringwidth

import (
	"testing"
)

func TestStringWidth(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		eastAsianWidth     bool
		strictEmojiNeutral bool
		expected           int
	}{
		// Basic ASCII characters
		{"empty string", "", false, false, 0},
		{"single ASCII", "a", false, false, 1},
		{"multiple ASCII", "hello", false, false, 5},
		{"ASCII with spaces", "hello world", false, false, 11},

		// Control characters (width 0)
		{"newline", "\n", false, false, 0},
		{"tab", "\t", false, false, 0},
		{"carriage return", "\r", false, false, 0},
		{"backspace", "\b", false, false, 0},

		// Mixed content
		{"ASCII with newline", "hello\nworld", false, false, 10},
		{"ASCII with tab", "hello\tworld", false, false, 10},

		// East Asian characters (should be in trie)
		{"CJK ideograph", "中", false, false, 2},
		{"CJK with ASCII", "hello中", false, false, 7},

		// Ambiguous characters
		{"ambiguous character", "★", false, false, 1},    // Default narrow
		{"ambiguous character EAW", "★", true, false, 2}, // East Asian wide

		// Emoji
		{"emoji", "😀", false, false, 2},       // Default emoji width
		{"emoji strict", "😀", false, true, 1}, // Strict emoji neutral

		// Invalid UTF-8 - the trie treats \xff as a valid character with default properties
		{"invalid UTF-8", "\xff", false, false, 1},
		{"partial UTF-8", "\xc2", false, false, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringWidth(tt.input, tt.eastAsianWidth, tt.strictEmojiNeutral)
			if result != tt.expected {
				t.Errorf("StringWidth(%q, %v, %v) = %d, want %d",
					tt.input, tt.eastAsianWidth, tt.strictEmojiNeutral, result, tt.expected)
			}
		})
	}
}

func TestStringWidthBytes(t *testing.T) {
	tests := []struct {
		name               string
		input              []byte
		eastAsianWidth     bool
		strictEmojiNeutral bool
		expected           int
	}{
		// Basic ASCII characters
		{"empty bytes", []byte{}, false, false, 0},
		{"single ASCII", []byte("a"), false, false, 1},
		{"multiple ASCII", []byte("hello"), false, false, 5},
		{"ASCII with spaces", []byte("hello world"), false, false, 11},

		// Control characters (width 0)
		{"newline", []byte("\n"), false, false, 0},
		{"tab", []byte("\t"), false, false, 0},
		{"carriage return", []byte("\r"), false, false, 0},
		{"backspace", []byte("\b"), false, false, 0},

		// Mixed content
		{"ASCII with newline", []byte("hello\nworld"), false, false, 10},
		{"ASCII with tab", []byte("hello\tworld"), false, false, 10},

		// East Asian characters (should be in trie)
		{"CJK ideograph", []byte("中"), false, false, 2},
		{"CJK with ASCII", []byte("hello中"), false, false, 7},

		// Ambiguous characters
		{"ambiguous character", []byte("★"), false, false, 1},    // Default narrow
		{"ambiguous character EAW", []byte("★"), true, false, 2}, // East Asian wide

		// Emoji
		{"emoji", []byte("😀"), false, false, 2},       // Default emoji width
		{"emoji strict", []byte("😀"), false, true, 1}, // Strict emoji neutral

		// Invalid UTF-8 - the trie treats \xff as a valid character with default properties
		{"invalid UTF-8", []byte{0xff}, false, false, 1},
		{"partial UTF-8", []byte{0xc2}, false, false, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringWidthBytes(tt.input, tt.eastAsianWidth, tt.strictEmojiNeutral)
			if result != tt.expected {
				t.Errorf("StringWidthBytes(%v, %v, %v) = %d, want %d",
					tt.input, tt.eastAsianWidth, tt.strictEmojiNeutral, result, tt.expected)
			}
		})
	}
}

func TestCalculateWidth(t *testing.T) {
	tests := []struct {
		name               string
		props              property
		eastAsianWidth     bool
		strictEmojiNeutral bool
		expected           int
	}{
		// Control characters
		{"control char", IsControlChar, false, false, 0},

		// Combining marks
		{"combining mark", IsCombiningMark, false, false, 0},

		// Zero width
		{"zero width", IsZeroWidth, false, false, 0},

		// East Asian Wide
		{"EAW fullwidth", EAW_Fullwidth, false, false, 2},
		{"EAW wide", EAW_Wide, false, false, 2},

		// East Asian Ambiguous
		{"EAW ambiguous default", EAW_Ambiguous, false, false, 1},
		{"EAW ambiguous EAW", EAW_Ambiguous, true, false, 2},

		// Emoji
		{"emoji default", IsEmoji, false, false, 2},
		{"emoji strict", IsEmoji, false, true, 1},

		// Default (no properties set)
		{"default", 0, false, false, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateWidth(tt.props, tt.eastAsianWidth, tt.strictEmojiNeutral)
			if result != tt.expected {
				t.Errorf("calculateWidth(%d, %v, %v) = %d, want %d",
					tt.props, tt.eastAsianWidth, tt.strictEmojiNeutral, result, tt.expected)
			}
		})
	}
}
