package comparison

import (
	"strings"
	"testing"

	"github.com/clipperhouse/displaywidth"
	"github.com/clipperhouse/displaywidth/testdata"
	"github.com/mattn/go-runewidth"
	"github.com/rivo/uniseg"
)

// TestCase represents a test case from the test_cases.txt file
type TestCase struct {
	Name  string
	Input string
}

// loadTestCases reads and parses test cases from test_cases.txt
func loadTestCases() ([]TestCase, int64, error) {
	file, err := testdata.TestCases()
	if err != nil {
		return nil, 0, err
	}

	var testCases []TestCase
	lines := strings.Split(string(file), "\n")

	for _, line := range lines {
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Handle special cases with colons (like "newline:", "tab:", etc.)
		if strings.HasSuffix(line, ":") {
			name := strings.TrimSuffix(line, ":")
			var input string

			switch name {
			case "newline":
				input = "\n"
			case "tab":
				input = "\t"
			case "carriage return":
				input = "\r"
			case "backspace":
				input = "\b"
			case "null":
				input = "\x00"
			case "del":
				input = "\x7f"
			case "Zero Width Space":
				input = "\u200b"
			case "Zero Width Joiner":
				input = "\u200d"
			case "Zero Width Non-Joiner":
				input = "\u200c"
			case "Empty string":
				input = ""
			case "Single space":
				input = " "
			case "Multiple spaces":
				input = "     "
			case "Tab and newline":
				input = "\t\n"
			case "Mixed whitespace":
				input = " \t \n "
			default:
				// For other cases, use the name as input
				input = name
			}

			testCases = append(testCases, TestCase{
				Name:  name,
				Input: input,
			})
		} else {
			// Regular test case - use the line as both name and input
			testCases = append(testCases, TestCase{
				Name:  line,
				Input: line,
			})
		}
	}

	totalBytes := 0
	for _, tc := range testCases {
		totalBytes += len(tc.Input)
	}

	return testCases, int64(totalBytes), nil
}

// BenchmarkStringDefault benchmarks our displaywidth package
func BenchmarkStringDefault(b *testing.B) {
	b.Run("clipperhouse/displaywidth", func(b *testing.B) {
		testCases, n, err := loadTestCases()
		if err != nil {
			b.Fatalf("Failed to load test cases: %v", err)
		}
		b.SetBytes(n)
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for _, tc := range testCases {
				// Test with default settings (eastAsianWidth=false, strictEmojiNeutral=false)
				_ = displaywidth.String(tc.Input)
			}
		}
	})

	b.Run("mattn/go-runewidth", func(b *testing.B) {
		testCases, n, err := loadTestCases()
		if err != nil {
			b.Fatalf("Failed to load test cases: %v", err)
		}
		b.SetBytes(n)
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for _, tc := range testCases {
				_ = runewidth.StringWidth(tc.Input)
			}
		}
	})

	b.Run("rivo/uniseg", func(b *testing.B) {
		testCases, n, err := loadTestCases()
		if err != nil {
			b.Fatalf("Failed to load test cases: %v", err)
		}
		b.SetBytes(n)
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for _, tc := range testCases {
				_ = uniseg.StringWidth(tc.Input)
			}
		}
	})
}

func BenchmarkString_EAW(b *testing.B) {
	options := displaywidth.Options{
		EastAsianWidth:     true,
		StrictEmojiNeutral: false,
	}

	condition := &runewidth.Condition{
		EastAsianWidth:     options.EastAsianWidth,
		StrictEmojiNeutral: options.StrictEmojiNeutral,
	}

	// Save original value and restore after benchmark
	originalEAAWidth := uniseg.EastAsianAmbiguousWidth
	defer func() {
		uniseg.EastAsianAmbiguousWidth = originalEAAWidth
	}()

	b.Run("clipperhouse/displaywidth", func(b *testing.B) {
		testCases, n, err := loadTestCases()
		if err != nil {
			b.Fatalf("Failed to load test cases: %v", err)
		}
		b.SetBytes(n)
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for _, tc := range testCases {
				// Test with East Asian Width enabled
				_ = options.String(tc.Input)
			}
		}
	})

	b.Run("mattn/go-runewidth", func(b *testing.B) {
		testCases, n, err := loadTestCases()
		if err != nil {
			b.Fatalf("Failed to load test cases: %v", err)
		}
		b.SetBytes(n)
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for _, tc := range testCases {
				_ = condition.StringWidth(tc.Input)
			}
		}
	})

	b.Run("rivo/uniseg", func(b *testing.B) {
		// Set EastAsianAmbiguousWidth to 2 to match the other libraries
		uniseg.EastAsianAmbiguousWidth = 2
		defer func() {
			uniseg.EastAsianAmbiguousWidth = 1
		}()

		testCases, n, err := loadTestCases()
		if err != nil {
			b.Fatalf("Failed to load test cases: %v", err)
		}
		b.SetBytes(n)
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for _, tc := range testCases {
				_ = uniseg.StringWidth(tc.Input)
			}
		}
	})
}

// BenchmarkString_StrictEmoji benchmarks our package with strict emoji neutral
func BenchmarkString_StrictEmoji(b *testing.B) {
	options := displaywidth.Options{
		EastAsianWidth:     false,
		StrictEmojiNeutral: true,
	}

	condition := &runewidth.Condition{
		EastAsianWidth:     options.EastAsianWidth,
		StrictEmojiNeutral: options.StrictEmojiNeutral,
	}

	b.Run("clipperhouse/displaywidth", func(b *testing.B) {
		testCases, n, err := loadTestCases()
		if err != nil {
			b.Fatalf("Failed to load test cases: %v", err)
		}
		b.SetBytes(n)
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for _, tc := range testCases {
				// Test with strict emoji neutral enabled
				_ = options.String(tc.Input)
			}
		}
	})

	b.Run("mattn/go-runewidth", func(b *testing.B) {
		testCases, n, err := loadTestCases()
		if err != nil {
			b.Fatalf("Failed to load test cases: %v", err)
		}
		b.SetBytes(n)
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for _, tc := range testCases {
				_ = condition.StringWidth(tc.Input)
			}
		}
	})

	b.Run("rivo/uniseg", func(b *testing.B) {
		// Note: rivo/uniseg doesn't have an equivalent to StrictEmojiNeutral
		// It uses emoji presentation properties directly
		testCases, n, err := loadTestCases()
		if err != nil {
			b.Fatalf("Failed to load test cases: %v", err)
		}
		b.SetBytes(n)
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for _, tc := range testCases {
				_ = uniseg.StringWidth(tc.Input)
			}
		}
	})
}

// BenchmarkString_ASCII benchmarks ASCII-only strings
func BenchmarkString_ASCII(b *testing.B) {
	asciiStrings := []string{
		"hello",
		"Hello World",
		"1234567890",
		"!@#$%^&*()",
		"This is a very long string with many characters to test performance of both implementations.",
	}

	n := 0
	for _, s := range asciiStrings {
		n += len(s)
	}

	b.Run("clipperhouse/displaywidth", func(b *testing.B) {
		b.SetBytes(int64(n))
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for _, s := range asciiStrings {
				_ = displaywidth.String(s)
			}
		}
	})

	b.Run("mattn/go-runewidth", func(b *testing.B) {
		b.SetBytes(int64(n))
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for _, s := range asciiStrings {
				_ = runewidth.StringWidth(s)
			}
		}
	})

	b.Run("rivo/uniseg", func(b *testing.B) {
		b.SetBytes(int64(n))
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for _, s := range asciiStrings {
				_ = uniseg.StringWidth(s)
			}
		}
	})
}

// BenchmarkString_Unicode benchmarks Unicode strings
func BenchmarkString_Unicode(b *testing.B) {
	unicodeStrings := []string{
		"café",
		"naïve",
		"résumé",
		"Zürich",
		"São Paulo",
		"中文",
		"こんにちは",
		"안녕하세요",
		"Hello 世界",
		"★ ☆ ♠ ♣ ♥ ♦",
		"° ± × ÷",
		"← → ↑ ↓",
	}

	n := 0
	for _, s := range unicodeStrings {
		n += len(s)
	}

	b.Run("clipperhouse/displaywidth", func(b *testing.B) {
		b.SetBytes(int64(n))
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, s := range unicodeStrings {
				_ = displaywidth.String(s)
			}
		}
	})

	b.Run("mattn/go-runewidth", func(b *testing.B) {
		b.SetBytes(int64(n))
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, s := range unicodeStrings {
				_ = runewidth.StringWidth(s)
			}
		}
	})

	b.Run("rivo/uniseg", func(b *testing.B) {
		b.SetBytes(int64(n))
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, s := range unicodeStrings {
				_ = uniseg.StringWidth(s)
			}
		}
	})
}

// BenchmarkStringWidth_Emoji benchmarks emoji strings
func BenchmarkStringWidth_Emoji(b *testing.B) {
	emojiStrings := []string{
		"😀 😁 😂 🤣 😃 😄 😅 😆 😉 😊",
		"🚀 🎉 🎊 🎈 🎁 🎂 🎃 🎄 🎆 🎇",
		"👨‍👩‍👧‍👦 👨‍💻 👩‍🔬 👨‍🎨 👩‍🚀",
		"🇺🇸 🇬🇧 🇫🇷 🇩🇪 🇯🇵 🇰🇷 🇨🇳",
		"Hello 世界! 😀",
		"👨‍💻 working on 🚀",
		"😀😁😂🤣😃😄😅😆😉😊😋😎😍😘🥰😗😙😚☺️🙂🤗🤩🤔🤨😐😑😶🙄😏😣😥😮🤐😯😪😫🥱😴😌😛😜😝🤤😒😓😔😕🙃🤑😲☹️🙁😖😞😟😤😢😭😦😧😨😩🤯😬😰😱🥵🥶😳🤪😵😡😠🤬😷🤒🤕🤢🤮🤧😇🤠🤡🥳🥴🥺🤥🤫🤭🧐🤓😈👿💀☠️👹👺🤖👽👾💩😺😸😹😻😼😽🙀😿😾",
	}

	n := 0
	for _, s := range emojiStrings {
		n += len(s)
	}

	b.Run("clipperhouse/displaywidth", func(b *testing.B) {
		b.SetBytes(int64(n))
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, s := range emojiStrings {
				_ = displaywidth.String(s)
			}
		}
	})

	b.Run("mattn/go-runewidth", func(b *testing.B) {
		b.SetBytes(int64(n))
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, s := range emojiStrings {
				_ = runewidth.StringWidth(s)
			}
		}
	})

	b.Run("rivo/uniseg", func(b *testing.B) {
		b.SetBytes(int64(n))
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, s := range emojiStrings {
				_ = uniseg.StringWidth(s)
			}
		}
	})
}

// BenchmarkString_Mixed benchmarks mixed content strings
func BenchmarkString_Mixed(b *testing.B) {
	mixedStrings := []string{
		"Hello 世界! 😀",
		"Price: $100.00 €85.50",
		"Math: ∑(x²) = ∞",
		"Emoji sequence: 👨‍💻 working on 🚀",
		"Mixed script: Hello 世界 안녕하세요 こんにちは",
		"This is a very long string with many characters to test performance of both implementations. It contains various character types including ASCII, Unicode, emoji, and special symbols. The purpose is to see how both packages handle longer strings and whether there are any performance differences or edge cases that emerge with more complex input.",
	}

	n := 0
	for _, s := range mixedStrings {
		n += len(s)
	}
	b.Run("clipperhouse/displaywidth", func(b *testing.B) {
		b.SetBytes(int64(n))
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, s := range mixedStrings {
				_ = displaywidth.String(s)
			}
		}
	})

	b.Run("mattn/go-runewidth", func(b *testing.B) {
		b.SetBytes(int64(n))
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, s := range mixedStrings {
				_ = runewidth.StringWidth(s)
			}
		}
	})

	b.Run("rivo/uniseg", func(b *testing.B) {
		b.SetBytes(int64(n))
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, s := range mixedStrings {
				_ = uniseg.StringWidth(s)
			}
		}
	})
}

// BenchmarkString_ControlChars benchmarks control characters
func BenchmarkString_ControlChars(b *testing.B) {
	controlStrings := []string{
		"\n",
		"\t",
		"\r",
		"\b",
		"\x00",
		"\x7f",
		"hello\nworld",
		"hello\tworld",
		" \t \n ",
	}

	n := 0
	for _, s := range controlStrings {
		n += len(s)
	}

	b.Run("clipperhouse/displaywidth", func(b *testing.B) {
		b.SetBytes(int64(n))
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, s := range controlStrings {
				_ = displaywidth.String(s)
			}
		}
	})

	b.Run("mattn/go-runewidth", func(b *testing.B) {
		b.SetBytes(int64(n))
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, s := range controlStrings {
				_ = runewidth.StringWidth(s)
			}
		}
	})

	b.Run("rivo/uniseg", func(b *testing.B) {
		b.SetBytes(int64(n))
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, s := range controlStrings {
				_ = uniseg.StringWidth(s)
			}
		}
	})
}

// BenchmarkRuneDefault benchmarks rune width calculation
func BenchmarkRuneDefault(b *testing.B) {
	testRunes := []rune{
		// Control characters
		'\x00', '\t', '\n', '\r', '\x7F',
		// ASCII printable
		' ', 'a', 'Z', '0', '9', '!', '@', '~',
		// Latin extended
		'é', 'ñ', 'ö',
		// East Asian Wide
		'中', '文', 'あ', 'ア', '가', '한',
		// Fullwidth
		'Ａ', 'Ｚ', '０', '９', '！', '　',
		// Ambiguous
		'★', '°', '±', '§', '©', '®',
		// Emoji
		'😀', '😁', '😍', '🤔', '🚀', '🎉', '🔥', '👍',
		// Mathematical symbols
		'∞', '∑', '∫', '√',
		// Currency
		'$', '€', '£', '¥',
		// Box drawing
		'─', '│', '┌',
		// Special Unicode
		'•', '…', '—',
	}

	n := 0
	for _, r := range testRunes {
		n += len(string(r))
	}
	b.SetBytes(int64(n))

	b.Run("clipperhouse/displaywidth", func(b *testing.B) {
		b.SetBytes(int64(n))
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, r := range testRunes {
				_ = displaywidth.Rune(r)
			}
		}
	})

	b.Run("mattn/go-runewidth", func(b *testing.B) {
		b.SetBytes(int64(n))
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, r := range testRunes {
				_ = runewidth.RuneWidth(r)
			}
		}
	})
}

// BenchmarkRuneWidth_EAW benchmarks rune width with East Asian Width option
func BenchmarkRuneWidth_EAW(b *testing.B) {
	options := displaywidth.Options{
		EastAsianWidth:     true,
		StrictEmojiNeutral: false,
	}

	condition := &runewidth.Condition{
		EastAsianWidth:     options.EastAsianWidth,
		StrictEmojiNeutral: options.StrictEmojiNeutral,
	}

	// Focus on ambiguous characters that are affected by EastAsianWidth
	testRunes := []rune{
		'★', '°', '±', '§', '©', '®',
		'中', '文', 'あ', 'ア', '가', '한',
		'Ａ', 'Ｚ', '０', '９',
		'😀', '🚀', '🎉',
	}

	n := 0
	for _, r := range testRunes {
		n += len(string(r))
	}

	b.Run("clipperhouse/displaywidth", func(b *testing.B) {
		b.SetBytes(int64(n))
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, r := range testRunes {
				_ = options.Rune(r)
			}
		}

	})

	b.Run("mattn/go-runewidth", func(b *testing.B) {
		b.SetBytes(int64(n))
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, r := range testRunes {
				_ = condition.RuneWidth(r)
			}
		}
	})
}

// BenchmarkRuneWidth_ASCII benchmarks ASCII rune width calculation
func BenchmarkRuneWidth_ASCII(b *testing.B) {
	asciiRunes := []rune{
		' ', 'a', 'b', 'c', 'x', 'y', 'z',
		'A', 'B', 'C', 'X', 'Y', 'Z',
		'0', '1', '2', '3', '8', '9',
		'!', '@', '#', '$', '%', '^', '&', '*', '(', ')',
	}

	n := 0
	for _, r := range asciiRunes {
		n += len(string(r))
	}
	b.SetBytes(int64(n))

	b.Run("clipperhouse/displaywidth", func(b *testing.B) {
		b.SetBytes(int64(n))
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, r := range asciiRunes {
				_ = displaywidth.Rune(r)
			}
		}
	})

	b.Run("mattn/go-runewidth", func(b *testing.B) {
		b.SetBytes(int64(n))
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, r := range asciiRunes {
				_ = runewidth.RuneWidth(r)
			}
		}
	})
}
