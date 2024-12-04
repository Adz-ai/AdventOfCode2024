package main

import (
	"aoc2024/utility"
	"log"
)

// Position represents a 2D coordinate in the grid with row and column values.
type Position struct {
	row, col int
}

// Direction represents a movement vector in the grid with row (dr) and column (dc) deltas.
type Direction struct {
	dr, dc int
}

// Pattern represents a configuration of M and S characters around a center A.
// The pattern forms an X shape with the specified runes in each corner position.
type Pattern struct {
	topLeft     rune
	topRight    rune
	bottomLeft  rune
	bottomRight rune
}

// checkPattern verifies if the characters around the center position match
// the specified pattern configuration.
//
// Parameters:
//   - grid: The input grid of characters
//   - center: The center Position to check around (must be 'A')
//   - pattern: The Pattern configuration to match against
//
// Returns:
//   - bool: true if the pattern matches, false otherwise
func checkPattern(grid []string, center Position, pattern Pattern) bool {
	return rune(grid[center.row-1][center.col-1]) == pattern.topLeft &&
		rune(grid[center.row-1][center.col+1]) == pattern.topRight &&
		rune(grid[center.row+1][center.col-1]) == pattern.bottomLeft &&
		rune(grid[center.row+1][center.col+1]) == pattern.bottomRight
}

// matchWordInDirection checks if a word matches in the grid starting from a position
// and moving in a specified direction.
//
// Parameters:
//   - pos: Starting Position in the grid
//   - dir: Direction to check for the word
//   - word: Slice of runes representing the word to match
//   - grid: The input grid of characters
//
// Returns:
//   - bool: true if the word matches, false otherwise
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

// isValidPosition checks if a position is within the bounds of the grid.
//
// Parameters:
//   - pos: Position to validate
//   - grid: The input grid to check against
//
// Returns:
//   - bool: true if the position is valid, false otherwise
func isValidPosition(pos Position, grid []string) bool {
	return pos.row >= 0 && pos.row < len(grid) &&
		pos.col >= 0 && pos.col < len(grid[pos.row])
}

// getChar safely retrieves a character from the grid at the specified position.
//
// Parameters:
//   - pos: Position to retrieve the character from
//   - grid: The input grid of characters
//
// Returns:
//   - rune: The character at the position, or 0 if the position is invalid
func getChar(pos Position, grid []string) rune {
	if !isValidPosition(pos, grid) {
		return 0
	}
	return rune(grid[pos.row][pos.col])
}

func part2(input []string) int {
	patterns := []Pattern{
		{'M', 'M', 'S', 'S'},
		{'S', 'S', 'M', 'M'},
		{'M', 'S', 'M', 'S'},
		{'S', 'M', 'S', 'M'},
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
