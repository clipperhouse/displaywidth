package displaywidth

import "testing"

func TestSetExternalWidths(t *testing.T) {
	// Clear at end so other tests aren't affected.
	defer SetExternalWidths(nil)

	// Baseline: no external map.
	SetExternalWidths(nil)
	baseHeart := String("❤️")
	baseThumbsup := String("👍")
	if baseThumbsup == 0 {
		t.Fatalf("baseline thumbsup should be > 0")
	}

	// Install external override that contradicts the spec values.
	SetExternalWidths(map[string]int{
		"❤️": 1, // spec/PR says 2 (VS16 promotion); override to 1
		"👍": 5, // spec says 2; override to 5 to verify override is used
	})

	if got := String("❤️"); got != 1 {
		t.Errorf("with external override, String(❤️) = %d, want 1", got)
	}
	if got := String("👍"); got != 5 {
		t.Errorf("with external override, String(👍) = %d, want 5", got)
	}

	// Strings not in the override fall back to spec.
	if got := String("🔥"); got != 2 {
		t.Errorf("fallback String(🔥) = %d, want 2 (no override)", got)
	}

	// Clearing returns to baseline.
	SetExternalWidths(nil)
	if got := String("❤️"); got != baseHeart {
		t.Errorf("after clear, String(❤️) = %d, want baseline %d", got, baseHeart)
	}
	if got := String("👍"); got != baseThumbsup {
		t.Errorf("after clear, String(👍) = %d, want baseline %d", got, baseThumbsup)
	}
}

func TestExternalWidthsAffectsBytes(t *testing.T) {
	defer SetExternalWidths(nil)

	SetExternalWidths(map[string]int{"👍": 4})

	if got := Bytes([]byte("👍")); got != 4 {
		t.Errorf("Bytes(👍) with override = %d, want 4", got)
	}
}

func TestExternalWidthsDoesNotAffectASCII(t *testing.T) {
	defer SetExternalWidths(nil)

	// Even if the override map says "a" is 99, ASCII fast paths should
	// return the standard width — the override is checked only inside
	// graphemeWidth, which doesn't run for printable ASCII.
	SetExternalWidths(map[string]int{"a": 99})

	if got := String("a"); got != 1 {
		t.Errorf("String(\"a\") = %d, want 1 (ASCII bypass)", got)
	}
	if got := String("hello"); got != 5 {
		t.Errorf("String(\"hello\") = %d, want 5 (ASCII bypass)", got)
	}
}

func TestExternalWidthsMixedContent(t *testing.T) {
	defer SetExternalWidths(nil)

	SetExternalWidths(map[string]int{
		"❤️": 1,
		"👍": 4,
	})

	// "abc❤️def👍" → 3 + 1 + 3 + 4 = 11
	if got := String("abc❤️def👍"); got != 11 {
		t.Errorf("mixed content = %d, want 11", got)
	}
}
