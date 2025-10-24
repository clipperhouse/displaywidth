package comparison

import (
	"testing"
	"unicode"

	"github.com/clipperhouse/displaywidth"
	"github.com/mattn/go-runewidth"
)

// Investigate differences between displaywidth and go-runewidth.
// Not meant to be a real test, more of a tool, but good to
// fail if something changes in the future.
func TestCompatibility(t *testing.T) {
	t.Skip("skipping compatibility test")
	if unicode.Version < "15" {
		// We only care about Unicode 15 and above,
		// which I believe was Go version 1.21.
		return
	}

	for r := rune(0); r <= unicode.MaxRune; r++ {
		w1 := displaywidth.Rune(r)
		w2 := runewidth.RuneWidth(r)
		if w1 != w2 {
			if unicode.Is(unicode.Mn, r) {
				// these are in the trie, known
				// we will return width 0,
				// go-runewidth may return width 1
				continue
			}
			if unicode.Is(unicode.Cf, r) {
				// these are in the trie, known
				// we will return width 0,
				// go-runewidth may return width 1
				continue
			}
			if unicode.Is(unicode.Mc, r) {
				// these are deliberately excluded from the trie, known
				// we will return width 1,
				// go-runewidth may return width 0
				continue
			}
			t.Errorf("%#x: runewidth is %d, displaywidth is %d, difference is %d", r, w2, w1, w2-w1)
		}
	}
}
