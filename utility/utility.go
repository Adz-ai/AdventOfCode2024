package utility

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func ParseTextFile(day, filename string) ([]string, error) {
	pathToWorkingDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	inputPath := filepath.Join(pathToWorkingDir, day, fmt.Sprintf("%s.txt", filename))
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func SliceOfStringsToInt(input []string) []int {
	var result []int
	for _, line := range input {
		i, err := strconv.Atoi(line)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, i)
	}
	return result
}
