package displaywidth

import "testing"

var csOptions = Options{ControlSequences: true}

// Inputs for benchmarking truncation with trailing escape sequence preservation
var (
	// Short colored text with reset
	shortANSI = "\x1b[31mhello world\x1b[0m"
	// Multiple stacked SGR sequences
	stackedANSI = "\x1b[1m\x1b[31m\x1b[42mhello world, this is some longer text\x1b[0m"
	// Many interleaved color changes
	interleavedANSI = "hello \x1b[31mworld \x1b[32mfoo \x1b[33mbar \x1b[34mbaz \x1b[35mqux \x1b[36mend\x1b[0m"
	// Plain text (no escape sequences) — baseline
	plainText = "hello world, this is some plain text without escapes"
)

func BenchmarkTruncateString(b *testing.B) {
	benchmarks := []struct {
		name    string
		input   string
		options Options
	}{
		{"plain/default", plainText, defaultOptions},
		{"plain/ControlSequences", plainText, csOptions},
		{"short_ANSI/default", shortANSI, defaultOptions},
		{"short_ANSI/ControlSequences", shortANSI, csOptions},
		{"stacked_ANSI/ControlSequences", stackedANSI, csOptions},
		{"interleaved_ANSI/ControlSequences", interleavedANSI, csOptions},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = bm.options.TruncateString(bm.input, 5, "...")
			}
		})
	}
}

var tail = []byte("...")

func BenchmarkTruncateBytes(b *testing.B) {
	benchmarks := []struct {
		name    string
		input   []byte
		options Options
	}{
		{"plain/default", []byte(plainText), defaultOptions},
		{"plain/ControlSequences", []byte(plainText), csOptions},
		{"short_ANSI/default", []byte(shortANSI), defaultOptions},
		{"short_ANSI/ControlSequences", []byte(shortANSI), csOptions},
		{"stacked_ANSI/ControlSequences", []byte(stackedANSI), csOptions},
		{"interleaved_ANSI/ControlSequences", []byte(interleavedANSI), csOptions},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = bm.options.TruncateBytes(bm.input, 5, tail)
			}
		})
	}
}
