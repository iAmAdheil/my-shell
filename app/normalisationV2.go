package main

import (
	"strings"
)

func HandleNormalisationV2(str string) []string {
	args := []string{}
	temp := []byte{}

	i := 0
	for i < len(str) {
		chr := str[i]

		if chr == ' ' {
			if len(temp) > 0 {
				arg := string(temp)
				args = append(args, arg)
				// clear character buffer
				temp = []byte{}
			}
			// move to next valid character
			j := i
			for ; j < len(str); j++ {
				if str[j] != ' ' {
					break
				}
			}
			i = j
			args = append(args, " ")
		} else if chr == '\\' {
			temp = append(temp, str[i+1])
			i += 2
		} else if chr == '\'' {
			j := i + 1
			for ; j < len(str); j++ {
				if str[j] == '\'' {
					break
				}
				temp = append(temp, str[j])
			}
			arg := string(temp)
			args = append(args, arg)

			i = j + 1
			temp = []byte{}
		} else if chr == '"' {
			j := i + 1
			for ; j < len(str); j++ {
				if str[j] == '"' {
					break
				} else if str[j] == '\\' {
					temp = append(temp, str[j+1])
					// only increment once as loop will also increment j
					j += 1
				} else {
					temp = append(temp, str[j])
				}
			}
			arg := string(temp)
			args = append(args, arg)

			i = j + 1
			temp = []byte{}
		} else {
			temp = append(temp, chr)
			i += 1
		}
	}

	// add the last arg after the string has ended
	if len(temp) > 0 {
		arg := string(temp)
		args = append(args, arg)
	}

	return args
}

func GetArgs(str string) []string {
	fArgs := []string{}
	args := HandleNormalisationV2(str)
	temp := []string{}

	for _, arg := range args {
		if arg == " " {
			fArgs = append(fArgs, strings.Join(temp, ""))
			// clear temp once its contents have been added to fArgs
			temp = []string{}
		} else {
			temp = append(temp, arg)
		}
	}
	// add the last item added to temp
	if len(temp) > 0 {
		fArgs = append(fArgs, strings.Join(temp, ""))
	}

	return fArgs
}
