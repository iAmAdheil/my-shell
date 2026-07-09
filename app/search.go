package main

import (
	"os"
	"slices"
	"strings"
)

func Search(line string) (string, []string) {
	p := strings.Split(line, " ")

	switch len(p) {
	case 0:
		return "", []string{}
	case 1:
		txt := p[0]
		suggs := SearchPath(txt)
		if len(suggs) == 0 {
			// no changes in string and no suggs
			return "", suggs
		} else if len(suggs) > 1 {
			pref := GetComPrefix(suggs)
			// if no common prefix, just print after next tab
			if len(pref) <= len(txt) {
				// no changes in string, but suggs available
				return "", suggs
			}
			// use common pref string, suggs available
			return pref, suggs
		}
		return suggs[0] + " ", suggs

	default:
		// text that needs to be completed
		txt := p[len(p)-1]
		var (
			fText string
			suggs []SearchRes
		)

		if len(txt) == 0 {
			suggs = SearchDir("", "")
			if len(suggs) == 0 {
				return "", []string{}
			} else if len(suggs) > 1 {
				return "", NormaliseSuggs(suggs)
			}

			sug := suggs[0]
			sugTxt := sug.Name
			if sug.IsDir {
				sugTxt += "/"
			}

			p[len(p)-1] = sugTxt

			fText = strings.Join(p, " ")
			if !sug.IsDir {
				fText += " "
			}
		} else if strings.Contains(txt, "/") {
			dirs := strings.Split(txt, "/")

			ctxt := dirs[len(dirs)-1]                    // text to be completed
			pth := strings.Join(dirs[:len(dirs)-1], "/") // path excluding txt to be completed

			suggs = SearchDir(ctxt, pth)
			if len(suggs) == 0 {
				return "", []string{}
			} else if len(suggs) > 1 {
				nsuggs := NormaliseSuggs(suggs)

				pref := GetComPrefix(nsuggs)
				// if no common prefix, just print after next tab
				if len(pref) <= len(ctxt) {
					// no changes in string, but suggs available
					return "", nsuggs
				}

				p[len(p)-1] = pth + "/" + pref
				fText = strings.Join(p, " ")
				return fText, nsuggs
			}

			sug := suggs[0]
			sugTxt := sug.Name
			if sug.IsDir {
				sugTxt += "/"
			}

			p[len(p)-1] = pth + "/" + sugTxt

			fText = strings.Join(p, " ")
			if !sug.IsDir {
				fText += " "
			}

		} else {
			suggs = SearchDir(txt, "")
			if len(suggs) == 0 {
				return "", []string{}
			} else if len(suggs) > 1 {
				nsuggs := NormaliseSuggs(suggs)

				pref := GetComPrefix(nsuggs)
				// if no common prefix, just print after next tab
				if len(pref) <= len(txt) {
					// no changes in string, but suggs available
					return "", nsuggs
				}

				p[len(p)-1] = pref
				fText = strings.Join(p, " ")
				return fText, nsuggs
			}

			sug := suggs[0]
			sugTxt := sug.Name
			if sug.IsDir {
				sugTxt += "/"
			}

			p[len(p)-1] = sugTxt

			fText = strings.Join(p, " ")
			if !sug.IsDir {
				fText += " "
			}
		}

		return fText, NormaliseSuggs(suggs)
	}
}

func NormaliseSuggs(sgs []SearchRes) []string {
	var nsgs []string
	for _, v := range sgs {
		if v.IsDir {
			nsgs = append(nsgs, v.Name+"/")
		} else {
			nsgs = append(nsgs, v.Name)
		}
	}

	slices.Sort(nsgs)

	return nsgs
}

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

	slices.Sort(matches)

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

	if len(word) == 0 {
		for _, v := range dirnames {
			matches = append(matches, SearchRes{
				Name:  v,
				IsDir: true,
			})
		}
		for _, v := range filenames {
			matches = append(matches, SearchRes{
				Name:  v,
				IsDir: false,
			})
		}

		return matches
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
