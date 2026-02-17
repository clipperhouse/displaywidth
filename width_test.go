package displaywidth

import (
	"bytes"
	"testing"
)

var defaultOptions = Options{}

var eawOptions = Options{EastAsianWidth: true}

func TestStringWidth(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		options  Options
		expected int
	}{
		// Basic ASCII characters
		{"empty string", "", defaultOptions, 0},
		{"single ASCII", "a", defaultOptions, 1},
		{"multiple ASCII", "hello", defaultOptions, 5},
		{"ASCII with spaces", "hello world", defaultOptions, 11},

		// Control characters (width 0)
		{"newline", "\n", defaultOptions, 0},
		{"tab", "\t", defaultOptions, 0},
		{"carriage return", "\r", defaultOptions, 0},
		{"backspace", "\b", defaultOptions, 0},

		// Mixed content
		{"ASCII with newline", "hello\nworld", defaultOptions, 10},
		{"ASCII with tab", "hello\tworld", defaultOptions, 10},

		// East Asian characters (should be in trie)
		{"CJK ideograph", "中", defaultOptions, 2},
		{"CJK with ASCII", "hello中", defaultOptions, 7},

		// Ambiguous characters
		{"ambiguous character", "★", defaultOptions, 1}, // Default narrow
		{"ambiguous character EAW", "★", eawOptions, 2}, // East Asian wide

		// Emoji
		{"emoji", "😀", defaultOptions, 2},          // Default emoji width
		{"checkered flag", "🏁", defaultOptions, 2}, // U+1F3C1 chequered flag is an emoji, width 2

		// Invalid UTF-8 - the trie treats \xff as a valid character with default properties
		{"invalid UTF-8", "\xff", defaultOptions, 1},
		{"partial UTF-8", "\xc2", defaultOptions, 1},

		// Variation selectors - VS16 (U+FE0F) requests emoji, VS15 (U+FE0E) is a no-op per Unicode TR51
		{"☺ text default", "☺", defaultOptions, 1},      // U+263A has text presentation by default
		{"☺️ emoji with VS16", "☺️", defaultOptions, 2}, // VS16 forces emoji presentation (width 2)
		{"⌛ emoji default", "⌛", defaultOptions, 2},     // U+231B has emoji presentation by default
		{"⌛︎ with VS15", "⌛︎", defaultOptions, 2},       // VS15 is a no-op, width remains 2
		{"❤ text default", "❤", defaultOptions, 1},      // U+2764 has text presentation by default
		{"❤️ emoji with VS16", "❤️", defaultOptions, 2}, // VS16 forces emoji presentation (width 2)
		{"✂ text default", "✂", defaultOptions, 1},      // U+2702 has text presentation by default
		{"✂️ emoji with VS16", "✂️", defaultOptions, 2}, // VS16 forces emoji presentation (width 2)
		{"keycap 1️⃣", "1️⃣", defaultOptions, 2},        // Keycap sequence: 1 + VS16 + U+20E3 (always width 2)
		{"keycap #️⃣", "#️⃣", defaultOptions, 2},        // Keycap sequence: # + VS16 + U+20E3 (always width 2)

		// Flags (regional indicator pairs form a single grapheme, always width 2 per TR51)
		{"flag US", "🇺🇸", defaultOptions, 2},
		{"flag JP", "🇯🇵", defaultOptions, 2},
		{"text with flags", "Go 🇺🇸🚀", defaultOptions, 3 + 2 + 2},

		// Partial ASCII optimization tests (8+ byte ASCII runs)
		{"ASCII 8 bytes then emoji", "12345678😀", defaultOptions, 8 + 2},
		{"ASCII 16 bytes then CJK", "1234567890abcdef中", defaultOptions, 16 + 2},
		{"emoji then ASCII 8 bytes", "😀12345678", defaultOptions, 2 + 8},
		{"CJK then ASCII 16 bytes", "中1234567890abcdef", defaultOptions, 2 + 16},
		{"ASCII-emoji-ASCII sandwich", "12345678😀abcdefgh", defaultOptions, 8 + 2 + 8},
		{"short ASCII then emoji", "hello😀", defaultOptions, 5 + 2},
		{"emoji-short ASCII-emoji", "😀abc😀", defaultOptions, 2 + 3 + 2},
		{"long mixed", "Hello World! 你好世界 12345678 emoji: 🎉🎊", defaultOptions, 42}, // 13 + 9 + 9 + 7 + 4

		// ASCII with embedded control characters
		{"ASCII with null in middle", "hello\x00world", defaultOptions, 10}, // 5 + 0 + 5
		{"ASCII with DEL in middle", "hello\x7Fworld", defaultOptions, 10},  // 5 + 0 + 5
		{"ASCII with multiple controls", "a\x00b\tc\nd", defaultOptions, 4}, // 1 + 0 + 1 + 0 + 1 + 0 + 1

		// Alternating short ASCII/non-ASCII sequences
		{"alternating ASCII-CJK", "a中b文c", defaultOptions, 7},   // 1 + 2 + 1 + 2 + 1
		{"alternating CJK-ASCII", "中a文b字c", defaultOptions, 9},  // 2 + 1 + 2 + 1 + 2 + 1
		{"single char alternating", "a😀b🎉c", defaultOptions, 7}, // 1 + 2 + 1 + 2 + 1
		{"rapid alternation", "aあbいcうd", defaultOptions, 10},    // 1 + 2 + 1 + 2 + 1 + 2 + 1
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

var controlSequences = Options{ControlSequences: true}
var controlSequences8Bit = Options{ControlSequences8Bit: true}
var controlSequencesBoth = Options{ControlSequences: true, ControlSequences8Bit: true}

func TestAnsiEscapeSequences(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		options  Options
		expected int
	}{
		// ANSI escape sequences (ECMA-48) should be zero width when parsed as single graphemes
		{"SGR red", "\x1b[31m", controlSequences, 0},
		{"SGR reset", "\x1b[0m", controlSequences, 0},
		{"SGR bold", "\x1b[1m", controlSequences, 0},
		{"SGR 256-color", "\x1b[38;5;196m", controlSequences, 0},
		{"SGR truecolor", "\x1b[38;2;255;0;0m", controlSequences, 0},
		{"cursor up", "\x1b[A", controlSequences, 0},
		{"cursor position", "\x1b[10;20H", controlSequences, 0},
		{"erase in display", "\x1b[2J", controlSequences, 0},

		// ANSI escape sequences mixed with visible text
		{"red hello", "\x1b[31mhello\x1b[0m", controlSequences, 5},
		{"bold world", "\x1b[1mworld\x1b[0m", controlSequences, 5},
		{"colored CJK", "\x1b[31m中文\x1b[0m", controlSequences, 4},
		{"colored emoji", "\x1b[31m😀\x1b[0m", controlSequences, 2},
		{"nested SGR", "\x1b[1m\x1b[31mhi\x1b[0m", controlSequences, 2},

		// CR+LF as a multi-byte C0-led grapheme (zero width)
		{"CRLF", "\r\n", controlSequences, 0},
		{"text with CRLF", "hello\r\nworld", controlSequences, 10},

		// Without ControlSequences, ESC is zero width but the rest of the sequence is visible
		{"bare ESC default options", "\x1b", defaultOptions, 0},
		{"SGR red default options", "\x1b[31m", defaultOptions, 4},
		{"red hello default options", "\x1b[31mhello\x1b[0m", defaultOptions, 12},

		// ControlSequences should not regress width for strings with no escape sequences
		{"plain ASCII with option", "hello", controlSequences, 5},
		{"plain ASCII spaces with option", "hello world", controlSequences, 11},
		{"CJK with option", "中文", controlSequences, 4},
		{"emoji with option", "😀", controlSequences, 2},
		{"flag with option", "🇺🇸", controlSequences, 2},
		{"mixed with option", "hello中文😀", controlSequences, 5 + 4 + 2},
		{"ambiguous with option", "★", controlSequences, 1},
		{"combining mark with option", "é", controlSequences, 1},
		{"control chars with option", "\t\n", controlSequences, 0},
		{"empty with option", "", controlSequences, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.options.String(tt.input)
			if result != tt.expected {
				t.Errorf("String(%q) = %d, want %d", tt.input, result, tt.expected)
			}

			result = tt.options.Bytes([]byte(tt.input))
			if result != tt.expected {
				t.Errorf("Bytes(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

func TestAnsiEscapeSequences8Bit(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		options  Options
		expected int
	}{
		// 8-bit C1 CSI sequences should be zero width
		{"C1 CSI red", "\x9B31m", controlSequences8Bit, 0},
		{"C1 CSI reset", "\x9B0m", controlSequences8Bit, 0},
		{"C1 CSI bold", "\x9B1m", controlSequences8Bit, 0},
		{"C1 CSI multi-param", "\x9B1;2;3m", controlSequences8Bit, 0},
		{"C1 CSI cursor up", "\x9BA", controlSequences8Bit, 0},

		// 8-bit C1 OSC/DCS/SOS/APC with C1 ST terminator
		{"C1 OSC with ST", "\x9D0;Title\x9C", controlSequences8Bit, 0},
		{"C1 OSC with BEL", "\x9D0;Title\x07", controlSequences8Bit, 0},
		{"C1 DCS with ST", "\x90qpayload\x9C", controlSequences8Bit, 0},
		{"C1 SOS with ST", "\x98hello\x9C", controlSequences8Bit, 0},
		{"C1 APC with ST", "\x9Fdata\x9C", controlSequences8Bit, 0},

		// Standalone C1 controls (single byte, no body)
		{"C1 IND", "\x84", controlSequences8Bit, 0},
		{"C1 NEL", "\x85", controlSequences8Bit, 0},

		// 8-bit sequences mixed with visible text
		{"C1 CSI red hello", "\x9B31mhello\x9B0m", controlSequences8Bit, 5},
		{"C1 CSI colored CJK", "\x9B31m中文\x9B0m", controlSequences8Bit, 4},
		{"C1 CSI colored emoji", "\x9B31m😀\x9B0m", controlSequences8Bit, 2},
		{"C1 CSI nested", "\x9B1m\x9B31mhi\x9B0m", controlSequences8Bit, 2},

		// Without ControlSequences8Bit, C1 bytes have width per asciiWidth (1 for >= 0x80)
		{"C1 CSI default options", "\x9B31m", defaultOptions, 4},

		// 8-bit option should not regress plain text
		{"plain ASCII with 8-bit option", "hello", controlSequences8Bit, 5},
		{"CJK with 8-bit option", "中文", controlSequences8Bit, 4},
		{"emoji with 8-bit option", "😀", controlSequences8Bit, 2},
		{"empty with 8-bit option", "", controlSequences8Bit, 0},

		// Both options enabled
		{"both: 7-bit SGR", "\x1b[31mhello\x1b[0m", controlSequencesBoth, 5},
		{"both: 8-bit CSI", "\x9B31mhello\x9B0m", controlSequencesBoth, 5},
		{"both: mixed 7 and 8-bit", "\x1b[31mhello\x9B0m", controlSequencesBoth, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.options.String(tt.input)
			if result != tt.expected {
				t.Errorf("String(%q) = %d, want %d", tt.input, result, tt.expected)
			}

			result = tt.options.Bytes([]byte(tt.input))
			if result != tt.expected {
				t.Errorf("Bytes(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

// TestAnsiEscapeSequencesIndependence verifies that the 7-bit and 8-bit options
// are strictly independent: enabling one must NOT cause the other's sequences
// to be treated as escape sequences.
func TestAnsiEscapeSequencesIndependence(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		options  Options
		expected int
		desc     string
	}{
		// 7-bit only: C1 bytes must NOT be treated as escape sequences.
		// \x9B31m is 4 visible chars (0x9B has width 1, '3' '1' 'm' each width 1)
		{
			name:     "7-bit on, 8-bit input C1 CSI",
			input:    "\x9B31m",
			options:  controlSequences,
			expected: 4,
			desc:     "C1 CSI should not be recognized when only 7-bit is enabled",
		},
		{
			name:     "7-bit on, 8-bit input standalone C1",
			input:    "\x84",
			options:  controlSequences,
			expected: 1,
			desc:     "Standalone C1 byte should have width 1 when only 7-bit is enabled",
		},
		{
			name:     "7-bit on, 8-bit input C1 with text",
			input:    "\x9B31mhello\x9B0m",
			options:  controlSequences,
			expected: 4 + 5 + 3,
			desc:     "C1 CSI sequences should contribute visible width when only 7-bit is enabled",
		},

		// 8-bit only: 7-bit ESC sequences must NOT be treated as escape sequences.
		// \x1b[31m is: ESC (width 0) + '[' (1) + '3' (1) + '1' (1) + 'm' (1) = 4
		{
			name:     "8-bit on, 7-bit input SGR",
			input:    "\x1b[31m",
			options:  controlSequences8Bit,
			expected: 4,
			desc:     "7-bit SGR should not be recognized when only 8-bit is enabled",
		},
		{
			name:     "8-bit on, 7-bit input SGR with text",
			input:    "\x1b[31mhello\x1b[0m",
			options:  controlSequences8Bit,
			expected: 4 + 5 + 3,
			desc:     "7-bit SGR should contribute visible width when only 8-bit is enabled",
		},

		// Both enabled: both kinds should be zero-width
		{
			name:     "both on, 7-bit SGR",
			input:    "\x1b[31m",
			options:  controlSequencesBoth,
			expected: 0,
			desc:     "7-bit SGR should be zero-width when both are enabled",
		},
		{
			name:     "both on, 8-bit CSI",
			input:    "\x9B31m",
			options:  controlSequencesBoth,
			expected: 0,
			desc:     "C1 CSI should be zero-width when both are enabled",
		},
		{
			name:     "both on, mixed sequences with text",
			input:    "\x1b[31mhello\x9B0m",
			options:  controlSequencesBoth,
			expected: 5,
			desc:     "Mixed 7-bit and 8-bit sequences should both be zero-width",
		},

		// Neither enabled: both kinds contribute visible width
		{
			name:     "neither, 7-bit SGR",
			input:    "\x1b[31m",
			options:  defaultOptions,
			expected: 4,
			desc:     "7-bit SGR should contribute visible width when neither is enabled",
		},
		{
			name:     "neither, 8-bit CSI",
			input:    "\x9B31m",
			options:  defaultOptions,
			expected: 4,
			desc:     "C1 CSI should contribute visible width when neither is enabled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.options.String(tt.input)
			if result != tt.expected {
				t.Errorf("String(%q) = %d, want %d (%s)", tt.input, result, tt.expected, tt.desc)
			}

			result = tt.options.Bytes([]byte(tt.input))
			if result != tt.expected {
				t.Errorf("Bytes(%q) = %d, want %d (%s)", tt.input, result, tt.expected, tt.desc)
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
		{"null char", '\x00', defaultOptions, 0},
		{"bell", '\x07', defaultOptions, 0},
		{"backspace", '\x08', defaultOptions, 0},
		{"tab", '\t', defaultOptions, 0},
		{"newline", '\n', defaultOptions, 0},
		{"carriage return", '\r', defaultOptions, 0},
		{"escape", '\x1B', defaultOptions, 0},
		{"delete", '\x7F', defaultOptions, 0},

		// Combining marks - when tested standalone as runes, they have width 0
		// (In actual strings with grapheme clusters, they combine and have width 0)
		{"combining grave accent", '\u0300', defaultOptions, 0},
		{"combining acute accent", '\u0301', defaultOptions, 0},
		{"combining diaeresis", '\u0308', defaultOptions, 0},
		{"combining tilde", '\u0303', defaultOptions, 0},

		// Zero width characters
		{"zero width space", '\u200B', defaultOptions, 0},
		{"zero width non-joiner", '\u200C', defaultOptions, 0},
		{"zero width joiner", '\u200D', defaultOptions, 0},

		// ASCII printable (width 1)
		{"space", ' ', defaultOptions, 1},
		{"letter a", 'a', defaultOptions, 1},
		{"letter Z", 'Z', defaultOptions, 1},
		{"digit 0", '0', defaultOptions, 1},
		{"digit 9", '9', defaultOptions, 1},
		{"exclamation", '!', defaultOptions, 1},
		{"at sign", '@', defaultOptions, 1},
		{"tilde", '~', defaultOptions, 1},

		// Latin extended (width 1)
		{"latin e with acute", 'é', defaultOptions, 1},
		{"latin n with tilde", 'ñ', defaultOptions, 1},
		{"latin o with diaeresis", 'ö', defaultOptions, 1},

		// East Asian Wide characters
		{"CJK ideograph", '中', defaultOptions, 2},
		{"CJK ideograph", '文', defaultOptions, 2},
		{"hiragana a", 'あ', defaultOptions, 2},
		{"katakana a", 'ア', defaultOptions, 2},
		{"hangul syllable", '가', defaultOptions, 2},
		{"hangul syllable", '한', defaultOptions, 2},

		// Fullwidth characters
		{"fullwidth A", 'Ａ', defaultOptions, 2},
		{"fullwidth Z", 'Ｚ', defaultOptions, 2},
		{"fullwidth 0", '０', defaultOptions, 2},
		{"fullwidth 9", '９', defaultOptions, 2},
		{"fullwidth exclamation", '！', defaultOptions, 2},
		{"fullwidth space", '　', defaultOptions, 2},

		// Ambiguous characters - default narrow
		{"black star default", '★', defaultOptions, 1},
		{"degree sign default", '°', defaultOptions, 1},
		{"plus-minus default", '±', defaultOptions, 1},
		{"section sign default", '§', defaultOptions, 1},
		{"copyright sign default", '©', defaultOptions, 1},
		{"registered sign default", '®', defaultOptions, 1},

		// Ambiguous characters - EastAsianWidth wide
		{"black star EAW", '★', eawOptions, 2},
		{"degree sign EAW", '°', eawOptions, 2},
		{"plus-minus EAW", '±', eawOptions, 2},
		{"section sign EAW", '§', eawOptions, 2},
		{"copyright sign EAW", '©', eawOptions, 1}, // Not in ambiguous category
		{"registered sign EAW", '®', eawOptions, 2},

		// Emoji (width 2)
		{"grinning face", '😀', defaultOptions, 2},
		{"grinning face with smiling eyes", '😁', defaultOptions, 2},
		{"smiling face with heart-eyes", '😍', defaultOptions, 2},
		{"thinking face", '🤔', defaultOptions, 2},
		{"rocket", '🚀', defaultOptions, 2},
		{"party popper", '🎉', defaultOptions, 2},
		{"fire", '🔥', defaultOptions, 2},
		{"thumbs up", '👍', defaultOptions, 2},
		{"red heart", '❤', defaultOptions, 1},      // Text presentation by default
		{"checkered flag", '🏁', defaultOptions, 2}, // U+1F3C1 chequered flag

		// Mathematical symbols
		{"infinity", '∞', defaultOptions, 1},
		{"summation", '∑', defaultOptions, 1},
		{"integral", '∫', defaultOptions, 1},
		{"square root", '√', defaultOptions, 1},

		// Currency symbols
		{"dollar", '$', defaultOptions, 1},
		{"euro", '€', defaultOptions, 1},
		{"pound", '£', defaultOptions, 1},
		{"yen", '¥', defaultOptions, 1},

		// Box drawing characters
		{"box light horizontal", '─', defaultOptions, 1},
		{"box light vertical", '│', defaultOptions, 1},
		{"box light down and right", '┌', defaultOptions, 1},

		// Special Unicode characters
		{"bullet", '•', defaultOptions, 1},
		{"ellipsis", '…', defaultOptions, 1},
		{"em dash", '—', defaultOptions, 1},
		{"left single quote", '\u2018', defaultOptions, 1},
		{"right single quote", '\u2019', defaultOptions, 1},

		// Test edge cases with options disabled
		{"ambiguous EAW disabled", '★', defaultOptions, 1},

		// Variation selectors (note: Rune() tests single runes without VS, not sequences)
		{"☺ U+263A text default", '☺', defaultOptions, 1},    // Has text presentation by default
		{"⌛ U+231B emoji default", '⌛', defaultOptions, 2},   // Has emoji presentation by default
		{"❤ U+2764 text default", '❤', defaultOptions, 1},    // Has text presentation by default
		{"✂ U+2702 text default", '✂', defaultOptions, 1},    // Has text presentation by default
		{"VS16 U+FE0F alone", '\ufe0f', defaultOptions, 0},   // Variation selectors are zero-width by themselves
		{"VS15 U+FE0E alone", '\ufe0e', defaultOptions, 0},   // Variation selectors are zero-width by themselves
		{"keycap U+20E3 alone", '\u20e3', defaultOptions, 0}, // Combining enclosing keycap is zero-width alone
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

// TestEmojiPresentation verifies correct width behavior for characters with different
// Emoji_Presentation property values according to TR51 conformance
func TestEmojiPresentation(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantDefault  int
		wantWithVS16 int
		wantWithVS15 int
		description  string
	}{
		// Characters with Extended_Pictographic=Yes AND Emoji_Presentation=Yes
		// Should have width 2 by default (emoji presentation)
		// VS15 is a no-op per Unicode TR51 - it requests text presentation but doesn't change width
		{
			name:         "Watch (EP=Yes, EmojiPres=Yes)",
			input:        "\u231A",
			wantDefault:  2,
			wantWithVS16: 2,
			wantWithVS15: 2, // VS15 is a no-op, width remains 2
			description:  "⌚ U+231A has default emoji presentation",
		},
		{
			name:         "Hourglass (EP=Yes, EmojiPres=Yes)",
			input:        "\u231B",
			wantDefault:  2,
			wantWithVS16: 2,
			wantWithVS15: 2, // VS15 is a no-op, width remains 2
			description:  "⌛ U+231B has default emoji presentation",
		},
		{
			name:         "Fast-forward (EP=Yes, EmojiPres=Yes)",
			input:        "\u23E9",
			wantDefault:  2,
			wantWithVS16: 2,
			wantWithVS15: 2, // VS15 is a no-op, width remains 2
			description:  "⏩ U+23E9 has default emoji presentation",
		},
		{
			name:         "Alarm Clock (EP=Yes, EmojiPres=Yes)",
			input:        "\u23F0",
			wantDefault:  2,
			wantWithVS16: 2,
			wantWithVS15: 2, // VS15 is a no-op, width remains 2
			description:  "⏰ U+23F0 has default emoji presentation",
		},
		{
			name:         "Soccer Ball (EP=Yes, EmojiPres=Yes)",
			input:        "\u26BD",
			wantDefault:  2,
			wantWithVS16: 2,
			wantWithVS15: 2, // VS15 is a no-op, width remains 2
			description:  "⚽ U+26BD has default emoji presentation",
		},
		{
			name:         "Anchor (EP=Yes, EmojiPres=Yes)",
			input:        "\u2693",
			wantDefault:  2,
			wantWithVS16: 2,
			wantWithVS15: 2, // VS15 is a no-op, width remains 2
			description:  "⚓ U+2693 has default emoji presentation",
		},

		// Characters with Extended_Pictographic=Yes BUT Emoji_Presentation=No
		// Should have width 1 by default (text presentation)
		{
			name:         "Star of David (EP=Yes, EmojiPres=No)",
			input:        "\u2721",
			wantDefault:  1,
			wantWithVS16: 2,
			wantWithVS15: 1,
			description:  "✡ U+2721 has default text presentation",
		},
		{
			name:         "Hammer and Pick (EP=Yes, EmojiPres=No)",
			input:        "\u2692",
			wantDefault:  1,
			wantWithVS16: 2,
			wantWithVS15: 1,
			description:  "⚒ U+2692 has default text presentation",
		},
		{
			name:         "Gear (EP=Yes, EmojiPres=No)",
			input:        "\u2699",
			wantDefault:  1,
			wantWithVS16: 2,
			wantWithVS15: 1,
			description:  "⚙ U+2699 has default text presentation",
		},
		{
			name:         "Star and Crescent (EP=Yes, EmojiPres=No)",
			input:        "\u262A",
			wantDefault:  1,
			wantWithVS16: 2,
			wantWithVS15: 1,
			description:  "☪ U+262A has default text presentation",
		},
		{
			name:         "Infinity (EP=Yes, EmojiPres=No)",
			input:        "\u267E",
			wantDefault:  1,
			wantWithVS16: 2,
			wantWithVS15: 1,
			description:  "♾ U+267E has default text presentation",
		},
		{
			name:         "Recycling Symbol (EP=Yes, EmojiPres=No)",
			input:        "\u267B",
			wantDefault:  1,
			wantWithVS16: 2,
			wantWithVS15: 1,
			description:  "♻ U+267B has default text presentation",
		},

		// Characters with Emoji=Yes but NOT Extended_Pictographic
		// These are typically ASCII characters like # that can become emoji with VS16
		{
			name:         "Hash Sign (Emoji=Yes, EP=No)",
			input:        "\u0023",
			wantDefault:  1,
			wantWithVS16: 2,
			wantWithVS15: 1,
			description:  "# U+0023 has default text presentation",
		},
		{
			name:         "Asterisk (Emoji=Yes, EP=No)",
			input:        "\u002A",
			wantDefault:  1,
			wantWithVS16: 2,
			wantWithVS15: 1,
			description:  "* U+002A has default text presentation",
		},
		{
			name:         "Digit Zero (Emoji=Yes, EP=No)",
			input:        "\u0030",
			wantDefault:  1,
			wantWithVS16: 2,
			wantWithVS15: 1,
			description:  "0 U+0030 has default text presentation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test default width (no variation selector)
			gotDefault := String(tt.input)
			if gotDefault != tt.wantDefault {
				t.Errorf("String(%q) default = %d, want %d (%s)",
					tt.input, gotDefault, tt.wantDefault, tt.description)
			}

			// Test with VS16 (U+FE0F) for emoji presentation
			inputWithVS16 := tt.input + "\uFE0F"
			gotWithVS16 := String(inputWithVS16)
			if gotWithVS16 != tt.wantWithVS16 {
				t.Errorf("String(%q) with VS16 = %d, want %d (%s)",
					tt.input, gotWithVS16, tt.wantWithVS16, tt.description)
			}

			// Test with VS15 (U+FE0E) - VS15 is a no-op per Unicode TR51
			// It requests text presentation but does not affect width calculation
			inputWithVS15 := tt.input + "\uFE0E"
			gotWithVS15 := String(inputWithVS15)
			if gotWithVS15 != tt.wantWithVS15 {
				t.Errorf("String(%q) with VS15 = %d, want %d (%s)",
					tt.input, gotWithVS15, tt.wantWithVS15, tt.description)
			}
		})
	}
}

// TestEmojiPresentationRune tests the Rune() function specifically
func TestEmojiPresentationRune(t *testing.T) {
	tests := []struct {
		name string
		r    rune
		want int
		desc string
	}{
		// Emoji_Presentation=Yes
		{name: "Watch", r: '\u231A', want: 2, desc: "⌚ has default emoji presentation"},
		{name: "Alarm Clock", r: '\u23F0', want: 2, desc: "⏰ has default emoji presentation"},
		{name: "Soccer Ball", r: '\u26BD', want: 2, desc: "⚽ has default emoji presentation"},

		// Emoji_Presentation=No (but Extended_Pictographic=Yes)
		{name: "Star of David", r: '\u2721', want: 1, desc: "✡ has default text presentation"},
		{name: "Hammer and Pick", r: '\u2692', want: 1, desc: "⚒ has default text presentation"},
		{name: "Gear", r: '\u2699', want: 1, desc: "⚙ has default text presentation"},
		{name: "Infinity", r: '\u267E', want: 1, desc: "♾ has default text presentation"},

		// Not Extended_Pictographic
		{name: "Hash Sign", r: '#', want: 1, desc: "# is regular ASCII"},
		{name: "Asterisk", r: '*', want: 1, desc: "* is regular ASCII"},
		{name: "Digit Zero", r: '0', want: 1, desc: "0 is regular ASCII"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Rune(tt.r)
			if got != tt.want {
				t.Errorf("Rune(%U) = %d, want %d (%s)", tt.r, got, tt.want, tt.desc)
			}
		})
	}
}

// TestComplexEmojiSequences tests width of complex emoji sequences
func TestComplexEmojiSequences(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
		desc  string
	}{
		{
			name:  "Keycap sequence #️⃣",
			input: "#\uFE0F\u20E3",
			want:  2,
			desc:  "# + VS16 + combining enclosing keycap",
		},
		{
			name:  "Keycap sequence 0️⃣",
			input: "0\uFE0F\u20E3",
			want:  2,
			desc:  "0 + VS16 + combining enclosing keycap",
		},
		{
			name:  "Flag sequence 🇺🇸 (Regional Indicator pair)",
			input: "\U0001F1FA\U0001F1F8",
			want:  2,
			desc:  "US flag (RI pair)",
		},
		{
			name:  "Single Regional Indicator 🇺",
			input: "\U0001F1FA",
			want:  2,
			desc:  "U (RI)",
		},
		{
			name:  "ZWJ sequence 👨‍👩‍👧",
			input: "\U0001F468\u200D\U0001F469\u200D\U0001F467",
			want:  2,
			desc:  "Family emoji (man + ZWJ + woman + ZWJ + girl)",
		},
		{
			name:  "Skin tone modifier 👍🏽",
			input: "\U0001F44D\U0001F3FD",
			want:  2,
			desc:  "Thumbs up with medium skin tone",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := String(tt.input)
			if got != tt.want {
				t.Errorf("String(%q) = %d, want %d (%s)",
					tt.input, got, tt.want, tt.desc)
			}
		})
	}
}

// TestMixedContent tests width of strings with mixed emoji and text
func TestMixedContent(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
		desc  string
	}{
		{
			name:  "Text with emoji-presentation emoji",
			input: "Hi\u231AWorld",
			want:  9, // "Hi" (2) + ⌚ (2) + "World" (5)
			desc:  "Text with watch emoji (emoji presentation)",
		},
		{
			name:  "Text with text-presentation emoji",
			input: "Hi\u2721Go",
			want:  5, // "Hi" (2) + ✡ (1) + "Go" (2)
			desc:  "Text with star of David (text presentation)",
		},
		{
			name:  "Text with text-presentation emoji + VS16",
			input: "Hi\u2721\uFE0FGo",
			want:  6, // "Hi" (2) + ✡️ (2) + "Go" (2)
			desc:  "Text with star of David (forced emoji presentation)",
		},
		{
			name:  "Multiple emojis",
			input: "⌚⚽⚓",
			want:  6, // All three have Emoji_Presentation=Yes
			desc:  "Watch, soccer ball, anchor",
		},
		{
			name:  "Mixed presentation",
			input: "⌚⚙⚓",
			want:  5, // ⌚(2) + ⚙(1) + ⚓(2)
			desc:  "Watch (emoji), gear (text), anchor (emoji)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := String(tt.input)
			if got != tt.want {
				t.Errorf("String(%q) = %d, want %d (%s)",
					tt.input, got, tt.want, tt.desc)
			}
		})
	}
}

// TestTR51Conformance verifies key TR51 conformance requirements
func TestTR51Conformance(t *testing.T) {
	t.Run("C1: Default Emoji Presentation", func(t *testing.T) {
		// Characters with Emoji_Presentation=Yes should display as emoji by default (width 2)
		emojiPresentationChars := []rune{
			'\u231A', // ⌚ watch
			'\u231B', // ⌛ hourglass
			'\u23F0', // ⏰ alarm clock
			'\u26BD', // ⚽ soccer ball
			'\u2693', // ⚓ anchor
		}

		for _, r := range emojiPresentationChars {
			got := Rune(r)
			if got != 2 {
				t.Errorf("Rune(%U) = %d, want 2 (should have default emoji presentation)", r, got)
			}
		}
	})

	t.Run("C1: Default Text Presentation", func(t *testing.T) {
		// Characters with Emoji_Presentation=No should display as text by default (width 1)
		textPresentationChars := []rune{
			'\u2721', // ✡ star of David
			'\u2692', // ⚒ hammer and pick
			'\u2699', // ⚙ gear
			'\u267E', // ♾ infinity
			'\u267B', // ♻ recycling symbol
		}

		for _, r := range textPresentationChars {
			got := Rune(r)
			if got != 1 {
				t.Errorf("Rune(%U) = %d, want 1 (should have default text presentation)", r, got)
			}
		}
	})

	t.Run("C2: VS15 is a no-op for width calculation", func(t *testing.T) {
		// VS15 (U+FE0E) requests text presentation but does not affect width per Unicode TR51.
		// The width should be the same as the base character.
		emojiWithVS15 := []struct {
			char string
			base string
		}{
			{"\u231A\uFE0E", "\u231A"}, // ⌚︎ watch with VS15
			{"\u26BD\uFE0E", "\u26BD"}, // ⚽︎ soccer ball with VS15
			{"\u2693\uFE0E", "\u2693"}, // ⚓︎ anchor with VS15
		}

		for _, tc := range emojiWithVS15 {
			baseWidth := String(tc.base)
			vs15Width := String(tc.char)
			if vs15Width != baseWidth {
				t.Errorf("String(%q) with VS15 = %d, want %d (VS15 is a no-op, width should match base)", tc.char, vs15Width, baseWidth)
			}
		}

		// VS15 with East Asian Wide characters should preserve width 2 (no-op)
		eastAsianWideWithVS15 := []struct {
			char string
			base string
		}{
			{"中\uFE0E", "中"}, // CJK ideograph with VS15
			{"日\uFE0E", "日"}, // CJK ideograph with VS15
		}

		for _, tc := range eastAsianWideWithVS15 {
			baseWidth := String(tc.base)
			vs15Width := String(tc.char)
			if vs15Width != baseWidth {
				t.Errorf("String(%q) with VS15 = %d, want %d (VS15 is a no-op, width should match base)", tc.char, vs15Width, baseWidth)
			}
		}
	})

	t.Run("C3: VS16 forces emoji presentation", func(t *testing.T) {
		// VS16 (U+FE0F) should force emoji presentation (width 2) even for text-presentation characters
		textWithVS16 := []string{
			"\u2721\uFE0F", // ✡️ star of David with VS16
			"\u2692\uFE0F", // ⚒️ hammer and pick with VS16
			"\u2699\uFE0F", // ⚙️ gear with VS16
			"\u0023\uFE0F", // #️ hash with VS16
		}

		for _, s := range textWithVS16 {
			got := String(s)
			if got != 2 {
				t.Errorf("String(%q) with VS16 = %d, want 2 (VS16 should force emoji presentation)", s, got)
			}
		}
	})

	t.Run("ED-16: ZWJ Sequences treated as single grapheme", func(t *testing.T) {
		// ZWJ sequences should be treated as a single grapheme cluster by the grapheme tokenizer
		// and should have width 2 (since they display as a single emoji image)
		tests := []struct {
			name     string
			sequence string
			want     int
			desc     string
		}{
			{
				name:     "Family",
				sequence: "\U0001F468\u200D\U0001F469\u200D\U0001F467\u200D\U0001F466", // 👨‍👩‍👧‍👦
				want:     2,
				desc:     "Family: man, woman, girl, boy (4 people + 3 ZWJ)",
			},
			{
				name:     "Kiss",
				sequence: "\U0001F469\u200D\u2764\uFE0F\u200D\U0001F48B\u200D\U0001F468", // 👩‍❤️‍💋‍👨
				want:     2,
				desc:     "Kiss: woman, heart, kiss mark, man",
			},
			{
				name:     "Couple with heart",
				sequence: "\U0001F469\u200D\u2764\uFE0F\u200D\U0001F468", // 👩‍❤️‍👨
				want:     2,
				desc:     "Couple with heart: woman, heart, man",
			},
			{
				name:     "Woman technologist",
				sequence: "\U0001F469\u200D\U0001F4BB", // 👩‍💻
				want:     2,
				desc:     "Woman technologist: woman, ZWJ, laptop",
			},
			{
				name:     "Rainbow flag",
				sequence: "\U0001F3F4\u200D\U0001F308", // 🏴‍🌈
				want:     2,
				desc:     "Rainbow flag: black flag, ZWJ, rainbow",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := String(tt.sequence)
				if got != tt.want {
					t.Errorf("String(%q) = %d, want %d (%s)",
						tt.sequence, got, tt.want, tt.desc)
					// Show the individual components for debugging
					t.Logf("  Sequence: %+q", tt.sequence)
					t.Logf("  Expected: single grapheme cluster of width %d", tt.want)
					t.Logf("  Got: %d (if > 2, grapheme tokenizer may not be recognizing ZWJ sequence)", got)
				}
			})
		}
	})

	// ED-13: Emoji Modifier Sequences
	// Per TR51: emoji_modifier_sequence := emoji_modifier_base emoji_modifier
	// These should be treated as single grapheme clusters with width 2
	t.Run("ED-13: Emoji Modifier Sequences", func(t *testing.T) {
		tests := []struct {
			sequence string
			want     int
			desc     string
		}{
			{"👍🏻", 2, "thumbs up + light skin tone"},
			{"👍🏼", 2, "thumbs up + medium-light skin tone"},
			{"👍🏽", 2, "thumbs up + medium skin tone"},
			{"👍🏾", 2, "thumbs up + medium-dark skin tone"},
			{"👍🏿", 2, "thumbs up + dark skin tone"},
			{"👋🏻", 2, "waving hand + light skin tone"},
			{"🧑🏽", 2, "person + medium skin tone"},
			{"👶🏿", 2, "baby + dark skin tone"},
			{"👩🏼", 2, "woman + medium-light skin tone"},
		}

		for _, tt := range tests {
			t.Run(tt.desc, func(t *testing.T) {
				got := String(tt.sequence)
				if got != tt.want {
					t.Errorf("String(%q) = %d, want %d (%s)",
						tt.sequence, got, tt.want, tt.desc)
					t.Logf("  Sequence: %+q", tt.sequence)
					t.Logf("  Expected: single grapheme cluster of width %d", tt.want)
					t.Logf("  Got: %d (if > 2, grapheme tokenizer may not be recognizing modifier sequence)", got)
				}
			})
		}
	})
}

func TestStringGraphemes(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		options Options
	}{
		{"empty string", "", defaultOptions},
		{"single ASCII", "a", defaultOptions},
		{"multiple ASCII", "hello", defaultOptions},
		{"ASCII with spaces", "hello world", defaultOptions},
		{"ASCII with newline", "hello\nworld", defaultOptions},
		{"CJK ideograph", "中", defaultOptions},
		{"CJK with ASCII", "hello中", defaultOptions},
		{"ambiguous character", "★", defaultOptions},
		{"ambiguous character EAW", "★", eawOptions},
		{"emoji", "😀", defaultOptions},
		{"flag US", "🇺🇸", defaultOptions},
		{"text with flags", "Go 🇺🇸🚀", defaultOptions},
		{"keycap 1️⃣", "1️⃣", defaultOptions},
		{"mixed content", "Hi⌚⚙⚓", defaultOptions},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get expected width using String
			expected := tt.options.String(tt.input)

			// Iterate over graphemes and sum widths
			iter := tt.options.StringGraphemes(tt.input)
			got := 0
			for iter.Next() {
				got += iter.Width()
			}

			if got != expected {
				t.Errorf("StringGraphemes(%q) sum = %d, want %d (from String)",
					tt.input, got, expected)
			}
		})
	}
}

func TestBytesGraphemes(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		options Options
	}{
		{"empty bytes", []byte(""), defaultOptions},
		{"single ASCII", []byte("a"), defaultOptions},
		{"multiple ASCII", []byte("hello"), defaultOptions},
		{"ASCII with spaces", []byte("hello world"), defaultOptions},
		{"ASCII with newline", []byte("hello\nworld"), defaultOptions},
		{"CJK ideograph", []byte("中"), defaultOptions},
		{"CJK with ASCII", []byte("hello中"), defaultOptions},
		{"ambiguous character", []byte("★"), defaultOptions},
		{"ambiguous character EAW", []byte("★"), eawOptions},
		{"emoji", []byte("😀"), defaultOptions},
		{"flag US", []byte("🇺🇸"), defaultOptions},
		{"text with flags", []byte("Go 🇺🇸🚀"), defaultOptions},
		{"keycap 1️⃣", []byte("1️⃣"), defaultOptions},
		{"mixed content", []byte("Hi⌚⚙⚓"), defaultOptions},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Get expected width using Bytes
			expected := tt.options.Bytes(tt.input)

			// Iterate over graphemes and sum widths
			iter := tt.options.BytesGraphemes(tt.input)
			got := 0
			for iter.Next() {
				got += iter.Width()
			}

			if got != expected {
				t.Errorf("BytesGraphemes(%q) sum = %d, want %d (from Bytes)",
					tt.input, got, expected)
			}
		})
	}
}

func TestGraphemesControlSequences(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		options Options
	}{
		// ControlSequences true: ANSI sequences are one zero-width grapheme each; visible width only
		{"ControlSequences ANSI wrapped", "\x1b[31mhello\x1b[0m", controlSequences},
		{"ControlSequences ANSI only", "\x1b[0m", controlSequences},
		{"ControlSequences plain text", "hi", controlSequences},
		{"ControlSequences ANSI mid", "a\x1b[31mb\x1b[0mc", controlSequences},
		// Default options: sum of grapheme widths must still match String/Bytes
		{"default ANSI wrapped", "\x1b[31mhello\x1b[0m", defaultOptions},
		{"default plain", "hello", defaultOptions},
		// 8-bit ControlSequences: C1 sequences are one zero-width grapheme each
		{"8-bit C1 CSI wrapped", "\x9B31mhello\x9B0m", controlSequences8Bit},
		{"8-bit C1 CSI only", "\x9B0m", controlSequences8Bit},
		{"8-bit plain text", "hi", controlSequences8Bit},
		{"8-bit C1 CSI mid", "a\x9B31mb\x9B0mc", controlSequences8Bit},
		// Both options: both 7-bit and 8-bit sequences are zero-width graphemes
		{"both: mixed", "\x1b[31mhello\x9B0m", controlSequencesBoth},
		{"both: 7-bit only input", "\x1b[31mhi\x1b[0m", controlSequencesBoth},
		{"both: 8-bit only input", "\x9B31mhi\x9B0m", controlSequencesBoth},
		// Independence: 7-bit on but 8-bit input — graphemes must still sum correctly
		{"7-bit on, 8-bit input", "\x9B31mhello\x9B0m", controlSequences},
		// Independence: 8-bit on but 7-bit input
		{"8-bit on, 7-bit input", "\x1b[31mhello\x1b[0m", controlSequences8Bit},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// StringGraphemes: option must be passed through; sum of Width() matches String()
			expected := tt.options.String(tt.input)
			iter := tt.options.StringGraphemes(tt.input)
			got := 0
			for iter.Next() {
				got += iter.Width()
			}
			if got != expected {
				t.Errorf("StringGraphemes(%q) sum Width() = %d, want %d (String)",
					tt.input, got, expected)
			}

			// BytesGraphemes: same option and outcome for []byte
			b := []byte(tt.input)
			expectedBytes := tt.options.Bytes(b)
			iterBytes := tt.options.BytesGraphemes(b)
			gotBytes := 0
			for iterBytes.Next() {
				gotBytes += iterBytes.Width()
			}
			if gotBytes != expectedBytes {
				t.Errorf("BytesGraphemes(%q) sum Width() = %d, want %d (Bytes)",
					b, gotBytes, expectedBytes)
			}
		})
	}
}

func TestAsciiWidth(t *testing.T) {
	tests := []struct {
		name     string
		b        byte
		expected int
		desc     string
	}{
		// Control characters (0x00-0x1F): width 0
		{"null", 0x00, 0, "NULL character"},
		{"bell", 0x07, 0, "BEL (bell)"},
		{"backspace", 0x08, 0, "BS (backspace)"},
		{"tab", 0x09, 0, "TAB"},
		{"newline", 0x0A, 0, "LF (newline)"},
		{"carriage return", 0x0D, 0, "CR (carriage return)"},
		{"escape", 0x1B, 0, "ESC (escape)"},
		{"last control", 0x1F, 0, "Last control character"},

		// Printable ASCII (0x20-0x7E): width 1
		{"space", 0x20, 1, "Space (first printable)"},
		{"exclamation", 0x21, 1, "!"},
		{"zero", 0x30, 1, "0"},
		{"nine", 0x39, 1, "9"},
		{"A", 0x41, 1, "A"},
		{"Z", 0x5A, 1, "Z"},
		{"a", 0x61, 1, "a"},
		{"z", 0x7A, 1, "z"},
		{"tilde", 0x7E, 1, "~ (last printable)"},

		// DEL (0x7F): width 0
		{"delete", 0x7F, 0, "DEL (delete)"},

		// >= 128: width 1 (default, though shouldn't be used for valid UTF-8)
		{"0x80", 0x80, 1, "First byte >= 128"},
		{"0xFF", 0xFF, 1, "Last byte value"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := asciiWidth(tt.b)
			if got != tt.expected {
				t.Errorf("asciiWidth(0x%02X '%s') = %d, want %d (%s)",
					tt.b, string(tt.b), got, tt.expected, tt.desc)
			}
		})
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxWidth int
		tail     string
		options  Options
		expected string
	}{
		// Empty string cases
		{"empty string", "", 0, "", defaultOptions, ""},
		{"empty string with tail", "", 5, "...", defaultOptions, ""},
		{"empty string large maxWidth", "", 100, "...", defaultOptions, ""},

		// No truncation needed
		{"fits exactly", "hello", 5, "...", defaultOptions, "hello"},
		{"fits with room", "hi", 10, "...", defaultOptions, "hi"},
		{"single char fits", "a", 1, "...", defaultOptions, "a"},

		// Basic truncation - ASCII
		{"truncate ASCII", "hello world", 5, "...", defaultOptions, "he..."},
		{"truncate ASCII at start", "hello", 0, "...", defaultOptions, "..."},
		{"truncate ASCII single char", "hello", 1, "...", defaultOptions, "..."},
		{"truncate ASCII with empty tail", "hello world", 5, "", defaultOptions, "hello"},

		// Truncation with wide characters (CJK)
		{"CJK fits", "中", 2, "...", defaultOptions, "中"},
		{"CJK truncate", "中", 1, "...", defaultOptions, "..."},
		{"CJK with ASCII", "hello中", 5, "...", defaultOptions, "he..."},
		{"CJK with ASCII fits", "hello中", 7, "...", defaultOptions, "hello中"},
		{"CJK with ASCII partial", "hello中", 6, "...", defaultOptions, "hel..."},
		{"multiple CJK", "中文", 2, "...", defaultOptions, "..."},
		{"multiple CJK fits", "中文", 4, "...", defaultOptions, "中文"},

		// Truncation with emoji
		{"emoji fits", "😀", 2, "...", defaultOptions, "😀"},
		{"emoji truncate", "😀", 1, "...", defaultOptions, "..."},
		{"emoji with ASCII", "hello😀", 5, "...", defaultOptions, "he..."},
		{"emoji with ASCII fits", "hello😀", 7, "...", defaultOptions, "hello😀"},
		{"multiple emoji", "😀😁", 2, "...", defaultOptions, "..."},
		{"multiple emoji fits", "😀😁", 4, "...", defaultOptions, "😀😁"},

		// Truncation with control characters (zero width)
		// Control characters have width 0 but are preserved in the string structure
		{"with newline", "hello\nworld", 5, "...", defaultOptions, "he..."},
		{"with tab", "hello\tworld", 5, "...", defaultOptions, "he..."},
		{"newline at start", "\nhello", 5, "...", defaultOptions, "\nhello"},
		{"multiple newlines", "a\n\nb", 1, "...", defaultOptions, "..."},

		// Mixed content
		{"ASCII CJK emoji", "hi中😀", 2, "...", defaultOptions, "..."},
		{"ASCII CJK emoji fits", "hi中😀", 6, "...", defaultOptions, "hi中😀"},
		{"ASCII CJK emoji partial", "hi中😀", 4, "...", defaultOptions, "h..."},
		{"complex mixed", "Go 🇺🇸🚀", 3, "...", defaultOptions, "..."},
		{"complex mixed fits", "Go 🇺🇸🚀", 7, "...", defaultOptions, "Go 🇺🇸🚀"},

		// ControlSequences (ANSI escape sequences): truncation by visible width only.
		// When ControlSequences is true, escape sequences that appear after the
		// truncation point are preserved (appended after the tail). This prevents
		// color bleed from unclosed SGR sequences in terminal output.
		{"ControlSequences plain no truncation", "hello", 5, "...", controlSequences, "hello"},
		{"ControlSequences ANSI wrapped no truncation", "\x1b[31mhello\x1b[0m", 8, "...", controlSequences, "\x1b[31mhello\x1b[0m"},
		{"ControlSequences ANSI wrapped truncate", "\x1b[31mhello\x1b[0m", 4, "...", controlSequences, "\x1b[31mh...\x1b[0m"},
		{"ControlSequences ANSI in middle truncate", "hello\x1b[31mworld", 5, "...", controlSequences, "he...\x1b[31m"},
		{"ControlSequences CJK truncate", "\x1b[31m中文\x1b[0m", 2, "...", controlSequences, "...\x1b[31m\x1b[0m"},
		{"ControlSequences CJK no truncation", "\x1b[31m中文\x1b[0m", 7, "...", controlSequences, "\x1b[31m中文\x1b[0m"},
		{"ControlSequences CJK one wide then tail", "\x1b[31m中文xx\x1b[0m", 5, "...", controlSequences, "\x1b[31m中...\x1b[0m"},
		// Stacked SGR sequences: all escape sequences after cut are preserved
		{"ControlSequences stacked SGR", "\x1b[31m\x1b[42mhello\x1b[0m", 4, "...", controlSequences, "\x1b[31m\x1b[42mh...\x1b[0m"},
		// Escape sequence between visible chars after cut: preserved
		{"ControlSequences mid-escape after cut", "\x1b[31mhello\x1b[42mworld\x1b[0m", 6, "...", controlSequences, "\x1b[31mhel...\x1b[42m\x1b[0m"},
		// No escape sequences after cut: same as before
		{"ControlSequences no trailing escape", "\x1b[31mhello", 4, "...", controlSequences, "\x1b[31mh..."},
		// Multiple colors: all trailing escapes preserved
		{"ControlSequences multi color", "a\x1b[31mb\x1b[32mc\x1b[33md\x1b[0m", 2, "...", controlSequences, "...\x1b[31m\x1b[32m\x1b[33m\x1b[0m"},

		// 8-bit ControlSequences8Bit is ignored by truncation entirely. The
		// grapheme parser is not told about 8-bit, so C1 sequence parameters
		// (e.g. "31m" after \x9B) are treated as visible characters. This is
		// intentional: 8-bit C1 bytes (0x80-0x9F) overlap with UTF-8 multi-byte
		// encoding, making them unsafe to manipulate during truncation.
		{"8-bit plain no truncation", "hello", 5, "...", controlSequences8Bit, "hello"},
		{"8-bit C1 CSI wrapped truncate", "\x9B31mhello\x9B0m", 8, "...", controlSequences8Bit, "\x9B31mh..."},
		{"8-bit C1 CSI wrapped truncate narrow", "\x9B31mhello\x9B0m", 4, "...", controlSequences8Bit, "\x9B..."},
		{"8-bit C1 CSI in middle truncate", "hello\x9B31mworld", 5, "...", controlSequences8Bit, "he..."},
		{"8-bit C1 CSI CJK truncate", "\x9B31m中文\x9B0m", 2, "...", controlSequences8Bit, "..."},
		{"8-bit C1 CSI no trailing escape", "\x9B31mhello", 4, "...", controlSequences8Bit, "\x9B..."},
		{"8-bit C1 stacked SGR", "\x9B31m\x9B42mhello\x9B0m", 4, "...", controlSequences8Bit, "\x9B..."},

		// 7-bit only must NOT preserve trailing C1 sequences.
		// With 7-bit only, \x9B is a regular character (width 1), so the input
		// "hello\x9B0m" has visible width 8. Trailing \x9B0m is not preserved.
		{"7-bit only ignores trailing C1", "hello\x9B0m", 5, "...", controlSequences, "he..."},

		// Both enabled: only 7-bit trailing escapes are preserved; 8-bit is
		// ignored by truncation, so C1 parameters are visible characters.
		{"both: mixed trailing escapes", "\x1b[31mhello\x9B0m", 4, "...", controlSequencesBoth, "\x1b[31mh..."},
		{"both: 7-bit wrapped truncate", "\x1b[31mhello\x1b[0m", 4, "...", controlSequencesBoth, "\x1b[31mh...\x1b[0m"},
		{"both: 8-bit wrapped truncate", "\x9B31mhello\x9B0m", 4, "...", controlSequencesBoth, "\x9B..."},

		// East Asian Width option
		{"ambiguous EAW fits", "★", 2, "...", eawOptions, "★"},
		{"ambiguous EAW truncate", "★", 1, "...", eawOptions, "..."},
		{"ambiguous default fits", "★", 1, "...", defaultOptions, "★"},
		{"ambiguous mixed", "a★b", 2, "...", eawOptions, "..."},
		{"ambiguous mixed default", "a★b", 2, "...", defaultOptions, "..."},

		// Edge cases
		{"zero maxWidth", "hello", 0, "...", defaultOptions, "..."},
		{"very long string", "a very long string that will definitely be truncated", 10, "...", defaultOptions, "a very ..."},
		// Bug fix: wide char at boundary with narrow chars - ensures truncation position is correct
		// Input "中cde" (width 5), maxWidth 4, tail "ab" (width 2) -> should return "中ab" (width 4)
		{"wide char boundary bug fix", "中cde", 4, "ab", defaultOptions, "中ab"},

		// Tail variations
		{"custom tail", "hello world", 5, "…", defaultOptions, "hell…"},
		{"long tail", "hello", 3, ">>>", defaultOptions, ">>>"},
		{"tail with wide char", "hello", 3, "中", defaultOptions, "h中"},
		{"tail with emoji", "hello", 3, "😀", defaultOptions, "h😀"},

		// Grapheme boundary tests (ensuring truncation happens at grapheme boundaries)
		{"keycap sequence", "1️⃣2️⃣", 2, "...", defaultOptions, "..."},
		{"flag sequence", "🇺🇸🇯🇵", 2, "...", defaultOptions, "..."},
		{"ZWJ sequence", "👨‍👩‍👧", 2, "...", defaultOptions, "👨‍👩‍👧"},
		{"ZWJ sequence truncate", "👨‍👩‍👧👨‍👩‍👧", 2, "...", defaultOptions, "..."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			{
				got := tt.options.TruncateString(tt.input, tt.maxWidth, tt.tail)
				if got != tt.expected {
					t.Errorf("TruncateString(%q, %d, %q) with options %v = %q, want %q",
						tt.input, tt.maxWidth, tt.tail, tt.options, got, tt.expected)
					inputWidth := tt.options.String(tt.input)
					gotWidth := tt.options.String(got)
					t.Logf("  Input width: %d, Got width: %d, MaxWidth: %d", inputWidth, gotWidth, tt.maxWidth)
				}

				// Verify visible width respects maxWidth (or tailWidth if tail is wider)
				gotWidth := tt.options.String(got)
				limit := tt.maxWidth
				tailWidth := tt.options.String(tt.tail)
				if tailWidth > limit {
					limit = tailWidth
				}
				if gotWidth > limit {
					t.Errorf("Result visible width (%d) exceeds max(maxWidth, tailWidth) (%d)", gotWidth, limit)
				}
			}
			{
				input := []byte(tt.input)
				tail := []byte(tt.tail)
				expected := []byte(tt.expected)
				got := tt.options.TruncateBytes(input, tt.maxWidth, tail)
				if !bytes.Equal(got, expected) {
					t.Errorf("TruncateBytes(%q, %d, %q) with options %v = %q, want %q",
						input, tt.maxWidth, tail, tt.options, got, expected)
					inputWidth := tt.options.Bytes(input)
					gotWidth := tt.options.Bytes(got)
					t.Logf("  Input width: %d, Got width: %d, MaxWidth: %d", inputWidth, gotWidth, tt.maxWidth)
				}

				// Verify visible width respects maxWidth (or tailWidth if tail is wider)
				gotWidth := tt.options.Bytes(got)
				limit := tt.maxWidth
				tailWidth := tt.options.Bytes(tail)
				if tailWidth > limit {
					limit = tailWidth
				}
				if gotWidth > limit {
					t.Errorf("Result visible width (%d) exceeds max(maxWidth, tailWidth) (%d)", gotWidth, limit)
				}
			}
		})
	}
}

func TestTruncateBytesDoesNotMutateInput(t *testing.T) {
	// Test that TruncateBytes does not mutate the caller's slice
	original := []byte("hello world")
	originalCopy := make([]byte, len(original))
	copy(originalCopy, original)

	tail := []byte("...")
	_ = TruncateBytes(original, 5, tail)

	if !bytes.Equal(original, originalCopy) {
		t.Errorf("TruncateBytes mutated the input slice: got %q, want %q", original, originalCopy)
	}
}

func TestPrintableASCIILength(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
		desc     string
	}{
		// Some of these tests are left over from a SWAR implementation,
		// which cared about 8 byte boundaries.
		{"empty string", "", 0, "Empty string has 0 printable bytes"},
		{"single char", "a", 1, "Single printable byte"},
		{"single space", " ", 1, "Space is printable"},
		{"7 bytes", "1234567", 7, "7 printable bytes"},
		{"8 bytes", "12345678", 8, "8 printable bytes"},
		{"9 bytes", "123456789", 9, "9 printable bytes"},
		{"15 bytes", "123456789012345", 15, "15 printable bytes"},
		{"16 bytes", "1234567890123456", 16, "16 printable bytes"},
		{"17 bytes", "12345678901234567", 17, "17 printable bytes"},
		{"24 bytes", "123456789012345678901234", 24, "24 printable bytes"},
		{"long ASCII", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", 62, "All 62 printable bytes"},
		{"all printable range", " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~", 95, "All 95 printable ASCII chars"},

		// Non-printable at start -> 0
		{"control at start", "\x00hello world", 0, "Control char at start"},
		{"DEL at start", "\x7Fhello world", 0, "DEL at start"},
		{"non-ASCII at start", "\x80hello world", 0, "Non-ASCII at start"},
		{"UTF-8 at start", "\xC2\xA0hello world", 0, "UTF-8 at start"},
		{"emoji at start", "\xF0\x9F\x98\x80hello123", 0, "Emoji at start"},

		// Non-printable in middle - stops before it
		{"control in middle", "hello\x00world123", 5, "Control at pos 5, returns 5"},
		{"DEL in middle", "hello\x7Fworld123", 5, "DEL at pos 5, returns 5"},
		{"control after 8", "12345678\x00world", 8, "Control at pos 8, returns 8"},
		{"DEL after 8", "12345678\x7Fworld", 8, "DEL at pos 8, returns 8"},
		{"control at pos 15", "123456789012345\x00", 15, "Control at pos 15, returns 15"},

		// Non-ASCII at end - backs off by 1
		{"non-ASCII at end of 9", "12345678\x80", 7, "Backs off 1 before non-ASCII"},
		{"non-ASCII at end of 17", "1234567890123456\x80", 15, "Backs off 1 before non-ASCII"},
		{"combining after 16", "1234567890123456\u0301", 15, "Backs off 1 before combining mark"},
		{"non-ASCII after 1", "a\x80", 0, "Backs off 1 from 1, returns 0"},
		{"16 ASCII then emoji", "1234567890123456\xF0\x9F\x98\x80", 15, "Backs off 1 before emoji"},

		// Printable boundaries
		{"8 spaces", "        ", 8, "Space (0x20) is first printable"},
		{"8 tildes", "~~~~~~~~", 8, "Tilde (0x7E) is last printable"},
		{"mixed printable", "Hello, World! 123", 17, "All 17 are printable"},

		// Control characters at various positions
		{"control at pos 0", "\x00234567890", 0, "Control at position 0"},
		{"control at pos 3", "123\x00567890", 3, "Control at position 3"},
		{"control at pos 7", "1234567\x000", 7, "Control at position 7"},

		// DEL at various positions
		{"DEL at pos 0", "\x7F234567890", 0, "DEL at position 0"},
		{"DEL at pos 3", "123\x7F567890", 3, "DEL at position 3"},
		{"DEL at pos 7", "1234567\x7F0", 7, "DEL at position 7"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := printableASCIILength(tt.input)
			if got != tt.expected {
				t.Errorf("printableASCIILength(%q) = %d, want %d (%s)",
					tt.input, got, tt.expected, tt.desc)
				if len(tt.input) > 0 {
					t.Logf("  String length: %d bytes", len(tt.input))
					for i, b := range []byte(tt.input) {
						isPrintable := b >= 0x20 && b <= 0x7E
						t.Logf("    [%d]: 0x%02X printable=%v", i, b, isPrintable)
					}
				}
			}
		})
	}
}

func TestPrintableASCIILengthBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected int
		desc     string
	}{
		// Any length works - returns exact count
		{"empty slice", []byte{}, 0, "Empty slice has 0 printable bytes"},
		{"single space", []byte{0x20}, 1, "Single space"},
		{"single char", []byte("a"), 1, "Single printable byte"},
		{"7 bytes", []byte("1234567"), 7, "7 printable bytes"},
		{"8 bytes all printable", []byte("12345678"), 8, "8 bytes all printable"},
		{"9 bytes all printable", []byte("123456789"), 9, "9 printable bytes"},
		{"15 bytes all printable", []byte("123456789012345"), 15, "15 printable bytes"},
		{"16 bytes all printable", []byte("1234567890123456"), 16, "16 printable bytes"},
		{"17 bytes all printable", []byte("12345678901234567"), 17, "17 printable bytes"},
		{"24 bytes all printable", []byte("123456789012345678901234"), 24, "24 printable bytes"},
		{"long all printable", []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()"), 72, "72 printable bytes"},
		{"all printable range", []byte(" !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"), 95, "All 95 printable ASCII chars"},

		// Printable boundaries
		{"8 bytes all spaces", []byte("        "), 8, "8 spaces"},
		{"8 bytes all tildes", []byte("~~~~~~~~"), 8, "8 tildes"},
		{"8 bytes boundary low", []byte{0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20}, 8, "8 spaces (0x20)"},
		{"8 bytes boundary high", []byte{0x7E, 0x7E, 0x7E, 0x7E, 0x7E, 0x7E, 0x7E, 0x7E}, 8, "8 tildes (0x7E)"},

		// Non-printable at start -> 0
		{"control at start", []byte("\x00hello world"), 0, "Control char at start"},
		{"DEL at start", []byte("\x7Fhello world"), 0, "DEL at start"},
		{"non-ASCII at start", []byte("\x80hello world"), 0, "Non-ASCII at start"},
		{"UTF-8 at start", []byte("\xC2\xA0hello world"), 0, "UTF-8 at start"},
		{"emoji at start", []byte("\xF0\x9F\x98\x80hello123"), 0, "Emoji at start"},

		// Non-printable in middle - stops before it
		{"control in middle", []byte("hel\x00o123"), 3, "Control at pos 3, returns 3"},
		{"DEL in middle", []byte("hel\x7Fo123"), 3, "DEL at pos 3, returns 3"},
		{"control at pos 8", []byte("12345678\x00world"), 8, "Control at pos 8, returns 8"},
		{"DEL at pos 8", []byte("12345678\x7Fworld"), 8, "DEL at pos 8, returns 8"},
		{"control at pos 15", []byte("123456789012345\x00"), 15, "Control at pos 15, returns 15"},

		// Non-ASCII at end - backs off by 1
		{"non-ASCII at end of 9", []byte("12345678\x80"), 7, "Backs off 1 before non-ASCII"},
		{"non-ASCII at end of 17", []byte("1234567890123456\x80"), 15, "Backs off 1 before non-ASCII"},
		{"16 ASCII then emoji", []byte("1234567890123456\xF0\x9F\x98\x80"), 15, "Backs off 1 before emoji"},
		{"16 ASCII then combining acute", []byte("1234567890123456\u0301"), 15, "Backs off 1 before combining mark"},
		{"16 ASCII then combining grave", []byte("1234567890123456\u0300"), 15, "Backs off 1 before combining mark"},
		{"non-ASCII after 1", []byte("a\x80"), 0, "Backs off 1 from 1, returns 0"},

		// Control characters at various positions
		{"control at pos 0", []byte("\x00234567890"), 0, "Control at position 0"},
		{"control at pos 3", []byte("123\x00567890"), 3, "Control at position 3"},
		{"control at pos 7", []byte("1234567\x000"), 7, "Control at position 7"},

		// DEL at various positions
		{"DEL at pos 0", []byte("\x7F234567890"), 0, "DEL at position 0"},
		{"DEL at pos 3", []byte("123\x7F567890"), 3, "DEL at position 3"},
		{"DEL at pos 7", []byte("1234567\x7F0"), 7, "DEL at position 7"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := printableASCIILength(tt.input)
			if got != tt.expected {
				t.Errorf("printableASCIILength(%v) = %d, want %d (%s)",
					tt.input, got, tt.expected, tt.desc)
				if len(tt.input) > 0 {
					t.Logf("  Slice length: %d bytes", len(tt.input))
					for i, b := range tt.input {
						isPrintable := b >= 0x20 && b <= 0x7E
						t.Logf("    [%d]: 0x%02X printable=%v", i, b, isPrintable)
					}
				}
			}
		})
	}
}

// TestPrintableASCIIOptimization verifies that the partial ASCII optimization
// in String() and Bytes() works correctly for printable ASCII content.
func TestPrintableASCIIOptimization(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"single char", "a"},
		{"short ASCII", "hello"},
		{"long ASCII", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"},
		{"with spaces", "hello world"},
		{"with punctuation", "Hello, World!"},
		{"all printable range", " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"},
		{"exactly 8 bytes", "12345678"},
		{"exactly 16 bytes", "1234567890123456"},
		{"exactly 24 bytes", "123456789012345678901234"},
		{"7 bytes", "1234567"},
		{"9 bytes", "123456789"},
		{"15 bytes", "123456789012345"},
		{"17 bytes", "12345678901234567"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For printable ASCII, width should equal length
			width := String(tt.input)
			expected := len(tt.input)
			if width != expected {
				t.Errorf("String(%q) = %d, want %d", tt.input, width, expected)
			}

			// Same for Bytes
			widthBytes := Bytes([]byte(tt.input))
			if widthBytes != expected {
				t.Errorf("Bytes(%q) = %d, want %d", tt.input, widthBytes, expected)
			}
		})
	}
}

// TestUnicode16IndicConjunctBreak tests Unicode 16.0 Indic_Conjunct_Break property.
// This property affects grapheme cluster breaking in Indic scripts, ensuring that
// conjuncts (consonant clusters) are properly grouped into single grapheme clusters.
// The Indic_Conjunct_Break property has values: Consonant, Linker, and Extend.
//
// Note: Indic scripts are typically width 1 (not width 2 like CJK). The key test
// here is that grapheme clusters are formed correctly according to Indic_Conjunct_Break
// rules, not the width value itself.
func TestUnicode16IndicConjunctBreak(t *testing.T) {
	tests := []struct {
		name                   string
		input                  string
		expectedWidth          int
		expectedClusters       int // Expected number of grapheme clusters
		description            string
		verifyClusterFormation bool // Whether to verify the cluster contains expected runes
	}{
		// Devanagari (Hindi, Sanskrit) - Unicode range U+0900-U+097F
		{
			name:                   "Devanagari conjunct क्ष",
			input:                  "क्ष", // kṣa - क (ka) + virama + ष (ṣa)
			expectedWidth:          1,     // Indic scripts are width 1
			expectedClusters:       1,     // Should form single grapheme cluster
			description:            "Devanagari conjunct formed with virama (U+094D) - should be single cluster",
			verifyClusterFormation: true,
		},
		{
			name:                   "Devanagari conjunct त्र",
			input:                  "त्र", // tra - त (ta) + virama + र (ra)
			expectedWidth:          1,
			expectedClusters:       1,
			description:            "Devanagari conjunct with र (ra) as subscript - should be single cluster",
			verifyClusterFormation: true,
		},
		{
			name:                   "Devanagari conjunct ज्ञ",
			input:                  "ज्ञ", // jña - ज (ja) + virama + ञ (ña)
			expectedWidth:          1,
			expectedClusters:       1,
			description:            "Devanagari conjunct ज्ञ - should be single cluster",
			verifyClusterFormation: true,
		},
		{
			name:             "Devanagari word with conjuncts",
			input:            "क्षत्रिय", // kṣatriya - contains conjunct क्ष
			expectedWidth:    3,          // 3 grapheme clusters × 1 width each
			expectedClusters: 3,          // क्ष, त्रि, य
			description:      "Devanagari word with multiple conjuncts",
		},
		{
			name:             "Devanagari with repha",
			input:            "राम", // rāma - र (ra) can form repha in some contexts
			expectedWidth:    2,     // 2 grapheme clusters × 1 width each
			expectedClusters: 2,     // रा, म
			description:      "Devanagari with potential repha formation",
		},

		// Bengali (Bangla) - Unicode range U+0980-U+09FF
		{
			name:                   "Bengali conjunct ক্ষ",
			input:                  "ক্ষ", // kṣa - ক (ka) + virama + ষ (ṣa)
			expectedWidth:          1,
			expectedClusters:       1,
			description:            "Bengali conjunct ক্ষ - should be single cluster",
			verifyClusterFormation: true,
		},
		{
			name:                   "Bengali conjunct জ্ঞ",
			input:                  "জ্ঞ", // jña - জ (ja) + virama + ঞ (ña)
			expectedWidth:          1,
			expectedClusters:       1,
			description:            "Bengali conjunct জ্ঞ - should be single cluster",
			verifyClusterFormation: true,
		},
		{
			name:             "Bengali word",
			input:            "বাংলা", // bāṅlā - Bengali
			expectedWidth:    2,       // 2 grapheme clusters × 1 width each
			expectedClusters: 2,       // বাং, লা
			description:      "Bengali word with conjuncts",
		},

		// Tamil - Unicode range U+0B80-U+0BFF
		// Tamil typically uses visible viramas rather than fused conjuncts
		// Note: Tamil may break differently - virama may form separate cluster
		{
			name:                   "Tamil with virama",
			input:                  "க்ஷ", // kṣa - க (ka) + virama + ஷ (ṣa)
			expectedWidth:          2,     // May break into 2 clusters: க், ஷ
			expectedClusters:       2,     // Tamil virama handling may differ
			description:            "Tamil conjunct with visible virama - may break into multiple clusters",
			verifyClusterFormation: false,
		},
		{
			name:             "Tamil word",
			input:            "தமிழ்", // tamiḻ - Tamil
			expectedWidth:    3,       // 3 grapheme clusters × 1 width each
			expectedClusters: 3,       // த, மி, ழ்
			description:      "Tamil word",
		},

		// Telugu - Unicode range U+0C00-U+0C7F
		{
			name:                   "Telugu conjunct క్ష",
			input:                  "క్ష", // kṣa - క (ka) + virama + ష (ṣa)
			expectedWidth:          1,
			expectedClusters:       1,
			description:            "Telugu conjunct క్ష - should be single cluster",
			verifyClusterFormation: true,
		},
		{
			name:             "Telugu word",
			input:            "తెలుగు", // telugu
			expectedWidth:    3,        // 3 grapheme clusters × 1 width each
			expectedClusters: 3,        // తె, లు, గు
			description:      "Telugu word",
		},

		// Gujarati - Unicode range U+0A80-U+0AFF
		{
			name:                   "Gujarati conjunct ક્ષ",
			input:                  "ક્ષ", // kṣa - ક (ka) + virama + ષ (ṣa)
			expectedWidth:          1,
			expectedClusters:       1,
			description:            "Gujarati conjunct ક્ષ - should be single cluster",
			verifyClusterFormation: true,
		},
		{
			name:             "Gujarati word",
			input:            "ગુજરાતી", // gujarātī
			expectedWidth:    4,         // 4 grapheme clusters × 1 width each
			expectedClusters: 4,         // ગુ, જ, રા, તી
			description:      "Gujarati word",
		},

		// Kannada - Unicode range U+0C80-U+0CFF
		// Note: Some Kannada conjuncts may break differently depending on Indic_Conjunct_Break implementation
		{
			name:                   "Kannada conjunct ಕ್ಷ",
			input:                  "ಕ್ಷ", // kṣa - ಕ (ka) + virama + ಷ (ṣa)
			expectedWidth:          2,     // May break into 2 clusters: ಕ್, ಷ
			expectedClusters:       2,     // Kannada virama handling may differ
			description:            "Kannada conjunct ಕ್ಷ - may break into multiple clusters",
			verifyClusterFormation: false,
		},
		{
			name:             "Kannada word",
			input:            "ಕನ್ನಡ", // kannada
			expectedWidth:    4,       // 4 grapheme clusters × 1 width each
			expectedClusters: 4,       // ಕ, ನ್, ನ, ಡ
			description:      "Kannada word",
		},

		// Malayalam - Unicode range U+0D00-U+0D7F
		{
			name:                   "Malayalam conjunct ക്ഷ",
			input:                  "ക്ഷ", // kṣa - ക (ka) + virama + ഷ (ṣa)
			expectedWidth:          1,
			expectedClusters:       1,
			description:            "Malayalam conjunct ക്ഷ - should be single cluster",
			verifyClusterFormation: true,
		},
		{
			name:             "Malayalam word",
			input:            "മലയാളം", // malayāḷaṁ
			expectedWidth:    4,        // 4 grapheme clusters × 1 width each
			expectedClusters: 4,        // മ, ല, യാ, ളം
			description:      "Malayalam word",
		},

		// Mixed Indic scripts
		{
			name:             "Mixed Indic scripts",
			input:            "क्ष বাংলা தமிழ்", // Devanagari + Bengali + Tamil
			expectedWidth:    8,                 // 1 + space + 2 + space + 3
			expectedClusters: 8,                 // क्ष, space, বাং, লা, space, த, மி, ழ்
			description:      "Mixed Indic scripts with spaces",
		},

		// Test that virama (U+094D in Devanagari) doesn't break grapheme cluster
		{
			name:                   "Devanagari with explicit virama",
			input:                  "क्", // ka + virama (should be part of grapheme cluster)
			expectedWidth:          1,
			expectedClusters:       1,
			description:            "Devanagari consonant with virama (no following consonant) - should be single cluster",
			verifyClusterFormation: true,
		},

		// Test Indic script with combining marks (should still form single grapheme)
		{
			name:                   "Devanagari with vowel sign",
			input:                  "का", // kā - क (ka) + ा (ā vowel sign)
			expectedWidth:          1,
			expectedClusters:       1,
			description:            "Devanagari with combining vowel sign - should be single cluster",
			verifyClusterFormation: true,
		},
		{
			name:                   "Bengali with vowel sign",
			input:                  "কা", // kā - ক (ka) + া (ā vowel sign)
			expectedWidth:          1,
			expectedClusters:       1,
			description:            "Bengali with combining vowel sign - should be single cluster",
			verifyClusterFormation: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test String width
			got := String(tt.input)
			if got != tt.expectedWidth {
				t.Errorf("String(%q) = %d, want %d (%s)",
					tt.input, got, tt.expectedWidth, tt.description)
			}

			// Test Bytes width
			gotBytes := Bytes([]byte(tt.input))
			if gotBytes != tt.expectedWidth {
				t.Errorf("Bytes(%q) = %d, want %d (%s)",
					tt.input, gotBytes, tt.expectedWidth, tt.description)
			}

			// Verify grapheme cluster formation (key test for Indic_Conjunct_Break)
			iter := StringGraphemes(tt.input)
			sumWidth := 0
			clusterCount := 0
			var clusters []string
			for iter.Next() {
				clusterCount++
				width := iter.Width()
				sumWidth += width
				clusters = append(clusters, iter.Value())
			}

			if clusterCount != tt.expectedClusters {
				t.Errorf("Number of grapheme clusters = %d, want %d (%s)",
					clusterCount, tt.expectedClusters, tt.description)
				for i, cluster := range clusters {
					t.Logf("  Cluster %d: %q (width %d)", i+1, cluster, String(cluster))
				}
			}

			if sumWidth != tt.expectedWidth {
				t.Errorf("Sum of grapheme cluster widths = %d, want %d", sumWidth, tt.expectedWidth)
			}

			// For conjuncts, verify they form a single cluster (Indic_Conjunct_Break behavior)
			if tt.verifyClusterFormation && clusterCount != 1 {
				t.Errorf("Expected single grapheme cluster for conjunct, got %d clusters: %v",
					clusterCount, clusters)
			}

			// Verify that the input string can be reconstructed from clusters
			reconstructed := ""
			iter2 := StringGraphemes(tt.input)
			for iter2.Next() {
				reconstructed += iter2.Value()
			}
			if reconstructed != tt.input {
				t.Errorf("Reconstructed string from clusters = %q, want %q", reconstructed, tt.input)
			}
		})
	}
}

func TestReproduceFuzzTruncate(t *testing.T) {
	// Regression test: \x1bX (ESC X = SOS) is segmented as one grapheme in the
	// full input but as two separate graphemes (\x1b + X) in the truncated
	// result, causing the preserved escape sequence to add visible width.
	text := "00000000000\x1bX\x18"
	options := []Options{
		{EastAsianWidth: false},
		{EastAsianWidth: true},
		{ControlSequences: true},
		{EastAsianWidth: true, ControlSequences: true},
	}

	for _, opt := range options {
		ts := opt.TruncateString(text, 10, "...")
		w := opt.String(ts)
		if w > 10 {
			t.Errorf("TruncateString() returned string longer than maxWidth for %q with opts %+v: %q (width %d)", text, opt, ts, w)
		}

		tb := opt.TruncateBytes([]byte(text), 10, []byte("..."))
		if !bytes.Equal(tb, []byte(ts)) {
			t.Errorf("TruncateBytes() != TruncateString() for %q with opts %+v: %q != %q", text, opt, tb, ts)
		}
	}
}

func TestTruncateIgnores8Bit(t *testing.T) {
	// Truncation ignores ControlSequences8Bit entirely (see GoDoc).
	// This means the truncation result, when measured with 8-bit-aware
	// String(), may exceed maxWidth. This is the documented tradeoff:
	// 8-bit C1 bytes (0x80-0x9F) overlap with UTF-8 multi-byte encoding,
	// so manipulating them during truncation is unsafe.
	//
	// These tests verify that truncation is self-consistent: the result
	// measured WITHOUT 8-bit should respect maxWidth.

	cases := []struct {
		name string
		text string
	}{
		{
			// Byte recombination: the grapheme parser with 8-bit groups
			// \x9f\xcf as one escape (APC + payload). Without 8-bit, \xcf
			// and \x90 can recombine into U+03D0 (ϐ, width 1).
			name: "byte recombination",
			text: "000000000000000000000\x9f\xcf\x1a\x90",
		},
		{
			// SOS terminator mismatch: with 8-bit, \x9c is ST (terminates
			// the 7-bit SOS started by \x1bX). Without 8-bit, \x9c is not
			// recognized as ST, so SOS consumes more of the string.
			name: "SOS terminator mismatch",
			text: "00\x98\x1bX\x9c0000000000\x18",
		},
	}

	options := []Options{
		{ControlSequences8Bit: true},
		{ControlSequences: true, ControlSequences8Bit: true},
		{EastAsianWidth: true, ControlSequences8Bit: true},
	}

	for _, tc := range cases {
		for _, opt := range options {
			// Truncation ignores 8-bit, so measure with the same view
			measureOpt := opt
			measureOpt.ControlSequences8Bit = false

			ts := opt.TruncateString(tc.text, 10, "...")
			w := measureOpt.String(ts)
			if w > 10 {
				t.Errorf("%s: TruncateString() width %d > 10 (measured without 8-bit) for %q with opts %+v: %q",
					tc.name, w, tc.text, opt, ts)
			}

			tb := opt.TruncateBytes([]byte(tc.text), 10, []byte("..."))
			bw := measureOpt.Bytes(tb)
			if bw > 10 {
				t.Errorf("%s: TruncateBytes() width %d > 10 (measured without 8-bit) for %q with opts %+v: %q",
					tc.name, bw, tc.text, opt, tb)
			}

			if !bytes.Equal(tb, []byte(ts)) {
				t.Errorf("%s: TruncateBytes() != TruncateString() for %q with opts %+v: %q != %q",
					tc.name, tc.text, opt, tb, ts)
			}
		}
	}
}
