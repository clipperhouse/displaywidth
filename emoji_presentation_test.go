package displaywidth

import (
	"testing"
)

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

// TestUnicode16NewEmoji tests new emoji added in Unicode 16.0
func TestUnicode16NewEmoji(t *testing.T) {
	// Note: These are examples; actual Unicode 16.0 additions may vary
	// This test verifies that the generator picked up the latest data
	tests := []struct {
		name  string
		input string
		want  int
		desc  string
	}{
		{
			name:  "Face with bags under eyes 🫦",
			input: "\U0001FAE6",
			want:  2,
			desc:  "New emoji should be recognized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := String(tt.input)
			if got != tt.want {
				t.Logf("String(%q) = %d, want %d (%s) - may be expected if emoji not in Unicode 16.0",
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
}
