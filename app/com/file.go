package com

import (
	"log"
	"os"
)

// mode == 0 -> read only
// mode == 1 -> truncate
// mode == 2 -> append
func OpenFile(filepath string, mode int) (*os.File, error) {
	var perm int = os.O_RDONLY | os.O_CREATE
	switch mode {
	case 1:
		perm = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	case 2:
		perm = os.O_WRONLY | os.O_CREATE | os.O_APPEND
	}

	file, err := os.OpenFile(filepath, perm, 0644)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
		return nil, err
	}

	return file, nil
}
