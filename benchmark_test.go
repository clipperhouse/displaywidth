package stringwidth

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/mattn/go-runewidth"
)

// TestCase represents a test case from the test_cases.txt file
type TestCase struct {
	Name  string
	Input string
}

// loadTestCases reads and parses test cases from test_cases.txt
func loadTestCases() ([]TestCase, error) {
	file, err := os.Open("test_cases.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var testCases []TestCase
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

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

	return testCases, scanner.Err()
}

// BenchmarkStringWidth_OurPackage benchmarks our stringwidth package
func BenchmarkStringWidth_OurPackage(b *testing.B) {
	testCases, err := loadTestCases()
	if err != nil {
		b.Fatalf("Failed to load test cases: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, tc := range testCases {
			// Test with default settings (eastAsianWidth=false, strictEmojiNeutral=false)
			_ = processStringWidth(tc.Input, false, false)
		}
	}

	// Set bytes for throughput calculation
	totalBytes := 0
	for _, tc := range testCases {
		totalBytes += len(tc.Input)
	}
	b.SetBytes(int64(totalBytes))
}

// BenchmarkStringWidth_GoRunewidth benchmarks go-runewidth package
func BenchmarkStringWidth_GoRunewidth(b *testing.B) {
	testCases, err := loadTestCases()
	if err != nil {
		b.Fatalf("Failed to load test cases: %v", err)
	}

	condition := runewidth.NewCondition()
	condition.EastAsianWidth = false
	condition.StrictEmojiNeutral = false

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, tc := range testCases {
			_ = condition.StringWidth(tc.Input)
		}
	}

	// Set bytes for throughput calculation
	totalBytes := 0
	for _, tc := range testCases {
		totalBytes += len(tc.Input)
	}
	b.SetBytes(int64(totalBytes))
}

// BenchmarkStringWidth_OurPackage_EAW benchmarks our package with East Asian Width enabled
func BenchmarkStringWidth_OurPackage_EAW(b *testing.B) {
	testCases, err := loadTestCases()
	if err != nil {
		b.Fatalf("Failed to load test cases: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, tc := range testCases {
			// Test with East Asian Width enabled
			_ = processStringWidth(tc.Input, true, false)
		}
	}

	// Set bytes for throughput calculation
	totalBytes := 0
	for _, tc := range testCases {
		totalBytes += len(tc.Input)
	}
	b.SetBytes(int64(totalBytes))
}

// BenchmarkStringWidth_GoRunewidth_EAW benchmarks go-runewidth with East Asian Width enabled
func BenchmarkStringWidth_GoRunewidth_EAW(b *testing.B) {
	testCases, err := loadTestCases()
	if err != nil {
		b.Fatalf("Failed to load test cases: %v", err)
	}

	condition := runewidth.NewCondition()
	condition.EastAsianWidth = true
	condition.StrictEmojiNeutral = false

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, tc := range testCases {
			_ = condition.StringWidth(tc.Input)
		}
	}

	// Set bytes for throughput calculation
	totalBytes := 0
	for _, tc := range testCases {
		totalBytes += len(tc.Input)
	}
	b.SetBytes(int64(totalBytes))
}

// BenchmarkStringWidth_OurPackage_StrictEmoji benchmarks our package with strict emoji neutral
func BenchmarkStringWidth_OurPackage_StrictEmoji(b *testing.B) {
	testCases, err := loadTestCases()
	if err != nil {
		b.Fatalf("Failed to load test cases: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, tc := range testCases {
			// Test with strict emoji neutral enabled
			_ = processStringWidth(tc.Input, false, true)
		}
	}

	// Set bytes for throughput calculation
	totalBytes := 0
	for _, tc := range testCases {
		totalBytes += len(tc.Input)
	}
	b.SetBytes(int64(totalBytes))
}

// BenchmarkStringWidth_GoRunewidth_StrictEmoji benchmarks go-runewidth with strict emoji neutral
func BenchmarkStringWidth_GoRunewidth_StrictEmoji(b *testing.B) {
	testCases, err := loadTestCases()
	if err != nil {
		b.Fatalf("Failed to load test cases: %v", err)
	}

	condition := runewidth.NewCondition()
	condition.EastAsianWidth = false
	condition.StrictEmojiNeutral = true

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, tc := range testCases {
			_ = condition.StringWidth(tc.Input)
		}
	}

	// Set bytes for throughput calculation
	totalBytes := 0
	for _, tc := range testCases {
		totalBytes += len(tc.Input)
	}
	b.SetBytes(int64(totalBytes))
}

// BenchmarkStringWidth_ASCII benchmarks ASCII-only strings
func BenchmarkStringWidth_ASCII(b *testing.B) {
	asciiStrings := []string{
		"hello",
		"Hello World",
		"1234567890",
		"!@#$%^&*()",
		"This is a very long string with many characters to test performance of both implementations.",
	}

	b.Run("OurPackage", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, s := range asciiStrings {
				_ = processStringWidth(s, false, false)
			}
		}

		totalBytes := 0
		for _, s := range asciiStrings {
			totalBytes += len(s)
		}
		b.SetBytes(int64(totalBytes))
	})

	b.Run("GoRunewidth", func(b *testing.B) {
		condition := runewidth.NewCondition()
		condition.EastAsianWidth = false
		condition.StrictEmojiNeutral = false

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, s := range asciiStrings {
				_ = condition.StringWidth(s)
			}
		}

		totalBytes := 0
		for _, s := range asciiStrings {
			totalBytes += len(s)
		}
		b.SetBytes(int64(totalBytes))
	})
}

// BenchmarkStringWidth_Unicode benchmarks Unicode strings
func BenchmarkStringWidth_Unicode(b *testing.B) {
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

	b.Run("OurPackage", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, s := range unicodeStrings {
				_ = processStringWidth(s, false, false)
			}
		}

		totalBytes := 0
		for _, s := range unicodeStrings {
			totalBytes += len(s)
		}
		b.SetBytes(int64(totalBytes))
	})

	b.Run("GoRunewidth", func(b *testing.B) {
		condition := runewidth.NewCondition()
		condition.EastAsianWidth = false
		condition.StrictEmojiNeutral = false

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, s := range unicodeStrings {
				_ = condition.StringWidth(s)
			}
		}

		totalBytes := 0
		for _, s := range unicodeStrings {
			totalBytes += len(s)
		}
		b.SetBytes(int64(totalBytes))
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

	b.Run("OurPackage", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, s := range emojiStrings {
				_ = processStringWidth(s, false, false)
			}
		}

		totalBytes := 0
		for _, s := range emojiStrings {
			totalBytes += len(s)
		}
		b.SetBytes(int64(totalBytes))
	})

	b.Run("GoRunewidth", func(b *testing.B) {
		condition := runewidth.NewCondition()
		condition.EastAsianWidth = false
		condition.StrictEmojiNeutral = false

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, s := range emojiStrings {
				_ = condition.StringWidth(s)
			}
		}

		totalBytes := 0
		for _, s := range emojiStrings {
			totalBytes += len(s)
		}
		b.SetBytes(int64(totalBytes))
	})
}

// BenchmarkStringWidth_Mixed benchmarks mixed content strings
func BenchmarkStringWidth_Mixed(b *testing.B) {
	mixedStrings := []string{
		"Hello 世界! 😀",
		"Price: $100.00 €85.50",
		"Math: ∑(x²) = ∞",
		"Emoji sequence: 👨‍💻 working on 🚀",
		"Mixed script: Hello 世界 안녕하세요 こんにちは",
		"This is a very long string with many characters to test performance of both implementations. It contains various character types including ASCII, Unicode, emoji, and special symbols. The purpose is to see how both packages handle longer strings and whether there are any performance differences or edge cases that emerge with more complex input.",
	}

	b.Run("OurPackage", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, s := range mixedStrings {
				_ = processStringWidth(s, false, false)
			}
		}

		totalBytes := 0
		for _, s := range mixedStrings {
			totalBytes += len(s)
		}
		b.SetBytes(int64(totalBytes))
	})

	b.Run("GoRunewidth", func(b *testing.B) {
		condition := runewidth.NewCondition()
		condition.EastAsianWidth = false
		condition.StrictEmojiNeutral = false

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, s := range mixedStrings {
				_ = condition.StringWidth(s)
			}
		}

		totalBytes := 0
		for _, s := range mixedStrings {
			totalBytes += len(s)
		}
		b.SetBytes(int64(totalBytes))
	})
}

// BenchmarkStringWidth_ControlChars benchmarks control characters
func BenchmarkStringWidth_ControlChars(b *testing.B) {
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

	b.Run("OurPackage", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, s := range controlStrings {
				_ = processStringWidth(s, false, false)
			}
		}

		totalBytes := 0
		for _, s := range controlStrings {
			totalBytes += len(s)
		}
		b.SetBytes(int64(totalBytes))
	})

	b.Run("GoRunewidth", func(b *testing.B) {
		condition := runewidth.NewCondition()
		condition.EastAsianWidth = false
		condition.StrictEmojiNeutral = false

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, s := range controlStrings {
				_ = condition.StringWidth(s)
			}
		}

		totalBytes := 0
		for _, s := range controlStrings {
			totalBytes += len(s)
		}
		b.SetBytes(int64(totalBytes))
	})
}
