package main

import (
	"os"
	"testing"
)

// use this to test random things

func TestRandom(t *testing.T) {
	filenames := []string{}

	dir, err := os.Getwd()
	if err != nil {
		return
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		filenames = append(filenames, entry.Name())
	}
}
