package main

import (
	"aoc2024/utility"
	"log"
)

const UP = 0

type Position struct {
	x, y int
}

var directions = [4][2]int{
	{0, -1}, // UP
	{1, 0},  // RIGHT
	{0, 1},  // DOWN
	{-1, 0}, // LEFT
}

var width, height int

// convertToByteGrid converts a slice of strings into a 2D byte grid for efficient processing
func convertToByteGrid(lines []string) [][]byte {
	grid := make([][]byte, len(lines))
	for i, line := range lines {
		grid[i] = []byte(line)
	}
	return grid
}

// findStart locates the guard's starting position (marked by '^') in the grid
// Returns a Position with x,y coordinates
func findStart(grid [][]byte) Position {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if grid[y][x] == '^' {
				return Position{x, y}
			}
		}
	}
	return Position{}
}

// getPath simulates the guard's movement and returns a slice of all positions visited
// The guard follows these rules:
// - If there's an obstacle ahead, turn right
// - Otherwise, move forward
// - Stop when reaching the grid boundary
func getPath(grid [][]byte, start Position) []Position {
	visited := make([]bool, width*height)
	var path []Position

	pos := start
	dir := UP

	for {
		idx := pos.y*width + pos.x
		if !visited[idx] {
			visited[idx] = true
			path = append(path, pos)
		}

		nextPos := Position{
			x: pos.x + directions[dir][0],
			y: pos.y + directions[dir][1],
		}

		if nextPos.x < 0 || nextPos.x >= width || nextPos.y < 0 || nextPos.y >= height {
			return path
		}

		if grid[nextPos.y][nextPos.x] == '#' {
			dir = (dir + 1) & 3
		} else {
			pos = nextPos
		}
	}
}

// countLoopPositions counts how many positions, when blocked, would cause the guard
// to enter an infinite loop. Tests each position in the initial path by temporarily
// placing an obstacle and checking for a loop.
func countLoopPositions(grid [][]byte, start Position, initialPath []Position) int {
	count := 0
	for _, pos := range initialPath {
		if grid[pos.y][pos.x] == '.' {
			grid[pos.y][pos.x] = '#'
			if hasLoop(grid, start) {
				count++
			}
			grid[pos.y][pos.x] = '.'
		}
	}
	return count
}

// hasLoop checks if the guard's path contains a loop by tracking visited positions
// and their entry direction using bit flags. Returns true if the same position
// is visited in the same direction twice.
func hasLoop(grid [][]byte, start Position) bool {
	visited := make([]uint8, width*height)
	pos := start
	dir := UP

	for {
		idx := pos.y*width + pos.x
		dirBit := uint8(1 << dir)

		if visited[idx]&dirBit != 0 {
			return true
		}
		visited[idx] |= dirBit

		nextPos := Position{
			x: pos.x + directions[dir][0],
			y: pos.y + directions[dir][1],
		}

		if nextPos.x < 0 || nextPos.x >= width || nextPos.y < 0 || nextPos.y >= height {
			return false
		}

		if grid[nextPos.y][nextPos.x] == '#' {
			dir = (dir + 1) & 3
		} else {
			pos = nextPos
		}
	}
}

// solve is a helper function that handles the common setup for both parts:
// converting input to a grid, setting dimensions, finding start position,
// and running the provided solver function
func solve(input []string, solver func(grid [][]byte, startPos Position) int) int {
	grid := convertToByteGrid(input)
	height, width = len(grid), len(grid[0])
	startPos := findStart(grid)
	return solver(grid, startPos)
}

func part1(input []string) int {
	return solve(input, func(grid [][]byte, startPos Position) int {
		return len(getPath(grid, startPos))
	})
}

func part2(input []string) int {
	return solve(input, func(grid [][]byte, startPos Position) int {
		return countLoopPositions(grid, startPos, getPath(grid, startPos))
	})
}

func main() {
	input, err := utility.ParseTextFile("input")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(part1(input))
	log.Println(part2(input))
}
