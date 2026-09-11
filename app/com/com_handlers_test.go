package com

import (
	"bytes"
	"testing"
)

// nopWriteCloser makes a bytes.Buffer usable as the Out field of a Com.
type nopWriteCloser struct {
	*bytes.Buffer
}

func (nopWriteCloser) Close() error { return nil }

func TestHandleEcho(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"no argument", []string{}, "\n"},
		{"one argument", []string{"hi"}, "hi\n"},
		{"only the -e flag", []string{"-e"}, "\n"},
		{"the -e flag with a newline", []string{"-e", `a\nb`}, "a\nb\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			com := &Com{Main: "echo", Args: tt.args, Out: nopWriteCloser{buf}}

			if err := com.HandleEcho(); err != nil {
				t.Fatalf("HandleEcho() returned %v", err)
			}

			if got := buf.String(); got != tt.want {
				t.Errorf("HandleEcho() wrote %q, want %q", got, tt.want)
			}
		})
	}
}
