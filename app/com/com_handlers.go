package com

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

func (com *Com) Stop() {
	// wait only unblocks when both the ends for the write end of a pipe close
	// when previous proc's wait unblocks (their own pipes close) and com.Close runs -> which closes the pipe
	if com.IsBgProc {
		fmt.Printf("[1] %v\n", com.Proc.Process.Pid)
		return
	}

	com.Proc.Wait()
	if com.Close {
		err := com.Out.Close()
		if err != nil {
			panic(err)
		}
	}
}

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
	case "history":
		com.HandleHistory()
	case "jobs":
		com.HandleJobs()
	default:
		exePath := GetBinaryPath(com.Main)
		if len(exePath) > 0 {
			err := com.RunBinary()
			if err != nil {
				// cat: nonexistent: No such file or directory
				fmt.Printf("%s\n", err)
			}
		} else {
			fmt.Printf("%s: command not found\n", com.Main)
		}
	}
}

func (com *Com) HandleExit() {
	os.Exit(0)
}

func (com *Com) HandlePwd() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(com.Out, "%s\n", err)
		return
	}
	fmt.Fprintf(com.Out, "%s\n", dir)
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
		fmt.Fprintf(com.Out, "cd: %s: No such file or directory\n", path)
	}
}

var History []string

func (com *Com) HandleHistory() error {
	var lines []string = History

	if len(com.Args) > 0 {
		switch {
		case com.Args[0] == "-r":
			filename := com.Args[1]
			if len(filename) == 0 {
				goto def
			}

			return HandleHistoryRead(filename)
		case com.Args[0] == "-w":
			filename := com.Args[1]
			if len(filename) == 0 {
				goto def
			}

			return HandleHistoryWrite(filename)
		case com.Args[0] == "-a":
			filename := com.Args[1]
			if len(filename) == 0 {
				goto def
			}

			return HandleHistoryAppend(filename)
		default:
			c, err := strconv.Atoi(com.Args[0])
			if err != nil || c < 0 {
				goto def
			}

			for i, line := range lines {
				if i >= len(lines)-c {
					fmt.Fprintf(com.Out, "%v %s\n", i+1, line)
				}

			}

			return nil
		}
	}

def:
	for i, line := range lines {
		fmt.Fprintf(com.Out, "%v %s\n", i+1, line)
	}

	return nil
}

func (com *Com) HandleEcho() error {
	wg := &sync.WaitGroup{}

	var txt string

	if com.Args[0] == "-e" {
		txt = strings.Join(com.Args[1:], " ")
		txt = strings.ReplaceAll(txt, `\n`, "\n")
	} else {
		txt = strings.Join(com.Args, " ")
	}

	var (
		out     string = txt
		fileout string = txt
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
	if (com.Redirect == 0 || com.Redirect == 2) || com.Close {
		fmt.Fprintf(com.Out, "%s\n", out)
	}

	return nil
}

func (com *Com) HandleType() {
	if len(com.Args) == 0 {
		return
	}

	m := com.Args[0]
	switch m {
	case "exit", "echo", "type", "pwd", "cd", "history", "jobs":
		fmt.Fprintf(com.Out, "%s is a shell builtin\n", m)

	default:
		exePath := GetBinaryPath(m)
		if len(exePath) > 0 {
			fmt.Fprintf(com.Out, "%s is %s\n", m, exePath)
		} else {
			fmt.Fprintf(com.Out, "%s: not found\n", m)
		}
	}
}

func (com *Com) HandleJobs() {}

// redirect == 1 -> stdout
// redirect == 2 -> stderr
func (com *Com) RunBinary() error {
	com.Proc.Stdin = com.In

	// if redirect == 1 -> stdout to file and print stderr message
	// if redirect == 2 -> stderr to file and print stdout message,
	if len(com.OutFilePath) > 0 {
		file, err := OpenFile(com.OutFilePath, com.Mode)
		if err != nil {
			panic(err)
		}

		switch com.Redirect {
		case 1:
			com.Proc.Stdout = file
			com.Proc.Stderr = com.Out
		case 2:
			com.Proc.Stdout = com.Out
			com.Proc.Stderr = file
		}
	} else {
		com.Proc.Stdout = com.Out
		com.Proc.Stderr = com.Out
	}

	if err := com.Proc.Start(); err != nil {
		return err
	}

	return nil
}
