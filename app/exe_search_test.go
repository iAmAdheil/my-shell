package main

import (
	"fmt"
	"testing"
)

// search for all executable exes in path startign with test string

func TestExeSearch(t *testing.T) {
	testStr := ".py"
	matches := SearchPath(testStr)
	fmt.Println(matches)
}
