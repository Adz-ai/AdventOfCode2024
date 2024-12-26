package main

import (
	"aoc2024/utility"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// dependency represents a logical gate operation
type dependency struct {
	w1, w2 string
	op     string
}

func parseInput(input []string) (map[string]int8, map[string]dependency) {
	instrRegex := regexp.MustCompile(`([a-z0-9]*) ([A-Z]*) ([a-z0-9]*) -> ([a-z0-9]*)`)
	wireValueRegex := regexp.MustCompile(`([a-zA-Z0-9]*): ([0-9])`)

	value := make(map[string]int8)
	dependencies := make(map[string]dependency)
	parsingValues := true

	for _, line := range input {
		line = strings.TrimSpace(line)
		if line == "" {
			parsingValues = false
			continue
		}

		if parsingValues {
			matches := wireValueRegex.FindStringSubmatch(line)
			w := matches[1]
			v := int8(matches[2][0] - '0')
			value[w] = v
		} else {
			matches := instrRegex.FindStringSubmatch(line)
			w := matches[4]
			op := matches[2]
			w1, w2 := matches[1], matches[3]

			dependencies[w] = dependency{
				w1: w1,
				w2: w2,
				op: op,
			}
		}
	}

	return value, dependencies
}

func partOne(value map[string]int8, dependencies map[string]dependency) (res uint64) {
	var resolve func(string) int8

	resolve = func(curr string) int8 {
		if v, ok := value[curr]; ok {
			return v
		}

		d := dependencies[curr]
		v1 := resolve(d.w1)
		v2 := resolve(d.w2)

		switch d.op {
		case "XOR":
			value[curr] = v1 ^ v2
		case "AND":
			value[curr] = v1 & v2
		case "OR":
			value[curr] = v1 | v2
		}

		return value[curr]
	}

	// Resolve all dependencies
	for n := range dependencies {
		resolve(n)
	}

	// Build result from z-wires
	for n, v := range value {
		if n[0] == 'z' {
			temp, _ := strconv.Atoi(n[1:])
			res |= uint64(v) << temp
		}
	}

	return
}

func isXOrY(wire string) bool {
	temp, _ := strconv.Atoi(wire[1:])
	return (wire[0] == 'x' || wire[0] == 'y') && temp != 0
}

func partTwo(dependencies map[string]dependency) string {
	temp := make(map[string]bool)

	for w, d := range dependencies {
		// Rule 1: Check z-wires that aren't XOR gates (except z45)
		if w[0] == 'z' {
			val, _ := strconv.Atoi(w[1:])
			if d.op != "XOR" && val != 45 {
				temp[w] = true
			}
		} else if !isXOrY(d.w1) && !isXOrY(d.w2) && d.w1[0] != d.w2[0] && d.op == "XOR" {
			temp[w] = true
		}

		// Rule 2: Check XOR gates with x/y inputs from different series
		if d.op == "XOR" && isXOrY(d.w1) && isXOrY(d.w2) && d.w1[0] != d.w2[0] {
			isValid := false
			for _, dp := range dependencies {
				if dp.op == "XOR" && (dp.w1 == w || dp.w2 == w) {
					isValid = true
				}
			}
			if !isValid {
				temp[w] = true
			}
		}

		// Rule 3: Check AND gates with x/y inputs from different series
		if d.op == "AND" && isXOrY(d.w1) && isXOrY(d.w2) && d.w1[0] != d.w2[0] {
			isValid := false
			for _, dp := range dependencies {
				if dp.op == "OR" && (dp.w1 == w || dp.w2 == w) {
					isValid = true
				}
			}
			if !isValid {
				temp[w] = true
			}
		}
	}

	ans := make([]string, 0, len(temp))
	for w := range temp {
		ans = append(ans, w)
	}
	sort.Strings(ans)

	return strings.Join(ans, ",")
}

func main() {
	input, err := utility.ParseTextFile("input")
	if err != nil {
		log.Fatal(err)
	}
	value, dependencies := parseInput(input)
	log.Println("Part One:", partOne(value, dependencies))
	log.Println("Part Two:", partTwo(dependencies))
}
