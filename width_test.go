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
		{"emoji", "😀", Options{}, 2},                                // Default emoji width
		{"emoji strict", "😀", Options{StrictEmojiNeutral: true}, 2}, // Strict emoji neutral - only ambiguous emoji get width 1

		// Invalid UTF-8 - the trie treats \xff as a valid character with default properties
		{"invalid UTF-8", "\xff", Options{}, 1},
		{"partial UTF-8", "\xc2", Options{}, 1},
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

func TestCalculateWidth(t *testing.T) {
	tests := []struct {
		name     string
		props    property
		options  Options
		expected int
	}{
		// Control characters
		{"control char", _ControlChar, Options{}, 0},

		// Combining marks
		{"combining mark", _CombiningMark, Options{}, 0},

		// Zero width
		{"zero width", _ZeroWidth, Options{}, 0},

		// East Asian Wide
		{"EAW fullwidth", _EAW_Fullwidth, Options{}, 2},
		{"EAW wide", _EAW_Wide, Options{}, 2},

		// East Asian Ambiguous
		{"EAW ambiguous default", _EAW_Ambiguous, Options{}, 1},
		{"EAW ambiguous EAW", _EAW_Ambiguous, Options{EastAsianWidth: true}, 2},

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

func TestComparisonWithRunewidth(t *testing.T) {
	testCases := []struct {
		name    string
		input   string
		options Options
	}{
		// Basic ASCII
		{"empty string", "", Options{}},
		{"single ASCII", "a", Options{}},
		{"multiple ASCII", "hello", Options{}},
		{"ASCII with spaces", "hello world", Options{}},
		{"numbers", "1234567890", Options{}},
		{"symbols", "!@#$%^&*()", Options{}},

		// Control characters
		{"newline", "\n", Options{}},
		{"tab", "\t", Options{}},
		{"carriage return", "\r", Options{}},
		{"backspace", "\b", Options{}},
		{"null", "\x00", Options{}},
		{"del", "\x7f", Options{}},

		// Latin characters with diacritics
		{"cafe", "café", Options{}},
		{"naive", "naïve", Options{}},
		{"resume", "résumé", Options{}},
		{"zurich", "Zürich", Options{}},
		{"sao paulo", "São Paulo", Options{}},

		// East Asian characters
		{"chinese", "中文", Options{}},
		{"japanese", "こんにちは", Options{}},
		{"korean", "안녕하세요", Options{}},
		{"mixed", "Hello 世界", Options{}},

		// Fullwidth characters
		{"fullwidth A", "Ａ", Options{}},
		{"fullwidth 1", "１", Options{}},
		{"fullwidth !", "！", Options{}},

		// Ambiguous characters
		{"star", "★", Options{}},
		{"star EAW", "★", Options{EastAsianWidth: true}},
		{"degree", "°", Options{}},
		{"degree EAW", "°", Options{EastAsianWidth: true}},
		{"plus minus", "±", Options{}},
		{"plus minus EAW", "±", Options{EastAsianWidth: true}},

		// Emoji
		{"grinning face", "😀", Options{}},
		{"grinning face strict", "😀", Options{StrictEmojiNeutral: true}},
		{"rocket", "🚀", Options{}},
		{"rocket strict", "🚀", Options{StrictEmojiNeutral: true}},
		{"party popper", "🎉", Options{}},
		{"party popper strict", "🎉", Options{StrictEmojiNeutral: true}},

		// Complex emoji sequences
		{"family", "👨‍👩‍👧‍👦", Options{}},
		{"family strict", "👨‍👩‍👧‍👦", Options{StrictEmojiNeutral: true}},
		{"technologist", "👨‍💻", Options{}},
		{"technologist strict", "👨‍💻", Options{StrictEmojiNeutral: true}},

		// Mixed content
		{"hello world emoji", "Hello 世界! 😀", Options{}},
		{"price", "Price: $100.00 €85.50", Options{}},
		{"math", "Math: ∑(x²) = ∞", Options{}},
		{"emoji sequence", "👨‍💻 working on 🚀", Options{}},

		// Edge cases
		{"single space", " ", Options{}},
		{"multiple spaces", "     ", Options{}},
		{"tab and newline", "\t\n", Options{}},
		{"mixed whitespace", " \t \n ", Options{}},

		// Long string
		{"long string", "This is a very long string with many characters to test performance of both implementations. It contains various character types including ASCII, Unicode, emoji, and special symbols.", Options{}},

		// Many emoji
		{"many emoji", "😀😁😂🤣😃😄😅😆😉😊😋😎😍😘🥰😗😙😚☺️🙂🤗🤩🤔🤨😐😑😶🙄😏😣😥😮🤐😯😪😫🥱😴😌😛😜😝🤤😒😓😔😕🙃🤑😲☹️🙁😖😞😟😤😢😭😦😧😨😩🤯😬😰😱🥵🥶😳🤪😵😡😠🤬😷🤒🤕🤢🤮🤧😇🤠🤡🥳🥴🥺🤥🤫🤭🧐🤓😈👿💀☠️👹👺🤖👽👾💩😺😸😹😻😼😽🙀😿😾", Options{}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test our implementation
			ourResult := tc.options.String(tc.input)

			// Test go-runewidth using Condition
			condition := runewidth.NewCondition()
			condition.EastAsianWidth = tc.options.EastAsianWidth
			condition.StrictEmojiNeutral = tc.options.StrictEmojiNeutral
			goRunewidthResult := condition.StringWidth(tc.input)

			// Compare results
			if ourResult != goRunewidthResult {
				t.Errorf("StringWidth mismatch for %q (eastAsianWidth=%v, strictEmojiNeutral=%v):\n"+
					"  Our result: %d\n"+
					"  go-runewidth result: %d\n"+
					"  Difference: %d",
					tc.input, tc.options.EastAsianWidth, tc.options.StrictEmojiNeutral,
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
				options := Options{
					EastAsianWidth:     e,
					StrictEmojiNeutral: s,
				}
				condition := &runewidth.Condition{
					EastAsianWidth:     options.EastAsianWidth,
					StrictEmojiNeutral: options.StrictEmojiNeutral,
				}

				for _, line := range lines {
					w1 := options.String(line)
					w2 := condition.StringWidth(line)
					if w1 != w2 {
						t.Errorf("TestCases mismatch for %q (eastAsianWidth=%v, strictEmojiNeutral=%v):\n"+
							"  displaywidth result: %d\n"+
							"  go-runewidth result: %d\n"+
							"  Difference: %d",
							line, options.EastAsianWidth, options.StrictEmojiNeutral, w1, w2, w1-w2)
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
				options := Options{
					EastAsianWidth:     e,
					StrictEmojiNeutral: s,
				}
				condition := &runewidth.Condition{
					EastAsianWidth:     options.EastAsianWidth,
					StrictEmojiNeutral: options.StrictEmojiNeutral,
				}

				for _, word := range words {
					w1 := options.String(word)
					w2 := condition.StringWidth(word)
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
			w1 := String(string(char))
			w2 := runewidth.RuneWidth(char)

			t.Logf("Character: %c (U+%04X)", char, char)
			t.Logf("Our properties: %d", props)
			t.Logf("Our width: %d", w1)
			t.Logf("go-runewidth width: %d", w2)

			// For now, just log the differences - we'll fix them next
			if w1 != w2 {
				t.Logf("DIFFERENCE: Our width %d != go-runewidth width %d", w1, w2)
			}
		})
	}
}
