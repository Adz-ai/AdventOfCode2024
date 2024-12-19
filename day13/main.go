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

// parseCoordinates extracts X and Y values from a string like "X+10, Y+5"
func parseCoordinates(s string) (x, y int64) {
	parts := strings.Split(s, "+")
	if len(parts) >= 3 {
		x, _ = strconv.ParseInt(strings.TrimRight(parts[1], ", Y"), 10, 64)
		y, _ = strconv.ParseInt(parts[2], 10, 64)
	}
	return x, y
}

// parsePrize extracts X and Y values from a string like "X=10, Y=5"
func parsePrize(s string) (x, y int64) {
	parts := strings.Split(s, "=")
	if len(parts) >= 3 {
		x, _ = strconv.ParseInt(strings.TrimRight(parts[1], ", Y"), 10, 64)
		y, _ = strconv.ParseInt(parts[2], 10, 64)
	}
	return x, y
}

// parseSystem extracts a complete system from three consecutive input lines
func parseSystem(lines []string, start int) (sys system, nextIndex int) {
	if start >= len(lines) {
		return sys, start
	}

	for i := start; i < len(lines) && i < start+3; i++ {
		line := strings.TrimSpace(lines[i])
		switch {
		case strings.HasPrefix(line, "Button A:"):
			sys.a1, sys.a2 = parseCoordinates(line)
		case strings.HasPrefix(line, "Button B:"):
			sys.b1, sys.b2 = parseCoordinates(line)
		case strings.HasPrefix(line, "Prize:"):
			sys.c1, sys.c2 = parsePrize(line)
		}
	}
	return sys, start + 3
}

// cramersRule solves the system using Cramer's rule
func (s system) cramersRule() (a, b int64) {
	det := s.a1*s.b2 - s.a2*s.b1
	if det == 0 {
		return 0, 0
	}

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
	var total int64
	offset := int64(0)
	if p2 {
		offset = 10000000000000
	}

	for i := 0; i < len(input); {
		// Skip empty lines
		if i >= len(input) || input[i] == "" {
			i++
			continue
		}

		// Parse system and update index
		sys, nextIndex := parseSystem(input, i)
		i = nextIndex

		// Apply offset for part 2
		if p2 {
			sys.c1 += offset
			sys.c2 += offset
		}

		// Calculate and add solution
		if a, b := sys.cramersRule(); a > 0 || b > 0 {
			total += 3*a + b
		}
	}
	return total
}

func main() {
	lines, err := utility.ParseTextFile("input")
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println(solve(lines, false))
	log.Println(solve(lines, true))
}
