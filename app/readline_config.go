package main

import (
	"fmt"
	"os"

	"github.com/chzyer/readline"
)

type BellNoMatch struct {
	inner *readline.PrefixCompleter
}

func (bnm *BellNoMatch) Do(line []rune, pos int) ([][]rune, int) {
	newLine, offset := bnm.inner.Do(line, pos)
	if len(newLine) == 0 {
		fmt.Fprint(os.Stderr, "\a")
	}
	return newLine, offset
}

func GetConfig() *readline.Config {
	bnm := &BellNoMatch{
		inner: completer,
	}
	return &readline.Config{
		// Prompt:          "\033[31m»\033[0m ",
		Prompt:              "$ ",
		HistoryFile:         "/tmp/readline.tmp",
		AutoComplete:        bnm,
		InterruptPrompt:     "^C",
		EOFPrompt:           "exit",
		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	}
}

// func usage(w io.Writer) {
// 	io.WriteString(w, "commands:\n")
// 	io.WriteString(w, completer.Tree("    "))
// }

// Function constructor - constructs new function for listing given directory
// func listFiles(path string) func(string) []string {
// 	return func(line string) []string {
// 		names := make([]string, 0)
// 		files, _ := ioutil.ReadDir(path)
// 		for _, f := range files {
// 			names = append(names, f.Name())
// 		}
// 		return names
// 	}
// }

var completer = readline.NewPrefixCompleter(
	readline.PcItem("echo"),
	readline.PcItem("exit"),

// readline.PcItem("mode",
//
//	readline.PcItem("vi"),
//	readline.PcItem("emacs"),
//
// ),
// readline.PcItem("login"),
// readline.PcItem("say",
//
//	readline.PcItemDynamic(listFiles("./"),
//		readline.PcItem("with",
//			readline.PcItem("following"),
//			readline.PcItem("items"),
//		),
//	),
//	readline.PcItem("hello"),
//	readline.PcItem("bye"),
//
// ),
// readline.PcItem("setprompt"),
// readline.PcItem("setpassword"),
// readline.PcItem("bye"),
// readline.PcItem("help"),
// readline.PcItem("go",
//
//	readline.PcItem("build", readline.PcItem("-o"), readline.PcItem("-v")),
//	readline.PcItem("install",
//		readline.PcItem("-v"),
//		readline.PcItem("-vv"),
//		readline.PcItem("-vvv"),
//	),
//	readline.PcItem("test"),
//
// ),
// readline.PcItem("sleep"),
)

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}
