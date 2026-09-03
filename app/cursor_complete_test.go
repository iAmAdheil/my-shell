package main

import (
	"os"
	"testing"

	"github.com/chzyer/readline"
)

// completion must act on the word that ends at the cursor.
// the text on the right of the cursor must stay where it is.

func TestDoCompletesWordAtCursor(t *testing.T) {
	dir := t.TempDir()
	for _, name := range []string{"alpha.txt", "beta.txt"} {
		if err := os.WriteFile(dir+"/"+name, nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	t.Chdir(dir)

	bnm := &BellNoMatch{
		inner: readline.NewPrefixCompleter(
			readline.PcItem("echo"),
			readline.PcItem("exit"),
		),
	}

	tests := []struct {
		name   string
		line   string
		pos    int
		insert string
		offset int
	}{
		{"cursor at the end", "cat alph", 8, "a.txt ", 4},
		{"cursor inside the word", "cat alph", 6, "pha.txt ", 2},
		{"cursor before more text", "cat alph zzz", 8, "a.txt ", 4},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, offset := bnm.Do([]rune(tc.line), tc.pos)

			if len(got) != 1 || string(got[0]) != tc.insert {
				t.Errorf("Do(%q, %d) inserts %q, want [%q]", tc.line, tc.pos, got, tc.insert)
			}
			if offset != tc.offset {
				t.Errorf("Do(%q, %d) offset = %d, want %d", tc.line, tc.pos, offset, tc.offset)
			}
		})
	}
}

// a word that matches nothing must leave the line alone.
func TestDoOnNoMatch(t *testing.T) {
	t.Chdir(t.TempDir())

	bnm := &BellNoMatch{inner: readline.NewPrefixCompleter(readline.PcItem("echo"))}

	got, _ := bnm.Do([]rune("cat zzz alph"), 7)
	if len(got) != 0 {
		t.Errorf("Do inserts %q, want no completion", got)
	}
}
