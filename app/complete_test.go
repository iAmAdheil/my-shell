package main

import (
	"fmt"
	"testing"
)

func TestComplete(t *testing.T) {
	words := []string{"echo", "exit"}
	txt := "ech"
	root := InitTrie(words)
	suggs := root.Complete(txt)
	fmt.Println(suggs)
}
