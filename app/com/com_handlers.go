package com

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func (com *Com) Run() {
	switch com.Main {
	case "":
	case "exit":
		com.HandleExit()
	case "echo":
		com.HandleEcho()
	case "type":
		com.HandleType()
	case "pwd":
		com.HandlePwd()
	case "cd":
		com.HandleCd()
	default:
		exePath := GetBinaryPath(com.Main)
		if len(exePath) > 0 {
			err := com.RunBinary()
			if err != nil {
				// cat: nonexistent: No such file or directory
				fmt.Printf("%s\n\r", err)
			}
		} else {
			fmt.Printf("%s: command not found\n\r", com.Main)
		}
	}

	if com.Close {
		err := com.Out.Close()
		if err != nil {
			panic(err)
		}
	}
}

func (com *Com) HandleExit() {
	os.Exit(0)
}

func (com *Com) HandlePwd() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(com.Out, "%s\n\r", err)
		return
	}
	fmt.Fprintf(com.Out, "%s\n\r", dir)
}

func (com *Com) HandleCd() {
	if len(com.Args) == 0 {
		return
	}

	path := com.Args[0]
	if path == "~" {
		path = os.Getenv("HOME")
	}

	if err := os.Chdir(path); err != nil {
		fmt.Fprintf(com.Out, "cd: %s: No such file or directory\n\r", path)
	}
}

func (com *Com) HandleEcho() error {
	wg := &sync.WaitGroup{}
	var (
		out     string = strings.Join(com.Args, " ")
		fileout string = strings.Join(com.Args, " ")
	)
	// stderr does not have any output, but file is created
	if com.Redirect == 2 {
		fileout = ""
	}

	if len(com.OutFilePath) > 0 {
		wg.Add(1)

		r := bytes.NewBufferString(fileout)
		s := bufio.NewScanner(r)
		HandleFileOut(com.OutFilePath, s, wg, com.Mode)

		wg.Wait()
	}
	// print when either no redirect or redirect stderr
	if com.Redirect == 0 || com.Redirect == 2 {
		fmt.Fprintf(com.Out, "%s\n\r", out)
	}

	return nil
}

func (com *Com) HandleType() {
	if len(com.Args) == 0 {
		return
	}

	m := com.Args[0]
	switch m {
	case "exit", "echo", "type", "pwd", "cd":
		fmt.Fprintf(com.Out, "%s is a shell builtin\n\r", m)

	default:
		exePath := GetBinaryPath(m)
		if len(exePath) > 0 {
			fmt.Fprintf(com.Out, "%s is %s\n\r", m, exePath)
		} else {
			fmt.Fprintf(com.Out, "%s: not found\n\r", m)
		}
	}
}

// redirect == 1 -> stdout
// redirect == 2 -> stderr
func (com *Com) RunBinary() error {
	proc := exec.Command(com.Main, com.Args...)

	proc.Stdin = com.In

	// if redirect == 1 -> stdout to file and print stderr message
	// if redirect == 2 -> stderr to file and print stdout message,
	if len(com.OutFilePath) > 0 {
		switch com.Redirect {
		case 1:
			file, err := OpenFile(com.OutFilePath, com.Mode)
			if err != nil {
				panic(err)
			}
			defer file.Close()

			proc.Stdout = file
			proc.Stderr = com.Out

		case 2:
			file, err := OpenFile(com.OutFilePath, com.Mode)
			if err != nil {
				panic(err)
			}
			defer file.Close()

			proc.Stdout = com.Out
			proc.Stderr = file
		}
	} else {
		proc.Stdout = com.Out
		proc.Stderr = com.Out
	}

	if err := proc.Start(); err != nil {
		return err
	}

	// ignore err from proc, handled by stderr pipe
	proc.Wait()
	return nil
}
