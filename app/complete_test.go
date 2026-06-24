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

func TestCommonPrefix(t *testing.T) {
	words := []string{"xyz_foo", "xyz_foo_bar", "xyz_foo_bar_baz"}
	pref := GetComPrefix(words)
	fmt.Println(pref)
}
