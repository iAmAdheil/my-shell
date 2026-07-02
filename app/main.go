package main

import (
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/chzyer/readline"
	"github.com/codecrafters-io/shell-starter-go/app/com"
)

func GetComms(txt string) []string {
	txt = strings.TrimSpace(txt)
	return strings.Split(txt, " | ")
}

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
		txt, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(txt) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		comms := GetComms(txt)

		var in io.Reader = nil
		var running []*com.Com

		for i, _ := range comms {
			ct := comms[i]
			main, args := GetComm(ct) // normalise args, and extract main query

			pr, pw, err := os.Pipe()
			if err != nil {
				panic(err)
			}

			var (
				out   io.WriteCloser = pw
				close bool           = true
			)
			// last command prints to terminal
			if i == len(comms)-1 {
				out = os.Stdout
				close = false
			}

			var (
				outFilePath string // output file path
				// print or redirect to file (stdout or stderr)
				// 0 -> normal print
				// 1 -> stdout to file
				// 2 -> stderr to file
				redirect int = 0
				// 0 -> overwrite
				// 1 -> append
				mode int = 0 // append or overwrite
			)

			// filter out args without the redirect args
			args = RedirectFilter(args, &outFilePath, &redirect, &mode)

			com := &com.Com{
				Main:        main,
				Args:        args,
				Proc:        exec.Command(main, args...),
				In:          in,
				Out:         out,
				OutFilePath: outFilePath,
				Redirect:    redirect,
				Mode:        mode,
				Close:       close,
			}

			com.Run()
			// pass current com's pr to next com,
			// to read whatever is added via pw
			in = pr
			running = append(running, com)
		}

		for _, com := range running {
			com.Stop()
		}

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
	// exit:
}
