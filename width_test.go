package displaywidth

import (
	"testing"
)

func TestStringWidth(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		options  Options
		expected int
	}{
		// Basic ASCII characters
		{"empty string", "", Options{}, 0},
		{"single ASCII", "a", Options{}, 1},
		{"multiple ASCII", "hello", Options{}, 5},
		{"ASCII with spaces", "hello world", Options{}, 11},

		// Control characters (width 0)
		{"newline", "\n", Options{}, 0},
		{"tab", "\t", Options{}, 0},
		{"carriage return", "\r", Options{}, 0},
		{"backspace", "\b", Options{}, 0},

		// Mixed content
		{"ASCII with newline", "hello\nworld", Options{}, 10},
		{"ASCII with tab", "hello\tworld", Options{}, 10},

		// East Asian characters (should be in trie)
		{"CJK ideograph", "中", Options{}, 2},
		{"CJK with ASCII", "hello中", Options{}, 7},

		// Ambiguous characters
		{"ambiguous character", "★", Options{}, 1},                         // Default narrow
		{"ambiguous character EAW", "★", Options{EastAsianWidth: true}, 2}, // East Asian wide

		// Emoji
		{"emoji", "😀", Options{}, 2},                                  // Default emoji width
		{"emoji strict", "😀", Options{StrictEmojiNeutral: true}, 2},   // Strict emoji neutral - only ambiguous emoji get width 1
		{"checkered flag", "🏁", Options{StrictEmojiNeutral: true}, 2}, // U+1F3C1 chequered flag is an emoji, width 2

		// Invalid UTF-8 - the trie treats \xff as a valid character with default properties
		{"invalid UTF-8", "\xff", Options{}, 1},
		{"partial UTF-8", "\xc2", Options{}, 1},

		// Variation selectors - VS16 (U+FE0F) requests emoji, VS15 (U+FE0E) requests text
		{"☺ text default", "☺", Options{}, 1},      // U+263A has text presentation by default
		{"☺️ emoji with VS16", "☺️", Options{}, 2}, // VS16 forces emoji presentation (width 2)
		{"⌛ emoji default", "⌛", Options{}, 2},     // U+231B has emoji presentation by default
		{"⌛︎ text with VS15", "⌛︎", Options{}, 1},  // VS15 forces text presentation (width 1)
		{"❤ text default", "❤", Options{}, 1},      // U+2764 has text presentation by default
		{"❤️ emoji with VS16", "❤️", Options{}, 2}, // VS16 forces emoji presentation (width 2)
		{"✂ text default", "✂", Options{}, 1},      // U+2702 has text presentation by default
		{"✂️ emoji with VS16", "✂️", Options{}, 2}, // VS16 forces emoji presentation (width 2)
		{"keycap 1️⃣", "1️⃣", Options{}, 2},        // Keycap sequence: 1 + VS16 + U+20E3 (always width 2)
		{"keycap #️⃣", "#️⃣", Options{}, 2},        // Keycap sequence: # + VS16 + U+20E3 (always width 2)

		// Flags (regional indicator pairs form a single grapheme, width 1 by default, width 2 with StrictEmojiNeutral=true)
		{"flag US", "🇺🇸", Options{}, 1},
		{"flag JP", "🇯🇵", Options{}, 1},
		{"text with flags", "Go 🇺🇸🚀", Options{}, 3 + 1 + 2},
		{"flag US strict", "🇺🇸", Options{StrictEmojiNeutral: true}, 2},
		{"flag JP strict", "🇯🇵", Options{StrictEmojiNeutral: true}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.options.String(tt.input)
			if result != tt.expected {
				t.Errorf("StringWidth(%q, %v) = %d, want %d",
					tt.input, tt.options, result, tt.expected)
			}

			b := []byte(tt.input)
			result = tt.options.Bytes(b)
			if result != tt.expected {
				t.Errorf("BytesWidth(%q, %v) = %d, want %d",
					b, tt.options, result, tt.expected)
			}
		})
	}
}

func TestRuneWidth(t *testing.T) {
	tests := []struct {
		name     string
		input    rune
		options  Options
		expected int
	}{
		// Control characters (width 0)
		{"null char", '\x00', Options{}, 0},
		{"bell", '\x07', Options{}, 0},
		{"backspace", '\x08', Options{}, 0},
		{"tab", '\t', Options{}, 0},
		{"newline", '\n', Options{}, 0},
		{"carriage return", '\r', Options{}, 0},
		{"escape", '\x1B', Options{}, 0},
		{"delete", '\x7F', Options{}, 0},

		// Combining marks - when tested standalone as runes, they have width 0
		// (In actual strings with grapheme clusters, they combine and have width 0)
		{"combining grave accent", '\u0300', Options{}, 0},
		{"combining acute accent", '\u0301', Options{}, 0},
		{"combining diaeresis", '\u0308', Options{}, 0},
		{"combining tilde", '\u0303', Options{}, 0},

		// Zero width characters
		{"zero width space", '\u200B', Options{}, 0},
		{"zero width non-joiner", '\u200C', Options{}, 0},
		{"zero width joiner", '\u200D', Options{}, 0},

		// ASCII printable (width 1)
		{"space", ' ', Options{}, 1},
		{"letter a", 'a', Options{}, 1},
		{"letter Z", 'Z', Options{}, 1},
		{"digit 0", '0', Options{}, 1},
		{"digit 9", '9', Options{}, 1},
		{"exclamation", '!', Options{}, 1},
		{"at sign", '@', Options{}, 1},
		{"tilde", '~', Options{}, 1},

		// Latin extended (width 1)
		{"latin e with acute", 'é', Options{}, 1},
		{"latin n with tilde", 'ñ', Options{}, 1},
		{"latin o with diaeresis", 'ö', Options{}, 1},

		// East Asian Wide characters
		{"CJK ideograph", '中', Options{}, 2},
		{"CJK ideograph", '文', Options{}, 2},
		{"hiragana a", 'あ', Options{}, 2},
		{"katakana a", 'ア', Options{}, 2},
		{"hangul syllable", '가', Options{}, 2},
		{"hangul syllable", '한', Options{}, 2},

		// Fullwidth characters
		{"fullwidth A", 'Ａ', Options{}, 2},
		{"fullwidth Z", 'Ｚ', Options{}, 2},
		{"fullwidth 0", '０', Options{}, 2},
		{"fullwidth 9", '９', Options{}, 2},
		{"fullwidth exclamation", '！', Options{}, 2},
		{"fullwidth space", '　', Options{}, 2},

		// Ambiguous characters - default narrow
		{"black star default", '★', Options{}, 1},
		{"degree sign default", '°', Options{}, 1},
		{"plus-minus default", '±', Options{}, 1},
		{"section sign default", '§', Options{}, 1},
		{"copyright sign default", '©', Options{}, 1},
		{"registered sign default", '®', Options{}, 1},

		// Ambiguous characters - EastAsianWidth wide
		{"black star EAW", '★', Options{EastAsianWidth: true}, 2},
		{"degree sign EAW", '°', Options{EastAsianWidth: true}, 2},
		{"plus-minus EAW", '±', Options{EastAsianWidth: true}, 2},
		{"section sign EAW", '§', Options{EastAsianWidth: true}, 2},
		{"copyright sign EAW", '©', Options{EastAsianWidth: true}, 1}, // Not in ambiguous category
		{"registered sign EAW", '®', Options{EastAsianWidth: true}, 2},

		// Emoji (width 2)
		{"grinning face", '😀', Options{}, 2},
		{"grinning face with smiling eyes", '😁', Options{}, 2},
		{"smiling face with heart-eyes", '😍', Options{}, 2},
		{"thinking face", '🤔', Options{}, 2},
		{"rocket", '🚀', Options{}, 2},
		{"party popper", '🎉', Options{}, 2},
		{"fire", '🔥', Options{}, 2},
		{"thumbs up", '👍', Options{}, 2},
		{"red heart", '❤', Options{}, 1},      // Text presentation by default
		{"checkered flag", '🏁', Options{}, 2}, // U+1F3C1 chequered flag

		// Emoji with StrictEmojiNeutral
		{"grinning face strict", '😀', Options{StrictEmojiNeutral: true}, 2},
		{"rocket strict", '🚀', Options{StrictEmojiNeutral: true}, 2},
		{"party popper strict", '🎉', Options{StrictEmojiNeutral: true}, 2},

		// Emoji with both options
		{"grinning face both", '😀', Options{EastAsianWidth: true, StrictEmojiNeutral: true}, 2},
		{"rocket both", '🚀', Options{EastAsianWidth: true, StrictEmojiNeutral: true}, 2},

		// Mathematical symbols
		{"infinity", '∞', Options{}, 1},
		{"summation", '∑', Options{}, 1},
		{"integral", '∫', Options{}, 1},
		{"square root", '√', Options{}, 1},

		// Currency symbols
		{"dollar", '$', Options{}, 1},
		{"euro", '€', Options{}, 1},
		{"pound", '£', Options{}, 1},
		{"yen", '¥', Options{}, 1},

		// Box drawing characters
		{"box light horizontal", '─', Options{}, 1},
		{"box light vertical", '│', Options{}, 1},
		{"box light down and right", '┌', Options{}, 1},

		// Special Unicode characters
		{"bullet", '•', Options{}, 1},
		{"ellipsis", '…', Options{}, 1},
		{"em dash", '—', Options{}, 1},
		{"left single quote", '\u2018', Options{}, 1},
		{"right single quote", '\u2019', Options{}, 1},

		// Test edge cases with both options disabled
		{"ambiguous both disabled", '★', Options{EastAsianWidth: false, StrictEmojiNeutral: false}, 1},
		{"ambiguous strict only", '★', Options{EastAsianWidth: false, StrictEmojiNeutral: true}, 1},

		// Variation selectors (note: Rune() tests single runes without VS, not sequences)
		{"☺ U+263A text default", '☺', Options{}, 1},    // Has text presentation by default
		{"⌛ U+231B emoji default", '⌛', Options{}, 2},   // Has emoji presentation by default
		{"❤ U+2764 text default", '❤', Options{}, 1},    // Has text presentation by default
		{"✂ U+2702 text default", '✂', Options{}, 1},    // Has text presentation by default
		{"VS16 U+FE0F alone", '\ufe0f', Options{}, 0},   // Variation selectors are zero-width by themselves
		{"VS15 U+FE0E alone", '\ufe0e', Options{}, 0},   // Variation selectors are zero-width by themselves
		{"keycap U+20E3 alone", '\u20e3', Options{}, 0}, // Combining enclosing keycap is zero-width alone
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.options.Rune(tt.input)
			if result != tt.expected {
				t.Errorf("RuneWidth(%q, %v) = %d, want %d",
					tt.input, tt.options, result, tt.expected)
			}
		})
	}
}

func TestCalculateWidth(t *testing.T) {
	tests := []struct {
		name     string
		props    property
		options  Options
		expected int
	}{ // Zero width
		{"zero width", _ZeroWidth, Options{}, 0},

		// East Asian Wide
		{"EAW fullwidth", _East_Asian_Full_Wide, Options{}, 2},
		{"EAW wide", _East_Asian_Full_Wide, Options{}, 2},

		// East Asian Ambiguous
		{"EAW ambiguous default", _East_Asian_Ambiguous, Options{}, 1},
		{"EAW ambiguous EAW", _East_Asian_Ambiguous, Options{EastAsianWidth: true}, 2},

		// Emoji
		// {"emoji default", _Emoji, false, false, 2},
		// {"emoji strict", _Emoji, false, true, 2}, // Only ambiguous emoji get width 1 in strict mode

		// Default (no properties set)
		{"default", 0, Options{}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.props.width(tt.options)
			if result != tt.expected {
				t.Errorf("calculateWidth(%d, %v) = %d, want %d",
					tt.props, tt.options, result, tt.expected)
			}
		})
	}
}
