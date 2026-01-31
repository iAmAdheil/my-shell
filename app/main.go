package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ExtractEchoTxt(s string) string {
	return s[5:]
}

func ExtractTypeTxt(s string) string {
	return s[5:]
}

func CommandRecog(c string) string {
	if c == "exit" {
		return "exit"
	} else if c[:4] == "echo" {
		return "echo"
	} else if c[:4] == "type" {
		return "type"
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
		case "type":
			innerCom := text[5:]
			switch CommandRecog(innerCom) {
			case "exit", "echo", "type":
				fmt.Println(innerCom, "is a shell builtin")
			default:
				fmt.Println(innerCom + ": not found")
			}
		default:
			fmt.Println(text + ": command not found")
		}
	}
}
