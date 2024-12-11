package main

import (
	"aoc2024/utility"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func readData(lines []string) []int {
	var stones []int
	for _, line := range lines {
		chars := strings.Fields(line)
		for _, char := range chars {
			stone, _ := strconv.Atoi(char)
			stones = append(stones, stone)
		}
	}
	return stones
}

func blink(stones []int) []int {
	var newStones []int
	for j := 0; j < len(stones); j++ {
		if stones[j] == 0 {
			newStones = append(newStones, 1)
		} else {
			charStone := strconv.Itoa(stones[j])
			if len(charStone)%2 != 0 {
				newStones = append(newStones, stones[j]*2024)
			} else {
				num1, _ := strconv.Atoi(charStone[:len(charStone)/2])
				num2, _ := strconv.Atoi(charStone[len(charStone)/2:])
				newStones = append(newStones, num1, num2)
			}
		}
	}
	return newStones
}

func blinkPartTwo(stones []int, iterations int) map[int]int {
	stoneCounts := make(map[int]int)
	for _, stone := range stones {
		stoneCounts[stone]++
	}

	for i := 0; i < iterations; i++ {
		nextStoneCounts := make(map[int]int)
		for stone, count := range stoneCounts {
			if stone == 0 {
				nextStoneCounts[1] += count
			} else {
				charStone := strconv.Itoa(stone)
				if len(charStone)%2 != 0 {
					nextStoneCounts[stone*2024] += count
				} else {
					num1, _ := strconv.Atoi(charStone[:len(charStone)/2])
					num2, _ := strconv.Atoi(charStone[len(charStone)/2:])
					nextStoneCounts[num1] += count
					nextStoneCounts[num2] += count
				}
			}
		}
		stoneCounts = nextStoneCounts
	}
	return stoneCounts
}

func part1(lines []string) int {
	stones := readData(lines)
	for i := 0; i < 25; i++ {
		stones = blink(stones)
	}
	return len(stones)
}

func part2(lines []string) int {
	stones := readData(lines)
	stoneCounts := blinkPartTwo(stones, 75)
	counts := 0
	for _, count := range stoneCounts {
		counts += count
	}
	return counts
}

func main() {
	lines, err := utility.ParseTextFile("input")
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		return
	}
	log.Println(part1(lines))
	log.Println(part2(lines))
}
