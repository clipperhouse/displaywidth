package displaywidth

import "testing"

// U+1ACF is a Unicode 17 nonspacing mark. Mc spacing marks are deliberately
// excluded from the generator's zero-width category and remain width 1.
func TestStandaloneCombiningMarkRegression(t *testing.T) {
	for _, tc := range []struct {
		name string
		text string
		want int
	}{
		{"new nonspacing mark", "\u1acf", 0},
		{"old nonspacing mark", "\u0591", 0},
		{"spacing mark stays visible", "\u0cf3", 1},
		{"reclassified spacing mark stays visible", "\U0001171e", 1},
		{"attached new mark", "e\u1acf", 1},
		{"new mark after control", "\n\u1acf", 0},
		{"ASCII control", "a", 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := String(tc.text); got != tc.want {
				t.Errorf("String(%q) = %d, want %d", tc.text, got, tc.want)
			}
			if got := Bytes([]byte(tc.text)); got != tc.want {
				t.Errorf("Bytes(%q) = %d, want %d", tc.text, got, tc.want)
			}
		})
	}
}
