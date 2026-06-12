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

func GetComm(com string) (string, []string) {
	commParts := SplitComm(com)
	if len(commParts) == 0 {
		return "", commParts
	}
	return commParts[0], commParts[1:]
}

func GetBinaryPath(filename string) string {
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

func RunBinary(file string, args []string, outFile string) error {
	proc := exec.Command(file, args...)

	stdout, err := proc.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := proc.StderrPipe()
	if err != nil {
		return err
	}

	if err := proc.Start(); err != nil {
		return err
	}

	if len(outFile) > 0 {
		if err := HandleBinaryFileOut(outFile, stdout, stderr); err != nil {
			return err
		}
	} else {
		if err := HandleBinaryOut(stdout, stderr); err != nil {
			return err
		}
	}

	errScanner := bufio.NewScanner(stderr)
	var errStr string
	for errScanner.Scan() {
		errStr += errScanner.Text()
	}

	if err := errScanner.Err(); err != nil {
		fmt.Errorf("reading from pipe failed: %s", err)
	}

	if len(errStr) > 0 {
		return fmt.Errorf("%s", errStr)
	}

	if err := proc.Wait(); err != nil {
		return err
	}

	return nil
}

func HandleExit() {
	os.Exit(0)
}

func HandleEcho(args []string) error {
	if len(args) >= 2 && (args[len(args)-2] == ">" || args[len(args)-2] == "1>") {
		filepath := args[len(args)-1]
		args = args[:len(args)-2]
		HandleFileOut(strings.Join(args, " "), filepath)
	} else {
		fmt.Println(strings.Join(args, " "))
	}

	return nil
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
		exePath := GetBinaryPath(com)

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

	exePath := GetBinaryPath(main)

	var outFilePath string
	if len(exePath) > 0 {

		if len(args) >= 2 && (args[len(args)-2] == ">" || args[len(args)-2] == "1>") {
			outFilePath = args[len(args)-1]
			args = args[:len(args)-2]
		}

		err := RunBinary(main, args, outFilePath)
		if err != nil {
			// cat: nonexistent: No such file or directory
			fmt.Printf("%s\n", err)
		}
	} else {
		fmt.Printf("%s: command not found\n", main)
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

		main, args := GetComm(com)

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
}
