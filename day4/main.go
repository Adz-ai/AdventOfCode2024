package main

import (
	"aoc2024/utility"
	"log"
)

// Position represents a 2D coordinate
type Position struct {
	row, col int
}

// Direction represents a movement vector
type Direction struct {
	dr, dc int
}

// Pattern represents a configuration of M and S around center A
type Pattern struct {
	topLeft     rune
	topRight    rune
	bottomLeft  rune
	bottomRight rune
}

func checkPattern(grid []string, center Position, pattern Pattern) bool {
	return rune(grid[center.row-1][center.col-1]) == pattern.topLeft &&
		rune(grid[center.row-1][center.col+1]) == pattern.topRight &&
		rune(grid[center.row+1][center.col-1]) == pattern.bottomLeft &&
		rune(grid[center.row+1][center.col+1]) == pattern.bottomRight
}

func part2(input []string) int {
	patterns := []Pattern{
		{'M', 'M', 'S', 'S'}, // M.M.S.S pattern
		{'S', 'S', 'M', 'M'}, // S.S.M.M pattern
		{'M', 'S', 'M', 'S'}, // M.S.M.S pattern
		{'S', 'M', 'S', 'M'}, // S.M.S.M pattern
	}

	count := 0
	for i := 1; i < len(input)-1; i++ {
		for j := 1; j < len(input[0])-1; j++ {
			if rune(input[i][j]) != 'A' {
				continue
			}

			center := Position{i, j}
			for _, pattern := range patterns {
				if checkPattern(input, center, pattern) {
					count++
				}
			}
		}
	}
	return count
}

func matchWordInDirection(pos Position, dir Direction, word []rune, grid []string) bool {
	for k := 0; k < len(word); k++ {
		currentPos := Position{
			row: pos.row + k*dir.dr,
			col: pos.col + k*dir.dc,
		}
		if !isValidPosition(currentPos, grid) || getChar(currentPos, grid) != word[k] {
			return false
		}
	}
	return true
}

func isValidPosition(pos Position, grid []string) bool {
	return pos.row >= 0 && pos.row < len(grid) &&
		pos.col >= 0 && pos.col < len(grid[pos.row])
}

func getChar(pos Position, grid []string) rune {
	if !isValidPosition(pos, grid) {
		return 0
	}
	return rune(grid[pos.row][pos.col])
}

func part1(input []string) int {
	directions := []Direction{
		{0, 1}, {1, 0}, {1, 1}, {1, -1},
		{0, -1}, {-1, 0}, {-1, -1}, {-1, 1},
	}
	word := []rune("XMAS")
	count := 0

	for r := 0; r < len(input); r++ {
		for c := 0; c < len(input[r]); c++ {
			startPos := Position{r, c}
			for _, dir := range directions {
				if matchWordInDirection(startPos, dir, word, input) {
					count++
				}
			}
		}
	}
	return count
}

func main() {
	input, err := utility.ParseTextFile("input")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(part1(input))
	log.Println(part2(input))
}
