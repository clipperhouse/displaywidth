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
	CombiningMarks map[rune]bool   // From Go stdlib
	ZeroWidthChars map[rune]bool   // Special zero-width characters
}

// CharProperties represents the properties of a character as bit flags
type CharProperties uint16

const (
	// East Asian Width properties
	EAW_Fullwidth CharProperties = 1 << iota // F
	EAW_Wide                                 // W
	EAW_Halfwidth                            // H
	EAW_Narrow                               // Na
	EAW_Neutral                              // N
	EAW_Ambiguous                            // A

	// General categories
	IsCombiningMark          // Mn, Mc, Me
	IsControlChar            // C0, C1, DEL
	IsZeroWidth              // ZWSP, ZWJ, ZWNJ, etc.
	IsEmoji                  // Emoji base characters
	IsEmojiModifier          // Emoji modifiers
	IsEmojiVariationSelector // Emoji variation selectors
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
	if err := downloadFile("https://unicode.org/Public/UCD/latest/ucd/EastAsianWidth.txt", eawFile); err != nil {
		return nil, fmt.Errorf("failed to download EastAsianWidth.txt: %v", err)
	}
	if err := parseEastAsianWidth(eawFile, data); err != nil {
		return nil, fmt.Errorf("failed to parse EastAsianWidth.txt: %v", err)
	}

	// Download and parse emoji-data.txt
	emojiFile := filepath.Join(dataDir, "emoji-data.txt")
	if err := downloadFile("https://unicode.org/Public/UCD/latest/ucd/emoji-data.txt", emojiFile); err != nil {
		fmt.Printf("Warning: failed to download emoji-data.txt: %v\n", err)
		fmt.Println("Continuing with basic emoji detection from Go stdlib...")
	} else {
		if err := parseEmojiData(emojiFile, data); err != nil {
			fmt.Printf("Warning: failed to parse emoji-data.txt: %v\n", err)
			fmt.Println("Continuing with basic emoji detection from Go stdlib...")
		}
	}

	// Extract data from Go stdlib
	extractStdlibData(data)

	// Add special zero-width characters
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

// parseEmojiData parses the emoji-data.txt file
func parseEmojiData(filename string, data *UnicodeData) error {
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
		property := strings.TrimSpace(parts[1])

		// We're interested in emoji properties
		if !strings.Contains(property, "Emoji") {
			continue
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
				data.EmojiData[r] = true
			}
		} else {
			// Single codepoint
			codepoint, err := strconv.ParseInt(rangeStr, 16, 32)
			if err != nil {
				continue
			}
			data.EmojiData[rune(codepoint)] = true
		}
	}

	return scanner.Err()
}

// extractStdlibData extracts character properties from Go's unicode package
func extractStdlibData(data *UnicodeData) {
	// Extract control characters
	for r := rune(0); r <= 0x1F; r++ {
		data.ControlChars[r] = true
	}
	data.ControlChars[0x7F] = true // DEL
	for r := rune(0x80); r <= 0x9F; r++ {
		data.ControlChars[r] = true // C1 controls
	}

	// Extract combining marks
	for r := rune(0); r <= unicode.MaxRune; r++ {
		if unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Mc, r) || unicode.Is(unicode.Me, r) {
			data.CombiningMarks[r] = true
		}
	}

	// Extract basic emoji data from Go stdlib
	// This is a simplified approach - we'll detect common emoji ranges
	emojiRanges := []*unicode.RangeTable{
		unicode.So, // Symbol, other (includes some emoji)
	}

	for _, table := range emojiRanges {
		for r := rune(0); r <= unicode.MaxRune; r++ {
			if unicode.Is(table, r) {
				// Basic emoji detection - characters in certain ranges
				if (r >= 0x1F600 && r <= 0x1F64F) || // Emoticons
					(r >= 0x1F300 && r <= 0x1F5FF) || // Misc Symbols and Pictographs
					(r >= 0x1F680 && r <= 0x1F6FF) || // Transport and Map
					(r >= 0x1F1E0 && r <= 0x1F1FF) || // Regional indicators
					(r >= 0x2600 && r <= 0x26FF) || // Misc symbols
					(r >= 0x2700 && r <= 0x27BF) { // Dingbats
					data.EmojiData[r] = true
				}
			}
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
	if eaw, exists := data.EastAsianWidth[r]; exists {
		switch eaw {
		case "F":
			props |= EAW_Fullwidth
		case "W":
			props |= EAW_Wide
		case "H":
			props |= EAW_Halfwidth
		case "Na":
			props |= EAW_Narrow
		case "N":
			props |= EAW_Neutral
		case "A":
			props |= EAW_Ambiguous
		}

	}

	// General categories
	if data.CombiningMarks[r] {
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
