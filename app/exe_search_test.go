package main

import (
	"fmt"
	"testing"
)

func TestExeSearch(t *testing.T) {
	testStr := ".py"
	matches := searchPath(testStr)
	fmt.Println(matches)
}
