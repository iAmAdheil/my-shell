package main

import (
	"bytes"
	"fmt"
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
