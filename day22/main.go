package main

import (
	"aoc2024/utility"
	"log"
	"strconv"
)

type bananaPrice struct {
	num    int
	change int
}

type seq struct {
	a, b, c, d int
}

// Get price and change from a secret number
func getNumAndChange(input, previous int) bananaPrice {
	currentNum := input % 10
	return bananaPrice{
		num:    currentNum,
		change: currentNum - previous,
	}
}

func mix(value, secretNum int) int {
	return value ^ secretNum
}

func prune(value int) int {
	return value % 16777216
}

func findSecretNumber(input int) int {
	input = prune(mix(64*input, input))
	input = prune(mix(input/32, input))
	input = prune(mix(input*2048, input))
	return input
}

func findMaxNumberOfBananas(b [][]bananaPrice) int {
	seqMap := make(map[seq][]int)

	// For each buyer
	for i := range b {
		// Initial sliding window of changes
		s := []int{b[i][0].change, b[i][1].change, b[i][2].change}

		// Slide window through all changes
		for j := 3; j < len(b[i]); j++ {
			s = append(s, b[i][j].change)
			sequence := seq{s[0], s[1], s[2], s[3]}

			// Initialize sequence in map if not exists
			if _, ok := seqMap[sequence]; !ok {
				seqMap[sequence] = make([]int, len(b))
			}

			// Store first occurrence of price for this buyer
			if seqMap[sequence][i] == 0 {
				seqMap[sequence][i] = b[i][j].num
			}

			// Move window forward
			s = s[1:]
		}
	}

	// Find sequence that gives maximum bananas
	maxBananas := 0
	for _, prices := range seqMap {
		sum := 0
		for _, price := range prices {
			sum += price
		}
		if sum > maxBananas {
			maxBananas = sum
		}
	}

	return maxBananas
}

func solve(input []string) (int, int) {
	var sum int
	priceData := make([][]bananaPrice, len(input))

	// Process each buyer
	for i, line := range input {
		num, _ := strconv.Atoi(line)
		prev := num % 10 // Initial price

		// Generate 2000 numbers and their changes
		priceData[i] = make([]bananaPrice, 2000)
		for j := 0; j < 2000; j++ {
			num = findSecretNumber(num)
			priceData[i][j] = getNumAndChange(num, prev)
			prev = num % 10

			if j == 1999 { // Last number for part 1
				sum += num
			}
		}
	}

	return sum, findMaxNumberOfBananas(priceData)
}

func main() {
	input, err := utility.ParseTextFile("input")
	if err != nil {
		log.Fatal(err)
	}
	part1, part2 := solve(input)
	log.Printf("Part 1: %d\n", part1)
	log.Printf("Part 2: %d\n", part2)
}
