package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var (
	comments   []string = []string{"#", "!"}
	delimiters []string = []string{":", "="}
	input      string
	output     string
)

type property struct {
	Key   string
	Value string
}

func main() {
	args := os.Args
	if len(args) < 2 {
		return
	}

	input := args[1]

	json, err := json.Marshal(getProperties(input))
	if err == nil {
		result := strings.ReplaceAll(string(json), "\\\\", "\\")
		fmt.Printf(result)
	} else {
		fmt.Printf(err.Error())
	}
}
