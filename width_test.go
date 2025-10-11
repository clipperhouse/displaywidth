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
		{"red heart", '❤', Options{}, 1}, // Text presentation by default

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
	sample, err := testdata.Sample()
	if err != nil {
		t.Fatalf("Failed to load sample: %v", err)
	}

	eastAsianWidth := []bool{false, true}
	strictEmojiNeutral := []bool{false, true}

	t.Run("TestCases", func(t *testing.T) {
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
	})
	t.Run("Sample", func(t *testing.T) {
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
	})
	t.Run("SampleRunes", func(t *testing.T) {
		runes := []rune(string(sample))

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

				for _, r := range runes {
					// Skip Unicode tag characters (U+E0020-U+E007F) - we correctly
					// treat them as zero-width per Unicode Cf category, while
					// go-runewidth treats them as width 1. This is a known difference
					// where we are more correct.
					if r >= 0xE0020 && r <= 0xE007F {
						continue
					}

					w1 := options.Rune(r)
					w2 := condition.RuneWidth(r)
					if w1 != w2 {
						t.Errorf("Sample mismatch for %q (eastAsianWidth=%v, strictEmojiNeutral=%v):\n"+
							"  displaywidth result: %d\n"+
							"  go-runewidth result: %d\n"+
							"  Difference: %d",
							r, e, s, w1, w2, w1-w2)
					}
				}
			}
		}
	})
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
