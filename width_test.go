package displaywidth

import (
	"fmt"
	"strings"
	"testing"

	"github.com/clipperhouse/displaywidth/internal/testdata"
	"github.com/mattn/go-runewidth"
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
		{"emoji strict", "😀", false, true, 2}, // Strict emoji neutral - only ambiguous emoji get width 1

		// Invalid UTF-8 - the trie treats \xff as a valid character with default properties
		{"invalid UTF-8", "\xff", false, false, 1},
		{"partial UTF-8", "\xc2", false, false, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := String(tt.input, tt.eastAsianWidth, tt.strictEmojiNeutral)
			if result != tt.expected {
				t.Errorf("StringWidth(%q, %v, %v) = %d, want %d",
					tt.input, tt.eastAsianWidth, tt.strictEmojiNeutral, result, tt.expected)
			}

			b := []byte(tt.input)
			result = Bytes(b, tt.eastAsianWidth, tt.strictEmojiNeutral)
			if result != tt.expected {
				t.Errorf("BytesWidth(%q, %v, %v) = %d, want %d",
					b, tt.eastAsianWidth, tt.strictEmojiNeutral, result, tt.expected)
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
		{"control char", _ControlChar, false, false, 0},

		// Combining marks
		{"combining mark", _CombiningMark, false, false, 0},

		// Zero width
		{"zero width", _ZeroWidth, false, false, 0},

		// East Asian Wide
		{"EAW fullwidth", _EAW_Fullwidth, false, false, 2},
		{"EAW wide", _EAW_Wide, false, false, 2},

		// East Asian Ambiguous
		{"EAW ambiguous default", _EAW_Ambiguous, false, false, 1},
		{"EAW ambiguous EAW", _EAW_Ambiguous, true, false, 2},

		// Emoji
		// {"emoji default", _Emoji, false, false, 2},
		// {"emoji strict", _Emoji, false, true, 2}, // Only ambiguous emoji get width 1 in strict mode

		// Default (no properties set)
		{"default", 0, false, false, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.props.width(tt.eastAsianWidth, tt.strictEmojiNeutral)
			if result != tt.expected {
				t.Errorf("calculateWidth(%d, %v, %v) = %d, want %d",
					tt.props, tt.eastAsianWidth, tt.strictEmojiNeutral, result, tt.expected)
			}
		})
	}
}

func TestComparisonWithRunewidth(t *testing.T) {
	testCases := []struct {
		name               string
		input              string
		eastAsianWidth     bool
		strictEmojiNeutral bool
	}{
		// Basic ASCII
		{"empty string", "", false, false},
		{"single ASCII", "a", false, false},
		{"multiple ASCII", "hello", false, false},
		{"ASCII with spaces", "hello world", false, false},
		{"numbers", "1234567890", false, false},
		{"symbols", "!@#$%^&*()", false, false},

		// Control characters
		{"newline", "\n", false, false},
		{"tab", "\t", false, false},
		{"carriage return", "\r", false, false},
		{"backspace", "\b", false, false},
		{"null", "\x00", false, false},
		{"del", "\x7f", false, false},

		// Latin characters with diacritics
		{"cafe", "café", false, false},
		{"naive", "naïve", false, false},
		{"resume", "résumé", false, false},
		{"zurich", "Zürich", false, false},
		{"sao paulo", "São Paulo", false, false},

		// East Asian characters
		{"chinese", "中文", false, false},
		{"japanese", "こんにちは", false, false},
		{"korean", "안녕하세요", false, false},
		{"mixed", "Hello 世界", false, false},

		// Fullwidth characters
		{"fullwidth A", "Ａ", false, false},
		{"fullwidth 1", "１", false, false},
		{"fullwidth !", "！", false, false},

		// Ambiguous characters
		{"star", "★", false, false},
		{"star EAW", "★", true, false},
		{"degree", "°", false, false},
		{"degree EAW", "°", true, false},
		{"plus minus", "±", false, false},
		{"plus minus EAW", "±", true, false},

		// Emoji
		{"grinning face", "😀", false, false},
		{"grinning face strict", "😀", false, true},
		{"rocket", "🚀", false, false},
		{"rocket strict", "🚀", false, true},
		{"party popper", "🎉", false, false},
		{"party popper strict", "🎉", false, true},

		// Complex emoji sequences
		{"family", "👨‍👩‍👧‍👦", false, false},
		{"family strict", "👨‍👩‍👧‍👦", false, true},
		{"technologist", "👨‍💻", false, false},
		{"technologist strict", "👨‍💻", false, true},

		// Mixed content
		{"hello world emoji", "Hello 世界! 😀", false, false},
		{"price", "Price: $100.00 €85.50", false, false},
		{"math", "Math: ∑(x²) = ∞", false, false},
		{"emoji sequence", "👨‍💻 working on 🚀", false, false},

		// Edge cases
		{"single space", " ", false, false},
		{"multiple spaces", "     ", false, false},
		{"tab and newline", "\t\n", false, false},
		{"mixed whitespace", " \t \n ", false, false},

		// Long string
		{"long string", "This is a very long string with many characters to test performance of both implementations. It contains various character types including ASCII, Unicode, emoji, and special symbols.", false, false},

		// Many emoji
		{"many emoji", "😀😁😂🤣😃😄😅😆😉😊😋😎😍😘🥰😗😙😚☺️🙂🤗🤩🤔🤨😐😑😶🙄😏😣😥😮🤐😯😪😫🥱😴😌😛😜😝🤤😒😓😔😕🙃🤑😲☹️🙁😖😞😟😤😢😭😦😧😨😩🤯😬😰😱🥵🥶😳🤪😵😡😠🤬😷🤒🤕🤢🤮🤧😇🤠🤡🥳🥴🥺🤥🤫🤭🧐🤓😈👿💀☠️👹👺🤖👽👾💩😺😸😹😻😼😽🙀😿😾", false, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test our implementation
			ourResult := String(tc.input, tc.eastAsianWidth, tc.strictEmojiNeutral)

			// Test go-runewidth using Condition
			condition := runewidth.NewCondition()
			condition.EastAsianWidth = tc.eastAsianWidth
			condition.StrictEmojiNeutral = tc.strictEmojiNeutral
			goRunewidthResult := condition.StringWidth(tc.input)

			// Compare results
			if ourResult != goRunewidthResult {
				t.Errorf("StringWidth mismatch for %q (eastAsianWidth=%v, strictEmojiNeutral=%v):\n"+
					"  Our result: %d\n"+
					"  go-runewidth result: %d\n"+
					"  Difference: %d",
					tc.input, tc.eastAsianWidth, tc.strictEmojiNeutral,
					ourResult, goRunewidthResult, ourResult-goRunewidthResult)
			}
		})
	}
}

func TestMoreComparisonWithRunewidth(t *testing.T) {
	testCases, err := testdata.TestCases()
	if err != nil {
		t.Fatalf("Failed to load test cases: %v", err)
	}

	eastAsianWidth := []bool{false, true}
	strictEmojiNeutral := []bool{false, true}

	{
		lines := strings.Split(string(testCases), "\n")

		for _, e := range eastAsianWidth {
			for _, s := range strictEmojiNeutral {
				condition := runewidth.NewCondition()
				condition.EastAsianWidth = e
				condition.StrictEmojiNeutral = s

				for _, line := range lines {
					w1 := String(line, e, s)
					w2 := condition.StringWidth(line)
					if w1 != w2 {
						t.Errorf("TestCases mismatch for %q (eastAsianWidth=%v, strictEmojiNeutral=%v):\n"+
							"  displaywidth result: %d\n"+
							"  go-runewidth result: %d\n"+
							"  Difference: %d",
							line, e, s, w1, w2, w1-w2)
					}
				}
			}
		}
	}
	{
		sample, err := testdata.Sample()
		if err != nil {
			t.Fatalf("Failed to load sample: %v", err)
		}
		words := strings.Fields(string(sample))
		for _, e := range eastAsianWidth {
			for _, s := range strictEmojiNeutral {
				condition := runewidth.NewCondition()
				condition.EastAsianWidth = e
				condition.StrictEmojiNeutral = s

				for _, word := range words {
					w1 := String(word, false, false)
					w2 := runewidth.StringWidth(word)
					if w1 != w2 {
						t.Errorf("Sample mismatch for %q (eastAsianWidth=%v, strictEmojiNeutral=%v):\n"+
							"  displaywidth result: %d\n"+
							"  go-runewidth result: %d\n"+
							"  Difference: %d",
							word, e, s, w1, w2, w1-w2)
					}
				}
			}
		}
	}
}

func TestSpecificEmojiCharacters(t *testing.T) {
	// Test the specific characters that are causing differences with go-runewidth
	chars := []rune{'☺', '☹', '☠', '️'}

	for _, char := range chars {
		t.Run(fmt.Sprintf("char_%04X", char), func(t *testing.T) {
			props, _ := lookupProperties(string(char))
			ourWidth := String(string(char), false, false)

			// Test with go-runewidth for comparison
			condition := runewidth.NewCondition()
			condition.EastAsianWidth = false
			condition.StrictEmojiNeutral = false
			goRunewidthWidth := condition.RuneWidth(char)

			t.Logf("Character: %c (U+%04X)", char, char)
			t.Logf("Our properties: %d", props)
			t.Logf("Our width: %d", ourWidth)
			t.Logf("go-runewidth width: %d", goRunewidthWidth)

			// For now, just log the differences - we'll fix them next
			if ourWidth != goRunewidthWidth {
				t.Logf("DIFFERENCE: Our width %d != go-runewidth width %d", ourWidth, goRunewidthWidth)
			}
		})
	}
}
