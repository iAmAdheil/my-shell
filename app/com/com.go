package com

import "io"

type Com struct {
	Main        string
	Args        []string
	Stdin       io.WriteCloser
	Stdout      io.ReadCloser
	OutFilePath string // output file path
	// print or redirect to file (stdout or stderr)
	// 0 -> normal print
	// 1 -> stdout to file
	// 2 -> stderr to file
	Redirect int
	// 0 -> overwrite
	// 1 -> append
	Mode int // append or overwrite
}
