// "hello"  "script's"  test""example -> test on this
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func IsExecAny(mode os.FileMode) bool {
	// The 0111 octal bitmask checks the execute bits for owner, group, and others.
	return mode&0111 != 0
}

func GetCommand(com string) (string, []string) {
	comParts := GetNormalCom(com)
	if len(comParts) == 0 {
		return "", comParts
	}
	return comParts[0], comParts[1:]
}

func GetExePath(filename string) string {
	path := os.Getenv("PATH")
	if len(path) == 0 {
		fmt.Println("no 'PATH' env variable set")
	}

	dirs := strings.Split(path, ":")

	for _, dir := range dirs {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, file := range files {
			if file.Name() == filename {
				fullPath := dir + "/" + filename
				fileInfo, err := os.Stat(fullPath)
				if err != nil {
					break
				}
				if IsExecAny(fileInfo.Mode().Perm()) { // check file permissions
					return fullPath
				} else {
					break
				}
			}
		}
	}
	return ""
}

func RunExe(filename string, args []string) error {
	proc := exec.Command(filename, args...)
	stdout, err := proc.StdoutPipe()
	if err != nil {
		return err
	}
	if err := proc.Start(); err != nil {
		return err
	}
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if err := proc.Wait(); err != nil {
		return err
	}
	return nil
}

func HandleExit() {
	os.Exit(0)
}

func HandleEcho(args []string) {
	fmt.Println(strings.Join(args, " "))
}

func HandleType(args []string) {
	// normalised args
	if len(args) == 0 {
		return
	}
	com := args[0]
	switch com {
	case "exit", "echo", "type", "pwd", "cd":
		fmt.Println(com, "is a shell builtin")
	default:
		exePath := GetExePath(com)

		if len(exePath) > 0 {
			fmt.Println(com + " is " + exePath)
		} else {
			fmt.Println(com + ": not found")
		}
	}
}

func HandlePwd() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(dir)
}

func HandleCd(args []string) {
	if len(args) == 0 {
		return
	}
	path := args[0]
	if path == "~" {
		path = os.Getenv("HOME")
	}
	if err := os.Chdir(path); err != nil {
		fmt.Println("cd:", path+": No such file or directory")
	}
}

func HandleDefault(main string, args []string) {
	if len(main) == 0 && len(args) == 0 {
		return
	}

	exePath := GetExePath(main)

	if len(exePath) > 0 {
		err := RunExe(main, args)
		if err != nil {
			fmt.Println("execution failed:", err)
		}
	} else {
		fmt.Println(main + ": command not found")
	}
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

		main, args := GetCommand(com)

		switch main {
		case "exit":
			HandleExit()
		case "echo":
			HandleEcho(args)
		case "type":
			HandleType(args)
		case "pwd":
			HandlePwd()
		case "cd":
			HandleCd(args)
		default:
			HandleDefault(main, args)
		}
	}

	//TESTING ENV
	// str := "\"/tmp/ant/f 60\" \"/tmp/ant/f   54\" \"/tmp/ant/f's98\""
	// res := GetArgs(str)
	// for i, v := range res {
	// 	fmt.Println("index:", i, "value:", v)
	// }
}
