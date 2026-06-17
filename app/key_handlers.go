package main

import (
	"bytes"
	"fmt"
	"unicode/utf8"
)

func Return(txt *bytes.Buffer) {
	fmt.Printf("\n\r")
	com := txt.String()
	fmt.Printf("command: %s\n\r", com)
	txt.Reset()

	main, args := GetComm(com)
	var (
		outFilePath string
		redirect    int = 0
		mode        int = 0
	)
	args = RedirectFilter(args, &mode, &redirect, &outFilePath)

	switch main {
	case "exit":
		HandleExit()
	case "echo":
		HandleEcho(args, outFilePath, redirect, mode)
	case "type":
		HandleType(args)
	case "pwd":
		HandlePwd()
	case "cd":
		HandleCd(args)
	default:
		HandleDefault(main, args, outFilePath, redirect, mode)
	}

	fmt.Printf("$ ")
}

var words = []string{"echo", "exit"}

func Tab(txt string) string {
	root := InitTrie(words)
	suggs := root.Complete(txt)
	if len(suggs) > 0 {
		return suggs[0]
	}
	return ""
}

func Delete(txt *bytes.Buffer) {
	fmt.Print("\b \b")
	if txt.Len() > 0 {
		b := txt.Bytes()
		_, size := utf8.DecodeLastRune(b) // size = bytes in the final rune
		txt.Truncate(txt.Len() - size)
	}
}

func Clear(txt *bytes.Buffer) {
	// clear everything from start of "$ " to end
	n := utf8.RuneCountInString(txt.String())
	if n > 0 {
		fmt.Printf("\033[%dD", n) // move cursor left n columns
	}
	fmt.Print("\033[K")
	txt.Reset()
}
