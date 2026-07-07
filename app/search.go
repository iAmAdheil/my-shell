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

type SearchRes struct {
	Name  string
	IsDir bool
}

// search path for executable exes that match the passed word arg
func SearchDir(word string, path string) []SearchRes {
	matches := []SearchRes{}
	filenames := []string{}
	dirnames := []string{}

	if len(path) == 0 {
		dir, err := os.Getwd()
		if err != nil {
			return []SearchRes{}
		}
		path = dir
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return []SearchRes{}
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			filenames = append(filenames, entry.Name())
		} else {
			dirnames = append(dirnames, entry.Name())
		}
	}

	froot := InitTrie(filenames)
	fmatches := froot.Complete(word)

	droot := InitTrie(dirnames)
	dmatches := droot.Complete(word)

	for _, v := range dmatches {
		matches = append(matches, SearchRes{
			Name:  v,
			IsDir: true,
		})
	}
	for _, v := range fmatches {
		matches = append(matches, SearchRes{
			Name:  v,
			IsDir: false,
		})
	}

	return matches
}
