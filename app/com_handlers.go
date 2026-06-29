package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"slices"
	"strings"
	"sync"
)

func HandleExit() {
	os.Exit(0)
}

func HandlePwd() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Printf("%s\n\r", err)
	}
	fmt.Printf("%s\n\r", dir)
}

func HandleCd(args []string) {
	if len(args) == 0 {
		return
	}
	path := args[0]
	if path == "~" {
		path = os.Getenv("HOME")
	}
	if err := os.Chdir(path); err != nil {
		fmt.Printf("cd: %s: No such file or directory\n\r", path)
	}
}

func HandleEcho(args []string, filepath string, redirect int, mode int) error {
	wg := &sync.WaitGroup{}
	var (
		out     string = strings.Join(args, " ")
		fileout string = strings.Join(args, " ")
	)
	if redirect == 2 {
		fileout = ""
	}
	if len(filepath) > 0 {
		wg.Add(1)
		r := bytes.NewBufferString(fileout)
		s := bufio.NewScanner(r)
		HandleFileOut(filepath, s, wg, mode)
		wg.Wait()
	}
	if redirect == 0 || redirect == 2 {
		fmt.Printf("%s\n\r", out)
	}

	return nil
}

func HandleType(args []string) {
	// normalised args
	if len(args) == 0 {
		return
	}
	com := args[0]
	switch com {
	case "exit", "echo", "type", "pwd", "cd":
		fmt.Printf("%s is a shell builtin\n\r", com)
	default:
		exePath := GetBinaryPath(com)

		if len(exePath) > 0 {
			fmt.Printf("%s is %s\n\r", com, exePath)
		} else {
			fmt.Printf("%s: not found\n\r", com)
		}
	}
}

func GetBinaryPath(filename string) string {
	path := os.Getenv("PATH")
	if len(path) == 0 {
		fmt.Printf("no 'PATH' env variable set\n\r")
	}

	dirs := strings.Split(path, ":")

	for _, dir := range dirs {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, file := range files {
			if file.Name() == filename {
				fullPath := dir + "/" + filename
				fileInfo, err := os.Stat(fullPath)
				if err != nil {
					break
				}
				if IsExecAny(fileInfo.Mode().Perm()) { // check file permissions
					return fullPath
				} else {
					break
				}
			}
		}
	}
	return ""
}

// redirect == 1 -> stdout
// redirect == 2 -> stderr
func RunBinary(file string, args []string, outFile string, redirect int, mode int) error {
	proc := exec.Command(file, args...)

	stdout, err := proc.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := proc.StderrPipe()
	if err != nil {
		return err
	}

	if err := proc.Start(); err != nil {
		return err
	}

	var (
		wg         *sync.WaitGroup = &sync.WaitGroup{}
		errstrch   chan string     = make(chan string)
		outScanner *bufio.Scanner  = bufio.NewScanner(stdout)
		errScanner *bufio.Scanner  = bufio.NewScanner(stderr)
	)

	// if redirect == 1 -> stdout to file and print stderr message
	// from chan and ignore the err from proc.wait
	// if redirect == 2 -> stderr to file and print stdout message,
	// ignore the err from proc.wait
	if len(outFile) > 0 {
		switch redirect {
		case 1:
			wg.Add(1)

			go HandleFileOut(outFile, outScanner, wg, mode)
			go HandlePrintOut(errScanner, errstrch, true)

			errstr := <-errstrch
			if len(errstr) > 0 {
				// return fmt.Errorf("%s\n\r", errstr)
				return errors.New(errstr)
			}

			wg.Wait()

		case 2:
			wg.Add(1)
			// to finish scanning from scanner before wait is called
			// and the internal stdout reader gets closed
			done := make(chan struct{})
			go HandleFileOut(outFile, errScanner, wg, mode)
			go func() {
				HandlePrintOut(outScanner, nil, false)
				close(done)
			}()
			<-done
			wg.Wait()
		}
	} else {
		// to finish scanning from scanner before wait is called
		// and the internal stdout reader gets closed
		done := make(chan struct{})
		go func() {
			HandlePrintOut(outScanner, nil, false)
			close(done)
		}()
		go HandlePrintOut(errScanner, errstrch, true)

		<-done
		errstr := <-errstrch
		if len(errstr) > 0 {
			// return fmt.Errorf("%s\n\r", errstr)
			return errors.New(errstr)
		}
	}

	// ignore err from proc, handled by stderr pipe
	proc.Wait()
	return nil
}

func HandleDualComm(c1 string, args1 []string, c2 string, args2 []string) {
	exeP1 := GetBinaryPath(c1)
	exeP2 := GetBinaryPath(c2)

	if len(exeP1) == 0 || len(exeP2) == 0 {
		return
	}

	proc1 := exec.Command(c1, args1...)
	proc2 := exec.Command(c2, args2...)

	stdout1, err := proc1.StdoutPipe()
	if err != nil {
		return
	}
	stdin2, err := proc2.StdinPipe()
	if err != nil {
		return
	}
	stdout2, err := proc2.StdoutPipe()
	if err != nil {
		return
	}

	if err := proc1.Start(); err != nil {
		return
	}
	if err := proc2.Start(); err != nil {
		return
	}

	outScanner2 := bufio.NewScanner(stdout2)
	done := make(chan struct{})
	go func() {
		HandlePrintOut(outScanner2, nil, false)
		close(done)
	}()

	if _, err := io.Copy(stdin2, stdout1); err != nil {
		return
	}

	stdin2.Close()

	<-done
	proc1.Wait()
	proc2.Wait()
}

func HandleDefault(main string, args []string, filepath string, redirect int, mode int) {
	if len(main) == 0 {
		return
	}

	if slices.Contains(args, "|") {
		idx := slices.Index(args, "|")
		if len(args)-1 == idx {
			return
		}
		c1 := main
		args1 := args[0:idx]
		c2 := args[idx+1]
		args2 := args[idx+2:]
		HandleDualComm(c1, args1, c2, args2)
		return
	}

	exePath := GetBinaryPath(main)

	if len(exePath) > 0 {
		err := RunBinary(main, args, filepath, redirect, mode)
		if err != nil {
			// cat: nonexistent: No such file or directory
			fmt.Printf("%s\n\r", err)
		}
	} else {
		fmt.Printf("%s: command not found\n\r", main)
	}
}
