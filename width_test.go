package displaywidth

import (
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
