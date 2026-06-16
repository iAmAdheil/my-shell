package main

import (
	"os"
	"strings"

	"golang.org/x/term"
)

// handles args and vars to handle output redirect
func RedirectFilter(args []string, mode *int, redirect *int, outFilePath *string) []string {
	if len(args) >= 2 {
		fileArg := args[len(args)-2]
		if strings.Contains(fileArg, ">") || strings.Contains(fileArg, "1>") || strings.Contains(fileArg, "2>") {
			if strings.Count(fileArg, ">") == 2 {
				*mode = 1
			}

			if fileArg[0:1] == ">" || fileArg[0:2] == "1>" {
				*redirect = 1
			} else if fileArg[0:2] == "2>" {
				*redirect = 2
			}

			*outFilePath = args[len(args)-1]
			args = args[:len(args)-2]
		}
	}

	return args
}

// switches terminal from cooked to raw mode
func EnableRaw() *term.State {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}

	return oldState
}

// checks if a binary has executable permissions
func IsExecAny(mode os.FileMode) bool {
	// The 0111 octal bitmask checks the execute bits for owner, group, and others.
	return mode&0111 != 0
}
