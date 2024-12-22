package main

import (
	"aoc2024/utility"
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
	result := mix(secret, secret*64)
	result = prune(result)

	result = mix(result, result/32)
	result = prune(result)

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
			log.Fatalf("Error parsing input: %v\n", err)
		}

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
