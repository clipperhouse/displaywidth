package displaywidth

import "testing"

func TestVariationSelectorsAffectWidth(t *testing.T) {
	// U+263A WHITE SMILING FACE: text default; VS16 -> emoji
	if got := String("\u263a"); got != 1 {
		t.Fatalf("expected U+263A width 1, got %d", got)
	}
	if got := String("\u263a\ufe0f"); got != 2 {
		t.Fatalf("expected U+263A+VS16 width 2, got %d", got)
	}

	// U+231B HOURGLASS: emoji default; VS15 -> text
	if got := String("\u231b"); got != 2 {
		t.Fatalf("expected U+231B width 2, got %d", got)
	}
	if got := String("\u231b\ufe0e"); got != 1 {
		t.Fatalf("expected U+231B+VS15 width 1, got %d", got)
	}

	// Keycap sequence: 1 + VS16 (FE0F) + COMBINING ENCLOSING KEYCAP (20E3)
	if got := String("1\ufe0f\u20e3"); got != 2 {
		t.Fatalf("expected keycap width 2, got %d", got)
	}
}
