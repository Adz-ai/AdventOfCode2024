package main

import (
	"aoc2024/utility"
	"fmt"
	"log"
	"regexp"
	"strconv"
)

func part2(input []string) int {
	pattern := `(mul\((\d+),(\d+)\)|do\(\)|don't\(\))`
	re := regexp.MustCompile(pattern)
	enabled := true
	total := 0
	for _, line := range input {
		matches := re.FindAllStringSubmatch(line, -1)
		for _, match := range matches {
			if match[0] == "do()" {
				enabled = true
			} else if match[0] == "don't()" {
				enabled = false
			} else if match[1] != "" && enabled {
				x, err := strconv.Atoi(match[2])
				if err != nil {
					log.Fatal(err)
				}
				y, err := strconv.Atoi(match[3])
				if err != nil {
					log.Fatal(err)
				}
				total = total + (x * y)
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
			x, err := strconv.Atoi(match[1])
			if err != nil {
				log.Fatal(err)
			}
			y, err := strconv.Atoi(match[2])
			if err != nil {
				log.Fatal(err)
			}
			total = total + (x * y)
		}
	}
	return total
}

func main() {
	input, err := utility.ParseTextFile("day3", "input")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(part1(input))
	fmt.Println(part2(input))
}
