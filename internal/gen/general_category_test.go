package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseGeneralCategories(t *testing.T) {
	path := filepath.Join(t.TempDir(), "categories.txt")
	text := "# fixture\n\n1ACF..1AD0 ; Mn # new marks\n20DD ; Me\n200D ; Cf\n2028 ; Zl\n2029 ; Zp\n0CF3 ; Mc\n0041 ; Lu\n"
	if err := os.WriteFile(path, []byte(text), 0644); err != nil {
		t.Fatal(err)
	}
	data := &UnicodeData{CombiningMarks: map[rune]bool{}, ZeroWidthChars: map[rune]bool{}}
	if err := parseGeneralCategories(path, data); err != nil {
		t.Fatal(err)
	}
	for _, r := range []rune{0x1ACF, 0x1AD0, 0x20DD} {
		if !data.CombiningMarks[r] {
			t.Errorf("missing combining mark %U", r)
		}
	}
	for _, r := range []rune{0x200D, 0x2028, 0x2029} {
		if !data.ZeroWidthChars[r] {
			t.Errorf("missing zero-width character %U", r)
		}
	}
	if len(data.CombiningMarks) != 3 || len(data.ZeroWidthChars) != 3 {
		t.Fatal("unexpected category entries, including possible Mc spacing marks")
	}
}

func TestParseGeneralCategoriesErrors(t *testing.T) {
	for _, text := range []string{"1ACF Mn", "xyz ; Mn", "110000 ; Cf", "1AD0..1ACF ; Mn", "-1 ; Me", "1..2..3 ; Zl", "1ACF ; Mn ; extra"} {
		t.Run(text, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "categories.txt")
			if err := os.WriteFile(path, []byte(text), 0644); err != nil {
				t.Fatal(err)
			}
			data := &UnicodeData{CombiningMarks: map[rune]bool{}, ZeroWidthChars: map[rune]bool{}}
			if err := parseGeneralCategories(path, data); err == nil {
				t.Fatal("expected malformed category data to fail")
			}
		})
	}
}

func TestPinnedGeneralCategories(t *testing.T) {
	data, err := ParseUnicodeData()
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range []rune{0x0591, 0x1ACF, 0x1ADD, 0x20DD} {
		if !data.CombiningMarks[r] {
			t.Errorf("missing Unicode 17 combining mark %U", r)
		}
	}
	for _, r := range []rune{0x0CF3, 0x1171E} {
		if data.CombiningMarks[r] {
			t.Fatalf("Mc %U must not be treated as a nonspacing mark", r)
		}
	}
}
