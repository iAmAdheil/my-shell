package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func ExtractEchoTxt(s string) string {
	return s[5:]
}

func CommandRecog(c string) string {
	if c == "exit" {
		return "exit"
	} else if c[:4] == "echo" {
		return "echo"
	}
	return "not found"
}

func main() {
	// TODO: Uncomment the code below to pass the first stage
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("$ ")
		// starts write and fills up reader buffer
		// once "enter" is pressed, stops, extracts text before delimiter and clears reader buffer
		in, _ := reader.ReadString('\n')
		text := strings.TrimSuffix(in, "\n")

		com := CommandRecog(text)

		switch com {
		case "exit":
			os.Exit(0)
		case "echo":
			fmt.Println(ExtractEchoTxt(text))
		default:
			fmt.Println(text + ": command not found")
		}
	}
}
