package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
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

			go HandleFileOut(outFile, errScanner, wg, mode)
			go HandlePrintOut(outScanner, nil, false)

			wg.Wait()
		}
	} else {
		go HandlePrintOut(outScanner, nil, false)
		go HandlePrintOut(errScanner, errstrch, true)

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

func HandleDefault(main string, args []string, filepath string, redirect int, mode int) {
	if len(main) == 0 {
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
