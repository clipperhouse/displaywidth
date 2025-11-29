package main

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/clipperhouse/displaywidth"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--help" {
		fmt.Println("Usage: terminal-test")
		fmt.Println("Tests the actual terminal display width of regional indicator symbols")
		fmt.Println("by using visual alignment tests.")
		fmt.Println()
		fmt.Println("Run from a terminal to see visual alignment results.")
		os.Exit(0)
	}

	fmt.Println("=== Terminal Regional Indicator Width Test ===")
	fmt.Printf("TERM: %s\n", os.Getenv("TERM"))
	fmt.Println()

	// Test characters
	singleRI := "🇦" // U+1F1E6
	pairRI := "🇺🇸"  // US flag
	regularEmoji := "😀"
	ascii := "a"
	cjk := "中"

	testCases := []struct {
		name  string
		char  string
		width int
	}{
		{"Single Regional Indicator", singleRI, displaywidth.String(singleRI)},
		{"Regional Indicator Pair (flag)", pairRI, displaywidth.String(pairRI)},
		{"Regular Emoji", regularEmoji, displaywidth.String(regularEmoji)},
		{"ASCII", ascii, displaywidth.String(ascii)},
		{"CJK", cjk, displaywidth.String(cjk)},
	}

	fmt.Println("Package calculated widths:")
	for _, tc := range testCases {
		fmt.Printf("  %s (%s): %d columns\n", tc.name, tc.char, tc.width)
	}
	fmt.Println()

	// Visual alignment tests
	fmt.Println("=== Visual Alignment Tests ===")
	fmt.Println("Check if the markers align correctly with the characters.")
	fmt.Println("If aligned: terminal width matches package calculation.")
	fmt.Println("If misaligned: terminal rendering differs from package.")
	fmt.Println()

	for _, tc := range testCases {
		visualTest(tc.char, tc.name, tc.width)
		fmt.Println()
	}

	fmt.Println("=== Summary ===")
	fmt.Println("Compare the visual alignment above.")
	fmt.Println("The '^' marker shows the START of the character.")
	fmt.Println("The 'x' marker shows the expected END position (start + width).")
	fmt.Println("If characters align with the markers, the package calculation is correct.")
}

// visualTest prints a visual alignment test
func visualTest(char string, label string, expectedWidth int) {
	fmt.Printf("--- %s: %s ---\n", label, char)

	// Print alignment markers with known-width characters
	marker := "0123456789"
	testLine := marker + char + marker
	fmt.Println(testLine)

	// Print caret marker at the start of the character
	caretStart := len(marker)
	caretLine := strings.Repeat(" ", caretStart) + "^ (start)" + strings.Repeat(" ", len(marker)-6)
	fmt.Println(caretLine)

	// Print expected end position marker
	expectedEnd := len(marker) + expectedWidth
	if expectedEnd <= len(testLine) {
		expectedLine := strings.Repeat(" ", expectedEnd) + "x (expected end, width=" + fmt.Sprintf("%d", expectedWidth) + ")"
		fmt.Println(expectedLine)
	}

	// Print information
	runeCount := utf8.RuneCountInString(char)
	fmt.Printf("  UTF-8 bytes: %d | Runes: %d | Package width: %d\n", len(char), runeCount, expectedWidth)
}
