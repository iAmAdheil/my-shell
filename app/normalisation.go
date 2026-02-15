package main

import "strings"

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
