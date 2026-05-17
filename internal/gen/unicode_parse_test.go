package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseEmojiVariationSequencesErrors(t *testing.T) {
	cases := []struct {
		name    string
		content string
		wantErr string
	}{
		{
			name:    "valid VS16 entry",
			content: "0023 FE0F ; emoji style;  # NUMBER SIGN\n",
			wantErr: "",
		},
		{
			name:    "valid VS15 entry is not an error (intentionally ignored)",
			content: "0023 FE0E ; text style;  # NUMBER SIGN\n",
			wantErr: "",
		},
		{
			name:    "blank and comment lines ignored",
			content: "# header\n\n  \n",
			wantErr: "",
		},
		{
			name:    "missing semicolons fails",
			content: "0023 FE0F\n",
			wantErr: "expected at least two ';'-separated fields",
		},
		{
			name:    "wrong field count in first part",
			content: "0023 FE0F EXTRA ; emoji style;  # bogus\n",
			wantErr: "expected '<base> <vs>'",
		},
		{
			name:    "invalid base hex",
			content: "GGGG FE0F ; emoji style;  # bogus\n",
			wantErr: "invalid base codepoint",
		},
		{
			name:    "invalid vs hex",
			content: "0023 ZZZZ ; emoji style;  # bogus\n",
			wantErr: "invalid variation selector codepoint",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "emoji-variation-sequences.txt")
			if err := os.WriteFile(path, []byte(tc.content), 0644); err != nil {
				t.Fatal(err)
			}
			data := &UnicodeData{VS16Eligible: make(map[rune]bool)}

			err := parseEmojiVariationSequences(path, data)
			if tc.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", tc.wantErr)
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("expected error containing %q, got %q", tc.wantErr, err.Error())
			}
		})
	}
}
