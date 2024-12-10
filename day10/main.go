package main

import (
	"aoc2024/utility"
	"log"
)

type Position struct {
	row, col int
}

type Grid [][]int

func parseGrid(lines []string) Grid {
	grid := make(Grid, len(lines))
	for i, line := range lines {
		grid[i] = make([]int, len(line))
		for j, ch := range line {
			grid[i][j] = int(ch - '0')
		}
	}
	return grid
}

func findTrailheads(grid Grid) []Position {
	var trailheads []Position
	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] == 0 {
				trailheads = append(trailheads, Position{i, j})
			}
		}
	}
	return trailheads
}

func isValidPosition(grid Grid, pos Position) bool {
	return pos.row >= 0 && pos.row < len(grid) && pos.col >= 0 && pos.col < len(grid[0])
}

func calculateTrailheadScore(grid Grid, start Position) int {
	reachableNines := make(map[Position]bool)

	var dfs func(pos Position, currentHeight int, path map[Position]bool)
	dfs = func(pos Position, currentHeight int, path map[Position]bool) {
		if grid[pos.row][pos.col] == 9 {
			reachableNines[pos] = true
			return
		}

		directions := []Position{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
		for _, dir := range directions {
			newPos := Position{pos.row + dir.row, pos.col + dir.col}

			if !isValidPosition(grid, newPos) {
				continue
			}

			newHeight := grid[newPos.row][newPos.col]
			if newHeight == currentHeight+1 && !path[newPos] {
				path[newPos] = true
				dfs(newPos, newHeight, path)
				delete(path, newPos)
			}
		}
	}

	initialPath := map[Position]bool{start: true}
	dfs(start, 0, initialPath)

	return len(reachableNines)
}

func calculateTrailheadRating(grid Grid, start Position) int {
	pathCount := 0

	var dfs func(pos Position, currentHeight int, path map[Position]bool)
	dfs = func(pos Position, currentHeight int, path map[Position]bool) {
		if grid[pos.row][pos.col] == 9 {
			pathCount++
			return
		}

		directions := []Position{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
		for _, dir := range directions {
			newPos := Position{pos.row + dir.row, pos.col + dir.col}

			if !isValidPosition(grid, newPos) {
				continue
			}

			newHeight := grid[newPos.row][newPos.col]
			if newHeight == currentHeight+1 && !path[newPos] {
				path[newPos] = true
				dfs(newPos, newHeight, path)
				delete(path, newPos)
			}
		}
	}

	initialPath := map[Position]bool{start: true}
	dfs(start, 0, initialPath)

	return pathCount
}

func part1(input []string) int {
	grid := parseGrid(input)
	totalScore := 0

	trailheads := findTrailheads(grid)

	for _, start := range trailheads {
		score := calculateTrailheadScore(grid, start)
		totalScore += score
	}

	return totalScore
}

func part2(input []string) int {
	grid := parseGrid(input)
	totalRating := 0

	trailheads := findTrailheads(grid)
	for _, start := range trailheads {
		rating := calculateTrailheadRating(grid, start)
		totalRating += rating
	}

	return totalRating
}

func main() {
	input, err := utility.ParseTextFile("input")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(part1(input))
	log.Println(part2(input))
}
