package com

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func CheckAlphaNumUS(word string) bool { // allow letters, numbers and underscore
	return regexp.MustCompile(`^[a-zA-Z0-9_]*$`).MatchString(word)
}

func CheckNum(word string) bool {
	return regexp.MustCompile(`^[0-9]*$`).MatchString(word)
}

func GetBinaryPath(filename string) string {
	path := os.Getenv("PATH")
	if len(path) == 0 {
		fmt.Printf("no 'PATH' env variable set\n\r")
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

// checks if a binary has executable permissions
func IsExecAny(mode os.FileMode) bool {
	// The 0111 octal bitmask checks the execute bits for owner, group, and others.
	return mode&0111 != 0
}
