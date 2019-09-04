package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func getFile(path string) (file *os.File) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(3)
	}

	return file
}

func getLinesFromFile(path string) (lines []string) {
	file := getFile(path)
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	file.Close()

	return lines
}

func getProperties(path string) (result []property) {
	lines := getLinesFromFile(path)
	for _, line := range lines {
		if line == "" {
			var p property
			p.Key, p.Value = "", ""

			result = append(result, p)

			continue
		}

		var p property
		p.Key, p.Value = getKeyAndValue(line)

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
	if len(value) == 0 {
		return ""
	}

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

func getKeyAndValue(line string) (key, value string) {
	delimiter := getDelimiter(line)
	escaped := isEscaped(line, delimiter)
	position := strings.Index(line, delimiter)
	if position == 0 {
		return delimiter, strings.TrimSpace(line[1:])
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
