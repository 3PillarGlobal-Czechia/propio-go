package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	PATH string = "m.properties"
)

var (
	comments   []string = []string{"#", "!"}
	delimiters []string = []string{":", "="}
)

type property struct {
	Line      int
	Delimiter string
	Escaped   bool
	Key       string
	Value     string
}

func main() {
	json, err := json.Marshal(getProperties())
	if err == nil {
		result := strings.ReplaceAll(string(json), "\\\\", "\\")
		fmt.Printf(result)
	} else {
		fmt.Printf(err.Error())
	}
}

func getFile() (file *os.File) {
	file, err := os.Open(PATH)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(3)
	}

	return file
}

func getLinesFromFile() (lines []string) {
	file := getFile()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	file.Close()

	return lines
}

func getProperties() (result []property) {
	lines := getLinesFromFile()
	for index, line := range lines {
		if line == "" {
			continue
		}

		var p property
		p.Line = index
		p.Delimiter = getDelimiter(line)
		p.Escaped = isEscaped(line, p.Delimiter)
		p.Key, p.Value = getKeyAndValue(line, p.Delimiter, p.Escaped)

		result = append(result, p)
	}

	return mergeMultilinesValue(result)
}

func mergeMultilinesValue(properties []property) (result []property) {
	skip := make(map[int]bool)
	merged := result
	for index, _ := range properties {
		i := index
		merged = append(merged, properties[i])
		for {
			i++
			if isMultiline(properties, i) {
				skip[i] = true
				merged[len(merged)-1].Value += properties[i].Value
				continue
			}

			break
		}
	}

	for index, _ := range merged {
		if skip[index] {
			fmt.Println(strconv.Itoa(index))
			continue
		}

		result = append(result, merged[index])
	}

	return result
}

func isMultiline(properties []property, index int) bool {
	if index-1 < 0 {
		return false
	}

	if getLastChar(properties[index-1].Value) == "\\" {
		return true
	}

	return false
}

func getLastChar(value string) string {
	return string(value[len(value)-1:])
}

func getDelimiter(line string) (delimiter string) {
	if isComment(line) {
		return string(line[0])
	}

	delimiterPosition := -1
	for _, value := range delimiters {
		position := strings.Index(line, string(value))
		if position == -1 {
			continue
		}
		if delimiterPosition == -1 || delimiterPosition > position {
			delimiterPosition = position
			delimiter = value
		}
	}

	if delimiterPosition == -1 {
		position := strings.Index(line, " ")
		if position > -1 {
			return " "
		}
	}

	return delimiter
}

func getKeyAndValue(line, delimiter string, escaped bool) (key, value string) {
	position := strings.Index(line, delimiter)
	if position == 0 {
		return "", strings.TrimSpace(line[1:])
	}

	if escaped {
		key = line[:position-1]
	} else {
		key = line[:position]
	}

	value = line[position+1:]

	return strings.TrimSuffix(key, " "), strings.TrimPrefix(value, " ")
}

func isEscaped(line, delimiter string) bool {
	if delimiter == "" {
		return false
	}

	position := strings.Index(line, delimiter)
	if position == 0 {
		return false
	}
	if string(line[position-1]) == "\\" {
		return true
	}

	return false
}

func isComment(line string) bool {
	for _, value := range comments {
		if string(line[0]) == value {
			return true
		}
	}

	return false
}
