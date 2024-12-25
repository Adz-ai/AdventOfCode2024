package main

import (
	"aoc2024/utility"
	"log"
)

func parseInput(input []string) ([][]int, [][]int) {
	var locks, keys [][]int
	var curr []string

	// Process each schema
	for i := 0; i < len(input); i++ {
		if input[i] == "" {
			if len(curr) > 0 {
				heights := processSchema(curr)
				if isLock(curr) {
					locks = append(locks, heights)
				} else {
					keys = append(keys, heights)
				}
				curr = nil
			}
			continue
		}
		curr = append(curr, input[i])
	}

	// Process last schema if exists
	if len(curr) > 0 {
		heights := processSchema(curr)
		if isLock(curr) {
			locks = append(locks, heights)
		} else {
			keys = append(keys, heights)
		}
	}

	return locks, keys
}

func isLock(schema []string) bool {
	count := 0
	for _, char := range schema[0] {
		if char == '#' {
			count++
		}
	}
	return count == 5
}

func processSchema(lines []string) []int {
	heights := make([]int, 5)

	// Count '#' characters in each column
	for _, line := range lines {
		for col, char := range line {
			if char == '#' {
				heights[col]++
			}
		}
	}

	// Subtract 1 from each height to account for the top/bottom row
	for i := range heights {
		heights[i]--
	}

	return heights
}

func countValidPairs(locks, keys [][]int) int {
	count := 0
	for _, lock := range locks {
		for _, key := range keys {
			check := make([]int, 5)
			for i := range lock {
				check[i] = lock[i] + key[i]
			}

			isValid := true
			for i := range check {
				if check[i] > 5 {
					isValid = false
					break
				}
			}

			if isValid {
				count++
			}
		}
	}
	return count
}

func part1(input []string) int {
	locks, keys := parseInput(input)
	return countValidPairs(locks, keys)
}

func main() {
	input, err := utility.ParseTextFile("input")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(part1(input))
	//log.Println(part2(input))
}
