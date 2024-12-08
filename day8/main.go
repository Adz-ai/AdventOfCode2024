package main

import (
	"aoc2024/utility"
	"log"
)

type Point struct {
	row, col int
}

// buildFrequencyMap constructs a map of frequencies to antenna positions
func buildFrequencyMap(lines []string) map[rune][]Point {
	freqToAntennas := make(map[rune][]Point)
	for row, line := range lines {
		for col, char := range line {
			if char != '.' {
				freqToAntennas[char] = append(freqToAntennas[char], Point{row, col})
			}
		}
	}
	return freqToAntennas
}

// isWithinBounds checks if a point is within the grid boundaries
func isWithinBounds(p Point, bounds Point) bool {
	return p.row >= 0 && p.col >= 0 && p.row < bounds.row && p.col < bounds.col
}

func part1(input []string) int {
	freqToAntennas := buildFrequencyMap(input)
	antinodes := make(map[Point]bool)
	bounds := Point{len(input), len(input[0])}

	for _, antennas := range freqToAntennas {
		for i := 0; i < len(antennas); i++ {
			for j := i + 1; j < len(antennas); j++ {
				an1 := Point{
					row: 2*antennas[i].row - antennas[j].row,
					col: 2*antennas[i].col - antennas[j].col,
				}

				an2 := Point{
					row: 2*antennas[j].row - antennas[i].row,
					col: 2*antennas[j].col - antennas[i].col,
				}

				if isWithinBounds(an1, bounds) {
					antinodes[an1] = true
				}
				if isWithinBounds(an2, bounds) {
					antinodes[an2] = true
				}
			}
		}
	}
	return len(antinodes)
}

func part2(input []string) int {
	freqToAntennas := buildFrequencyMap(input)
	antinodes := make(map[Point]bool)
	bounds := Point{len(input), len(input[0])}

	for _, antennas := range freqToAntennas {
		for i := 0; i < len(antennas); i++ {
			for j := i + 1; j < len(antennas); j++ {
				dRow := antennas[i].row - antennas[j].row
				dCol := antennas[i].col - antennas[j].col

				nextPoint := Point{
					row: antennas[i].row - dRow,
					col: antennas[i].col - dCol,
				}
				for isWithinBounds(nextPoint, bounds) {
					antinodes[nextPoint] = true
					nextPoint.row -= dRow
					nextPoint.col -= dCol
				}

				nextPoint = Point{
					row: antennas[j].row + dRow,
					col: antennas[j].col + dCol,
				}
				for isWithinBounds(nextPoint, bounds) {
					antinodes[nextPoint] = true
					nextPoint.row += dRow
					nextPoint.col += dCol
				}

				if len(antennas) > 1 {
					antinodes[antennas[i]] = true
					antinodes[antennas[j]] = true
				}
			}
		}
	}
	return len(antinodes)
}

func main() {
	input, err := utility.ParseTextFile("input")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(part1(input))
	log.Println(part2(input))
}
