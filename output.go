package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Line struct {
	Key   string
	Value string
}

func inputToLines(data string) []Line {
	var lines []Line
	err := json.Unmarshal([]byte(data), &lines)
	if err != nil {
		fmt.Println(err.Error())
	}

	return lines
}

func saveToFile(lines []Line, path string) {
	file, err := os.Create(path)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer file.Close()

	for _, line := range lines {
		_, err := file.WriteString(line.Key + "=" + line.Value + "\n")
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	file.Sync()
}

func output(data, path string) {
	saveToFile(inputToLines(data), path)
}
