package displaywidth

import (
	"strings"

	"github.com/clipperhouse/uax29/v2/graphemes"
)

// TruncateString truncates a string to the given maxWidth, and appends the
// given tail if the string is truncated.
//
// It ensures the visible width, including the width of the tail, is less than or
// equal to maxWidth.
//
// When [Options.ControlSequences] is true, ANSI escape sequences that appear
// after the truncation point are preserved in the output. This ensures that
// escape sequences such as SGR resets are not lost, preventing color bleed
// in terminal output.
func (options Options) TruncateString(s string, maxWidth int, tail string) string {
	maxWidthWithoutTail := maxWidth - options.String(tail)

	var pos, total int
	g := graphemes.FromString(s)
	g.AnsiEscapeSequences = options.ControlSequences
	g.AnsiEscapeSequences8Bit = options.ControlSequences8Bit

	for g.Next() {
		gw := graphemeWidth(g.Value(), options)
		if total+gw <= maxWidthWithoutTail {
			pos = g.End()
		}
		total += gw
		if total > maxWidth {
			if options.ControlSequences || options.ControlSequences8Bit {
				// Build result with trailing ANSI escape sequences preserved
				var b strings.Builder
				b.Grow(len(s) + len(tail)) // at most original + tail
				b.WriteString(s[:pos])
				b.WriteString(tail)

				rem := graphemes.FromString(s[pos:])
				rem.AnsiEscapeSequences = options.ControlSequences
				rem.AnsiEscapeSequences8Bit = options.ControlSequences8Bit

				for rem.Next() {
					v := rem.Value()
					// Only preserve escapes that measure as zero-width
					// on their own; some sequences (e.g. SOS) are only
					// valid in their original context.
					if len(v) > 0 && isEscapeLeader(v[0], options) && options.String(v) == 0 {
						b.WriteString(v)
					}
				}
				return b.String()
			}
			return s[:pos] + tail
		}
	}
	// No truncation
	return s
}

// TruncateString truncates a string to the given maxWidth, and appends the
// given tail if the string is truncated.
//
// It ensures the total width, including the width of the tail, is less than or
// equal to maxWidth.
func TruncateString(s string, maxWidth int, tail string) string {
	return DefaultOptions.TruncateString(s, maxWidth, tail)
}

// TruncateBytes truncates a []byte to the given maxWidth, and appends the
// given tail if the []byte is truncated.
//
// It ensures the visible width, including the width of the tail, is less than or
// equal to maxWidth.
//
// When [Options.ControlSequences] is true, ANSI escape sequences that appear
// after the truncation point are preserved in the output. This ensures that
// escape sequences such as SGR resets are not lost, preventing color bleed
// in terminal output.
func (options Options) TruncateBytes(s []byte, maxWidth int, tail []byte) []byte {
	maxWidthWithoutTail := maxWidth - options.Bytes(tail)

	var pos, total int
	g := graphemes.FromBytes(s)
	g.AnsiEscapeSequences = options.ControlSequences
	g.AnsiEscapeSequences8Bit = options.ControlSequences8Bit

	for g.Next() {
		gw := graphemeWidth(g.Value(), options)
		if total+gw <= maxWidthWithoutTail {
			pos = g.End()
		}
		total += gw
		if total > maxWidth {
			if options.ControlSequences || options.ControlSequences8Bit {
				// Build result with trailing ANSI escape sequences preserved
				result := make([]byte, 0, len(s)+len(tail)) // at most original + tail
				result = append(result, s[:pos]...)
				result = append(result, tail...)

				rem := graphemes.FromBytes(s[pos:])
				rem.AnsiEscapeSequences = options.ControlSequences
				rem.AnsiEscapeSequences8Bit = options.ControlSequences8Bit

				for rem.Next() {
					v := rem.Value()
					// Only preserve escapes that measure as zero-width
					// on their own; some sequences (e.g. SOS) are only
					// valid in their original context.
					if len(v) > 0 && isEscapeLeader(v[0], options) && options.Bytes(v) == 0 {
						result = append(result, v...)
					}
				}
				return result
			}
			result := make([]byte, 0, pos+len(tail))
			result = append(result, s[:pos]...)
			result = append(result, tail...)
			return result
		}
	}
	// No truncation
	return s
}

// TruncateBytes truncates a []byte to the given maxWidth, and appends the
// given tail if the []byte is truncated.
//
// It ensures the total width, including the width of the tail, is less than or
// equal to maxWidth.
func TruncateBytes(s []byte, maxWidth int, tail []byte) []byte {
	return DefaultOptions.TruncateBytes(s, maxWidth, tail)
}

// isEscapeLeader reports whether the byte is the leading byte of an
// escape sequence that is active for the given options: 7-bit ESC (0x1B)
// when ControlSequences is true, or 8-bit C1 (0x80-0x9F) when
// ControlSequences8Bit is true.
func isEscapeLeader(b byte, options Options) bool {
	return (options.ControlSequences && b == 0x1B) ||
		(options.ControlSequences8Bit && b >= 0x80 && b <= 0x9F)
}
