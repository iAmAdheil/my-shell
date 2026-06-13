package main

import (
	"bufio"
	"fmt"
	"sync"
)

func HandlePrintOut(s *bufio.Scanner, errstrch chan string, isErr bool) {
	// isErr -> collect only, no print
	var out string
	for s.Scan() {
		buf := s.Text()
		if !isErr {
			fmt.Println(buf)
		}
		out += buf
	}

	if err := s.Err(); err != nil {
		fmt.Errorf("reading from pipe failed: %s", err)
	}

	if errstrch != nil {
		errstrch <- out
	}
}

func HandleFileOut(filepath string, s *bufio.Scanner, wg *sync.WaitGroup) {
	if err := handleFileWrite(filepath, s, wg); err != nil {
		fmt.Errorf("file out failed: %s\n", err)
	}

}

func handleFileWrite(filepath string, s *bufio.Scanner, wg *sync.WaitGroup) error {
	f, err := OpenFile(filepath)
	if err != nil {
		return err
	}
	for s.Scan() {
		buf := s.Text()
		f.WriteString(buf)
		f.WriteString("\n")
	}

	if err := s.Err(); err != nil {
		return err
	}

	if wg != nil {
		wg.Done()
	}

	return nil
}
