package main

import (
	"aoc2024/utility"
	"log"
	"strconv"
	"strings"
)

// parseLine parses a single line into target value and numbers
func parseLine(line string) (int, []int) {
	parts := strings.Split(line, ": ")
	target, _ := strconv.Atoi(parts[0])

	numbersStr := strings.Fields(parts[1])
	numbers := make([]int, len(numbersStr))
	for i, numStr := range numbersStr {
		numbers[i], _ = strconv.Atoi(numStr)
	}

	return target, numbers
}

// isValid recursively checks if it's possible to reach the target
func isValid(target int, values []int, allowConcat bool) bool {
	if len(values) == 1 {
		return values[0] == target
	}

	if isValid(target, append([]int{values[0] + values[1]}, values[2:]...), allowConcat) {
		return true
	}

	if isValid(target, append([]int{values[0] * values[1]}, values[2:]...), allowConcat) {
		return true
	}

	if allowConcat {
		concat := strconv.Itoa(values[0]) + strconv.Itoa(values[1])
		concatVal, _ := strconv.Atoi(concat)
		if isValid(target, append([]int{concatVal}, values[2:]...), allowConcat) {
			return true
		}
	}

	return false
}

func part1(input []string) int {
	total := 0
	for _, line := range input {
		target, numbers := parseLine(line)
		if isValid(target, numbers, false) {
			total += target
		}
	}
	return total
}

func part2(input []string) int {
	total := 0
	for _, line := range input {
		target, numbers := parseLine(line)
		if isValid(target, numbers, true) {
			total += target
		}
	}
	return total
}

func main() {
	input, err := utility.ParseTextFile("input")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(part1(input))
	log.Println(part2(input))
}
