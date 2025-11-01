package comparison

import (
	"testing"

	"github.com/clipperhouse/displaywidth"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/uniseg"
)

func TestLibraryBehaviorComparison(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected map[string]int // library -> expected width
	}{
		// Basic ASCII
		{
			name:  "ASCII text",
			input: "Hello World",
			expected: map[string]int{
				"displaywidth_default":   11,
				"displaywidth_options{}": 11,
				"go-runewidth_default":   11,
				"uniseg_default":         11,
			},
		},

		// East Asian characters
		{
			name:  "CJK characters",
			input: "中文",
			expected: map[string]int{
				"displaywidth_default":   4,
				"displaywidth_options{}": 4,
				"go-runewidth_default":   4,
				"uniseg_default":         4,
			},
		},

		// Ambiguous characters (width depends on EastAsianWidth)
		{
			name:  "Ambiguous characters",
			input: "★°±",
			expected: map[string]int{
				"displaywidth_default":   3,
				"displaywidth_options{}": 3,
				"displaywidth_EAW":       6,
				"go-runewidth_default":   3,
				"go-runewidth_EAW":       6,
				"uniseg_default":         3,
				"uniseg_EAW":             5, // uniseg behavior is different
			},
		},

		// Emoji
		{
			name:  "Basic emoji",
			input: "😀🚀🎉",
			expected: map[string]int{
				"displaywidth_default":   6,
				"displaywidth_options{}": 6,
				"go-runewidth_default":   6,
				"uniseg_default":         6,
			},
		},

		// Unicode 16.0 new emojis
		{
			name:  "Unicode 16.0 emojis",
			input: "🫩🫆🪾🫜🪉🪏🫟", // Face with Bags Under Eyes, Fingerprint, Leafless Tree, Root Vegetable, Harp, Shovel, Splatter
			expected: map[string]int{
				"displaywidth_default":   14, // 2 per emoji (properly handles Unicode 16.0)
				"displaywidth_options{}": 14,
				"go-runewidth_default":   7, // go-runewidth may not fully support Unicode 16.0 yet (treats as width 1)
				"uniseg_default":         7,  // uniseg may not fully support Unicode 16.0 yet (treats as width 1)
			},
		},

		// Regional Indicator Pairs (flags) - the key difference
		{
			name:  "Flags",
			input: "🇺🇸🇯🇵🇬🇧",
			expected: map[string]int{
				"displaywidth_default":      6, // flags are always width 2 (modern standard)
				"displaywidth_options{}":    6, // same as default
				"go-runewidth_default":      3, // go-runewidth treats flags as width 1
				"go-runewidth_strict_false": 3,
				"go-runewidth_strict_true":  3, // go-runewidth always returns 1 for flags
				"uniseg_default":            6, // uniseg treats flags as width 2
			},
		},

		// Variation selectors
		{
			name:  "Variation selectors",
			input: "☺️⌛︎❤️",
			expected: map[string]int{
				"displaywidth_default":   5,
				"displaywidth_options{}": 5,
				"go-runewidth_default":   4,
				"uniseg_default":         5,
			},
		},

		// Keycap sequences
		{
			name:  "Keycap sequences",
			input: "1️⃣#️⃣",
			expected: map[string]int{
				"displaywidth_default":   4,
				"displaywidth_options{}": 4,
				"go-runewidth_default":   2,
				"uniseg_default":         2,
			},
		},

		// Mixed content
		{
			name:  "Mixed content",
			input: "Hello 世界! 😀🇺🇸",
			expected: map[string]int{
				"displaywidth_default":   16, // 6 + 4 + 2 + 2 + 2
				"displaywidth_options{}": 16, // same as default
				"go-runewidth_default":   15, // 6 + 4 + 2 + 2 + 1 (flags are width 1)
				"uniseg_default":         16, // 6 + 4 + 2 + 2 + 2
			},
		},

		// Control characters
		{
			name:  "Control characters",
			input: "hello\nworld\t",
			expected: map[string]int{
				"displaywidth_default":   10, // newline and tab are width 0
				"displaywidth_options{}": 10,
				"go-runewidth_default":   10,
				"uniseg_default":         10,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test displaywidth with default options
			displaywidthDefault := displaywidth.String(tc.input)
			if expected, ok := tc.expected["displaywidth_default"]; ok {
				if displaywidthDefault != expected {
					t.Errorf("displaywidth.String() = %d, want %d", displaywidthDefault, expected)
				}
			}

			// Test displaywidth with zero-value options (should behave same as default)
			displaywidthZero := displaywidth.String(tc.input)
			if expected, ok := tc.expected["displaywidth_options{}"]; ok {
				if displaywidthZero != expected {
					t.Errorf("displaywidth.String() = %d, want %d", displaywidthZero, expected)
				}
			}

			// Test displaywidth with EastAsianWidth=true
			displaywidthEAW := displaywidth.Options{EastAsianWidth: true}.String(tc.input)
			if expected, ok := tc.expected["displaywidth_EAW"]; ok {
				if displaywidthEAW != expected {
					t.Errorf("displaywidth.Options{EastAsianWidth: true}.String() = %d, want %d", displaywidthEAW, expected)
				}
			}

			// Test go-runewidth default
			goRunewidthDefault := runewidth.StringWidth(tc.input)
			if expected, ok := tc.expected["go-runewidth_default"]; ok {
				if goRunewidthDefault != expected {
					t.Errorf("runewidth.StringWidth() = %d, want %d", goRunewidthDefault, expected)
				}
			}

			// Test go-runewidth with StrictEmojiNeutral=false
			goRunewidthStrictFalse := (&runewidth.Condition{StrictEmojiNeutral: false}).StringWidth(tc.input)
			if expected, ok := tc.expected["go-runewidth_strict_false"]; ok {
				if goRunewidthStrictFalse != expected {
					t.Errorf("runewidth.Condition{StrictEmojiNeutral: false}.StringWidth() = %d, want %d", goRunewidthStrictFalse, expected)
				}
			}

			// Test go-runewidth with StrictEmojiNeutral=true
			goRunewidthStrictTrue := (&runewidth.Condition{StrictEmojiNeutral: true}).StringWidth(tc.input)
			if expected, ok := tc.expected["go-runewidth_strict_true"]; ok {
				if goRunewidthStrictTrue != expected {
					t.Errorf("runewidth.Condition{StrictEmojiNeutral: true}.StringWidth() = %d, want %d", goRunewidthStrictTrue, expected)
				}
			}

			// Test go-runewidth with EastAsianWidth=true
			goRunewidthEAW := (&runewidth.Condition{EastAsianWidth: true}).StringWidth(tc.input)
			if expected, ok := tc.expected["go-runewidth_EAW"]; ok {
				if goRunewidthEAW != expected {
					t.Errorf("runewidth.Condition{EastAsianWidth: true}.StringWidth() = %d, want %d", goRunewidthEAW, expected)
				}
			}

			// Test uniseg default
			unisegDefault := uniseg.StringWidth(tc.input)
			if expected, ok := tc.expected["uniseg_default"]; ok {
				if unisegDefault != expected {
					t.Errorf("uniseg.StringWidth() = %d, want %d", unisegDefault, expected)
				}
			}

			// Test uniseg with EastAsianWidth=true
			originalEAW := uniseg.EastAsianAmbiguousWidth
			uniseg.EastAsianAmbiguousWidth = 2
			unisegEAW := uniseg.StringWidth(tc.input)
			uniseg.EastAsianAmbiguousWidth = originalEAW
			if expected, ok := tc.expected["uniseg_EAW"]; ok {
				if unisegEAW != expected {
					t.Errorf("uniseg.StringWidth() with EastAsianAmbiguousWidth=2 = %d, want %d", unisegEAW, expected)
				}
			}
		})
	}
}

func TestFlagBehaviorDetailed(t *testing.T) {
	flags := []string{"🇺🇸", "🇯🇵", "🇬🇧", "🇫🇷", "🇩🇪"}

	t.Log("Flag behavior comparison:")
	t.Log("Library | displaywidth | go-runewidth (default) | go-runewidth (strict=false) | go-runewidth (strict=true) | uniseg")
	t.Log("--------|--------------|------------------------|----------------------------|----------------------------|-------")

	for _, flag := range flags {
		// displaywidth (always width 2, no StrictEmojiNeutral option)
		displaywidthDefault := displaywidth.String(flag)

		// go-runewidth
		goRunewidthDefault := runewidth.StringWidth(flag)
		goRunewidthStrictFalse := (&runewidth.Condition{StrictEmojiNeutral: false}).StringWidth(flag)
		goRunewidthStrictTrue := (&runewidth.Condition{StrictEmojiNeutral: true}).StringWidth(flag)

		// uniseg
		unisegDefault := uniseg.StringWidth(flag)

		t.Logf("%s | %d | %d | %d | %d | %d", flag, displaywidthDefault, goRunewidthDefault, goRunewidthStrictFalse, goRunewidthStrictTrue, unisegDefault)
	}
}
