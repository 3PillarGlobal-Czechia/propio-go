package main

import (
	"os"
)

var (
	comments   []string = []string{"#", "!"}
	delimiters []string = []string{":", "="}
)

type property struct {
	Key   string
	Value string
}

func main() {
	args := os.Args
	if len(args) == 2 {
		input(args[1])
	} else if len(args) == 3 {
		output(args[1], args[2])
	}
}

func isComment(line string) bool {
	if line == "" {
		return false
	}

	for _, value := range comments {
		if string(line[0]) == value {
			return true
		}
	}

	return false
}
