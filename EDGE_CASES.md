# Edge Cases and Special Sequences

## Overview

This document details edge cases, special character sequences, and compatibility considerations for string width measurement.

## Control Characters

### C0 Controls (0x00-0x1F)

| Code Point | Name | Width | Notes |
|------------|------|-------|-------|
| 0x00 | NULL | 0 | Should not advance cursor |
| 0x01-0x07 | SOH-STX-ETX-EOT-ENQ-ACK-BEL | 0 | Control characters |
| 0x08 | Backspace | 0 | Cursor movement, no width |
| 0x09 | Tab | 0 | Cursor movement, no width |
| 0x0A | Line Feed | 0 | Cursor movement, no width |
| 0x0B | Vertical Tab | 0 | Control character |
| 0x0C | Form Feed | 0 | Control character |
| 0x0D | Carriage Return | 0 | Cursor movement, no width |
| 0x0E-0x1F | SO-SI-DLE-DC1-DC2-DC3-DC4-NAK-SYN-ETB-CAN-EM-SUB-ESC-FS-GS-RS-US | 0 | Control characters |

### DEL and C1 Controls

| Code Point | Name | Width | Notes |
|------------|------|-------|-------|
| 0x7F | DEL | 0 | Delete character |
| 0x80-0x9F | C1 Controls | 0 | Extended control characters |

## Combining Characters

### Nonspacing Marks (Mn)

**Common Examples**:
- U+0300-U+036F: Combining Diacritical Marks
- U+1AB0-U+1AFF: Combining Diacritical Marks Extended
- U+1DC0-U+1DFF: Combining Diacritical Marks Supplement
- U+20D0-U+20FF: Combining Diacritical Marks for Symbols
- U+FE20-U+FE2F: Combining Half Marks

**Width**: Always 0

### Spacing Marks (Mc)

**Common Examples**:
- U+0900-U+097F: Devanagari (some characters)
- U+0B00-U+0B7F: Oriya (some characters)
- U+0C00-U+0C7F: Telugu (some characters)

**Width**: Always 0

### Enclosing Marks (Me)

**Common Examples**:
- U+20DD-U+20E0: Enclosing Marks
- U+20E2-U+20E4: Enclosing Marks
- U+A670-U+A672: Cyrillic Extended-B

**Width**: Always 0

## Zero-Width Characters

### Standard Zero-Width Characters

| Code Point | Name | Width | Usage |
|------------|------|-------|-------|
| U+200B | Zero Width Space | 0 | Invisible word separator |
| U+200C | Zero Width Non-Joiner | 0 | Prevents joining in Arabic/Persian |
| U+200D | Zero Width Joiner | 0 | Joins characters (emoji sequences) |
| U+2060 | Word Joiner | 0 | Invisible word separator |
| U+2061 | Function Application | 0 | Mathematical notation |
| U+FEFF | Zero Width No-Break Space | 0 | Byte order mark, invisible |

## East Asian Width Edge Cases

### Ambiguous Characters (A)

**Common Ambiguous Characters**:
- U+00A1: ¡ (Inverted Exclamation Mark)
- U+00A4: ¤ (Currency Sign)
- U+00A7: § (Section Sign)
- U+00A8: ¨ (Diaeresis)
- U+00AA: ª (Feminine Ordinal Indicator)
- U+00AD: Soft Hyphen
- U+00AE: ® (Registered Sign)
- U+00B0: ° (Degree Sign)
- U+00B1: ± (Plus-Minus Sign)
- U+00B2: ² (Superscript Two)
- U+00B3: ³ (Superscript Three)
- U+00B4: ´ (Acute Accent)
- U+00B6: ¶ (Pilcrow Sign)
- U+00B7: · (Middle Dot)
- U+00B8: ¸ (Cedilla)
- U+00B9: ¹ (Superscript One)
- U+00BA: º (Masculine Ordinal Indicator)
- U+00BC: ¼ (Vulgar Fraction One Quarter)
- U+00BD: ½ (Vulgar Fraction One Half)
- U+00BE: ¾ (Vulgar Fraction Three Quarters)

**Default Width**: 1 (narrow)
**East Asian Width**: 2 (wide)

### Neutral Characters (N)

**Common Neutral Characters**:
- Most ASCII printable characters (0x20-0x7E)
- Most Latin characters
- Most symbols and punctuation

**Width**: Always 1

## Emoji and Special Sequences

### Emoji Base Characters

**Width**: Generally 2

**Examples**:
- 😀 (U+1F600): Grinning Face
- 🎉 (U+1F389): Party Popper
- 🌟 (U+1F31F): Glowing Star
- 🚀 (U+1F680): Rocket

### Emoji Modifiers

**Width**: Always 0

**Examples**:
- U+1F3FB-U+1F3FF: Skin Tone Modifiers
- U+1F9B0-U+1F9B3: Hair Style Modifiers

### Emoji Variation Selectors

**Width**: Always 0

**Examples**:
- U+FE0E: Variation Selector-15 (text style)
- U+FE0F: Variation Selector-16 (emoji style)

### ZWJ Sequences

**Complex Emoji Sequences**:
- Family emoji: 👨‍👩‍👧‍👦 (U+1F468 U+200D U+1F469 U+200D U+1F467 U+200D U+1F466)
- Flag sequences: 🇺🇸 (U+1F1FA U+1F1F8)
- Skin tone sequences: 👋🏿 (U+1F44B U+1F3FF)

**Width Calculation**: Use the width of the base emoji (typically 2)

## Surrogate Pairs

### High Surrogates (U+D800-U+DBFF)
### Low Surrogates (U+DC00-U+DFFF)

**Width**: -1 (invalid when unpaired)

**Valid Surrogate Pairs**: Calculate width of the combined character
**Invalid Surrogates**: Return -1

## Private Use Areas

### Private Use Area (U+E000-U+F8FF)
### Supplementary Private Use Area A (U+F0000-U+FFFFD)
### Supplementary Private Use Area B (U+100000-U+10FFFD)

**Width**: 1 (default assumption)

## Format Characters

### Soft Hyphen (U+00AD)
**Width**: 1 (visible when line breaks)

### Word Joiner (U+2060)
**Width**: 0 (invisible)

### Function Application (U+2061)
**Width**: 0 (invisible, mathematical)

## Compatibility Considerations

### go-runewidth Behavior

1. **Ambiguous Width**: Defaults to width 1
2. **Emoji**: Treats most emoji as width 2
3. **Control Characters**: Returns 0 for most control characters
4. **Combining Marks**: Returns 0 for all combining marks

### wcwidth Behavior

1. **Ambiguous Width**: Defaults to width 1
2. **Emoji**: May not handle emoji correctly (depends on implementation)
3. **Control Characters**: Returns 0 for control characters
4. **Combining Marks**: Returns 0 for combining marks

## Test Cases

### Basic Width Tests

```go
// Control characters
assertWidth("", 0)           // Empty string
assertWidth("\x00", 0)       // NULL
assertWidth("\x08", 0)       // Backspace
assertWidth("\x09", 0)       // Tab
assertWidth("\x0A", 0)       // Line Feed
assertWidth("\x0D", 0)       // Carriage Return
assertWidth("\x7F", 0)       // DEL

// ASCII printable
assertWidth("A", 1)          // Latin A
assertWidth("!", 1)          // Exclamation mark
assertWidth(" ", 1)          // Space

// Combining marks
assertWidth("é", 1)          // e + combining acute
assertWidth("c\u0327", 1)    // c + combining cedilla
assertWidth("\u0300", 0)     // Combining grave accent

// East Asian characters
assertWidth("中", 2)          // Chinese character (Wide)
assertWidth("あ", 2)          // Hiragana (Wide)
assertWidth("Ａ", 2)          // Fullwidth A (Fullwidth)

// Emoji
assertWidth("😀", 2)         // Emoji base
assertWidth("👋🏿", 2)        // Emoji with skin tone modifier
assertWidth("👨‍👩‍👧‍👦", 2)    // Family emoji (ZWJ sequence)
```

### Edge Case Tests

```go
// Ambiguous characters
assertWidth("¡", 1)          // Inverted exclamation (Ambiguous)
assertWidth("°", 1)          // Degree sign (Ambiguous)
assertWidth("±", 1)          // Plus-minus (Ambiguous)

// Zero-width characters
assertWidth("\u200B", 0)     // Zero width space
assertWidth("\u200D", 0)     // Zero width joiner
assertWidth("\uFEFF", 0)     // Zero width no-break space

// Surrogate pairs
assertWidth("\uD83D\uDE00", 2)  // Valid surrogate pair (😀)
assertWidth("\uD83D", -1)       // Unpaired high surrogate
assertWidth("\uDE00", -1)       // Unpaired low surrogate

// Invalid UTF-8
assertWidth("\xFF", -1)      // Invalid UTF-8 sequence
assertWidth("\xC0\x80", -1)  // Overlong encoding
```

## Performance Considerations

### Large String Handling

1. **Memory Usage**: Process strings in chunks to avoid large memory allocations
2. **Cache**: Cache width lookups for frequently used characters
3. **SIMD**: Use vectorized operations for bulk processing

### Table Lookup Optimization

1. **Trie Structure**: Use compressed trie for efficient range lookups
2. **Binary Search**: Use binary search for large ranges
3. **Hash Tables**: Use hash tables for frequently accessed characters

## Error Handling

### Invalid Input

1. **Invalid UTF-8**: Return -1 for malformed sequences
2. **Unmapped Characters**: Default to width 1
3. **Surrogate Pairs**: Return -1 for unpaired surrogates

### Edge Case Handling

1. **Empty Strings**: Return 0
2. **Null Characters**: Return 0
3. **Very Long Strings**: Handle without memory issues
4. **Mixed Encodings**: Detect and handle encoding issues
