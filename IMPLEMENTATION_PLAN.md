# Implementation Plan

## Overview

This document outlines the implementation plan for the stringwidth package, following the approach used in the uax29 package with code-generated tries for table lookups.

## Architecture

### Core Components

1. **Width Calculator**: Main API for calculating string widths with boolean parameters
2. **Trie Generator**: Code generation for compressed trie with sparse character properties
3. **Trie Structure**: Efficient lookup structure for character properties (similar to uax29)
4. **Test Suite**: Comprehensive test coverage with go-runewidth compatibility tests

### Package Structure

```
stringwidth/
├── internal/
│   ├── gen/               # Code generation tools (following uax29 pattern)
│   │   ├── main.go        # Trie generation entry point
│   │   ├── unicode.go     # Unicode data parsing
│   │   └── trie.go        # Compressed trie generation
│   ├── trie/              # Generated trie structure
│   │   ├── trie.go        # Compressed trie with sparse properties
│   │   └── properties.go  # Character property bitmap definitions
│   ├── stringish/         # String/byte processing utilities
│   │   ├── stringish.go   # String and byte processing (similar to uax29)
│   │   └── utf8.go        # UTF-8 processing utilities
│   └── width/             # Core width calculation
│       ├── calculator.go  # Main width calculation logic
│       └── lookup.go      # Trie lookup implementation
├── stringwidth.go         # Public API
├── stringwidth_test.go    # Test suite
└── go.mod
```

## Implementation Phases

### Phase 1: Core Infrastructure

#### 1.1 Trie Generation System

**Files**: `internal/gen/` (using golang.org/x/text/internal/triegen)

**Responsibilities**:
- Download and parse Unicode data files
- Generate compressed trie with sparse character properties using triegen
- Create bitmap-based property lookup
- Follow uax29 generation pattern with triegen foundation

**Key Components**:
```go
// internal/gen/main.go
package main

import (
    "golang.org/x/text/internal/triegen"
    "unicode"
)

func main() {
    // Extract range tables from Go stdlib unicode package where helpful
    // Download external Unicode data files for properties not in stdlib
    // Parse EastAsianWidth.txt (not in stdlib)
    // Parse emoji-data.txt (not in stdlib)
    // Build trie using triegen from combined data sources
    // Generate compressed trie with sparse properties
    // Only include codepoints that differ from defaults
    // Output to internal/trie/trie.go
}

// internal/gen/unicode.go
type UnicodeData struct {
    EastAsianWidth map[rune]string  // From external file
    EmojiData map[rune]bool         // From external file
    // Other properties from Go stdlib range tables or external files
}

func ParseEastAsianWidth(filename string) (map[rune]string, error)
func ParseEmojiData(filename string) (map[rune]bool, error)
func ExtractStdlibRangeTables() map[rune]CharProperties  // Extract from unicode package range tables

// internal/gen/trie.go
func GenerateTrie(data UnicodeData) *triegen.Trie
func BuildPropertyBitmap(r rune, data UnicodeData) CharProperties
func WriteTrieGo(trie *triegen.Trie, filename string) error
```

#### 1.2 Trie Structure

**Files**: `internal/trie/trie.go`, `internal/trie/properties.go`

**Responsibilities**:
- Generate compressed trie for character property lookups
- Implement efficient range-based lookups with bitmap properties
- Support sparse data (only non-default codepoints)
- Use stringish interface for lookups

**Key Components**:
```go
// internal/trie/properties.go
type CharProperties uint8

const (
    EAW_Fullwidth CharProperties = 1 << iota  // F
    EAW_Wide                                  // W
    EAW_Halfwidth                             // H
    EAW_Narrow                                // Na
    EAW_Neutral                               // N
    EAW_Ambiguous                             // A
    IsCombiningMark                           // Mn, Mc, Me
    IsControlChar                             // C0, C1, DEL
    IsZeroWidth                               // ZWSP, ZWJ, ZWNJ, etc.
    IsEmoji                                   // Emoji base characters
    IsEmojiModifier                           // Emoji modifiers
    IsEmojiVariationSelector                  // Emoji variation selectors
)

// internal/trie/trie.go
type CompressedTrie struct {
    nodes []TrieNode
}

type TrieNode struct {
    start, end rune
    properties CharProperties
    children   []int
}

// Following uax29 pattern with stringish interface
func (t *CompressedTrie) Lookup(r rune) CharProperties
func (t *CompressedTrie) LookupStringish(s Stringish) CharProperties
func (t *CompressedTrie) LookupRange(start, end rune) CharProperties
```

### Phase 2: Width Calculation Engine

#### 2.1 Stringish Processing

**Files**: `internal/stringish/stringish.go`, `internal/stringish/utf8.go`

**Responsibilities**:
- Process strings and bytes directly (similar to uax29 stringish package)
- Extract code points from UTF-8 without full rune decoding
- Validate UTF-8 sequences
- Handle string/byte iteration efficiently

**Key Components**:
```go
// internal/stringish/stringish.go
type Stringish interface {
    String() string
    Bytes() []byte
}

// internal/stringish/utf8.go
func NextCodePoint(s []byte) (rune, int, error)
func IsValidUTF8(s []byte) bool
func CodePointCount(s []byte) int
func ProcessString(s string, fn func(rune) bool) bool
func ProcessBytes(b []byte, fn func(rune) bool) bool
```

#### 2.2 Width Calculator

**Files**: `internal/width/calculator.go`, `internal/width/lookup.go`

**Responsibilities**:
- Calculate character width based on Unicode properties from trie
- Handle special cases and edge cases
- Support boolean configuration parameters
- Process strings and bytes without exposing runes

**Key Components**:
```go
// internal/width/calculator.go
func StringWidth(s string, eastAsianWidth bool, strictEmojiNeutral bool) int
func StringWidthBytes(s []byte, eastAsianWidth bool, strictEmojiNeutral bool) int

// internal/width/lookup.go
func calculateWidth(r rune, props CharProperties, eastAsianWidth bool, strictEmojiNeutral bool) int
func getDefaultWidth(r rune) int
func isDefaultBehavior(r rune) bool
func processStringWidth(s string, eastAsianWidth bool, strictEmojiNeutral bool) int
func processBytesWidth(b []byte, eastAsianWidth bool, strictEmojiNeutral bool) int
```

### Phase 3: Public API

#### 3.1 Main Package

**Files**: `stringwidth.go`

**Responsibilities**:
- Provide clean public API with boolean parameters
- Match go-runewidth API signature
- Maintain backward compatibility

**Key Components**:
```go
// stringwidth.go
package stringwidth

import "github.com/clipperhouse/stringwidth/internal/width"

// Public API functions - match go-runewidth signature
func StringWidth(s string, eastAsianWidth bool, strictEmojiNeutral bool) int {
    return width.StringWidth(s, eastAsianWidth, strictEmojiNeutral)
}

func StringWidthBytes(s []byte, eastAsianWidth bool, strictEmojiNeutral bool) int {
    return width.StringWidthBytes(s, eastAsianWidth, strictEmojiNeutral)
}

// Convenience functions with default parameters
func StringWidthDefault(s string) int {
    return StringWidth(s, false, false)
}

func StringWidthBytesDefault(s []byte) int {
    return StringWidthBytes(s, false, false)
}
```

### Phase 4: Testing and Validation

#### 4.1 Test Suite

**Files**: `stringwidth_test.go`

**Responsibilities**:
- Comprehensive test coverage
- Compatibility tests with go-runewidth and wcwidth
- Performance benchmarks
- Edge case validation

**Key Components**:
```go
// stringwidth_test.go
package stringwidth

import (
    "testing"
    "github.com/mattn/go-runewidth"
)

func TestStringWidth(t *testing.T) {
    // Basic functionality tests
    // Edge case tests
    // Compatibility tests
}

func TestCompatibility(t *testing.T) {
    // Compare with go-runewidth
    // Compare with wcwidth
}

func BenchmarkStringWidth(b *testing.B) {
    // Performance benchmarks
}

func TestEdgeCases(t *testing.T) {
    // Control characters
    // Combining marks
    // Emoji sequences
    // Surrogate pairs
    // Invalid UTF-8
}
```

## Data Sources and Generation

### Unicode Data Files

1. **EastAsianWidth.txt**: Primary source for character width
2. **UnicodeData.txt**: General category and combining mark information
3. **emoji-data.txt**: Emoji properties and sequences
4. **DerivedGeneralCategory.txt**: Derived character categories

### Generation Process

1. **Evaluate Data Sources**: Check which Unicode categories can be extracted from Go's existing range tables
2. **Extract Stdlib Data**: Extract range table data from Go's `unicode` package where helpful
3. **Download External**: Fetch external Unicode data files for properties not available in stdlib
4. **Parse**: Extract relevant character properties from external files and stdlib range tables
5. **Build Trie**: Use triegen to create compressed trie from combined data sources
6. **Generate Code**: Use triegen.Gen() to output Go code to internal/trie/trie.go
7. **Validate**: Test generated trie against reference data

**Note**: We'll evaluate whether to use Go's existing range tables as data sources for trie generation, or download external Unicode data files. The goal is to save work by leveraging existing range table data where helpful. The triegen package handles the actual trie compression and code generation.

### Generated Trie

Following the [uax29 generation pattern](https://github.com/clipperhouse/uax29/tree/master/internal/gen) using [golang.org/x/text/internal/triegen](https://pkg.go.dev/golang.org/x/text/internal/triegen):

```go
// internal/trie/trie.go (generated by internal/gen using triegen)
package trie

// Generated by triegen - compressed trie containing only non-default codepoints
type stringWidthTrie struct {
    // triegen-generated trie structure
}

func newStringWidthTrie(x int) *stringWidthTrie {
    // triegen-generated constructor
}

func (t *stringWidthTrie) lookup(s []byte) (v uint8, sz int) {
    // triegen-generated lookup method
}

func (t *stringWidthTrie) lookupString(s string) (v uint8, sz int) {
    // triegen-generated string lookup
}

// Generated trie data
var stringWidthValues []uint8
var stringWidthIndex []uint16
var stringWidthTrieHandles []uint16

// Default behavior for codepoints not in trie:
// - Width: 1
// - EAW: Neutral (N)
// - Not combining, not control, not zero-width, not emoji
```

## Performance Optimization

### Trie Optimization

1. **Compression**: Use compressed trie to reduce memory usage
2. **Simple Lookups**: Direct trie traversal only, no range tables or binary search
3. **No Caching**: No cache for frequently accessed values
4. **No SIMD**: No vectorized operations for bulk processing

### Memory Management

1. **Read-only Tables**: Use read-only data structures
2. **Efficient Lookups**: Minimize memory allocations
3. **String Processing**: Process strings in-place when possible
4. **Trie-Only**: All lookups go through the compressed trie

## Compatibility Strategy

### go-runewidth Compatibility

1. **API Compatibility**: Match public API where possible
2. **Behavior Compatibility**: Match width calculation behavior
3. **Configuration Compatibility**: Support same configuration options

### wcwidth Compatibility

1. **Function Compatibility**: Match wcwidth function behavior
2. **Return Values**: Return same values for equivalent characters
3. **Error Handling**: Handle errors consistently

## Testing Strategy

### Unit Tests

1. **Basic Functionality**: Test core width calculation
2. **Edge Cases**: Test control characters, combining marks, emoji
3. **Configuration**: Test different configuration options
4. **Error Handling**: Test invalid input handling

### Integration Tests

1. **Compatibility Tests**: Compare with go-runewidth and wcwidth
2. **Unicode Tests**: Test against Unicode test data
3. **Performance Tests**: Benchmark against existing implementations

### Continuous Integration

1. **Automated Testing**: Run tests on every commit
2. **Unicode Updates**: Test with latest Unicode data
3. **Cross-platform**: Test on multiple platforms
4. **Performance Monitoring**: Track performance regressions

## Release Strategy

### Version 0.1.0 (Initial Release)

- Basic width calculation functionality
- Compatibility with go-runewidth
- Core test suite
- Documentation

### Version 0.2.0 (Performance Release)

- Optimized trie implementation
- Performance improvements
- Extended test coverage
- Benchmark suite

### Version 1.0.0 (Stable Release)

- Full compatibility with go-runewidth and wcwidth
- Complete test coverage
- Performance optimizations
- Production-ready implementation

## Maintenance

### Unicode Updates

1. **Regular Updates**: Update with new Unicode releases
2. **Automated Testing**: Test compatibility with new Unicode data
3. **Backward Compatibility**: Maintain compatibility with older Unicode versions

### Performance Monitoring

1. **Benchmark Tracking**: Monitor performance over time
2. **Memory Usage**: Track memory usage patterns
3. **Trie Optimization**: Focus on trie structure improvements

### Community Feedback

1. **Issue Tracking**: Monitor and respond to issues
2. **Feature Requests**: Evaluate and implement new features
3. **Documentation**: Maintain comprehensive documentation
