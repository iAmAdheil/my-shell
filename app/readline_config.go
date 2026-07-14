package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
	"github.com/codecrafters-io/shell-starter-go/app/com"
)

func handleCompleter(rline []rune) (string, []string) {
	line := string(rline)

	// use main to check for an existing entry
	// main & args to be passed to proc as arguments
	args := strings.Split(line, " ")
	main := args[0]

	exepath, ok := com.CEntries[strings.TrimSpace(main)]
	if !ok {
		return "", []string{}
	}

	cw := NewCustomWriter()
	n := len(args)

	arg1 := main
	var arg2 string = "" // word to be completed
	var arg3 string = "" // second last arg (before word to be completed)

	if n > 0 {
		arg2 = args[n-1]
	}
	if n > 1 {
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

	newLine := strings.Split(cw.buf.String(), "\n")
	nlLen := len(newLine)
	// remove last string if empty
	if len(newLine[nlLen-1]) == 0 {
		newLine = newLine[:nlLen-1]
	}

	if len(newLine) > 1 {
		var opts []string

		for _, opt := range newLine {
			opts = append(opts, opt)
		}

		pref := GetComPrefix(opts)
		// if no common prefix, just print after next tab
		if len(pref) == 0 || len(pref) <= len(arg2) {
			// no changes in string, but suggs available
			return "", opts
		}

		return pref[len(arg2):], opts
	} else if len(newLine) == 0 {
		return "", []string{}
	}

	nw := newLine[0]
	if len(nw) == 0 || len(nw) <= len(arg2) {
		return "", []string{}
	}

	return nw[len(arg2):] + " ", []string{}
}

func handleSetupEnvVars(line []rune, pos int) {
	lineStr := string(line)

	if err := os.Setenv("COMP_POINT", strconv.Itoa(pos)); err != nil {
		fmt.Printf("failed to setup comp_point env.")
	}
	if err := os.Setenv("COMP_LINE", lineStr); err != nil {
		fmt.Printf("failed to setup comp_line env.")
	}
}

// implements the readline.AutoCompleter interface
type BellNoMatch struct {
	inner *readline.PrefixCompleter
}

// offset -> offset on the left for new completion
func (bnm *BellNoMatch) Do(line []rune, pos int) ([][]rune, int) {
	handleSetupEnvVars(line, pos)

	compLine, _ := handleCompleter(line)
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

func PrintOpts(opts []string) {
	fmt.Printf("\n")
	for _, opt := range opts {
		fmt.Printf("%s  ", opt)
	}
	fmt.Printf("\n")
}

func (ml *MyListener) OnChange(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
	// obstructs logs from query execution
	// returning true prints prompt
	// return false when "return" is pressed
	if key == 13 || key == 10 {
		return line, pos, false
	}
	if checkPath && second && key == 9 {
		_, compsuggs := handleCompleter(line)
		if len(compsuggs) > 0 {
			PrintOpts(compsuggs)
			goto esc
		}

		_, suggs := Search(string(line))
		if len(suggs) > 0 {
			slices.Sort(suggs)
			PrintOpts(suggs)
			goto esc
		}
	}

	if !second && key == 9 {
		second = true
	}

	if key != 9 {
		checkPath = false
		second = false
	}
esc:
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
