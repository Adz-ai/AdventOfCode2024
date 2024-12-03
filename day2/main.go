package main

import (
	"aoc2024/utility"
	"fmt"
	"log"
	"strings"
)

func isSafe(report []int) bool {
	if len(report) < 2 {
		return true
	}

	increasing := report[1] > report[0]

	for i := 1; i < len(report); i++ {
		diff := report[i] - report[i-1]

		if diff == 0 || diff < -3 || diff > 3 {
			return false
		}

		if (increasing && diff < 0) || (!increasing && diff > 0) {
			return false
		}
	}

	return true
}

func part2(input []string) int {
	safeCounter := 0

	for _, line := range input {
		split := strings.Split(line, " ")
		splitInt := utility.SliceOfStringsToInt(split)

		if isSafe(splitInt) {
			safeCounter++
			continue
		}

		for i := 0; i < len(splitInt); i++ {
			modified := append([]int{}, splitInt[:i]...)
			modified = append(modified, splitInt[i+1:]...)

			if isSafe(modified) {
				safeCounter++
				break
			}
		}
	}

	return safeCounter
}

func part1(input []string) int {
	safeCounter := 0

	for _, line := range input {
		split := strings.Split(line, " ")
		splitInt := utility.SliceOfStringsToInt(split)

		if isSafe(splitInt) {
			safeCounter++
		}
	}

	return safeCounter
}

func main() {
	input, err := utility.ParseTextFile("input")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(part1(input))
	fmt.Println(part2(input))
}
