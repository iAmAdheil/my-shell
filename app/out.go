package main

import (
	"bufio"
	"fmt"
	"sync"
)

func HandlePrintOut(s *bufio.Scanner, errstrch chan string, isErr bool) {
	// isErr -> collect only, no print, transmit msg via channel
	var out string
	for s.Scan() {
		buf := s.Text()
		if !isErr {
			fmt.Printf("%s\n\r", buf)
		}
		out += buf
	}

	if err := s.Err(); err != nil {
		errstrch <- err.Error()
	}

	if errstrch != nil {
		errstrch <- out
	}
}

func HandleFileOut(filepath string, s *bufio.Scanner, wg *sync.WaitGroup, mode int) {
	if err := handleFileWrite(filepath, s, wg, mode); err != nil {
		panic(err)
		// fmt.Errorf("file out failed: %s\n\r", err)
	}

}

func handleFileWrite(filepath string, s *bufio.Scanner, wg *sync.WaitGroup, mode int) error {
	f, err := OpenFile(filepath, mode)
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
