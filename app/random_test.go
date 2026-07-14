package main

import (
	"fmt"
	"strings"
	"testing"
)

// use this to test random things

func TestRandom(t *testing.T) {
	str := "var="
	sec := strings.Split(str, "=")
	fmt.Println(sec)
}
