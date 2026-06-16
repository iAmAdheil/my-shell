package main

import (
	"bufio"
	"bytes"
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
		fmt.Println(err)
	}
	fmt.Println(dir)
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
		fmt.Println("cd:", path+": No such file or directory")
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
		fmt.Println(out)
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
		fmt.Println(com, "is a shell builtin")
	default:
		exePath := GetBinaryPath(com)

		if len(exePath) > 0 {
			fmt.Println(com + " is " + exePath)
		} else {
			fmt.Println(com + ": not found")
		}
	}
}

func GetBinaryPath(filename string) string {
	path := os.Getenv("PATH")
	if len(path) == 0 {
		fmt.Println("no 'PATH' env variable set")
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
		if redirect == 1 {
			wg.Add(1)

			go HandleFileOut(outFile, outScanner, wg, mode)
			go HandlePrintOut(errScanner, errstrch, true)

			errstr := <-errstrch
			if len(errstr) > 0 {
				return fmt.Errorf("%s", errstr)
			}

			wg.Wait()
		} else if redirect == 2 {
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
			return fmt.Errorf("%s", errstr)
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
			fmt.Printf("%s\n", err)
		}
	} else {
		fmt.Printf("%s: command not found\n", main)
	}
}
