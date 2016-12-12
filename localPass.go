package main

import (
	"bufio"
	"fmt"
	"os"
)

func checkLocalPass(filename string, pass string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening passwd file %q: %q", filename, err.Error())
	}
	defer file.Close()

	line, isPrefix, err := bufio.NewReader(file).ReadLine()
	if err != nil {
		return fmt.Errorf("error reading passwd file %q: %q", filename, err.Error())
	}
	if isPrefix {
		return fmt.Errorf("error reading passwd file %q: line too long", filename)
	}

	if string(line) != pass {
		return fmt.Errorf("password mismatch")
	}

	return nil
}
