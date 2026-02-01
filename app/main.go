package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func isExecAny(mode os.FileMode) bool {
	// The 0111 octal bitmask checks the execute bits for owner, group, and others.
	return mode&0111 != 0
}

func ExtractCmdTxt(t string) string {
	echoTxt := strings.SplitN(t, " ", 2)[1]
	return echoTxt
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
	case "cd":
		return "cd"
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
			fmt.Println(ExtractCmdTxt(text))
		case "type":
			innerCom := ExtractCmdTxt(text)
			switch CommandRecog(innerCom) {
			case "exit", "echo", "type", "pwd", "cd":
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
				fmt.Println(dir)
			}
			fmt.Println(dir)
		case "cd":
			path := ExtractCmdTxt(text)
			if err := os.Chdir(path); err != nil {
				fmt.Println("cd:", path+": No such file or directory")
			}
		default:
			comParts := strings.Split(text, " ")
			exePath := CheckExeExistance(comParts[0])
			if len(exePath) > 0 {
				err := ExecuteExe(comParts[0], comParts[1:])
				if err != nil {
					fmt.Println("Execution Failed:", err)
				}
			} else {
				fmt.Println(text + ": command not found")
			}
		}
	}
}
