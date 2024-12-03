package main

import (
	"aoc2024/utility"
	"fmt"
	"log"
	"sort"
	"strings"
)

func splitPuzzleInputToTwoSortedLines(input []string) ([]int, []int) {
	line1 := make([]string, 0)
	line2 := make([]string, 0)

	for _, line := range input {
		split := strings.Split(line, "  ")
		line1 = append(line1, strings.TrimSpace(split[0]))
		line2 = append(line2, strings.TrimSpace(split[1]))
	}
	line1Int := utility.SliceOfStringsToInt(line1)
	line2Int := utility.SliceOfStringsToInt(line2)

	sort.Ints(line1Int)
	sort.Ints(line2Int)
	return line1Int, line2Int
}

func countNumberOfTimesValuesIsInSlice(slice []int) map[int]int {
	uniqueMap := make(map[int]int)
	for _, value := range slice {
		uniqueMap[value]++
	}
	return uniqueMap
}

func compareAndGenerateSimilarityScore(l1 []int, m2 map[int]int) int {
	var similarityScore int
	for i := 0; i < len(l1); i++ {
		if v2, ok := m2[l1[i]]; ok {
			similarityScore += l1[i] * v2
		}
	}
	return similarityScore
}

func part2(input []string) int {
	l1, l2 := splitPuzzleInputToTwoSortedLines(input)
	m2 := countNumberOfTimesValuesIsInSlice(l2)
	return compareAndGenerateSimilarityScore(l1, m2)

}

func part1(input []string) int {
	l1, l2 := splitPuzzleInputToTwoSortedLines(input)

	var totalDistance int
	for i := 0; i < len(l1); i++ {
		distance := l1[i] - l2[i]
		if distance < 0 {
			distance *= -1
		}
		totalDistance += distance
	}

	return totalDistance
}

func main() {
	input, err := utility.ParseTextFile("day1", "input")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(part1(input))
	fmt.Println(part2(input))
}
