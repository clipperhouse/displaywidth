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
					// Show width information for debugging
					inputWidth := tt.options.String(tt.input)
					gotWidth := tt.options.String(got)
					t.Logf("  Input width: %d, Got width: %d, MaxWidth: %d", inputWidth, gotWidth, tt.maxWidth)
				}

				if len(got) >= len(tt.tail) && tt.tail != "" {
					truncatedPart := got[:len(got)-len(tt.tail)]
					truncatedWidth := tt.options.String(truncatedPart)
					if truncatedWidth > tt.maxWidth {
						t.Errorf("Truncated part width (%d) exceeds maxWidth (%d)", truncatedWidth, tt.maxWidth)
					}
				} else if tt.tail == "" {
					// If no tail, the result itself should fit within maxWidth
					gotWidth := tt.options.String(got)
					if gotWidth > tt.maxWidth {
						t.Errorf("Result width (%d) exceeds maxWidth (%d) when tail is empty", gotWidth, tt.maxWidth)
					}
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
					// Show width information for debugging
					inputWidth := tt.options.Bytes(input)
					gotWidth := tt.options.Bytes(got)
					t.Logf("  Input width: %d, Got width: %d, MaxWidth: %d", inputWidth, gotWidth, tt.maxWidth)
				}

				if len(got) >= len(tt.tail) && tt.tail != "" {
					truncatedPart := got[:len(got)-len(tt.tail)]
					truncatedWidth := tt.options.Bytes(truncatedPart)
					if truncatedWidth > tt.maxWidth {
						t.Errorf("Truncated part width (%d) exceeds maxWidth (%d)", truncatedWidth, tt.maxWidth)
					}
				} else if tt.tail == "" {
					// If no tail, the result itself should fit within maxWidth
					gotWidth := tt.options.Bytes(got)
					if gotWidth > tt.maxWidth {
						t.Errorf("Result width (%d) exceeds maxWidth (%d) when tail is empty", gotWidth, tt.maxWidth)
					}
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

func TestIsPrintableASCII(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
		desc     string
	}{
		// Empty string
		{"empty string", "", true, "Empty string is printable ASCII"},

		// Single byte cases
		{"single space", " ", true, "Space (0x20) is printable"},
		{"single exclamation", "!", true, "! (0x21) is printable"},
		{"single tilde", "~", true, "~ (0x7E) is printable"},
		{"single null", "\x00", false, "NULL (0x00) is control char"},
		{"single tab", "\t", false, "TAB (0x09) is control char"},
		{"single newline", "\n", false, "LF (0x0A) is control char"},
		{"single last control", "\x1F", false, "Last control (0x1F) is control char"},
		{"single DEL", "\x7F", false, "DEL (0x7F) is not printable"},
		{"single 0x80", "\x80", false, "0x80 is non-ASCII"},
		{"single 0xFF", "\xFF", false, "0xFF is non-ASCII"},

		// Two byte cases
		{"two printable", "ab", true, "Two printable ASCII"},
		{"two with control", "a\x00", false, "Contains control char"},
		{"two with DEL", "a\x7F", false, "Contains DEL"},
		{"two with non-ASCII", "a\x80", false, "Contains non-ASCII"},

		// Exactly 8 bytes (SWAR boundary)
		{"8 bytes all printable", "12345678", true, "8 bytes all printable"},
		{"8 bytes with space", "hello wo", true, "8 bytes with space"},
		{"8 bytes with control", "hello\x00", false, "8 bytes with control char"},
		{"8 bytes with DEL", "hello\x7F", false, "8 bytes with DEL"},
		{"8 bytes with non-ASCII", "hello\x80", false, "8 bytes with non-ASCII"},
		{"8 bytes all spaces", "        ", true, "8 spaces"},
		{"8 bytes all tildes", "~~~~~~~~", true, "8 tildes"},
		{"8 bytes boundary low", "\x20\x20\x20\x20\x20\x20\x20\x20", true, "8 spaces (0x20)"},
		{"8 bytes boundary high", "\x7E\x7E\x7E\x7E\x7E\x7E\x7E\x7E", true, "8 tildes (0x7E)"},

		// Just before 8 bytes (7 bytes)
		{"7 bytes all printable", "1234567", true, "7 bytes all printable"},
		{"7 bytes with control", "123456\x00", false, "7 bytes with control char"},
		{"7 bytes with DEL", "123456\x7F", false, "7 bytes with DEL"},

		// Just after 8 bytes (9 bytes)
		{"9 bytes all printable", "123456789", true, "9 bytes all printable"},
		{"9 bytes with control", "12345678\x00", false, "9 bytes with control char"},
		{"9 bytes with DEL", "12345678\x7F", false, "9 bytes with DEL"},
		{"9 bytes with non-ASCII", "12345678\x80", false, "9 bytes with non-ASCII"},

		// Exactly 16 bytes (two SWAR chunks)
		{"16 bytes all printable", "1234567890123456", true, "16 bytes all printable"},
		{"16 bytes with control", "123456789012345\x00", false, "16 bytes with control char"},
		{"16 bytes with DEL", "123456789012345\x7F", false, "16 bytes with DEL"},
		{"16 bytes with non-ASCII", "123456789012345\x80", false, "16 bytes with non-ASCII"},

		// Just before 16 bytes (15 bytes)
		{"15 bytes all printable", "123456789012345", true, "15 bytes all printable"},
		{"15 bytes with control", "12345678901234\x00", false, "15 bytes with control char"},

		// Just after 16 bytes (17 bytes)
		{"17 bytes all printable", "12345678901234567", true, "17 bytes all printable"},
		{"17 bytes with control", "1234567890123456\x00", false, "17 bytes with control char"},

		// Exactly 24 bytes (three SWAR chunks)
		{"24 bytes all printable", "123456789012345678901234", true, "24 bytes all printable"},
		{"24 bytes with control", "12345678901234567890123\x00", false, "24 bytes with control char"},

		// Long strings
		{"long all printable", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()", true, "Long string all printable"},
		{"long with control", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!\x00@#$%^&*()", false, "Long string with control char"},
		{"long with DEL", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!\x7F@#$%^&*()", false, "Long string with DEL"},
		{"long with non-ASCII", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!\x80@#$%^&*()", false, "Long string with non-ASCII"},

		// Edge cases: first and last bytes in range
		{"starts with space", " hello", true, "Starts with space (0x20)"},
		{"starts with tilde", "~hello", true, "Starts with tilde (0x7E)"},
		{"ends with space", "hello ", true, "Ends with space (0x20)"},
		{"ends with tilde", "hello~", true, "Ends with tilde (0x7E)"},

		// Edge cases: just outside range
		{"starts with 0x1F", "\x1Fhello", false, "Starts with last control (0x1F)"},
		{"starts with 0x7F", "\x7Fhello", false, "Starts with DEL (0x7F)"},
		{"starts with 0x80", "\x80hello", false, "Starts with non-ASCII (0x80)"},
		{"ends with 0x1F", "hello\x1F", false, "Ends with last control (0x1F)"},
		{"ends with 0x7F", "hello\x7F", false, "Ends with DEL (0x7F)"},
		{"ends with 0x80", "hello\x80", false, "Ends with non-ASCII (0x80)"},

		// All printable ASCII range
		{"all printable range", " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~", true, "All printable ASCII (0x20-0x7E)"},

		// Control characters at various positions
		{"control at start", "\x00hello", false, "Control char at start"},
		{"control in middle", "he\x00llo", false, "Control char in middle"},
		{"control at end", "hello\x00", false, "Control char at end"},
		{"control in 8-byte chunk", "123456\x0078", false, "Control char in first 8-byte chunk"},
		{"control in second 8-byte chunk", "12345678\x0090", false, "Control char in second 8-byte chunk"},
		{"control in remainder", "123456789012345\x006", false, "Control char in remainder"},

		// DEL at various positions
		{"DEL at start", "\x7Fhello", false, "DEL at start"},
		{"DEL in middle", "he\x7Fllo", false, "DEL in middle"},
		{"DEL at end", "hello\x7F", false, "DEL at end"},
		{"DEL in 8-byte chunk", "123456\x7F78", false, "DEL in first 8-byte chunk"},
		{"DEL in second 8-byte chunk", "12345678\x7F90", false, "DEL in second 8-byte chunk"},
		{"DEL in remainder", "123456789012345\x7F6", false, "DEL in remainder"},

		// Non-ASCII at various positions
		{"non-ASCII at start", "\x80hello", false, "Non-ASCII at start"},
		{"non-ASCII in middle", "he\x80llo", false, "Non-ASCII in middle"},
		{"non-ASCII at end", "hello\x80", false, "Non-ASCII at end"},
		{"non-ASCII in 8-byte chunk", "123456\x8078", false, "Non-ASCII in first 8-byte chunk"},
		{"non-ASCII in second 8-byte chunk", "12345678\x8090", false, "Non-ASCII in second 8-byte chunk"},
		{"non-ASCII in remainder", "123456789012345\x806", false, "Non-ASCII in remainder"},

		// UTF-8 sequences (should fail since they contain non-ASCII bytes)
		{"UTF-8 2-byte start", "\xC2\xA0", false, "UTF-8 2-byte sequence (non-breaking space)"},
		{"UTF-8 3-byte start", "\xE2\x82\xAC", false, "UTF-8 3-byte sequence (euro sign)"},
		{"UTF-8 4-byte start", "\xF0\x9F\x98\x80", false, "UTF-8 4-byte sequence (emoji)"},
		{"ASCII with UTF-8", "hello\xC2\xA0world", false, "ASCII with UTF-8 sequence"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isPrintableASCIIString(tt.input)
			if got != tt.expected {
				t.Errorf("isPrintableASCII(%q) = %v, want %v (%s)",
					tt.input, got, tt.expected, tt.desc)
				// Show byte details for debugging
				if len(tt.input) > 0 {
					t.Logf("  String length: %d bytes", len(tt.input))
					t.Logf("  Bytes: %v", []byte(tt.input))
					for i, b := range []byte(tt.input) {
						isPrintable := b >= 0x20 && b <= 0x7E
						t.Logf("    [%d]: 0x%02X '%c' printable=%v", i, b, b, isPrintable)
					}
				}
			}
		})
	}
}

func TestIsPrintableASCIIBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected bool
		desc     string
	}{
		// Empty slice
		{"empty slice", []byte{}, true, "Empty slice is printable ASCII"},

		// Single byte cases
		{"single space", []byte{0x20}, true, "Space (0x20) is printable"},
		{"single exclamation", []byte{0x21}, true, "! (0x21) is printable"},
		{"single tilde", []byte{0x7E}, true, "~ (0x7E) is printable"},
		{"single null", []byte{0x00}, false, "NULL (0x00) is control char"},
		{"single tab", []byte{0x09}, false, "TAB (0x09) is control char"},
		{"single newline", []byte{0x0A}, false, "LF (0x0A) is control char"},
		{"single last control", []byte{0x1F}, false, "Last control (0x1F) is control char"},
		{"single DEL", []byte{0x7F}, false, "DEL (0x7F) is not printable"},
		{"single 0x80", []byte{0x80}, false, "0x80 is non-ASCII"},
		{"single 0xFF", []byte{0xFF}, false, "0xFF is non-ASCII"},

		// Two byte cases
		{"two printable", []byte("ab"), true, "Two printable ASCII"},
		{"two with control", []byte("a\x00"), false, "Contains control char"},
		{"two with DEL", []byte("a\x7F"), false, "Contains DEL"},
		{"two with non-ASCII", []byte("a\x80"), false, "Contains non-ASCII"},

		// Exactly 8 bytes (SWAR boundary)
		{"8 bytes all printable", []byte("12345678"), true, "8 bytes all printable"},
		{"8 bytes with space", []byte("hello wo"), true, "8 bytes with space"},
		{"8 bytes with control", []byte("hello\x00"), false, "8 bytes with control char"},
		{"8 bytes with DEL", []byte("hello\x7F"), false, "8 bytes with DEL"},
		{"8 bytes with non-ASCII", []byte("hello\x80"), false, "8 bytes with non-ASCII"},
		{"8 bytes all spaces", []byte("        "), true, "8 spaces"},
		{"8 bytes all tildes", []byte("~~~~~~~~"), true, "8 tildes"},
		{"8 bytes boundary low", []byte{0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20, 0x20}, true, "8 spaces (0x20)"},
		{"8 bytes boundary high", []byte{0x7E, 0x7E, 0x7E, 0x7E, 0x7E, 0x7E, 0x7E, 0x7E}, true, "8 tildes (0x7E)"},

		// Just before 8 bytes (7 bytes)
		{"7 bytes all printable", []byte("1234567"), true, "7 bytes all printable"},
		{"7 bytes with control", []byte("123456\x00"), false, "7 bytes with control char"},
		{"7 bytes with DEL", []byte("123456\x7F"), false, "7 bytes with DEL"},

		// Just after 8 bytes (9 bytes)
		{"9 bytes all printable", []byte("123456789"), true, "9 bytes all printable"},
		{"9 bytes with control", []byte("12345678\x00"), false, "9 bytes with control char"},
		{"9 bytes with DEL", []byte("12345678\x7F"), false, "9 bytes with DEL"},
		{"9 bytes with non-ASCII", []byte("12345678\x80"), false, "9 bytes with non-ASCII"},

		// Exactly 16 bytes (two SWAR chunks)
		{"16 bytes all printable", []byte("1234567890123456"), true, "16 bytes all printable"},
		{"16 bytes with control", []byte("123456789012345\x00"), false, "16 bytes with control char"},
		{"16 bytes with DEL", []byte("123456789012345\x7F"), false, "16 bytes with DEL"},
		{"16 bytes with non-ASCII", []byte("123456789012345\x80"), false, "16 bytes with non-ASCII"},

		// Just before 16 bytes (15 bytes)
		{"15 bytes all printable", []byte("123456789012345"), true, "15 bytes all printable"},
		{"15 bytes with control", []byte("12345678901234\x00"), false, "15 bytes with control char"},

		// Just after 16 bytes (17 bytes)
		{"17 bytes all printable", []byte("12345678901234567"), true, "17 bytes all printable"},
		{"17 bytes with control", []byte("1234567890123456\x00"), false, "17 bytes with control char"},

		// Exactly 24 bytes (three SWAR chunks)
		{"24 bytes all printable", []byte("123456789012345678901234"), true, "24 bytes all printable"},
		{"24 bytes with control", []byte("12345678901234567890123\x00"), false, "24 bytes with control char"},

		// Long slices
		{"long all printable", []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()"), true, "Long slice all printable"},
		{"long with control", []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!\x00@#$%^&*()"), false, "Long slice with control char"},
		{"long with DEL", []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!\x7F@#$%^&*()"), false, "Long slice with DEL"},
		{"long with non-ASCII", []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!\x80@#$%^&*()"), false, "Long slice with non-ASCII"},

		// Edge cases: first and last bytes in range
		{"starts with space", []byte(" hello"), true, "Starts with space (0x20)"},
		{"starts with tilde", []byte("~hello"), true, "Starts with tilde (0x7E)"},
		{"ends with space", []byte("hello "), true, "Ends with space (0x20)"},
		{"ends with tilde", []byte("hello~"), true, "Ends with tilde (0x7E)"},

		// Edge cases: just outside range
		{"starts with 0x1F", []byte("\x1Fhello"), false, "Starts with last control (0x1F)"},
		{"starts with 0x7F", []byte("\x7Fhello"), false, "Starts with DEL (0x7F)"},
		{"starts with 0x80", []byte("\x80hello"), false, "Starts with non-ASCII (0x80)"},
		{"ends with 0x1F", []byte("hello\x1F"), false, "Ends with last control (0x1F)"},
		{"ends with 0x7F", []byte("hello\x7F"), false, "Ends with DEL (0x7F)"},
		{"ends with 0x80", []byte("hello\x80"), false, "Ends with non-ASCII (0x80)"},

		// All printable ASCII range
		{"all printable range", []byte(" !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"), true, "All printable ASCII (0x20-0x7E)"},

		// Control characters at various positions
		{"control at start", []byte("\x00hello"), false, "Control char at start"},
		{"control in middle", []byte("he\x00llo"), false, "Control char in middle"},
		{"control at end", []byte("hello\x00"), false, "Control char at end"},
		{"control in 8-byte chunk", []byte("123456\x0078"), false, "Control char in first 8-byte chunk"},
		{"control in second 8-byte chunk", []byte("12345678\x0090"), false, "Control char in second 8-byte chunk"},
		{"control in remainder", []byte("123456789012345\x006"), false, "Control char in remainder"},

		// DEL at various positions
		{"DEL at start", []byte("\x7Fhello"), false, "DEL at start"},
		{"DEL in middle", []byte("he\x7Fllo"), false, "DEL in middle"},
		{"DEL at end", []byte("hello\x7F"), false, "DEL at end"},
		{"DEL in 8-byte chunk", []byte("123456\x7F78"), false, "DEL in first 8-byte chunk"},
		{"DEL in second 8-byte chunk", []byte("12345678\x7F90"), false, "DEL in second 8-byte chunk"},
		{"DEL in remainder", []byte("123456789012345\x7F6"), false, "DEL in remainder"},

		// Non-ASCII at various positions
		{"non-ASCII at start", []byte("\x80hello"), false, "Non-ASCII at start"},
		{"non-ASCII in middle", []byte("he\x80llo"), false, "Non-ASCII in middle"},
		{"non-ASCII at end", []byte("hello\x80"), false, "Non-ASCII at end"},
		{"non-ASCII in 8-byte chunk", []byte("123456\x8078"), false, "Non-ASCII in first 8-byte chunk"},
		{"non-ASCII in second 8-byte chunk", []byte("12345678\x8090"), false, "Non-ASCII in second 8-byte chunk"},
		{"non-ASCII in remainder", []byte("123456789012345\x806"), false, "Non-ASCII in remainder"},

		// UTF-8 sequences (should fail since they contain non-ASCII bytes)
		{"UTF-8 2-byte start", []byte{0xC2, 0xA0}, false, "UTF-8 2-byte sequence (non-breaking space)"},
		{"UTF-8 3-byte start", []byte{0xE2, 0x82, 0xAC}, false, "UTF-8 3-byte sequence (euro sign)"},
		{"UTF-8 4-byte start", []byte{0xF0, 0x9F, 0x98, 0x80}, false, "UTF-8 4-byte sequence (emoji)"},
		{"ASCII with UTF-8", []byte("hello\xC2\xA0world"), false, "ASCII with UTF-8 sequence"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isPrintableASCIIBytes(tt.input)
			if got != tt.expected {
				t.Errorf("isPrintableASCIIBytes(%v) = %v, want %v (%s)",
					tt.input, got, tt.expected, tt.desc)
				// Show byte details for debugging
				if len(tt.input) > 0 {
					t.Logf("  Slice length: %d bytes", len(tt.input))
					t.Logf("  Bytes: %v", tt.input)
					for i, b := range tt.input {
						isPrintable := b >= 0x20 && b <= 0x7E
						t.Logf("    [%d]: 0x%02X '%c' printable=%v", i, b, b, isPrintable)
					}
				}
			}
		})
	}
}

// TestIsPrintableASCIIBytesOptimization verifies that the isPrintableASCIIBytes optimization
// in Bytes() works correctly by comparing results with and without the optimization path.
func TestIsPrintableASCIIBytesOptimization(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		desc  string
	}{
		{"empty", []byte{}, "Empty slice"},
		{"single char", []byte("a"), "Single ASCII character"},
		{"short ASCII", []byte("hello"), "Short ASCII slice"},
		{"long ASCII", []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"), "Long ASCII slice"},
		{"with spaces", []byte("hello world"), "ASCII with spaces"},
		{"with punctuation", []byte("Hello, World!"), "ASCII with punctuation"},
		{"all printable range", []byte(" !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"), "All printable ASCII"},
		{"exactly 8 bytes", []byte("12345678"), "Exactly 8 bytes (SWAR boundary)"},
		{"exactly 16 bytes", []byte("1234567890123456"), "Exactly 16 bytes (two SWAR chunks)"},
		{"exactly 24 bytes", []byte("123456789012345678901234"), "Exactly 24 bytes (three SWAR chunks)"},
		{"7 bytes", []byte("1234567"), "7 bytes (just before SWAR)"},
		{"9 bytes", []byte("123456789"), "9 bytes (just after SWAR)"},
		{"15 bytes", []byte("123456789012345"), "15 bytes (just before two SWAR)"},
		{"17 bytes", []byte("12345678901234567"), "17 bytes (just after two SWAR)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// The optimization should return len(s) for printable ASCII slices
			// This is tested indirectly through Bytes() which uses isPrintableASCIIBytes
			if isPrintableASCIIBytes(tt.input) {
				width := Bytes(tt.input)
				expected := len(tt.input)
				if width != expected {
					t.Errorf("Bytes(%q) = %d, want %d (optimization should return length for printable ASCII)",
						tt.input, width, expected)
				}
			} else {
				t.Errorf("isPrintableASCIIBytes(%q) = false, but input should be printable ASCII (%s)",
					tt.input, tt.desc)
			}
		})
	}
}

// TestIsPrintableASCIIOptimization verifies that the isPrintableASCII optimization
// in String() works correctly by comparing results with and without the optimization path.
func TestIsPrintableASCIIOptimization(t *testing.T) {
	tests := []struct {
		name  string
		input string
		desc  string
	}{
		{"empty", "", "Empty string"},
		{"single char", "a", "Single ASCII character"},
		{"short ASCII", "hello", "Short ASCII string"},
		{"long ASCII", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", "Long ASCII string"},
		{"with spaces", "hello world", "ASCII with spaces"},
		{"with punctuation", "Hello, World!", "ASCII with punctuation"},
		{"all printable range", " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~", "All printable ASCII"},
		{"exactly 8 bytes", "12345678", "Exactly 8 bytes (SWAR boundary)"},
		{"exactly 16 bytes", "1234567890123456", "Exactly 16 bytes (two SWAR chunks)"},
		{"exactly 24 bytes", "123456789012345678901234", "Exactly 24 bytes (three SWAR chunks)"},
		{"7 bytes", "1234567", "7 bytes (just before SWAR)"},
		{"9 bytes", "123456789", "9 bytes (just after SWAR)"},
		{"15 bytes", "123456789012345", "15 bytes (just before two SWAR)"},
		{"17 bytes", "12345678901234567", "17 bytes (just after two SWAR)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// The optimization should return len(s) for printable ASCII strings
			// This is tested indirectly through String() which uses isPrintableASCII
			if isPrintableASCIIString(tt.input) {
				width := String(tt.input)
				expected := len(tt.input)
				if width != expected {
					t.Errorf("String(%q) = %d, want %d (optimization should return length for printable ASCII)",
						tt.input, width, expected)
				}
			} else {
				t.Errorf("isPrintableASCII(%q) = false, but input should be printable ASCII (%s)",
					tt.input, tt.desc)
			}
		})
	}
}
