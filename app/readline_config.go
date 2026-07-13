package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/chzyer/readline"
	"github.com/codecrafters-io/shell-starter-go/app/com"
)

func handleCompleter(rline []rune) string {
	line := string(rline)

	// use main to check for an existing entry
	// main & args to be passed to proc as arguments
	main, args := GetComm(line)

	exepath, ok := com.CEntries[strings.TrimSpace(main)]
	if !ok {
		return ""
	}

	cw := NewCustomWriter()
	n := len(args)

	arg1 := main
	var arg2 string = "" // word to be completed
	var arg3 string = "" // second last arg (before word to be completed)

	if len(args) > 0 {
		arg2 = args[n-1]
	}
	if len(args) > 1 {
		arg3 = args[n-2]
	}

	cargs := []string{arg1, arg2, arg3}

	c := &com.Com{
		Proc: exec.Command(exepath, cargs...),
		In:   nil,
		Out:  cw,
	}

	c.RunBinary()
	c.Stop()

	newLine := strings.Trim(cw.buf.String(), "\n")
	if len(newLine) == 0 || len(newLine) <= len(arg2) {
		return ""
	}

	return newLine[len(arg2):] + " "
}

// implements the readline.AutoCompleter interface
type BellNoMatch struct {
	inner *readline.PrefixCompleter
}

// offset -> offset on the left for new completion
func (bnm *BellNoMatch) Do(line []rune, pos int) ([][]rune, int) {
	compLine := handleCompleter(line)
	if len(compLine) > 0 {
		nr := make([][]rune, 1)
		r := []rune(compLine)
		nr[0] = r

		return nr, len(line)
	}

	newLine, offset := bnm.inner.Do(line, pos)
	if len(newLine) > 0 {
		return newLine, offset
	}

	sug := getPathSugg(string(line))
	if len(sug) > 0 && len(sug) > len(string(line)) {
		nr := make([][]rune, 1)
		r := []rune(sug)
		nr[0] = r[len(line):]

		checkPath = false

		return nr, len(line)
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
		_, suggs := Search(string(line))
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
	res, _ := Search(line)
	return res
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
