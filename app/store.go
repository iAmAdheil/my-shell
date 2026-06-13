package main

import (
	"log"
	"os"
)

func OpenFile(filepath string, mode int) (*os.File, error) {
	var wPerm int
	if mode == 0 {
		wPerm = os.O_TRUNC
	} else if mode == 1 {
		wPerm = os.O_APPEND
	}

	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|wPerm, 0644)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
		return nil, err
	}

	return file, nil
}
