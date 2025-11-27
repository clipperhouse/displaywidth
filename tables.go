package displaywidth

// propertyWidths is a jump table of sorts, instead of a switch
var propertyWidths = [5]int{
	_Default:              1,
	_Zero_Width:           0,
	_East_Asian_Wide:      2,
	_East_Asian_Ambiguous: 1,
	_Emoji:                2,
}

// asciiWidths is a bitmask using 4 uint64's.
//
// Layout:
//   - asciiWidths[0]: bits 0-63   (ASCII 0x00-0x3F)
//   - asciiWidths[1]: bits 64-127 (ASCII 0x40-0x7F)
//   - asciiWidths[2]: bits 128-191 (ASCII 0x80-0xBF)
//   - asciiWidths[3]: bits 192-255 (ASCII 0xC0-0xFF)
var asciiWidths = [4]uint64{
	// Mask 0: 0x00-0x3F
	// 0x00-0x1F: all 0 (control characters)
	// 0x20-0x3F: all 1 (printable ASCII)
	0xFFFFFFFF00000000,
	// Mask 1: 0x40-0x7F
	// 0x40-0x7E: all 1 (printable ASCII)
	// 0x7F: 0 (DEL)
	0x7FFFFFFFFFFFFFFF,

	// >= 128 means you should not be using this table, because valid
	// single-byte UTF-8 is < 128. We will return a default value of
	// _Default in those cases, so as not to panic.

	// Mask 2: 0x80-0xBF
	// All 1 (>= 128)
	0xFFFFFFFFFFFFFFFF, // all bits set
	// Mask 3: 0xC0-0xFF
	// All 1 (>= 128)
	0xFFFFFFFFFFFFFFFF, // all bits set
}

// asciiWidth returns the width for a byte
func asciiWidth(b byte) int {
	// determine the uint64 mask for the byte
	mask := asciiWidths[b>>6]
	// determine which bit within the uint64 to use
	pos := b & 0x3F
	return int((mask >> pos) & 1)
}

// asciiProperty returns the property for a byte
func asciiProperty(b byte) property {
	// We can reuse (invert) asciiWidth because _Default happens to be 0,
	// and _Zero_Width happens to be 1.

	// determine the uint64 mask for the byte
	mask := asciiWidths[b>>6]
	// determine which bit within the uint64 to use
	pos := b & 0x3F
	// invert the mask and extract the bit
	return property((^mask >> pos) & 1)
}

var (
	// asciiProperty depends on _Default being 0 and _Zero_Width being 1.
	// Some compile-time checks.

	// If _Default != 0, out of bounds.
	_ = [1]int{}[_Default]

	// If _Zero_Width is 0, index is -1, out of bounds
	// If _Zero_Width is 1, index is 0, correct
	// If _Zero_Width is > 1, out of bounds
	_ = [1]int{}[_Zero_Width-1]
)
