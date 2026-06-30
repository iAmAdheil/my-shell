package main

import (
	"os"
	"strings"

	"golang.org/x/term"
)

// handles args and vars to handle output redirect
func RedirectFilter(args []string, outFilePath *string, redirect *int, mode *int) []string {
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

// search path for executable exes that match the passed word arg
func SearchPath(word string) []string {
	matches := []string{}
	path := os.Getenv("PATH")
	if len(path) == 0 {
		return matches
	}
	dirs := strings.Split(path, ":")
	for _, dir := range dirs {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		exes := []string{}
		for _, file := range files {
			fileFP := dir + "/" + file.Name()
			fileInfo, err := os.Stat(fileFP)
			if err != nil {
				break
			}
			if IsExecAny(fileInfo.Mode().Perm()) { // check file permissions
				exes = append(exes, file.Name())
			} else {
				break
			}
		}
		root := InitTrie(exes)
		cur_matches := root.Complete(word)
		matches = append(matches, cur_matches...)
	}
	return matches
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
