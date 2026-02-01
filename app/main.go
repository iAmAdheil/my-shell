package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func isExecAny(mode os.FileMode) bool {
	// The 0111 octal bitmask checks the execute bits for owner, group, and others.
	return mode&0111 != 0
}

func ExtractEchoTxt(t string) string {
	echoTxt := strings.SplitN(t, " ", 2)[1]
	return echoTxt
}

func ExtractTypeTxt(t string) string {
	typeTxt := strings.SplitN(t, " ", 2)[1]
	return typeTxt
}

func CommandRecog(t string) string {
	com := strings.Split(t, " ")[0]

	switch com {
	case "exit":
		return "exit"
	case "echo":
		return "echo"
	case "type":
		return "type"
	case "pwd":
		return "pwd"
	default:
		return "not builtin"
	}
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

func ExecuteExe(filename string, args []string) error {
	process := exec.Command(filename, args...)
	stdout, err := process.StdoutPipe()
	if err != nil {
		return err
	}
	if err := process.Start(); err != nil {
		return err
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		// Process the line immediately as it arrives
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if err := process.Wait(); err != nil {
		return err
	}
	return nil
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
			innerCom := strings.SplitN(text, " ", 2)[1]
			switch CommandRecog(innerCom) {
			case "exit", "echo", "type", "pwd":
				fmt.Println(innerCom, "is a shell builtin")
			default:
				exePath := CheckExeExistance(innerCom)

				if len(exePath) > 0 {
					fmt.Println(innerCom + " is " + exePath)
				} else {
					fmt.Println(innerCom + ": not found")
				}
			}
		case "pwd":
			dir, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(dir)
		default:
			comParts := strings.Split(text, " ")
			exePath := CheckExeExistance(comParts[0])
			if len(exePath) > 0 {
				err := ExecuteExe(comParts[0], comParts[1:])
				if err != nil {
					log.Fatal("Execution Failed:", err)
				}
			} else {
				fmt.Println(text + ": command not found")
			}
		}
	}
}
