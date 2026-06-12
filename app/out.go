package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
)

func HandleBinaryOut(stdout, stderr io.ReadCloser) error {
	wg := &sync.WaitGroup{}

	outScanner := bufio.NewScanner(stdout)
	// errScanner := bufio.NewScanner(stderr)

	wg.Add(1)

	go handlePrint(outScanner, wg)
	// go handlePrint(errScanner, wg)

	wg.Wait()

	return nil
}

func handlePrint(s *bufio.Scanner, wg *sync.WaitGroup) {
	for s.Scan() {
		fmt.Println(s.Text())
	}

	if err := s.Err(); err != nil {
		fmt.Errorf("reading from pipe failed: %s", err)
	}

	wg.Done()
}

func HandleBinaryFileOut(filepath string, stdout, stderr io.ReadCloser) error {
	wg := &sync.WaitGroup{}
	// errCh := make(chan string)
	file, err := OpenFile(filepath)
	if err != nil {
		return err
	}

	outScanner := bufio.NewScanner(stdout)
	// errScanner := bufio.NewScanner(stderr)

	wg.Add(1)

	go handleFileWrite(file, outScanner, nil, wg)
	// go handleFileWrite(file, errScanner, errCh, nil)

	// errStr := <-errCh
	wg.Wait()

	// if len(errStr) > 0 {
	// 	return fmt.Errorf("%s", errStr)
	// }
	return nil
}

func HandleFileOut(out string, filepath string) error {
	f, err := OpenFile(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := bytes.NewBufferString(out)
	s := bufio.NewScanner(buf)
	handleFileWrite(f, s, nil, nil)

	return nil
}

func handleFileWrite(w *os.File, s *bufio.Scanner, errCh chan string, wg *sync.WaitGroup) {
	// var out string
	for s.Scan() {
		buf := s.Text()
		w.WriteString(buf)
		w.WriteString("\n")
		// out += buf
	}

	if err := s.Err(); err != nil {
		fmt.Errorf("reading from pipe failed: %s", err)
	}

	// if errCh != nil {
	// 	errCh <- out
	// }
	if wg != nil {
		wg.Done()
	}
}
