package main

import (
	. "aoc2024/utility"
	"log"
)

// Find start and end positions
func findPosition(grid []string, char byte) [2]int {
	for i, line := range grid {
		for j, c := range line {
			if byte(c) == char {
				return [2]int{i, j}
			}
		}
	}
	return [2]int{-1, -1}
}

// Check if position is valid within grid
func isValid(i, j int, grid []string) bool {
	return i >= 0 && i < len(grid) && j >= 0 && j < len(grid[i])
}

// Check if char is valid path
func isValidChar(c byte) bool {
	return c == 'S' || c == 'E' || c == '.'
}

// Check if position exists in track map
func isInTrack(track map[[2]int]int, pos [2]int) bool {
	_, exists := track[pos]
	return exists
}

// Manhattan distance between two positions
func manhattanDistance(pos1, pos2 [2]int) int {
	return Abs(pos1[0]-pos2[0]) + Abs(pos1[1]-pos2[1])
}

// Find all possible cheat endpoints within maxDist
func findCheatEndpoints(pos [2]int, track map[[2]int]int, maxDist int) map[[2]int]struct{} {
	endpoints := make(map[[2]int]struct{})
	for di := -maxDist; di <= maxDist; di++ {
		maxJ := maxDist - Abs(di)
		for dj := -maxJ; dj <= maxJ; dj++ {
			newPos := [2]int{pos[0] + di, pos[1] + dj}
			if _, exists := track[newPos]; exists {
				endpoints[newPos] = struct{}{}
			}
		}
	}
	return endpoints
}

// Calculate if a cheat is valid and saves enough time
func isValidCheat(startPos, endPos [2]int, track map[[2]int]int, maxCheatLen int) bool {
	cheatDist := manhattanDistance(startPos, endPos)
	if cheatDist > maxCheatLen {
		return false
	}

	timeSaved := track[endPos] - track[startPos] - cheatDist
	return timeSaved >= 100
}

// Find next valid move in BFS
func findNextMove(cur [2]int, grid []string, track map[[2]int]int) ([2]int, bool) {
	directions := [][2]int{{-1, 0}, {0, -1}, {0, 1}, {1, 0}}
	i, j := cur[0], cur[1]

	for _, dir := range directions {
		newI, newJ := i+dir[0], j+dir[1]
		newPos := [2]int{newI, newJ}
		if isValid(newI, newJ, grid) && !isInTrack(track, newPos) && isValidChar(grid[newI][newJ]) {
			return newPos, true
		}
	}
	return [2]int{}, false
}

// Perform BFS to build track map
func buildTrackMap(grid []string, start, end [2]int) map[[2]int]int {
	track := make(map[[2]int]int)
	track[start] = 0
	cur := start
	curStep := 0

	for cur != end {
		curStep++
		nextPos, found := findNextMove(cur, grid, track)
		if !found {
			break
		}
		track[nextPos] = curStep
		cur = nextPos
	}

	return track
}

// Count valid cheats from track map
func countValidCheats(track map[[2]int]int, maxCheatLen int) int {
	count := 0
	for startPos := range track {
		endpoints := findCheatEndpoints(startPos, track, maxCheatLen)
		for endPos := range endpoints {
			if isValidCheat(startPos, endPos, track, maxCheatLen) {
				count++
			}
		}
	}
	return count
}

func solve(grid []string, maxCheatLen int) int {
	start := findPosition(grid, 'S')
	end := findPosition(grid, 'E')

	track := buildTrackMap(grid, start, end)
	return countValidCheats(track, maxCheatLen)
}

func part1(grid []string) int {
	return solve(grid, 2)
}

func part2(grid []string) int {
	return solve(grid, 20)
}

func main() {
	input, err := ParseTextFile("input")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(part1(input))
	log.Println(part2(input))
}
