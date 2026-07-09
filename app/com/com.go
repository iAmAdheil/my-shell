package com

import (
	"io"
	"os/exec"
)

type Com struct {
	Main        string
	Args        []string
	Proc        *exec.Cmd
	In          io.Reader
	Out         io.WriteCloser
	OutFilePath string // output file path
	// print or redirect to file (stdout or stderr)
	// 0 -> normal print
	// 1 -> stdout to file
	// 2 -> stderr to file
	Redirect int
	// 0 -> overwrite
	// 1 -> append
	Mode     int  // append or overwrite
	Close    bool // check to manually close the write end of pipe except when os.Stdout
	IsBgProc bool // flag to mark proc as a background running proc
}
