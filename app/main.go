// "hello"  "script's"  test""example -> test on this
package main

import (
	"bytes"
	"fmt"
	"os"

	"golang.org/x/term"
)

func GetComm(com string) (string, []string) {
	commParts := SplitComm(com)
	if len(commParts) == 0 {
		return "", commParts
	}
	return commParts[0], commParts[1:]
}

func main() {
	oldState := EnableRaw()
	// Ensure we restore the terminal state when finished
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	in := make([]byte, 0)
	txt := bytes.NewBuffer(in)

	fmt.Printf("$ ")

	b := make([]byte, 1)
	for {
		if _, err := os.Stdin.Read(b); err != nil {
			fmt.Printf("error: %s\n\r", err)
			break
		}

		if b[0] == 3 {
			break
		}
		if b[0] == 13 {
			Return(txt)
			continue
		}

		if err := txt.WriteByte(b[0]); err != nil {
			fmt.Printf("error: %s\n\r", err)
		}
		fmt.Printf("%s", string(b[0]))
	}
}

// tab => 9
// ctrl+c => 3
// \n => 13
