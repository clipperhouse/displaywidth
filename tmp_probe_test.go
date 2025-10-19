package displaywidth

import "testing"

// Temporary probe test to log widths for selected runes/strings.
// This test does not assert and will be removed after analysis.
func TestTmpProbe(t *testing.T) {
	cases := []struct{ label, s string }{
		{"desert_island", "🏝"},
		{"mount_fuji", "🗻"},
		{"beach_with_umbrella", "🏖"},
		{"white_smiling_face_vs16", "\u263a\ufe0f"},
		{"hourglass_vs15", "\u231b\ufe0e"},
		{"two_em_dash", "\u2e3a"},
		{"three_em_dash", "\u2e3b"},
	}

	for _, c := range cases {
		t.Logf("%s: %q width=%d", c.label, c.s, String(c.s))
	}
}
