package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func isExecAny(mode os.FileMode) bool {
	// The 0111 octal bitmask checks the execute bits for owner, group, and others.
	return mode&0111 != 0
}

func ExtractEchoTxt(s string) string {
	if len(s) <= 5 {
		return ""
	}
	return s[5:]
}

func ExtractTypeTxt(s string) string {
	if len(s) <= 5 {
		return ""
	}
	return s[5:]
}

func CommandRecog(c string) string {
	if c == "exit" {
		return "exit"
	} else if len(c) >= 4 && c[:4] == "echo" {
		return "echo"
	} else if len(c) >= 4 && c[:4] == "type" {
		return "type"
	}
	return "not found"
}

func CheckExeExistance(filename string) string {
	path := os.Getenv("PATH")
	dirs := strings.Split(path, ":")

	for _, dir := range dirs {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, file := range files {
			if file.Name() == filename {
				fullPath := dir + "/" + filename
				fileInfo, err := os.Stat(fullPath) // Replace "myprogram" with the file path
				if err != nil {
					break
				}
				if isExecAny(fileInfo.Mode().Perm()) {
					return fullPath
				} else {
					break
				}
			}
		}
	}
	return ""
}

func ExecuteExe(filename string, args []string, nArgs int) (string, string, error) {
	cmd := exec.Command(filename, args...)
	var stdoutbuf, stderrbuf bytes.Buffer
	cmd.Stdout = &stdoutbuf
	cmd.Stderr = &stderrbuf
	if err := cmd.Run(); err != nil {
		return "", "", err
	}
	return stdoutbuf.String(), stderrbuf.String(), nil
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
				exePath := CheckExeExistance(innerCom)

				if len(exePath) > 0 {
					fmt.Println(innerCom + " is " + exePath)
				} else {
					fmt.Println(innerCom + ": not found")
				}
			}
		default:
			comParts := strings.Split(text, " ")
			exePath := CheckExeExistance(comParts[0])
			if len(exePath) > 0 {
				logs, _, _ := ExecuteExe(comParts[0], comParts[1:], len(comParts)-1)
				fmt.Print(logs)
				// fmt.Println("stderr:", b)
			} else {
				fmt.Println(text + ": command not found")
			}
		}
	}
}
