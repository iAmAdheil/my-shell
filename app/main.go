// "hello"  "script's"  test""example -> test on this
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func GetComm(com string) (string, []string) {
	commParts := SplitComm(com)
	if len(commParts) == 0 {
		return "", commParts
	}
	return commParts[0], commParts[1:]
}

func main() {
	// TODO: Uncomment the code below to pass the first stage
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		// starts write and fills up reader buffer
		// once "enter" is pressed, stops, extracts text before delimiter and clears reader buffer
		in, _ := reader.ReadString('\n')
		com := strings.TrimSuffix(in, "\n")

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
	}
}
