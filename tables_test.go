package displaywidth

import (
	"testing"
)

// TestAsciiWidth verifies the bitmask values for specific
// control and printable ASCII characters.
func TestAsciiWidth(t *testing.T) {
	tests := []struct {
		name     string
		b        byte
		expected int
		desc     string
	}{
		// Control characters (0x00-0x1F): width 0
		{"null", 0x00, 0, "NULL character"},
		{"bell", 0x07, 0, "BEL (bell)"},
		{"backspace", 0x08, 0, "BS (backspace)"},
		{"tab", 0x09, 0, "TAB"},
		{"newline", 0x0A, 0, "LF (newline)"},
		{"carriage return", 0x0D, 0, "CR (carriage return)"},
		{"escape", 0x1B, 0, "ESC (escape)"},
		{"last control", 0x1F, 0, "Last control character"},

		// Printable ASCII (0x20-0x7E): width 1
		{"space", 0x20, 1, "Space (first printable)"},
		{"exclamation", 0x21, 1, "!"},
		{"zero", 0x30, 1, "0"},
		{"nine", 0x39, 1, "9"},
		{"A", 0x41, 1, "A"},
		{"Z", 0x5A, 1, "Z"},
		{"a", 0x61, 1, "a"},
		{"z", 0x7A, 1, "z"},
		{"tilde", 0x7E, 1, "~ (last printable)"},

		// DEL (0x7F): width 0
		{"delete", 0x7F, 0, "DEL (delete)"},

		// >= 128: width 1 (default, though shouldn't be used for valid UTF-8)
		{"0x80", 0x80, 1, "First byte >= 128"},
		{"0xFF", 0xFF, 1, "Last byte value"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := asciiWidth(tt.b)
			if got != tt.expected {
				t.Errorf("asciiWidth(0x%02X '%s') = %d, want %d (%s)",
					tt.b, string(tt.b), got, tt.expected, tt.desc)
			}
		})
	}
}

// TestAsciiProperty verifies the bitmask values for specific
// control and printable ASCII characters.
func TestAsciiProperty(t *testing.T) {
	tests := []struct {
		name     string
		b        byte
		expected property
		desc     string
	}{
		// Control characters (0x00-0x1F): _Zero_Width (1)
		{"null", 0x00, _Zero_Width, "NULL character"},
		{"bell", 0x07, _Zero_Width, "BEL (bell)"},
		{"backspace", 0x08, _Zero_Width, "BS (backspace)"},
		{"tab", 0x09, _Zero_Width, "TAB"},
		{"newline", 0x0A, _Zero_Width, "LF (newline)"},
		{"carriage return", 0x0D, _Zero_Width, "CR (carriage return)"},
		{"escape", 0x1B, _Zero_Width, "ESC (escape)"},
		{"last control", 0x1F, _Zero_Width, "Last control character"},

		// Printable ASCII (0x20-0x7E): _Default (0)
		{"space", 0x20, _Default, "Space (first printable)"},
		{"exclamation", 0x21, _Default, "!"},
		{"zero", 0x30, _Default, "0"},
		{"nine", 0x39, _Default, "9"},
		{"A", 0x41, _Default, "A"},
		{"Z", 0x5A, _Default, "Z"},
		{"a", 0x61, _Default, "a"},
		{"z", 0x7A, _Default, "z"},
		{"tilde", 0x7E, _Default, "~ (last printable)"},

		// DEL (0x7F): _Zero_Width (1)
		{"delete", 0x7F, _Zero_Width, "DEL (delete)"},

		// >= 128: _Default (0) (default, though shouldn't be used for valid UTF-8)
		{"0x80", 0x80, _Default, "First byte >= 128"},
		{"0xFF", 0xFF, _Default, "Last byte value"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := asciiProperty(tt.b)
			if got != tt.expected {
				t.Errorf("asciiProperty(0x%02X '%s') = %d, want %d (%s)",
					tt.b, string(tt.b), got, tt.expected, tt.desc)
			}
		})
	}
}

func TestAsciiPropertyEnums(t *testing.T) {
	// We need _Default to be 0, and _Zero_Width to be 1, in order for
	// asciiProperty to work as the inverse of asciiWidth.
	if _Default != 0 {
		t.Errorf("_Default = %d, want 0", _Default)
	}
	if _Zero_Width != 1 {
		t.Errorf("_Zero_Width = %d, want 1", _Zero_Width)
	}
}
