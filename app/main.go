package main

import (
	"fmt"
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

func GetComm(c string) (string, []string) {
	commParts := SplitComm(c)
	if len(commParts) == 0 {
		return "", commParts
	}

	commParts = com.HandleExpandVar(commParts)

	return commParts[0], commParts[1:]
}

func main() {
	// runs on startup
	// Init()

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

		com.History = append(com.History, txt)

		comms := GetComms(txt)

		// "&" applies to the entire command, not to one pipeline segment
		isBg := CheckBgComm(comms)

		// a command inherits the terminal and can change its settings.
		// a command that stops early leaves the settings changed, and
		// readline returns to the state that readline found. read the
		// state here, and put it back when the line finishes.
		ttyState := GetTtyState()

		var in *os.File = nil // read end of the previous command's pipe
		var running []*com.Com

		for i, _ := range comms {
			ct := comms[i]
			main, args := GetComm(ct) // normalise args, and extract main query

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
			args = HandleArgs(args)

			var (
				out   io.WriteCloser = os.Stdout
				close bool           = false
				pr    *os.File       // read end for the next command
			)
			// last command prints to terminal, and needs no pipe
			if i != len(comms)-1 {
				r, pw, err := os.Pipe()
				if err != nil {
					panic(err)
				}
				pr, out, close = r, pw, true
			}

			// the first command of a pipeline has no pipe to read from.
			// give it the terminal, so a program can ask for input
			// during execution (for example "y/n" or a password).
			// a background command must not read from the terminal,
			// because the prompt reads from the terminal at the same time.
			if i == 0 && !isBg {
				in = os.Stdin
			}

			// an *os.File that holds nil is not equal to nil after
			// assignment to an io.Reader, so convert it explicitly
			var cin io.Reader
			if in != nil {
				cin = in
			}

			c := &com.Com{
				Main:        main,
				Args:        args,
				Proc:        exec.Command(main, args...),
				In:          cin,
				Out:         out,
				OutFilePath: outFilePath,
				Redirect:    redirect,
				Mode:        mode,
				Close:       close,
				// the job records the last process of the pipeline
				IsBgProc: isBg && i == len(comms)-1,
			}

			c.Run()
			// the command holds its own copy of the read end, so the
			// shell must close the copy that the shell holds.
			// os.Stdin is not a pipe end, and the shell keeps it open
			if in != nil && in != os.Stdin {
				in.Close()
			}
			// pass current com's pr to next com,
			// to read whatever is added via pw
			in = pr
			running = append(running, c)

			if c.IsBgProc {
				com.Count++
				njob := &com.Job{
					Id:      com.Count,
					PId:     c.Proc.Process.Pid,
					Status:  "Running",
					ComText: c.Main + " " + strings.Join(c.Args, " "),
				}
				fmt.Printf("[%v] %v\n", njob.Id, njob.PId)
				com.Jobs = append(com.Jobs, njob)
			}
		}

		for _, c := range running {
			if c.IsBgProc {
				go func() {
					pid := c.Proc.Process.Pid
					c.Stop()
					com.UpdateJobStatus(pid)
				}()
			} else {
				c.Stop()
			}
		}

		RestoreTtyState(ttyState)

		com.HandleCompleteJobs()

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
