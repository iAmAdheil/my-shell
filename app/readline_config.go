package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"slices"
	"strings"

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
	if len(sug) > 0 && len(sug) > len(line) {
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
	second    bool = false
)

type MyListener struct{}

func (ml *MyListener) OnChange(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
	// obstructs logs from query execution
	// returning true prints prompt
	// return false when "return" is pressed
	if key == 13 || key == 10 {
		return line, pos, false
	}
	if checkPath && second && key == 9 {
		suggs := SearchPath(string(line))
		slices.Sort(suggs)
		if len(suggs) > 0 {
			fmt.Printf("\n")
			for _, sug := range suggs {
				fmt.Printf("%s  ", sug)
			}
			fmt.Printf("\n")
		}
	}

	if !second && key == 9 {
		second = true
	}

	if key != 9 {
		checkPath = false
		second = false
	}
	return line, pos, true
}

// runs on startup
func Init() {
	filename := os.Getenv("HISTFILE")
	if len(filename) == 0 {
		return
	}

	// loads history on start
	err := com.HandleHistoryRead(filename)
	if err != nil {
		panic(err)
	}
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

	histfile := os.Getenv("HISTFILE")
	if len(histfile) == 0 {
		histfile = com.HISTORY_FILE
	}

	// loads history on start
	err := com.HandleHistoryRead(histfile)
	if err != nil {
		panic(err)
	}

	return &readline.Config{
		// Prompt:          "\033[31m»\033[0m ",
		Prompt:          "$ ",
		Listener:        &MyListener{},
		AutoComplete:    bnm,
		HistoryFile:     histfile,
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
	}
	return r, true
}

// tab handler
func getPathSugg(line string) string {
	p := strings.Split(line, " ")

	switch len(p) {
	case 0:
		return ""
	case 1:
		txt := p[0]
		suggs := SearchPath(txt)
		if len(suggs) == 0 {
			return ""
		} else if len(suggs) > 1 {
			pref := GetComPrefix(suggs)
			// if no common prefix, just print after next tab
			if len(pref) <= len(txt) {
				return ""
			}
			return pref
		}
		return suggs[0] + " "

	default:
		// text that needs to be completed
		txt := p[len(p)-1]

		if strings.Contains(txt, "/") {
			dirs := strings.Split(txt, "/")

			ctxt := dirs[len(dirs)-1]
			pth := strings.Join(dirs[:len(dirs)-1], "/")

			suggs := SearchDir(ctxt, pth)
			if len(suggs) == 0 {
				return ""
			}
			p[len(p)-1] = pth + "/" + suggs[0]

		} else {
			suggs := SearchDir(txt, "")
			if len(suggs) == 0 {
				return ""
			}
			p[len(p)-1] = suggs[0]
		}

		fText := strings.Join(p, " ")
		return fText + " "
	}

	return ""
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
