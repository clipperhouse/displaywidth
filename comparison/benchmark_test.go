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

var (
	// Shared test data for benchmarks
	asciiTestStrings = []string{
		"hello",
		"Hello World",
		"1234567890",
		"!@#$%^&*()",
		"This is a very long string with many characters to test performance of both implementations.",
	}

	emojiTestStrings = []string{
		"😀 😁 😂 🤣 😃 😄 😅 😆 😉 😊",
		"🚀 🎉 🎊 🎈 🎁 🎂 🎃 🎄 🎆 🎇",
		"👨‍👩‍👧‍👦 👨‍💻 👩‍🔬 👨‍🎨 👩‍🚀",
		"🇺🇸 🇬🇧 🇫🇷 🇩🇪 🇯🇵 🇰🇷 🇨🇳",
		"Hello 世界! 😀",
		"👨‍💻 working on 🚀",
		"😀😁😂🤣😃😄😅😆😉😊😋😎😍😘🥰😗😙😚☺️🙂🤗🤩🤔🤨😐😑😶🙄😏😣😥😮🤐😯😪😫🥱😴😌😛😜😝🤤😒😓😔😕🙃🤑😲☹️🙁😖😞😟😤😢😭😦😧😨😩🤯😬😰😱🥵🥶😳🤪😵😡😠🤬😷🤒🤕🤢🤮🤧😇🤠🤡🥳🥴🥺🤥🤫🤭🧐🤓😈👿💀☠️👹👺🤖👽👾💩😺😸😹😻😼😽🙀😿😾",
	}
)

// BenchmarkString_Mixed benchmarks our displaywidth package
func BenchmarkString_Mixed(b *testing.B) {
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
				// Test with default settings (eastAsianWidth=false)
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

func BenchmarkString_EastAsian(b *testing.B) {
	options := displaywidth.Options{
		EastAsianWidth: true,
	}

	condition := &runewidth.Condition{
		EastAsianWidth: true,
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

// BenchmarkString_ASCII benchmarks ASCII-only strings
func BenchmarkString_ASCII(b *testing.B) {
	n := 0
	for _, s := range asciiTestStrings {
		n += len(s)
	}

	b.Run("clipperhouse/displaywidth", func(b *testing.B) {
		b.SetBytes(int64(n))
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for _, s := range asciiTestStrings {
				_ = displaywidth.String(s)
			}
		}
	})

	b.Run("mattn/go-runewidth", func(b *testing.B) {
		b.SetBytes(int64(n))
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for _, s := range asciiTestStrings {
				_ = runewidth.StringWidth(s)
			}
		}
	})

	b.Run("rivo/uniseg", func(b *testing.B) {
		b.SetBytes(int64(n))
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			for _, s := range asciiTestStrings {
				_ = uniseg.StringWidth(s)
			}
		}
	})
}

// BenchmarkString_Emoji benchmarks emoji strings
func BenchmarkString_Emoji(b *testing.B) {
	n := 0
	for _, s := range emojiTestStrings {
		n += len(s)
	}

	b.Run("clipperhouse/displaywidth", func(b *testing.B) {
		b.SetBytes(int64(n))
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, s := range emojiTestStrings {
				_ = displaywidth.String(s)
			}
		}
	})

	b.Run("mattn/go-runewidth", func(b *testing.B) {
		b.SetBytes(int64(n))
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, s := range emojiTestStrings {
				_ = runewidth.StringWidth(s)
			}
		}
	})

	b.Run("rivo/uniseg", func(b *testing.B) {
		b.SetBytes(int64(n))
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, s := range emojiTestStrings {
				_ = uniseg.StringWidth(s)
			}
		}
	})
}

// BenchmarkRune_Mixed benchmarks rune width calculation using test cases
func BenchmarkRune_Mixed(b *testing.B) {
	testCases, _, err := loadTestCases()
	if err != nil {
		b.Fatalf("Failed to load test cases: %v", err)
	}

	// Convert all strings to []rune
	var testRunes []rune
	n := 0
	for _, tc := range testCases {
		runes := []rune(tc.Input)
		testRunes = append(testRunes, runes...)
		n += len(tc.Input)
	}

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

// BenchmarkRune_EastAsian benchmarks rune width with East Asian Width option
func BenchmarkRune_EastAsian(b *testing.B) {
	options := displaywidth.Options{
		EastAsianWidth: true,
	}

	condition := &runewidth.Condition{
		EastAsianWidth: true,
	}

	testCases, _, err := loadTestCases()
	if err != nil {
		b.Fatalf("Failed to load test cases: %v", err)
	}

	// Convert all strings to []rune
	var testRunes []rune
	n := 0
	for _, tc := range testCases {
		runes := []rune(tc.Input)
		testRunes = append(testRunes, runes...)
		n += len(tc.Input)
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

// BenchmarkRune_ASCII benchmarks ASCII rune width calculation
func BenchmarkRune_ASCII(b *testing.B) {
	// Convert ASCII strings to []rune
	var asciiRunes []rune
	n := 0
	for _, s := range asciiTestStrings {
		runes := []rune(s)
		asciiRunes = append(asciiRunes, runes...)
		n += len(s)
	}

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

// BenchmarkRune_Emoji benchmarks emoji rune width calculation
func BenchmarkRune_Emoji(b *testing.B) {
	// Convert emoji strings to []rune
	var emojiRunes []rune
	n := 0
	for _, s := range emojiTestStrings {
		runes := []rune(s)
		emojiRunes = append(emojiRunes, runes...)
		n += len(s)
	}

	b.Run("clipperhouse/displaywidth", func(b *testing.B) {
		b.SetBytes(int64(n))
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, r := range emojiRunes {
				_ = displaywidth.Rune(r)
			}
		}
	})

	b.Run("mattn/go-runewidth", func(b *testing.B) {
		b.SetBytes(int64(n))
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, r := range emojiRunes {
				_ = runewidth.RuneWidth(r)
			}
		}
	})
}
