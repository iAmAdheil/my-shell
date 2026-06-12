package main

import (
	"log"
	"os"
)

func OpenFile(filename string) (*os.File, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
		return nil, err
	}

	return file, nil
}
