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

func ExtractCmdArgs(t string) string {
	echoTxt := strings.SplitN(t, " ", 2)[1]
	resultStr := HandleNormalisation(echoTxt)
	return strings.Join(resultStr, " ")
}

func CommandRecog(t string) string {
	com := strings.Split(t, " ")[0]

	switch com {
	case "exit", "echo", "type", "pwd", "cd":
		return com
	default:
		return "not builtin"
	}
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

func HandleEcho(com string) {
	fmt.Println(ExtractCmdArgs(com))
}

func HandleType(com string) {
	arg := ExtractCmdArgs(com)
	switch CommandRecog(arg) {
	case "exit", "echo", "type", "pwd", "cd":
		fmt.Println(arg, "is a shell builtin")
	default:
		exePath := GetExePath(arg)

		if len(exePath) > 0 {
			fmt.Println(arg + " is " + exePath)
		} else {
			fmt.Println(arg + ": not found")
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

func HandleCd(com string) {
	path := ExtractCmdArgs(com)
	if path == "~" {
		path = os.Getenv("HOME")
	}
	if err := os.Chdir(path); err != nil {
		fmt.Println("cd:", path+": No such file or directory")
	}
}

func HandleDefault(com string) {
	if len(com) == 0 {
		return
	}
	parts := strings.SplitN(com, " ", 2)
	main := parts[0]
	exePath := GetExePath(main)
	if len(exePath) > 0 {
		if len(parts) > 1 {
			args := parts[1]
			err := RunExe(main, HandleNormalisation(args))
			if err != nil {
				fmt.Println("execution failed:", err)
			}
		} else if len(parts) == 1 {
			err := RunExe(main, HandleNormalisation(""))
			if err != nil {
				fmt.Println("execution failed:", err)
			}
		}
	} else {
		fmt.Println(com + ": command not found")
	}
}

func HandleQuoteRemoval(argString string) string {
	runes := []rune(argString)
	clean := make([]rune, 0, len(runes))

	for i := 0; i < len(runes); i++ {
		if (runes[i] == '"' || runes[i] == '\'') && i+1 < len(runes) && (runes[i+1] == '"' || runes[i+1] == '\'') {
			i++
			continue
		}

		clean = append(clean, runes[i])
	}

	return string(clean)
}

func HandleNormalisation(argsString string) []string {
	if len(argsString) == 0 {
		return []string{}
	}

	spl := strings.Split(argsString, " ")
	temp := make([]string, 0)
	argsSl := make([]string, 0)

	i := 0
	for i < len(spl) {
		temp = []string{}
		if len(strings.TrimSpace(spl[i])) == 0 {
			i += 1
		} else if spl[i][0] == '"' {
			if spl[i][len(spl[i])-1] == '"' {
				argsSl = append(argsSl, HandleQuoteRemoval(spl[i][1:len(spl[i])-1]))
				i += 1
			} else {
				temp = append(temp, spl[i][1:])
				for j := i + 1; j < len(spl); j++ {
					if len(spl[j]) > 0 && spl[j][len(spl[j])-1] == '"' {
						temp = append(temp, spl[j][:len(spl[j])-1])
						arg := strings.Join(temp, " ")
						argsSl = append(argsSl, arg)
						i = j + 1
						break
					}
					temp = append(temp, spl[j])
				}
			}
		} else if spl[i][0] == '\'' {
			if spl[i][len(spl[i])-1] == '\'' {
				argsSl = append(argsSl, HandleQuoteRemoval(spl[i][1:len(spl[i])-1]))
				i += 1
			} else {
				temp = append(temp, spl[i][1:])
				for j := i + 1; j < len(spl); j++ {
					if len(spl[j]) > 0 && spl[j][len(spl[j])-1] == '\'' {
						temp = append(temp, spl[j][:len(spl[j])-1])
						arg := strings.Join(temp, " ")
						argsSl = append(argsSl, arg)
						i = j + 1
						break
					}
					temp = append(temp, spl[j])
				}
			}
		} else {
			argsSl = append(argsSl, HandleQuoteRemoval(spl[i]))
			i += 1
		}
	}

	return argsSl
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

		comPref := CommandRecog(com)

		switch comPref {
		case "exit":
			HandleExit()
		case "echo":
			HandleEcho(com)
		case "type":
			HandleType(com)
		case "pwd":
			HandlePwd()
		case "cd":
			HandleCd(com)
		default:
			HandleDefault(com)
		}
	}

	//TESTING ENV
	// str := "\"/tmp/ant/f 60\" \"/tmp/ant/f   54\" \"/tmp/ant/f's98\""
	// res := HandleNormalisation(str)
	// for i, v := range res {
	// 	fmt.Println("index:", i, "value:", v)
	// }
}
