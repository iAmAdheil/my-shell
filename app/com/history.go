package com

import (
	"bufio"
	"fmt"
)

func HandleHistoryWrite(filename string) error {
	file, err := OpenFile(filename, 1)
	if err != nil {
		return fmt.Errorf("history logs could not be opened: %v", err)
	}
	defer file.Close()

	var txt string = ""

	for _, line := range History {
		txt += line + "\n"
	}

	_, err = file.WriteString(txt)
	if err != nil {
		return fmt.Errorf("history logs could not be written to file: %v", err)
	}

	return nil
}

func HandleHistoryRead(filename string) error {
	file, err := OpenFile(filename, 0)
	if err != nil {
		return fmt.Errorf("history logs could not be opened: %v", err)
	}
	defer file.Close()

	sc := bufio.NewScanner(file)

	for sc.Scan() {
		txt := sc.Text()
		History = append(History, txt)
	}

	return nil
}
