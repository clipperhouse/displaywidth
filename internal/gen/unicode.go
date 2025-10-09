// Package unicode handles parsing of Unicode data files for string width calculation
package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

// UnicodeData contains all the parsed Unicode character properties
type UnicodeData struct {
	EastAsianWidth map[rune]string // From EastAsianWidth.txt
	EmojiData      map[rune]bool   // From emoji-data.txt
	AmbiguousData  map[rune]bool   // Ambiguous width characters from EastAsianWidth.txt
	ControlChars   map[rune]bool   // From Go stdlib
	CombiningMarks map[rune]bool   // From Go stdlib
	ZeroWidthChars map[rune]bool   // Special zero-width characters
}

// CharProperties represents the properties of a character as bit flags
type CharProperties uint8

const (
	// East Asian Width properties
	EAW_Fullwidth CharProperties = 1 << iota // F
	EAW_Wide                                 // W
	EAW_Ambiguous                            // A

	// General categories
	IsCombiningMark // Mn, Mc, Me
	IsControlChar   // C0, C1, DEL
	IsZeroWidth     // ZWSP, ZWJ, ZWNJ, etc.
	IsEmoji         // Emoji base characters
)

// ParseUnicodeData downloads and parses all required Unicode data files
func ParseUnicodeData() (*UnicodeData, error) {
	data := &UnicodeData{
		EastAsianWidth: make(map[rune]string),
		EmojiData:      make(map[rune]bool),
		AmbiguousData:  make(map[rune]bool),
		ControlChars:   make(map[rune]bool),
		CombiningMarks: make(map[rune]bool),
		ZeroWidthChars: make(map[rune]bool),
	}

	// Create data directory
	dataDir := "data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %v", err)
	}

	// Download and parse EastAsianWidth.txt
	eawFile := filepath.Join(dataDir, "EastAsianWidth.txt")
	if err := downloadFile("https://unicode.org/Public/UCD/15.1.0/ucd/EastAsianWidth.txt", eawFile); err != nil {
		return nil, fmt.Errorf("failed to download EastAsianWidth.txt: %v", err)
	}
	if err := parseEastAsianWidth(eawFile, data); err != nil {
		return nil, fmt.Errorf("failed to parse EastAsianWidth.txt: %v", err)
	}

	// Download and parse emoji-data.txt (use same version as go-runewidth for compatibility)
	emojiFile := filepath.Join(dataDir, "emoji-data.txt")
	if err := downloadFile("https://unicode.org/Public/15.1.0/ucd/emoji/emoji-data.txt", emojiFile); err != nil {
		fmt.Printf("Warning: failed to download emoji-data.txt: %v\n", err)
		fmt.Println("Continuing with basic emoji detection from Go stdlib...")
	} else {
		if err := parseEmojiData(emojiFile, data); err != nil {
			fmt.Printf("Warning: failed to parse emoji-data.txt: %v\n", err)
			fmt.Println("Continuing with basic emoji detection from Go stdlib...")
		}
	}

	extractStdlibData(data)
	extractAmbiguousChars(data)
	addZeroWidthChars(data)

	return data, nil
}

// downloadFile downloads a file from URL to local path
func downloadFile(url, filepath string) error {
	// Check if file already exists
	if _, err := os.Stat(filepath); err == nil {
		fmt.Printf("File %s already exists, skipping download\n", filepath)
		return nil
	}

	fmt.Printf("Downloading %s...\n", url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("Downloaded %s\n", filepath)
	return nil
}

// parseEastAsianWidth parses the EastAsianWidth.txt file
func parseEastAsianWidth(filename string, data *UnicodeData) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, ";")
		if len(parts) < 2 {
			continue
		}

		rangeStr := strings.TrimSpace(parts[0])
		widthStr := strings.TrimSpace(parts[1])

		// Remove comments from width string
		if commentIndex := strings.Index(widthStr, "#"); commentIndex != -1 {
			widthStr = strings.TrimSpace(widthStr[:commentIndex])
		}

		// Parse range
		if strings.Contains(rangeStr, "..") {
			// Range of codepoints
			rangeParts := strings.Split(rangeStr, "..")
			if len(rangeParts) != 2 {
				continue
			}
			start, err1 := strconv.ParseInt(rangeParts[0], 16, 32)
			end, err2 := strconv.ParseInt(rangeParts[1], 16, 32)
			if err1 != nil || err2 != nil {
				continue
			}
			for r := rune(start); r <= rune(end); r++ {
				data.EastAsianWidth[r] = widthStr
			}
		} else {
			// Single codepoint
			codepoint, err := strconv.ParseInt(rangeStr, 16, 32)
			if err != nil {
				continue
			}
			data.EastAsianWidth[rune(codepoint)] = widthStr
		}
	}

	return scanner.Err()
}

// parseEmojiData parses the emoji-data.txt file using the same logic as go-runewidth
func parseEmojiData(filename string, data *UnicodeData) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Skip until we find the Extended_Pictographic=No line (same as go-runewidth)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Extended_Pictographic=No") {
			break
		}
	}

	if scanner.Err() != nil {
		return scanner.Err()
	}

	// Parse the Extended_Pictographic=Yes ranges (same as go-runewidth)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		var r1, r2 rune
		n, err := fmt.Sscanf(line, "%x..%x ", &r1, &r2)
		if err != nil || n == 1 {
			n, err = fmt.Sscanf(line, "%x ", &r1)
			if err != nil || n != 1 {
				continue
			}
			r2 = r1
		}

		// Skip characters below 0xFF (same as go-runewidth)
		if r2 < 0xFF {
			continue
		}

		for r := r1; r <= r2; r++ {
			data.EmojiData[r] = true
		}
	}

	return scanner.Err()
}

// extractStdlibData extracts character properties from Go's unicode package
func extractStdlibData(data *UnicodeData) {
	// Extract control characters
	// Skip 0x00-0x1F and 0x7F as they're handled by the fast path in width.go
	// Only add C1 controls (0x80-0x9F) which are multi-byte in UTF-8
	for r := rune(0x80); r <= 0x9F; r++ {
		data.ControlChars[r] = true // C1 controls
	}

	// Extract combining marks
	for r := rune(0); r <= unicode.MaxRune; r++ {
		if unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Mc, r) || unicode.Is(unicode.Me, r) {
			data.CombiningMarks[r] = true
		}
	}
}

// extractAmbiguousChars extracts all characters marked as "A" (Ambiguous)
// in EastAsianWidth.txt and adds them to the AmbiguousData map
func extractAmbiguousChars(data *UnicodeData) {
	for r, width := range data.EastAsianWidth {
		if width == "A" {
			data.AmbiguousData[r] = true
		}
	}
}

// addZeroWidthChars adds special zero-width characters
func addZeroWidthChars(data *UnicodeData) {
	zeroWidthChars := []rune{
		0x200B, // Zero Width Space
		0x200C, // Zero Width Non-Joiner
		0x200D, // Zero Width Joiner
		0x2060, // Word Joiner
		0x2061, // Function Application
		0xFEFF, // Zero Width No-Break Space
	}

	for _, r := range zeroWidthChars {
		data.ZeroWidthChars[r] = true
	}
}

// BuildPropertyBitmap creates a CharProperties bitmap for a given rune
func BuildPropertyBitmap(r rune, data *UnicodeData) CharProperties {
	var props CharProperties

	// East Asian Width
	// Only store properties that affect width calculation
	if eaw, exists := data.EastAsianWidth[r]; exists {
		switch eaw {
		case "F":
			props |= EAW_Fullwidth
		case "W":
			props |= EAW_Wide
		case "A":
			props |= EAW_Ambiguous
			// H (Halfwidth), Na (Narrow), and N (Neutral) are not stored
			// as they all result in width 1 (default behavior)
		}
	}

	if data.CombiningMarks[r] && !isExceptionalCombiningMark(r) {
		props |= IsCombiningMark
	}
	if data.ControlChars[r] {
		props |= IsControlChar
	}
	if data.ZeroWidthChars[r] {
		props |= IsZeroWidth
	}

	if data.EmojiData[r] {
		props |= IsEmoji
	}

	return props
}

// isExceptionalCombiningMark removes certain combining marks that
// go-runewidth treats as regular characters. We believe this might be
// incorrect in go-runewidth, but would need to confirm. This is
// debt/expediency for now.
func isExceptionalCombiningMark(r rune) bool {
	// Thai combining marks (U+0E31-U+0E3A, U+0E47-U+0E4F)
	if (r >= 0x0E31 && r <= 0x0E3A) || (r >= 0x0E47 && r <= 0x0E4F) {
		return true
	}
	// Bengali combining marks (U+0982, U+09BC-U+09C4, U+09C7-U+09C8, U+09CB-U+09CD, U+09D7)
	if r == 0x0982 || (r >= 0x09BC && r <= 0x09C4) || (r >= 0x09C7 && r <= 0x09C8) || (r >= 0x09CB && r <= 0x09CD) || r == 0x09D7 {
		return true
	}
	// Devanagari combining marks (U+0900-U+0903, U+093A-U+093C, U+093E-U+0940, U+0941-U+0948, U+0949-U+094D, U+0951-U+0957, U+0962-U+0963)
	if (r >= 0x0900 && r <= 0x0903) || (r >= 0x093A && r <= 0x093C) || (r >= 0x093E && r <= 0x0940) || (r >= 0x0941 && r <= 0x0948) || (r >= 0x0949 && r <= 0x094D) || (r >= 0x0951 && r <= 0x0957) || (r >= 0x0962 && r <= 0x0963) {
		return true
	}
	// Arabic combining marks (U+064B-U+0655, U+0657-U+065E, U+0670, U+06D6-U+06DC, U+06DF-U+06E4, U+06E7-U+06E8, U+06EA-U+06ED)
	if (r >= 0x064B && r <= 0x0655) || (r >= 0x0657 && r <= 0x065E) || r == 0x0670 || (r >= 0x06D6 && r <= 0x06DC) || (r >= 0x06DF && r <= 0x06E4) || (r >= 0x06E7 && r <= 0x06E8) || (r >= 0x06EA && r <= 0x06ED) {
		return true
	}
	// Additional Devanagari combining marks that go-runewidth treats as regular characters
	if (r >= 0x09BE && r <= 0x09C4) || (r >= 0x09C7 && r <= 0x09C7) || (r >= 0x09CB && r <= 0x09CD) || (r >= 0x09D7 && r <= 0x09D7) {
		return true
	}
	// Variation selectors that go-runewidth treats as regular characters
	if r == 0xFE0F {
		return true
	}
	return false
}
