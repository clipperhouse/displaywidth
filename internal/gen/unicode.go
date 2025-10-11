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
	ControlChars   map[rune]bool   // From Go stdlib
	CombiningMarks map[rune]bool   // From Go stdlib (Mn, Me only - Mc excluded for proper width)
	ZeroWidthChars map[rune]bool   // Special zero-width characters
}

// property represents the properties of a character as bit flags
type property uint8

// PropertyDefinition describes a single character property flag
type PropertyDefinition struct {
	Name    string
	Comment string
}

// PropertyDefinitions is the single source of truth for all character properties.
// The order matters - it defines the bit positions (via iota).
var PropertyDefinitions = []PropertyDefinition{
	{"East_Asian_Full_Wide", "F, W"},
	{"East_Asian_Ambiguous", "A"},
	{"Emoji", "Emoji base characters"},
	{"ZeroWidth", "ZWSP, ZWJ, ZWNJ, etc."},
}

const (
	East_Asian_Full_Wide property = 1 << iota // F, W
	East_Asian_Ambiguous                      // A
	Emoji                                     // Emoji base characters
	ZeroWidth                                 // ZWSP, ZWJ, ZWNJ, etc.
)

// ParseUnicodeData downloads and parses all required Unicode data files
func ParseUnicodeData() (*UnicodeData, error) {
	data := &UnicodeData{
		EastAsianWidth: make(map[rune]string),
		EmojiData:      make(map[rune]bool),
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

	// Extract combining marks using range tables for efficiency
	// Mn: Nonspacing_Mark, Me: Enclosing_Mark
	// Note: Mc (Spacing Mark) characters are excluded so they get default width 1
	extractRunesFromRangeTable(unicode.Mn, data.CombiningMarks)
	extractRunesFromRangeTable(unicode.Me, data.CombiningMarks)

	// Cf (Other, format) is the official Unicode category for format characters
	// which are generally invisible and have zero width.
	extractRunesFromRangeTable(unicode.Cf, data.ZeroWidthChars)

	// Zl (Other, line separator) is the official Unicode category for line separator characters
	// which are generally invisible and have zero width.
	extractRunesFromRangeTable(unicode.Zl, data.ZeroWidthChars)

	// Zp (Other, paragraph separator) is the official Unicode category for paragraph separator characters
	// which are generally invisible and have zero width.
	extractRunesFromRangeTable(unicode.Zp, data.ZeroWidthChars)

	// Noncharacters (U+nFFFE and U+nFFFF)
	data.ZeroWidthChars[0xFFFE] = true
	data.ZeroWidthChars[0xFFFF] = true
}

// extractRunesFromRangeTable efficiently extracts all runes from a Unicode range table
func extractRunesFromRangeTable(table *unicode.RangeTable, target map[rune]bool) {
	// Iterate over 16-bit ranges
	for _, r16 := range table.R16 {
		for r := rune(r16.Lo); r <= rune(r16.Hi); r += rune(r16.Stride) {
			target[r] = true
		}
	}

	// Iterate over 32-bit ranges
	for _, r32 := range table.R32 {
		for r := rune(r32.Lo); r <= rune(r32.Hi); r += rune(r32.Stride) {
			target[r] = true
		}
	}
}

// BuildPropertyBitmap creates a properties bitmap for a given rune
func BuildPropertyBitmap(r rune, data *UnicodeData) property {
	var props property

	// East Asian Width
	// Only store properties that affect width calculation
	if eaw, exists := data.EastAsianWidth[r]; exists {
		switch eaw {
		case "F", "W":
			props |= East_Asian_Full_Wide
		case "A":
			props |= East_Asian_Ambiguous
			// H (Halfwidth), Na (Narrow), and N (Neutral) are not stored
			// as they all result in width 1 (default behavior)
		}
	}

	if data.CombiningMarks[r] {
		props |= ZeroWidth
	}
	if data.ControlChars[r] {
		props |= ZeroWidth
	}
	if data.ZeroWidthChars[r] {
		props |= ZeroWidth
	}

	if data.EmojiData[r] {
		props |= Emoji
	}

	return props
}
