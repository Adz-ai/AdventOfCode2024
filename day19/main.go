package main

import (
	"aoc2024/utility"
	"fmt"
	"log"
	"regexp"
	"strings"
)

func solve(input []string) (int, int) {
	// First line contains the patterns
	patterns := strings.Split(strings.ReplaceAll(input[0], " ", ""), ",")

	// Part 1: Build regex pattern
	regexStr := fmt.Sprintf("^(%s)+$", strings.Join(patterns, "|"))
	regex := regexp.MustCompile(regexStr)

	// Count matches in designs
	part1Count := 0
	part2Sum := 0

	for _, design := range input[2:] {
		if regex.MatchString(design) {
			part1Count++
			combinations := countCombinations(design, patterns)
			part2Sum += combinations
		}
	}

	return part1Count, part2Sum
}

func countCombinations(design string, patterns []string) int {
	return findAllCombinations(design, patterns, make(map[string]int))
}

func findAllCombinations(remaining string, patterns []string, memo map[string]int) int {
	if remaining == "" {
		return 1
	}

	if count, exists := memo[remaining]; exists {
		return count
	}

	total := 0
	for _, pattern := range patterns {
		if strings.HasPrefix(remaining, pattern) {
			total += findAllCombinations(remaining[len(pattern):], patterns, memo)
		}
	}

	memo[remaining] = total
	return total
}

func main() {
	input, err := utility.ParseTextFile("input")
	if err != nil {
		log.Fatal(err)
	}
	part1, part2 := solve(input)
	log.Println(part1)
	log.Println(part2)
}
