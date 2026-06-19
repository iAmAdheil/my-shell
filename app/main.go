// // "hello"  "script's"  test""example -> test on this
// package main

// import (
// 	"fmt"
// 	"os"

// 	"github.com/chzyer/readline"
// 	"golang.org/x/term"
// )

// func main() {
// 	oldState := EnableRaw()
// 	defer term.Restore(int(os.Stdin.Fd()), oldState)

// 	rl, err := readline.New("> ")
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer rl.Close()

// 	for {
// 		line, err := rl.Readline()
// 		if err != nil { // io.EOF
// 			break
// 		}
// 		fmt.Printf("%s\n\r", line)
// 	}

// 	// cl := NewComLineTxt()
// 	// ks := make([]byte, 5) // key stroke store

// 	// fmt.Printf("$ ")

// 	// for {
// 	// 	n, err := os.Stdin.Read(ks)
// 	// 	if err != nil {
// 	// 		fmt.Printf("error: %s\n\r", err)
// 	// 		break
// 	// 	}

// 	// 	// fmt.Printf("%v %v\n\r", n, ks)
// 	// 	com := cl.text.String()

// 	// 	if n == 1 && ks[0] == 3 { // ctrl+c
// 	// 		break
// 	// 	} else if n == 3 && ks[0] == 27 && ks[1] == 91 && ks[2] == 65 {
// 	// 		// break
// 	// 	} else if n == 3 && ks[0] == 27 && ks[1] == 91 && ks[2] == 66 {
// 	// 		// break
// 	// 	} else if n == 3 && ks[0] == 27 && ks[1] == 91 && ks[2] == 67 {
// 	// 		// break
// 	// 	} else if n == 3 && ks[0] == 27 && ks[1] == 91 && ks[2] == 68 {
// 	// 		// break
// 	// 	} else if n == 1 && ks[0] == 9 { // ctrl+c
// 	// 		sug := Tab(com)
// 	// 		if len(sug) == 0 {
// 	// 			continue
// 	// 		}
// 	// 		cl.Clear()
// 	// 		nb := []byte(sug)
// 	// 		cl.Write(nb, len(nb))
// 	// 	} else if n == 1 && ks[0] == 13 { // return (or enter)
// 	// 		Return(com)
// 	// 		cl.Clear()
// 	// 		fmt.Printf("$ ")
// 	// 	} else if n == 1 && ks[0] == 127 { // delete
// 	// 		cl.Delete()
// 	// 		continue
// 	// 	} else {
// 	// 		cl.Write(ks, n)
// 	// 	}
// 	// }
// }

// // tab => 9
// // ctrl+c => 3
// // \n => 13
// // delete => 127

package main

import (
	"io"
	"strings"

	"github.com/chzyer/readline"
)

func GetComm(com string) (string, []string) {
	commParts := SplitComm(com)
	if len(commParts) == 0 {
		return "", commParts
	}
	return commParts[0], commParts[1:]
}

func main() {
	l, err := readline.NewEx(GetConfig())
	if err != nil {
		panic(err)
	}
	defer l.Close()

	for {
		com, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(com) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		com = strings.TrimSpace(com)
		main, args := GetComm(com)
		var (
			outFilePath string
			redirect    int = 0
			mode        int = 0
		)
		// filter out args without the redirect
		args = RedirectFilter(args, &mode, &redirect, &outFilePath)
		switch main {
		case "":
		case "exit":
			HandleExit()
		case "echo":
			HandleEcho(args, outFilePath, redirect, mode)
		case "type":
			HandleType(args)
		case "pwd":
			HandlePwd()
		case "cd":
			HandleCd(args)
		default:
			HandleDefault(main, args, outFilePath, redirect, mode)
			// case strings.HasPrefix(line, "mode "):
			// 	switch line[5:] {
			// 	case "vi":
			// 		l.SetVimMode(true)
			// 	case "emacs":
			// 		l.SetVimMode(false)
			// 	default:
			// 		println("invalid mode:", line[5:])
			// 	}
			// case line == "mode":
			// 	if l.IsVimMode() {
			// 		println("current mode: vim")
			// 	} else {
			// 		println("current mode: emacs")
			// 	}
			// case line == "login":
			// 	pswd, err := l.ReadPassword("please enter your password: ")
			// 	if err != nil {
			// 		break
			// 	}
			// 	println("you enter:", strconv.Quote(string(pswd)))
			// case line == "help":
			// 	usage(l.Stderr())
			// case strings.HasPrefix(line, "setprompt"):
			// 	if len(line) <= 10 {
			// 		log.Println("setprompt <prompt>")
			// 		break
			// 	}
			// 	l.SetPrompt(line[10:])
			// case strings.HasPrefix(line, "say"):
			// 	line := strings.TrimSpace(line[3:])
			// 	if len(line) == 0 {
			// 		log.Println("say what?")
			// 		break
			// 	}
			// 	go func() {
			// 		for range time.Tick(time.Second) {
			// 			log.Println(line)
			// 		}
			// 	}()
			// case line == "bye":
			// 	goto exit
			// case line == "sleep":
			// 	log.Println("sleep 4 second")
			// 	time.Sleep(4 * time.Second)
			// default:
			// 	log.Println("you said:", strconv.Quote(line))
		}
	}
	// exit:
}
