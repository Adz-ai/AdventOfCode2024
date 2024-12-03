package main

import (
	"aoc2024/utility"
	"fmt"
	"log"
	"regexp"
	"strconv"
)

func processMultiplication(match []string) (int, error) {
	first, second := 1, 2
	if len(match) > 3 {
		first, second = 2, 3
	}

	x, err := strconv.Atoi(match[first])
	if err != nil {
		return 0, fmt.Errorf("failed to parse first number: %w", err)
	}

	y, err := strconv.Atoi(match[second])
	if err != nil {
		return 0, fmt.Errorf("failed to parse second number: %w", err)
	}

	return x * y, nil
}

func part2(input []string) int {
	pattern := `(mul\((\d+),(\d+)\)|do\(\)|don't\(\))`
	re := regexp.MustCompile(pattern)
	enabled := true
	total := 0
	for _, line := range input {
		matches := re.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			switch match[0] {
			case "do()":
				enabled = true
			case "don't()":
				enabled = false
			default:
				if match[1] != "" && enabled {
					if result, err := processMultiplication(match); err != nil {
						log.Fatal(err)
					} else {
						total += result
					}
				}
			}
		}
	}
	return total
}

func part1(input []string) int {
	pattern := `mul\((\d+),(\d+)\)`
	re := regexp.MustCompile(pattern)
	total := 0
	for _, line := range input {
		matches := re.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if result, err := processMultiplication(match); err != nil {
				log.Fatal(err)
			} else {
				total += result
			}
		}
	}
	return total
}

func main() {
	input, err := utility.ParseTextFile("input")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(part1(input))
	fmt.Println(part2(input))
}
