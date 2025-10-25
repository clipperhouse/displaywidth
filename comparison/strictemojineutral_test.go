package comparison

import (
	"testing"

	"github.com/clipperhouse/displaywidth"
	"github.com/mattn/go-runewidth"
)

// TestStrictEmojiNeutralBehavior documents the actual behavior of StrictEmojiNeutral
// in both displaywidth and go-runewidth
func TestStrictEmojiNeutralBehavior(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		description string
	}{
		// Regular emoji (always width 2)
		{"Regular emoji 😀", "😀", "Regular emoji (grinning face)"},
		{"Regular emoji 🚀", "🚀", "Regular emoji (rocket)"},
		{"Regular emoji 🎉", "🎉", "Regular emoji (party popper)"},
		{"Regular emoji 🔥", "🔥", "Regular emoji (fire)"},
		{"Multiple emoji", "😀🚀🎉", "Multiple regular emojis"},

		// Flags (regional indicator pairs)
		{"Flag 🇺🇸", "🇺🇸", "US flag (regional indicator pair)"},
		{"Flag 🇯🇵", "🇯🇵", "Japan flag (regional indicator pair)"},
		{"Flag 🇬🇧", "🇬🇧", "UK flag (regional indicator pair)"},
		{"Multiple flags", "🇺🇸🇯🇵🇬🇧", "Multiple flags"},

		// Single-codepoint flag emoji
		{"Checkered flag 🏁", "🏁", "Checkered flag (single codepoint)"},
		{"White flag 🏳", "🏳", "White flag (single codepoint)"},
		{"Black flag 🏴", "🏴", "Black flag (single codepoint)"},

		// Emoji with text/emoji presentation variation
		{"Red heart ❤", "❤", "Red heart (text presentation by default)"},
		{"Red heart ❤️", "❤️", "Red heart with VS16 (emoji presentation)"},
		{"Scissors ✂", "✂", "Scissors (text presentation by default)"},
		{"Scissors ✂️", "✂️", "Scissors with VS16 (emoji presentation)"},
		{"Smiling face ☺", "☺", "Smiling face (text presentation by default)"},
		{"Smiling face ☺️", "☺️", "Smiling face with VS16 (emoji presentation)"},

		// Keycap sequences
		{"Keycap 1️⃣", "1️⃣", "Keycap sequence (1 + VS16 + combining keycap)"},
		{"Keycap #️⃣", "#️⃣", "Keycap sequence (# + VS16 + combining keycap)"},
	}

	t.Log("Testing StrictEmojiNeutral behavior:")
	t.Log("")
	t.Log("Character | Description | displaywidth (default) | displaywidth (strict=false) | go-runewidth (default) | go-runewidth (strict=false)")
	t.Log("----------|-------------|------------------------|-----------------------------|-----------------------|----------------------------")

	for _, tc := range testCases {
		// displaywidth tests
		dwDefault := displaywidth.String(tc.input) // StrictEmojiNeutral=true (default)
		dwStrictFalse := displaywidth.Options{StrictEmojiNeutral: false}.String(tc.input)

		// go-runewidth tests
		grDefault := runewidth.StringWidth(tc.input) // StrictEmojiNeutral=true (default)
		grStrictFalse := (&runewidth.Condition{StrictEmojiNeutral: false}).StringWidth(tc.input)

		t.Logf("%s | %s | %d | %d | %d | %d",
			tc.input, tc.description, dwDefault, dwStrictFalse, grDefault, grStrictFalse)

		// Document the differences
		if dwDefault != grDefault {
			t.Logf("  ⚠️  Difference at default settings: displaywidth=%d, go-runewidth=%d", dwDefault, grDefault)
		}
		if dwStrictFalse != grStrictFalse {
			t.Logf("  ⚠️  Difference with strict=false: displaywidth=%d, go-runewidth=%d", dwStrictFalse, grStrictFalse)
		}
	}

	t.Log("")
	t.Log("Key Findings:")
	t.Log("1. Regular emoji (😀, 🚀, 🎉) are ALWAYS width 2 in both libraries, regardless of StrictEmojiNeutral")
	t.Log("2. Flags (🇺🇸) in displaywidth: width 2 with strict=true, width 1 with strict=false")
	t.Log("3. Flags (🇺🇸) in go-runewidth: width 1 ALWAYS, StrictEmojiNeutral has NO effect on flags")
	t.Log("4. StrictEmojiNeutral ONLY affects flags in displaywidth, not regular emojis")
}

// TestStrictEmojiNeutralWithEAW tests the interaction between StrictEmojiNeutral and EastAsianWidth
func TestStrictEmojiNeutralWithEAW(t *testing.T) {
	testCases := []struct {
		name  string
		input string
	}{
		{"Regular emoji 😀", "😀"},
		{"Flag 🇺🇸", "🇺🇸"},
		{"Ambiguous star ★", "★"},
		{"Red heart ❤", "❤"},
	}

	t.Log("Testing StrictEmojiNeutral with EastAsianWidth:")
	t.Log("")
	t.Log("Character | EAW=false,Strict=true | EAW=false,Strict=false | EAW=true,Strict=true | EAW=true,Strict=false")
	t.Log("----------|----------------------|------------------------|----------------------|-----------------------")

	for _, tc := range testCases {
		// displaywidth combinations
		dw1 := displaywidth.Options{EastAsianWidth: false, StrictEmojiNeutral: true}.String(tc.input)
		dw2 := displaywidth.Options{EastAsianWidth: false, StrictEmojiNeutral: false}.String(tc.input)
		dw3 := displaywidth.Options{EastAsianWidth: true, StrictEmojiNeutral: true}.String(tc.input)
		dw4 := displaywidth.Options{EastAsianWidth: true, StrictEmojiNeutral: false}.String(tc.input)

		t.Logf("displaywidth %s | %d | %d | %d | %d", tc.input, dw1, dw2, dw3, dw4)

		// go-runewidth combinations
		gr1 := (&runewidth.Condition{EastAsianWidth: false, StrictEmojiNeutral: true}).StringWidth(tc.input)
		gr2 := (&runewidth.Condition{EastAsianWidth: false, StrictEmojiNeutral: false}).StringWidth(tc.input)
		gr3 := (&runewidth.Condition{EastAsianWidth: true, StrictEmojiNeutral: true}).StringWidth(tc.input)
		gr4 := (&runewidth.Condition{EastAsianWidth: true, StrictEmojiNeutral: false}).StringWidth(tc.input)

		t.Logf("go-runewidth %s | %d | %d | %d | %d", tc.input, gr1, gr2, gr3, gr4)
		t.Log("")
	}
}

// TestRegularEmojiAreAlwaysWidth2 explicitly verifies that regular emojis
// are always width 2, regardless of StrictEmojiNeutral
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
		// displaywidth
		dw1 := displaywidth.Options{StrictEmojiNeutral: true}.String(emoji)
		dw2 := displaywidth.Options{StrictEmojiNeutral: false}.String(emoji)

		// go-runewidth
		gr1 := (&runewidth.Condition{StrictEmojiNeutral: true}).StringWidth(emoji)
		gr2 := (&runewidth.Condition{StrictEmojiNeutral: false}).StringWidth(emoji)

		if dw1 != 2 || dw2 != 2 || gr1 != 2 || gr2 != 2 {
			t.Errorf("%s: Expected width 2 in all cases, got displaywidth(strict=true)=%d, displaywidth(strict=false)=%d, go-runewidth(strict=true)=%d, go-runewidth(strict=false)=%d",
				emoji, dw1, dw2, gr1, gr2)
			allPass = false
		}
	}

	if allPass {
		t.Log("✅ All regular emojis have width 2 regardless of StrictEmojiNeutral setting")
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

// TestFlagsAreAffectedByStrictEmojiNeutral documents that flags behave differently
func TestFlagsAreAffectedByStrictEmojiNeutral(t *testing.T) {
	flags := []string{
		"🇺🇸", "🇯🇵", "🇬🇧", "🇫🇷", "🇩🇪",
		"🇨🇦", "🇦🇺", "🇧🇷", "🇮🇳", "🇨🇳",
	}

	t.Log("Verifying that flags ARE affected by StrictEmojiNeutral:")
	t.Log("")
	t.Log("Flag | displaywidth (strict=true) | displaywidth (strict=false) | go-runewidth (strict=true) | go-runewidth (strict=false)")
	t.Log("-----|----------------------------|-----------------------------|-----------------------------|------------------------------")

	for _, flag := range flags {
		// displaywidth
		dw1 := displaywidth.Options{StrictEmojiNeutral: true}.String(flag)
		dw2 := displaywidth.Options{StrictEmojiNeutral: false}.String(flag)

		// go-runewidth
		gr1 := (&runewidth.Condition{StrictEmojiNeutral: true}).StringWidth(flag)
		gr2 := (&runewidth.Condition{StrictEmojiNeutral: false}).StringWidth(flag)

		t.Logf("%s | %d | %d | %d | %d", flag, dw1, dw2, gr1, gr2)

		// Verify displaywidth behavior
		if dw1 != 2 {
			t.Errorf("%s: displaywidth with strict=true should be 2, got %d", flag, dw1)
		}
		if dw2 != 1 {
			t.Errorf("%s: displaywidth with strict=false should be 1, got %d", flag, dw2)
		}

		// Document go-runewidth behavior (always 1)
		if gr1 != 1 {
			t.Errorf("%s: go-runewidth with strict=true should be 1, got %d", flag, gr1)
		}
		if gr2 != 1 {
			t.Errorf("%s: go-runewidth with strict=false should be 1, got %d", flag, gr2)
		}
	}

	t.Log("")
	t.Log("Summary:")
	t.Log("- displaywidth: flags are width 2 with StrictEmojiNeutral=true (default), width 1 with false")
	t.Log("- go-runewidth: flags are ALWAYS width 1, regardless of StrictEmojiNeutral")
}
