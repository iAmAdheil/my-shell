package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"slices"

	"github.com/chzyer/readline"
	"github.com/codecrafters-io/shell-starter-go/app/com"
)

// implements the readline.AutoCompleter interface
type BellNoMatch struct {
	inner *readline.PrefixCompleter
}

func (bnm *BellNoMatch) Do(line []rune, pos int) ([][]rune, int) {
	newLine, offset := bnm.inner.Do(line, pos)
	// for i, v := range newLine {
	//  fmt.Printf("\nnewline %v: %s and %v\n", i, string(v), len(string(v)))
	// }

	if len(newLine) > 0 {
		return newLine, offset
	}
	sug := getPathSugg(string(line))
	if len(sug) > 0 {
		nr := make([][]rune, 1)
		r := []rune(sug)
		nr[0] = r[len(line):]
		return nr, offset
	}

	fmt.Fprint(os.Stderr, "\a")
	if !checkPath {
		checkPath = true
	}

	return newLine, offset
}

var (
	checkPath bool = false
)

type MyListener struct{}

func (ml *MyListener) OnChange(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
	// obstructs logs from query execution
	// returning true prints prompt
	// return false when "return" is pressed
	// fmt.Println("\n", prevTab, checkPath)
	// fmt.Println("\n", key)
	if key == 13 || key == 10 {
		return line, pos, false
	}
	if checkPath && key == 9 {
		suggs := listPathBinaries(string(line))
		slices.Sort(suggs)
		if len(suggs) > 0 {
			fmt.Printf("\n")
			for _, sug := range suggs {
				fmt.Printf("%s  ", sug)
			}
			fmt.Printf("\n")
		}
	} else {
		checkPath = false
	}
	return line, pos, true
}

func GetConfig() *readline.Config {
	var completer = readline.NewPrefixCompleter(
		readline.PcItem("echo"),
		readline.PcItem("exit"),
		// readline.PcItemDynamic(listPathBinaries()),
		// readline.PcItem("mode",
		//
		//  readline.PcItem("vi"),
		//  readline.PcItem("emacs"),
		//
		// ),
		// readline.PcItem("login"),
		// readline.PcItem("say",
		//  readline.PcItemDynamic(listFiles("./"),
		//      readline.PcItem("with",
		//          readline.PcItem("following"),
		//          readline.PcItem("items"),
		//      ),
		//  ),
		//  readline.PcItem("hello"),
		//  readline.PcItem("bye"),
		// ),
	// readline.PcItem("setprompt"),
	// readline.PcItem("setpassword"),
	// readline.PcItem("bye"),
	// readline.PcItem("help"),
	// readline.PcItem("go",
	//
	//  readline.PcItem("build", readline.PcItem("-o"), readline.PcItem("-v")),
	//  readline.PcItem("install",
	//      readline.PcItem("-v"),
	//      readline.PcItem("-vv"),
	//      readline.PcItem("-vvv"),
	//  ),
	//  readline.PcItem("test"),
	//
	// ),
	// readline.PcItem("sleep"),
	)

	bnm := &BellNoMatch{
		inner: completer,
	}

	return &readline.Config{
		// Prompt:          "\033[31m»\033[0m ",
		Prompt:          "$ ",
		Listener:        &MyListener{},
		AutoComplete:    bnm,
		HistoryFile:     com.HISTORY_FILE,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	}
}

// func usage(w io.Writer) {
//  io.WriteString(w, "commands:\n")
//  io.WriteString(w, completer.Tree("    "))
// }

// Function constructor - constructs new function for listing given directory
// func listFiles(path string) func(string) []string {
//  return func(line string) []string {
//      names := make([]string, 0)
//      files, _ := ioutil.ReadDir(path)
//      for _, f := range files {
//          names = append(names, f.Name())
//      }
//      return names
//  }
// }

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	case readline.CharTab:
		return r, true
	}
	return r, true
}

// tab handler
func getPathSugg(line string) string {
	suggs := SearchPath(line)
	if len(suggs) == 0 {
		return ""
	} else if len(suggs) > 1 {
		pref := GetComPrefix(suggs)
		if len(pref) <= len(line) {
			return ""
		}
		return pref
	}
	return suggs[0] + " "
}

// tab handler
func listPathBinaries(line string) []string {
	suggs := SearchPath(line)
	return suggs
}

func listFiles(path string) func(string) []string {
	return func(line string) []string {
		names := make([]string, 0)
		files, _ := ioutil.ReadDir(path)
		for _, f := range files {
			names = append(names, f.Name())
		}
		return names
	}
}
