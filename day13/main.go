package main

import (
	"aoc2024/utility"
	"log"
	"strconv"
	"strings"
)

type system struct {
	a1, a2 int64 // Button A's X and Y movements
	b1, b2 int64 // Button B's X and Y movements
	c1, c2 int64 // Prize X and Y coordinates
}

// solve uses Cramer's rule to find the solution to the system of equations
func (s system) cramersRule() (a, b int64) {
	// Calculate determinant for the system
	det := s.a1*s.b2 - s.a2*s.b1

	// Calculate numerators for a and b using Cramer's rule
	aNum := s.c1*s.b2 - s.b1*s.c2
	bNum := s.a1*s.c2 - s.c1*s.a2

	// Check if we have exact integer solutions
	if a = aNum / det; a < 0 || a*det != aNum {
		return 0, 0
	}
	if b = bNum / det; b < 0 || b*det != bNum {
		return 0, 0
	}

	return a, b
}

func solve(input []string, p2 bool) int64 {
	var offset int64 = 10000000000000
	var total int64
	var sys system

	for i := 0; i < len(input); {
		if i >= len(input) || input[i] == "" {
			i++
			continue
		}

		// Parse Button A
		if strings.Contains(input[i], "Button A:") {
			parts := strings.Split(input[i], "+")
			sys.a1, _ = strconv.ParseInt(strings.TrimRight(parts[1], ", Y"), 10, 64)
			sys.a2, _ = strconv.ParseInt(parts[2], 10, 64)
		}
		i++

		// Parse Button B
		if i < len(input) && strings.Contains(input[i], "Button B:") {
			parts := strings.Split(input[i], "+")
			sys.b1, _ = strconv.ParseInt(strings.TrimRight(parts[1], ", Y"), 10, 64)
			sys.b2, _ = strconv.ParseInt(parts[2], 10, 64)
		}
		i++

		// Parse Prize
		if i < len(input) && strings.Contains(input[i], "Prize:") {
			parts := strings.Split(input[i], "=")
			sys.c1, _ = strconv.ParseInt(strings.TrimRight(parts[1], ", Y"), 10, 64)
			sys.c2, _ = strconv.ParseInt(parts[2], 10, 64)

			// Offset coordinates for part 2
			if p2 {
				sys.c1 += offset
				sys.c2 += offset
			}
		}
		i++

		// Solve the system and add to total
		a, b := sys.cramersRule()
		if a > 0 || b > 0 { // Only add if solution exists
			total += 3*a + b
		}
	}
	return total
}

func part1(input []string) int64 {
	return solve(input, false)
}

func part2(input []string) int64 {
	return solve(input, true)
}
func main() {
	lines, err := utility.ParseTextFile("input")
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println(part1(lines))
	log.Println(part2(lines))
}
