package main

import (
	"os"
	"strings"
)

// search path for executable exes that match the passed word arg
func SearchPath(word string) []string {
	matches := []string{}
	path := os.Getenv("PATH")
	if len(path) == 0 {
		return matches
	}
	dirs := strings.Split(path, ":")
	for _, dir := range dirs {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		exes := []string{}
		for _, file := range files {
			fileFP := dir + "/" + file.Name()
			fileInfo, err := os.Stat(fileFP)
			if err != nil {
				break
			}
			if IsExecAny(fileInfo.Mode().Perm()) { // check file permissions
				exes = append(exes, file.Name())
			} else {
				break
			}
		}
		root := InitTrie(exes)
		cur_matches := root.Complete(word)
		matches = append(matches, cur_matches...)
	}
	return matches
}

// search path for executable exes that match the passed word arg
func SearchCurDir(word string) []string {
	filenames := []string{}

	dir, err := os.Getwd()
	if err != nil {
		return []string{}
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return []string{}
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			filenames = append(filenames, entry.Name())
		}
	}

	root := InitTrie(filenames)
	matches := root.Complete(word)

	return matches
}
