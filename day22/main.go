package main

import (
	"aoc2024/utility"
	"fmt"
	"log"
	"strconv"
)

const (
	MOD = 16777216
)

// mix performs a bitwise XOR operation between the secret and value
func mix(secret, value int) int {
	return secret ^ value
}

// prune performs modulo operation with MOD constant
func prune(secret int) int {
	return secret % MOD
}

// generateNext creates the next secret number based on the current one
func generateNext(secret int) int {
	// Step 1: Multiply by 64
	result := mix(secret, secret*64)
	result = prune(result)

	// Step 2: Divide by 32
	result = mix(result, result/32)
	result = prune(result)

	// Step 3: Multiply by 2048
	result = mix(result, result*2048)
	result = prune(result)

	return result
}

// generate2000th generates the 2000th number in the sequence
func generate2000th(initial int) int {
	current := initial
	for i := 0; i < 2000; i++ {
		current = generateNext(current)
	}
	return current
}

func part1(input []string) int {
	var sum int

	for _, line := range input {
		initial, err := strconv.Atoi(line)
		if err != nil {
			fmt.Printf("Error parsing input: %v\n", err)
			continue
		}

		// Generate 2000th number and add to sum
		final := generate2000th(initial)
		sum += final
	}

	return sum
}

func main() {
	input, err := utility.ParseTextFile("input")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(part1(input))
}
