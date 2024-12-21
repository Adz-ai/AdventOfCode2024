package main

import (
	. "aoc2024/utility"
	"log"
	"strconv"
	"strings"
)

type Coordinates struct {
	X, Y int
}

// getNumericValue extracts numeric value from code ignoring leading zeros
func getNumericValue(code string) int {
	numStr := strings.TrimLeft(code[:len(code)-1], "0")
	if numStr == "" {
		return 0
	}
	num, _ := strconv.Atoi(numStr)
	return num
}

func getPressesForNumericPad(input []string, start string, numericalMap map[string]Coordinates) []string {
	current := numericalMap[start]
	output := make([]string, 0)

	for _, char := range input {
		dest := numericalMap[char]
		diffX, diffY := dest.X-current.X, dest.Y-current.Y

		horizontal := make([]string, 0)
		vertical := make([]string, 0)

		for i := 0; i < Abs(diffX); i++ {
			if diffX >= 0 {
				horizontal = append(horizontal, ">")
			} else {
				horizontal = append(horizontal, "<")
			}
		}

		for i := 0; i < Abs(diffY); i++ {
			if diffY >= 0 {
				vertical = append(vertical, "^")
			} else {
				vertical = append(vertical, "v")
			}
		}

		// Prioritize movement based on current position
		switch {
		case current.Y == 0 && dest.X == 0:
			output = append(output, vertical...)
			output = append(output, horizontal...)
		case current.X == 0 && dest.Y == 0:
			output = append(output, horizontal...)
			output = append(output, vertical...)
		case diffX < 0:
			output = append(output, horizontal...)
			output = append(output, vertical...)
		case diffX >= 0:
			output = append(output, vertical...)
			output = append(output, horizontal...)
		}

		current = dest
		output = append(output, "A")
	}
	return output
}

func getPressesForDirectionalPad(input []string, start string, directionalMap map[string]Coordinates) []string {
	current := directionalMap[start]
	output := make([]string, 0)

	for _, char := range input {
		dest := directionalMap[char]
		diffX, diffY := dest.X-current.X, dest.Y-current.Y

		horizontal := make([]string, 0)
		vertical := make([]string, 0)

		for i := 0; i < Abs(diffX); i++ {
			if diffX >= 0 {
				horizontal = append(horizontal, ">")
			} else {
				horizontal = append(horizontal, "<")
			}
		}

		for i := 0; i < Abs(diffY); i++ {
			if diffY >= 0 {
				vertical = append(vertical, "^")
			} else {
				vertical = append(vertical, "v")
			}
		}

		// Prioritize movement based on current position
		if current.X == 0 && dest.Y == 1 {
			output = append(output, horizontal...)
			output = append(output, vertical...)
		} else if current.Y == 1 && dest.X == 0 {
			output = append(output, vertical...)
			output = append(output, horizontal...)
		} else if diffX < 0 {
			output = append(output, horizontal...)
			output = append(output, vertical...)
		} else if diffX >= 0 {
			output = append(output, vertical...)
			output = append(output, horizontal...)
		}
		current = dest
		output = append(output, "A")
	}
	return output
}

func getIndividualSteps(input []string) [][]string {
	output := make([][]string, 0)
	current := make([]string, 0)
	for _, char := range input {
		current = append(current, char)
		if char == "A" {
			output = append(output, current)
			current = []string{}
		}
	}
	return output
}

func getCountAfterRobots(input []string, maxRobots int, robot int, cache map[string][]int, directionalMap map[string]Coordinates) int { //nolint:lll
	key := strings.Join(input, "")
	if val, ok := cache[key]; ok {
		if val[robot-1] != 0 {
			return val[robot-1]
		}
	} else {
		cache[key] = make([]int, maxRobots)
	}

	seq := getPressesForDirectionalPad(input, "A", directionalMap)
	cache[key][0] = len(seq)

	if robot == maxRobots {
		return len(seq)
	}

	splitSeq := getIndividualSteps(seq)
	count := 0
	for _, s := range splitSeq {
		c := getCountAfterRobots(s, maxRobots, robot+1, cache, directionalMap)
		if _, ok := cache[strings.Join(s, "")]; !ok {
			cache[strings.Join(s, "")] = make([]int, maxRobots)
		}
		cache[strings.Join(s, "")][0] = c
		count += c
	}

	cache[key][robot-1] = count
	return count
}

func getSequence(input []string, numericalMap, directionalMap map[string]Coordinates, robots int) int {
	count := 0
	cache := make(map[string][]int)
	for _, line := range input {
		row := strings.Split(line, "")
		seq1 := getPressesForNumericPad(row, "A", numericalMap)
		num := getCountAfterRobots(seq1, robots, 1, cache, directionalMap)
		count += getNumericValue(line) * num
	}
	return count
}

func main() {
	// Initialize maps
	numericalMap := map[string]Coordinates{
		"A": {2, 0},
		"0": {1, 0},
		"1": {0, 1},
		"2": {1, 1},
		"3": {2, 1},
		"4": {0, 2},
		"5": {1, 2},
		"6": {2, 2},
		"7": {0, 3},
		"8": {1, 3},
		"9": {2, 3},
	}

	directionalMap := map[string]Coordinates{
		"A": {2, 1},
		"^": {1, 1},
		"<": {0, 0},
		"v": {1, 0},
		">": {2, 0},
	}

	input, err := ParseTextFile("input")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Part 1:", getSequence(input, numericalMap, directionalMap, 2))
	log.Println("Part 2:", getSequence(input, numericalMap, directionalMap, 25))
}
