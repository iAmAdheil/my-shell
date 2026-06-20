package main

import (
	"fmt"
	"testing"
)

func TestExeSearch(t *testing.T) {
	testStr := ".py"
	matches := SearchPath(testStr)
	fmt.Println(matches)
}
