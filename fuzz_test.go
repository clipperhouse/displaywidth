package displaywidth

import (
	"bytes"
	"testing"
	"unicode/utf8"

	"github.com/clipperhouse/displaywidth/testdata"
)

// FuzzBytes fuzzes the Bytes function with valid and invalid UTF-8.
func FuzzBytes(f *testing.F) {
	if testing.Short() {
		f.Skip("skipping fuzz test in short mode")
	}

	// Seed with multi-lingual text (paragraph-sized chunks)
	file, err := testdata.Sample()
	if err != nil {
		f.Fatal(err)
	}
	chunks := bytes.Split(file, []byte("\n"))
	for _, chunk := range chunks {
		f.Add(chunk)
	}

	// Seed with invalid UTF-8
	invalid, err := testdata.InvalidUTF8()
	if err != nil {
		f.Fatal(err)
	}
	chunks = bytes.Split(invalid, []byte("\n"))
	for _, chunk := range chunks {
		f.Add(chunk)
	}

	// Seed with test cases
	testCases, err := testdata.TestCases()
	if err != nil {
		f.Fatal(err)
	}
	chunks = bytes.Split(testCases, []byte("\n"))
	for _, chunk := range chunks {
		f.Add(chunk)
	}

	// Seed with random bytes
	for i := 0; i < 10; i++ {
		b, err := testdata.RandomBytes()
		if err != nil {
			f.Fatal(err)
		}
		f.Add(b)
	}

	// Seed with edge cases
	f.Add([]byte(""))               // empty
	f.Add([]byte("a"))              // single ASCII
	f.Add([]byte("\x00"))           // null byte
	f.Add([]byte("\t\n\r"))         // whitespace
	f.Add([]byte("🌍"))              // emoji
	f.Add([]byte("\u0301"))         // combining mark
	f.Add([]byte{0xff, 0xfe, 0xfd}) // invalid UTF-8

	f.Fuzz(func(t *testing.T, text []byte) {
		// Test with default options
		wb := Bytes(text)

		// Invariant: width should never be negative
		if wb < 0 {
			t.Errorf("Bytes() returned negative width for %q: %d", text, wb)
		}

		// Invariant: empty input should always return 0
		if len(text) == 0 && wb != 0 {
			t.Errorf("Bytes() returned non-zero width %d for empty input", wb)
		}

		// Invariant: for valid UTF-8, width should never exceed input length
		// (each byte is at most 1 column wide, some are 0, some multi-byte chars are 2)
		if utf8.Valid(text) {
			runeCount := utf8.RuneCount(text)
			if wb > len(text) {
				t.Errorf("Bytes() width %d exceeds byte length %d for valid UTF-8: %q", wb, len(text), text)
			}

			// Also shouldn't exceed rune count * 2 (max width per rune is 2)
			if wb > runeCount*2 {
				t.Errorf("Bytes() width %d exceeds rune count * 2 (%d) for %q", wb, runeCount*2, text)
			}

			// Consistency check: String() and Bytes() should agree on valid UTF-8
			ws := String(string(text))
			if wb != ws {
				t.Errorf("Bytes() returned %d but String() returned %d for %q", wb, ws, text)
			}
		}

		// Test with different options combinations
		options := []Options{
			{EastAsianWidth: false}, // default
			{EastAsianWidth: true},
		}

		for _, option := range options {
			wb := option.Bytes(text)

			// Same invariants apply
			if wb < 0 {
				t.Errorf("Bytes() with options %+v returned negative width for %q: %d", option, text, wb)
			}

			if len(text) == 0 && wb != 0 {
				t.Errorf("Bytes() with options %+v returned non-zero width %d for empty input", option, wb)
			}

			// Consistency check with String() for valid UTF-8
			if utf8.Valid(text) {
				ws := option.String(string(text))
				if wb != ws {
					t.Errorf("Bytes() returned %d but String() returned %d with options %+v for %q", wb, ws, option, text)
				}
			}
		}
	})
}

// FuzzRune fuzzes the Rune function.
func FuzzRune(f *testing.F) {
	if testing.Short() {
		f.Skip("skipping fuzz test in short mode")
	}

	// Seed with interesting runes
	seeds := []rune{
		0,        // null
		' ',      // space
		'A',      // ASCII
		'\t',     // tab
		'\n',     // newline
		'\u0000', // null
		'\u0301', // combining acute accent
		'\u00A0', // non-breaking space
		'\u2028', // line separator
		'\u2029', // paragraph separator
		'\uFEFF', // zero-width no-break space
		'\uFFFD', // replacement character
		'\uFFFE', // noncharacter
		'\uFFFF', // noncharacter
		'世',      // CJK
		'界',      // CJK
		'🌍',      // emoji
		'👨',      // emoji
		0xD800,   // surrogate (invalid)
		0xDFFF,   // surrogate (invalid)
		0x10FFFF, // max valid rune
	}

	for _, r := range seeds {
		f.Add(r)
	}

	f.Fuzz(func(t *testing.T, r rune) {
		// Test with default options
		wr := Rune(r)

		// Invariant: width should never be negative
		if wr < 0 {
			t.Errorf("Rune() returned negative width for %U (%c): %d", r, r, wr)
		}

		// Invariant: width should be 0, 1, or 2
		if wr > 2 {
			t.Errorf("Rune() returned invalid width for %U (%c): %d (expected 0, 1, or 2)", r, r, wr)
		}

		// Consistency check: compare with Bytes/String for valid runes
		if utf8.ValidRune(r) {
			var buf [4]byte
			n := utf8.EncodeRune(buf[:], r)

			wb := Bytes(buf[:n])
			if wr != wb {
				t.Errorf("Rune() returned %d but Bytes() returned %d for %U (%c)", wr, wb, r, r)
			}

			ws := String(string(r))
			if wr != ws {
				t.Errorf("Rune() returned %d but String() returned %d for %U (%c)", wr, ws, r, r)
			}
		}

		// Test with different options
		options := []Options{
			{EastAsianWidth: false}, // default
			{EastAsianWidth: true},
		}

		for _, option := range options {
			wr := option.Rune(r)

			// Same invariants apply
			if wr < 0 || wr > 2 {
				t.Errorf("Rune() with options %+v returned invalid width for %U (%c): %d", option, r, r, wr)
			}

			// Consistency check with Bytes/String for valid runes
			if utf8.ValidRune(r) {
				var buf [4]byte
				n := utf8.EncodeRune(buf[:], r)

				wb := option.Bytes(buf[:n])
				if wr != wb {
					t.Errorf("Rune() returned %d but Bytes() returned %d with options %+v for %U (%c)", wr, wb, option, r, r)
				}

				ws := option.String(string(r))
				if wr != ws {
					t.Errorf("Rune() returned %d but String() returned %d with options %+v for %U (%c)", wr, ws, option, r, r)
				}
			}
		}
	})
}
