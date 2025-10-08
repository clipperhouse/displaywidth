# String Width Measurement Specification

## Overview

This document defines the specification for measuring the display (monospace column) width of strings, compatible with existing implementations like `go-runewidth` and `wcwidth`. The specification is based on Unicode standards and terminal rendering conventions.

## Core Principles

1. **Byte-level processing**: Operate strictly on strings and `[]byte` without decoding runes
2. **Unicode compliance**: Follow Unicode East Asian Width (UAX #11) and POSIX wcwidth standards
3. **Terminal compatibility**: Match behavior of existing terminal width measurement libraries
4. **Performance**: Use code-generated tries for efficient table lookups

## Width Classification

### Return Values

Characters are classified into the following width categories:

- **Width 0**: Non-printing characters that don't advance the cursor
- **Width 1**: Standard single-cell characters
- **Width 2**: Double-wide characters
- **Width -1**: Invalid or non-printable characters (implementation specific)

### Character Categories

#### 1. Control Characters (Width 0)

**C0 Controls (0x00-0x1F)**:
- NULL (0x00): Width 0
- Backspace (0x08): Width 0
- Tab (0x09): Width 0
- Line Feed (0x0A): Width 0
- Carriage Return (0x0D): Width 0
- All other C0 controls: Width 0

**DEL (0x7F)**: Width 0

**C1 Controls (0x80-0x9F)**: Width 0

#### 2. Combining Characters (Width 0)

**Unicode Categories**:
- **Mn (Nonspacing_Mark)**: Width 0
- **Mc (Spacing_Mark)**: Width 0
- **Me (Enclosing_Mark)**: Width 0

**Special Zero-Width Characters**:
- Zero Width Space (U+200B): Width 0
- Zero Width Joiner (U+200D): Width 0
- Zero Width Non-Joiner (U+200C): Width 0
- Zero Width No-Break Space (U+FEFF): Width 0

#### 3. East Asian Width Property

Based on Unicode UAX #11 (East Asian Width):

- **F (Fullwidth)**: Width 2
  - CJK fullwidth characters
  - Fullwidth ASCII variants
  - Fullwidth punctuation

- **W (Wide)**: Width 2
  - CJK ideographs
  - CJK symbols and punctuation
  - Some emoji

- **H (Halfwidth)**: Width 1
  - Halfwidth CJK characters
  - Halfwidth katakana

- **Na (Narrow)**: Width 1
  - ASCII characters
  - Most Latin characters
  - Most symbols

- **N (Neutral)**: Width 1
  - Characters with neutral width
  - Default to single width

- **A (Ambiguous)**: Width 1 (default) or Width 2 (East Asian context)
  - Characters that can be either narrow or wide
  - Context-dependent width
  - Default to narrow (width 1) for Western contexts

#### 4. Special Cases

**Format Characters**:
- Soft Hyphen (U+00AD): Width 1
- Word Joiner (U+2060): Width 0
- Function Application (U+2061): Width 0

**Emoji**:
- Emoji base characters: Generally Width 2
- Emoji modifiers (skin tone): Width 0
- Emoji variation selectors: Width 0
- ZWJ sequences: Combined width calculation needed

**Private Use and Surrogates**:
- Private Use Area: Width 1 (default)
- Surrogate pairs: Width -1 (invalid)

## String Width Calculation Algorithm

### Basic Algorithm

1. **Initialize**: Set total width to 0
2. **Iterate**: Process each UTF-8 byte sequence
3. **Decode**: Extract Unicode code point from UTF-8
4. **Classify**: Determine width based on character properties
5. **Accumulate**: Add width to total
6. **Return**: Total width

### Byte-Level Processing

Following the uax29 stringish approach, the implementation should:

1. **UTF-8 Validation**: Ensure valid UTF-8 sequences
2. **Code Point Extraction**: Extract Unicode code points from UTF-8 bytes without full rune decoding
3. **Trie Lookup**: Use compressed trie for efficient character property lookup
4. **Stringish Processing**: Process strings and bytes directly, similar to uax29's stringish package
5. **No Rune Exposure**: Keep rune processing internal to the package

### Special Sequence Handling

**Grapheme Clusters**:
- Base character + combining marks: Use base character width
- Emoji sequences: Handle ZWJ sequences specially

**Emoji ZWJ Sequences**:
- Zero Width Joiner (U+200D) connects emoji
- Calculate combined width of the sequence
- Default to width 2 for emoji sequences

## Data Sources

### Unicode Character Database Files

1. **EastAsianWidth.txt**: Primary source for East Asian Width property
2. **UnicodeData.txt**: General category and combining mark information
3. **emoji-data.txt**: Emoji properties and sequences
4. **DerivedGeneralCategory.txt**: Derived character categories

### Trie-Based Table Generation

Following the uax29 approach, the implementation will generate a compressed trie containing only non-default codepoints:

1. **Trie Structure**: Compressed trie with bitmap/enum values for character properties
2. **Sparse Data**: Only include codepoints that differ from default behavior
3. **Bitmap Properties**: Use bit flags for multiple properties per codepoint
4. **Efficient Lookup**: Fast trie traversal for character property lookup

#### Trie Property Bitmap

```go
type CharProperties uint8

const (
    // East Asian Width properties
    EAW_Fullwidth CharProperties = 1 << iota  // F
    EAW_Wide                                  // W
    EAW_Halfwidth                             // H
    EAW_Narrow                                // Na
    EAW_Neutral                               // N
    EAW_Ambiguous                             // A

    // General categories
    IsCombiningMark                           // Mn, Mc, Me
    IsControlChar                             // C0, C1, DEL
    IsZeroWidth                               // ZWSP, ZWJ, ZWNJ, etc.
    IsEmoji                                   // Emoji base characters
    IsEmojiModifier                           // Emoji modifiers
    IsEmojiVariationSelector                  // Emoji variation selectors
)
```

#### Default Behavior (Not in Trie)

- **Default Width**: 1 (for most characters)
- **Default EAW**: Neutral (N) - width 1
- **Default Emoji**: Width 2
- **Default Combining**: Width 0
- **Default Control**: Width 0

## Stringish Processing

Following the [uax29 stringish package](https://github.com/clipperhouse/uax29/tree/master/internal/stringish) approach:

- **String Interface**: Provide a unified interface for processing both strings and byte slices
- **Internal Rune Processing**: Keep rune extraction and processing internal to the package
- **No Rune Exposure**: Public API only exposes string and []byte types
- **Efficient Iteration**: Process UTF-8 sequences without full rune decoding overhead
- **Byte-Level Control**: Maintain control over UTF-8 byte sequences for performance

### Stringish Implementation

Following the [uax29 graphemes trie.go](https://github.com/clipperhouse/uax29/blob/master/graphemes/trie.go) pattern:

```go
// internal/stringish/stringish.go
type Stringish interface {
    String() string
    Bytes() []byte
}

// Process strings and bytes without exposing runes
func ProcessString(s string, fn func(rune) bool) bool
func ProcessBytes(b []byte, fn func(rune) bool) bool

// Trie lookup using stringish interface
func (t *CompressedTrie) LookupStringish(s Stringish) CharProperties
```

## Configuration Options

### API Configuration

The package will accept two boolean parameters to match go-runewidth behavior:

```go
func StringWidth(s string, eastAsianWidth bool, strictEmojiNeutral bool) int
func StringWidthBytes(s []byte, eastAsianWidth bool, strictEmojiNeutral bool) int
```

**Note**: No `RuneWidth` function in the public API - all processing is done at the string/byte level internally.

### Configuration Parameters

- **eastAsianWidth**: When true, treat ambiguous width characters as wide (width 2)
- **strictEmojiNeutral**: When true, use strict emoji width calculation (some emoji become width 1)

### Default Behavior

- **eastAsianWidth = false**: Ambiguous characters default to width 1 (narrow)
- **strictEmojiNeutral = false**: Emoji use default width calculation (typically width 2)

## Compatibility Requirements

### go-runewidth Compatibility

- Match `runewidth.StringWidth(s, eastAsianWidth, strictEmojiNeutral)` behavior
- Handle ambiguous width characters consistently with `eastAsianWidth` parameter
- Support same emoji handling with `strictEmojiNeutral` parameter

### wcwidth Compatibility

- Match POSIX `wcwidth()` function behavior (equivalent to `eastAsianWidth=false, strictEmojiNeutral=false`)
- Return same values for equivalent characters
- Handle control characters consistently

## Implementation Notes

### Performance Considerations

1. **Trie Structure**: Use compressed trie for efficient lookups
2. **Range Tables**: Use binary search for large ranges
3. **Cache**: Cache frequently accessed width values
4. **SIMD**: Consider SIMD optimizations for bulk processing

### Error Handling

1. **Invalid UTF-8**: Return -1 for invalid sequences
2. **Unmapped Characters**: Default to width 1
3. **Surrogate Pairs**: Return -1 for unpaired surrogates

### Testing Requirements

1. **Unicode Test Cases**: Test all East Asian Width categories
2. **Edge Cases**: Test control characters, combining marks, emoji
3. **Compatibility Tests**: Verify behavior matches go-runewidth and wcwidth
4. **Performance Tests**: Benchmark against existing implementations

## Version Compatibility

### Unicode Version

- **Target**: Unicode 15.0 (latest stable)
- **Update Strategy**: Regular updates with new Unicode releases
- **Backward Compatibility**: Maintain compatibility with older Unicode versions

### Library Versions

- **go-runewidth**: Compatible with v0.0.15+
- **wcwidth**: Compatible with v0.2.6+
- **Testing**: Continuous integration with latest versions

## Future Considerations

### Extensibility

1. **Custom Width Tables**: Allow custom width definitions
2. **Plugin Architecture**: Support for additional character sets
3. **Runtime Configuration**: Dynamic width behavior changes

### Performance Improvements

1. **JIT Compilation**: Compile width tables at runtime
2. **Hardware Acceleration**: Use CPU-specific optimizations
3. **Parallel Processing**: Multi-threaded width calculation

## References

- [Unicode UAX #11: East Asian Width](https://unicode.org/reports/tr11/)
- [POSIX wcwidth specification](https://pubs.opengroup.org/onlinepubs/9699919799/functions/wcwidth.html)
- [go-runewidth implementation](https://github.com/mattn/go-runewidth)
- [wcwidth Python library](https://github.com/jquast/wcwidth)
- [Unicode Character Database](https://unicode.org/ucd/)
