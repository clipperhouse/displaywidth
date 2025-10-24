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
		{"emoji", "😀", Options{}, 2},          // Default emoji width
		{"checkered flag", "🏁", Options{}, 2}, // U+1F3C1 chequered flag is an emoji, width 2

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

		// Test edge cases with options disabled
		{"ambiguous EAW disabled", '★', Options{EastAsianWidth: false}, 1},

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
		{
			name:         "Watch (EP=Yes, EmojiPres=Yes)",
			input:        "\u231A",
			wantDefault:  2,
			wantWithVS16: 2,
			wantWithVS15: 1,
			description:  "⌚ U+231A has default emoji presentation",
		},
		{
			name:         "Hourglass (EP=Yes, EmojiPres=Yes)",
			input:        "\u231B",
			wantDefault:  2,
			wantWithVS16: 2,
			wantWithVS15: 1,
			description:  "⌛ U+231B has default emoji presentation",
		},
		{
			name:         "Fast-forward (EP=Yes, EmojiPres=Yes)",
			input:        "\u23E9",
			wantDefault:  2,
			wantWithVS16: 2,
			wantWithVS15: 1,
			description:  "⏩ U+23E9 has default emoji presentation",
		},
		{
			name:         "Alarm Clock (EP=Yes, EmojiPres=Yes)",
			input:        "\u23F0",
			wantDefault:  2,
			wantWithVS16: 2,
			wantWithVS15: 1,
			description:  "⏰ U+23F0 has default emoji presentation",
		},
		{
			name:         "Soccer Ball (EP=Yes, EmojiPres=Yes)",
			input:        "\u26BD",
			wantDefault:  2,
			wantWithVS16: 2,
			wantWithVS15: 1,
			description:  "⚽ U+26BD has default emoji presentation",
		},
		{
			name:         "Anchor (EP=Yes, EmojiPres=Yes)",
			input:        "\u2693",
			wantDefault:  2,
			wantWithVS16: 2,
			wantWithVS15: 1,
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

			// Test with VS15 (U+FE0E) for text presentation
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
			name:  "Flag sequence 🇺🇸",
			input: "\U0001F1FA\U0001F1F8",
			want:  2,
			desc:  "US flag (RI pair)",
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

	t.Run("C2: VS15 forces text presentation", func(t *testing.T) {
		// VS15 (U+FE0E) should force text presentation (width 1) even for emoji-presentation characters
		emojiWithVS15 := []string{
			"\u231A\uFE0E", // ⌚︎ watch with VS15
			"\u26BD\uFE0E", // ⚽︎ soccer ball with VS15
			"\u2693\uFE0E", // ⚓︎ anchor with VS15
		}

		for _, s := range emojiWithVS15 {
			got := String(s)
			if got != 1 {
				t.Errorf("String(%q) with VS15 = %d, want 1 (VS15 should force text presentation)", s, got)
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
}
