package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"sync"
)

func HandleOut(proc *exec.Cmd) error {
	wg := &sync.WaitGroup{}

	stdout, err := proc.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := proc.StderrPipe()
	if err != nil {
		return err
	}

	outScanner := bufio.NewScanner(stdout)
	errScanner := bufio.NewScanner(stderr)

	wg.Add(2)

	go handlePrint(outScanner, wg)
	go handlePrint(errScanner, wg)

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
