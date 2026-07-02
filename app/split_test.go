// TESTING ENV

package main

import (
	"fmt"
	"testing"
)

// test command split, (@iAmAdheil) needs rework here

func TestSplit(t *testing.T) {
	// str := "\"/tmp/ant/f 60\" \"/tmp/ant/f   54\" \"/tmp/ant/f's98\""
	str := "cat /tmp/foo/file | wc"
	res := SplitComm(str)
	for i, v := range res {
		fmt.Println("index:", i, "value:", v)
	}
}

func TestCommSplit(t *testing.T) {
	str := "echo hello > output.txt"
	com, args := GetComm(str)
	fmt.Printf("command: %s\n", com)
	for i, v := range args {
		fmt.Printf("arg %v: %s\n", i+1, v)
	}
}

func TestCommNormalise(t *testing.T) {
	com := `echo -e "1. grape banana\n2. apple pear\n3. mango pineapple" > "/tmp/pig/file-89"`
	res := NormaliseComm(com)
	fmt.Println(res)
}
