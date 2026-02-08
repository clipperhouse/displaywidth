package displaywidth

import (
	"bytes"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/clipperhouse/displaywidth/testdata"
)

// FuzzBytesAndString fuzzes the Bytes function with valid and invalid UTF-8.
func FuzzBytesAndString(f *testing.F) {
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
			{EastAsianWidth: false},
			{EastAsianWidth: true},
			{IgnoreControlSequences: true},
			{EastAsianWidth: true, IgnoreControlSequences: true},
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

			ws := option.String(string(text))
			if wb != ws {
				t.Errorf("Bytes() returned %d but String() returned %d with options %+v for %q", wb, ws, option, text)
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

		// Test with different options (Rune is per-rune, IgnoreControlSequences
		// doesn't affect single runes, but we include it for completeness)
		options := []Options{
			{EastAsianWidth: false},
			{EastAsianWidth: true},
			{IgnoreControlSequences: true},
			{EastAsianWidth: true, IgnoreControlSequences: true},
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

func FuzzTruncateStringAndBytes(f *testing.F) {
	if testing.Short() {
		f.Skip("skipping fuzz test in short mode")
	}

	// Seed with multi-lingual text (paragraph-sized chunks)
	file, err := testdata.Sample()
	if err != nil {
		f.Fatal(err)
	}
	fs := string(file)
	chunks := strings.Split(fs, "\n")
	for _, chunk := range chunks {
		f.Add(chunk)
	}

	// Seed with invalid UTF-8
	invalid, err := testdata.InvalidUTF8()
	if err != nil {
		f.Fatal(err)
	}
	fs = string(invalid)
	chunks = strings.Split(fs, "\n")
	for _, chunk := range chunks {
		f.Add(chunk)
	}

	// Seed with test cases
	testCases, err := testdata.TestCases()
	if err != nil {
		f.Fatal(err)
	}
	fs = string(testCases)
	chunks = strings.Split(fs, "\n")
	for _, chunk := range chunks {
		f.Add(chunk)
	}

	// Seed with random bytes
	for i := 0; i < 10; i++ {
		b, err := testdata.RandomBytes()
		if err != nil {
			f.Fatal(err)
		}
		f.Add(string(b))
	}

	// Seed with edge cases
	f.Add("")             // empty
	f.Add("a")            // single ASCII
	f.Add("\t\n\r")       // whitespace
	f.Add("🌍")            // emoji
	f.Add("\u0301")       // combining mark
	f.Add("\xff\xfe\xfd") // invalid UTF-8

	f.Fuzz(func(t *testing.T, text string) {
		// Test with default options
		ts := TruncateString(text, 10, "...")

		// Invariant: truncated string should be less than or equal to maxWidth
		if String(ts) > 10 {
			t.Errorf("TruncateString() returned string longer than maxWidth for %q: %q", text, ts)
		}

		// Invariant: truncated string should be less than or equal to maxWidth
		if len(ts) > len(text) {
			t.Errorf("TruncateString() returned string longer than original for %q: %q", text, ts)
		}

		tb := TruncateBytes([]byte(text), 10, []byte("..."))

		// Invariant: truncated bytes should be less than or equal to maxWidth
		if Bytes(tb) > 10 {
			t.Errorf("TruncateBytes() returned bytes longer than maxWidth for %q: %q", text, tb)
		}

		// Invariant: truncated bytes should be less than or equal to original
		if len(tb) > len(text) {
			t.Errorf("TruncateBytes() returned bytes longer than original for %q: %q", text, tb)
		}

		if !bytes.Equal(tb, []byte(ts)) {
			t.Errorf("TruncateBytes() returned bytes different from TruncateString() for %q: %q != %q", text, tb, ts)
		}

		// Test with different options
		options := []Options{
			{EastAsianWidth: false},
			{EastAsianWidth: true},
			{IgnoreControlSequences: true},
			{EastAsianWidth: true, IgnoreControlSequences: true},
		}

		for _, option := range options {
			ts := option.TruncateString(text, 10, "...")

			// Invariant: truncated string should be less than or equal to maxWidth
			if option.String(ts) > 10 {
				t.Errorf("TruncateString() returned string longer than maxWidth for %q: %q", text, ts)
			}

			tb := option.TruncateBytes([]byte(text), 10, []byte("..."))
			if !bytes.Equal(tb, []byte(ts)) {
				t.Errorf("TruncateBytes() returned bytes different from TruncateString() for %q: %q != %q", text, tb, ts)
			}
		}
	})
}

// FuzzControlSequences fuzzes strings containing ANSI/ECMA-48 escape sequences
// across all option combinations (EastAsianWidth x IgnoreControlSequences).
func FuzzControlSequences(f *testing.F) {
	if testing.Short() {
		f.Skip("skipping fuzz test in short mode")
	}

	// Seed with ANSI escape sequences
	f.Add([]byte("\x1b[31m"))                                  // SGR red
	f.Add([]byte("\x1b[0m"))                                   // SGR reset
	f.Add([]byte("\x1b[1m"))                                   // SGR bold
	f.Add([]byte("\x1b[38;5;196m"))                            // SGR 256-color
	f.Add([]byte("\x1b[38;2;255;0;0m"))                        // SGR truecolor
	f.Add([]byte("\x1b[A"))                                    // cursor up
	f.Add([]byte("\x1b[10;20H"))                               // cursor position
	f.Add([]byte("\x1b[2J"))                                   // erase in display
	f.Add([]byte("\x1b[31mhello\x1b[0m"))                      // red text
	f.Add([]byte("\x1b[1m\x1b[31mhi\x1b[0m"))                  // nested SGR
	f.Add([]byte("hello\x1b[31mworld\x1b[0m"))                 // ANSI mid-string
	f.Add([]byte("\x1b[31m中文\x1b[0m"))                         // colored CJK
	f.Add([]byte("\x1b[31m😀\x1b[0m"))                          // colored emoji
	f.Add([]byte("\x1b[31m🇺🇸\x1b[0m"))                         // colored flag
	f.Add([]byte("a\x1b[31mb\x1b[32mc\x1b[33md\x1b[0m"))       // multiple colors
	f.Add([]byte("\x1b[31m\x1b[42m\x1b[1mbold on red\x1b[0m")) // stacked SGR
	f.Add([]byte("\r\n"))                                      // CR+LF
	f.Add([]byte("hello\r\nworld"))                            // text with CRLF
	f.Add([]byte("\x1b"))                                      // bare ESC
	f.Add([]byte("\x1b["))                                     // incomplete sequence
	f.Add([]byte("\x1b[31"))                                   // incomplete SGR
	f.Add([]byte(""))                                          // empty
	f.Add([]byte("hello"))                                     // plain ASCII
	f.Add([]byte("中文"))                                        // plain CJK
	f.Add([]byte("😀"))                                         // plain emoji

	// Seed with multi-lingual text
	file, err := testdata.Sample()
	if err != nil {
		f.Fatal(err)
	}
	chunks := bytes.Split(file, []byte("\n"))
	for _, chunk := range chunks {
		f.Add(chunk)
	}

	allOptions := []Options{
		{},
		{EastAsianWidth: true},
		{IgnoreControlSequences: true},
		{EastAsianWidth: true, IgnoreControlSequences: true},
	}

	f.Fuzz(func(t *testing.T, text []byte) {
		for _, opt := range allOptions {
			wb := opt.Bytes(text)
			ws := opt.String(string(text))

			// Invariant: width is never negative
			if wb < 0 {
				t.Errorf("Bytes() with %+v returned negative width %d for %q", opt, wb, text)
			}

			// Invariant: String and Bytes agree
			if wb != ws {
				t.Errorf("Bytes()=%d != String()=%d with %+v for %q", wb, ws, opt, text)
			}

			// Invariant: empty input is always 0
			if len(text) == 0 && wb != 0 {
				t.Errorf("non-zero width %d for empty input with %+v", wb, opt)
			}

			// Invariant: sum of grapheme widths equals total width
			gIter := opt.BytesGraphemes(text)
			gSum := 0
			for gIter.Next() {
				gw := gIter.Width()
				if gw < 0 {
					t.Errorf("grapheme Width() < 0 with %+v for %q", opt, text)
				}
				gSum += gw
			}
			if gSum != wb {
				t.Errorf("sum of grapheme widths %d != Bytes() %d with %+v for %q", gSum, wb, opt, text)
			}

			// Same for StringGraphemes
			sgIter := opt.StringGraphemes(string(text))
			sgSum := 0
			for sgIter.Next() {
				sgSum += sgIter.Width()
			}
			if sgSum != ws {
				t.Errorf("sum of StringGraphemes widths %d != String() %d with %+v for %q", sgSum, ws, opt, text)
			}

			// Invariant: IgnoreControlSequences width <= default width
			// (escape sequences become 0 instead of their visible char widths)
			if opt.IgnoreControlSequences {
				noIgnore := Options{EastAsianWidth: opt.EastAsianWidth}
				wDefault := noIgnore.Bytes(text)
				if wb > wDefault {
					t.Errorf("IgnoreControlSequences width %d > default width %d with %+v for %q", wb, wDefault, opt, text)
				}
			}

			// Invariant: truncation respects maxWidth (accounting for the tail,
			// which is always appended and may itself exceed maxWidth)
			tail := "..."
			tailWidth := opt.String(tail)
			for _, maxWidth := range []int{0, 1, 3, 5, 10, 20} {
				ts := opt.TruncateString(string(text), maxWidth, tail)
				tsWidth := opt.String(ts)
				limit := maxWidth
				if tailWidth > limit {
					limit = tailWidth
				}
				if tsWidth > limit {
					t.Errorf("TruncateString() width %d > max(maxWidth, tailWidth) %d with %+v for %q -> %q",
						tsWidth, limit, opt, text, ts)
				}

				tb := opt.TruncateBytes(text, maxWidth, []byte(tail))
				if !bytes.Equal(tb, []byte(ts)) {
					t.Errorf("TruncateBytes() != TruncateString() with %+v for %q: %q != %q",
						opt, text, tb, ts)
				}
			}
		}
	})
}
