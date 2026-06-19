package main

import (
	"bytes"
	"fmt"
	"unicode/utf8"
)

// Handles all prints within the "$ " lines
type ComLine struct {
	cursor int
	text   *bytes.Buffer
}

func NewComLineTxt() *ComLine {
	in := make([]byte, 0)
	buf := bytes.NewBuffer(in)
	cl := &ComLine{
		cursor: 0,
		text:   buf,
	}
	return cl
}

func (cl *ComLine) Write(b []byte, n int) {
	n, err := cl.text.Write(b[:n])
	if err != nil {
		panic(err)
	}
	cl.cursor += n
	fmt.Printf("%s", b[:n])
}

// currently only deletes the last char
// (@iAmAdheil) pls fix it to delete char at position cursor
// once the cursor is free to move
func (cl *ComLine) Delete() {
	if cl.text.Len() > 0 {
		b := cl.text.Bytes()
		_, size := utf8.DecodeLastRune(b) // size = bytes in the final rune
		cl.text.Truncate(cl.text.Len() - size)
		fmt.Print("\b \b")
	}
}

// clear everything from start of "$ " to end
func (cl *ComLine) Clear() {
	cmd := cl.text.String()
	cl.text.Reset()
	n := utf8.RuneCountInString(cmd)
	if n > 0 {
		fmt.Printf("\033[%dD", n) // move cursor left n columns
	}
	fmt.Print("\033[K")
}
