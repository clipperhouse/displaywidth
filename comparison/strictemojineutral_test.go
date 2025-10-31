package comparison

import (
	"testing"

	"github.com/clipperhouse/displaywidth"
	"github.com/mattn/go-runewidth"
)

// TestRegularEmojiAreAlwaysWidth2 explicitly verifies that regular emojis
// are always width 2
func TestRegularEmojiAreAlwaysWidth2(t *testing.T) {
	// Regular emojis WITHOUT variation selectors (VS16)
	// go-runewidth handles VS16 differently, treating it as a separate character
	emojis := []string{
		"😀", "😁", "😂", "🤣", "😃", "😄", "😅", "😆",
		"🚀", "🎉", "🎊", "🎈", "🎁", "🎂",
		"👍", "👎", "👏", "🙏",
		"🔥", "💯", "✨", "⭐",
	}

	t.Log("Verifying that regular emojis are ALWAYS width 2:")
	t.Log("(Note: excludes emojis with VS16, which go-runewidth treats differently)")
	t.Log("")

	allPass := true
	for _, emoji := range emojis {
		// displaywidth (always width 2, no StrictEmojiNeutral option)
		dw := displaywidth.String(emoji)

		// go-runewidth (always width 2 regardless of StrictEmojiNeutral)
		gr1 := (&runewidth.Condition{StrictEmojiNeutral: true}).StringWidth(emoji)
		gr2 := (&runewidth.Condition{StrictEmojiNeutral: false}).StringWidth(emoji)

		if dw != 2 || gr1 != 2 || gr2 != 2 {
			t.Errorf("%s: Expected width 2 in all cases, got displaywidth=%d, go-runewidth(strict=true)=%d, go-runewidth(strict=false)=%d",
				emoji, dw, gr1, gr2)
			allPass = false
		}
	}

	if allPass {
		t.Log("✅ All regular emojis have width 2")
	}

	// Document the difference with VS16
	t.Log("")
	t.Log("Variation Selector Difference:")
	emojiWithVS16 := []string{"❤️", "✂️", "☺️"}
	for _, emoji := range emojiWithVS16 {
		dw := displaywidth.String(emoji)
		gr := runewidth.StringWidth(emoji)
		t.Logf("  %s: displaywidth=%d, go-runewidth=%d (go-runewidth treats VS16 as separate char)", emoji, dw, gr)
	}
}

// TestFlagsBehavior documents flag behavior across libraries
func TestFlagsBehavior(t *testing.T) {
	flags := []string{
		"🇺🇸", "🇯🇵", "🇬🇧", "🇫🇷", "🇩🇪",
		"🇨🇦", "🇦🇺", "🇧🇷", "🇮🇳", "🇨🇳",
	}

	t.Log("Flag behavior comparison:")
	t.Log("(displaywidth follows modern standards: flags are always width 2)")
	t.Log("")
	t.Log("Flag | displaywidth | go-runewidth (default) | go-runewidth (strict=false) | go-runewidth (strict=true)")
	t.Log("-----|-------------|------------------------|---------------------------|-------------------------")

	for _, flag := range flags {
		// displaywidth (always width 2, no StrictEmojiNeutral option)
		dw := displaywidth.String(flag)

		// go-runewidth (always width 1, regardless of StrictEmojiNeutral)
		grDefault := runewidth.StringWidth(flag)
		gr1 := (&runewidth.Condition{StrictEmojiNeutral: false}).StringWidth(flag)
		gr2 := (&runewidth.Condition{StrictEmojiNeutral: true}).StringWidth(flag)

		t.Logf("%s | %d | %d | %d | %d", flag, dw, grDefault, gr1, gr2)

		// Verify displaywidth behavior (always width 2)
		if dw != 2 {
			t.Errorf("%s: displaywidth should always be 2, got %d", flag, dw)
		}

		// Document go-runewidth behavior (always 1)
		if grDefault != 1 {
			t.Errorf("%s: go-runewidth default should be 1, got %d", flag, grDefault)
		}
		if gr1 != 1 {
			t.Errorf("%s: go-runewidth with strict=false should be 1, got %d", flag, gr1)
		}
		if gr2 != 1 {
			t.Errorf("%s: go-runewidth with strict=true should be 1, got %d", flag, gr2)
		}
	}

	t.Log("")
	t.Log("Summary:")
	t.Log("- displaywidth: flags are always width 2 (modern standard)")
	t.Log("- go-runewidth: flags are always width 1, regardless of StrictEmojiNeutral")
}
