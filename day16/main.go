package main

import (
	"aoc2024/utility"
	"image"
	"log"
)

type State struct {
	pos image.Point
	dir image.Point
}

type RouteState struct {
	state State
	path  map[image.Point]struct{}
	cost  int
}

type Grid struct {
	cells map[image.Point]rune
	start image.Point
}

func parseGrid(lines []string) Grid {
	grid := Grid{
		cells: make(map[image.Point]rune),
	}
	for y, line := range lines {
		for x, r := range line {
			pos := image.Point{X: x, Y: y}
			if r == 'S' {
				grid.start = pos
			}
			grid.cells[pos] = r
		}
	}
	return grid
}

func getNeighbors(state State) []State {
	var moves []State

	// Move forward
	moves = append(moves, State{
		pos: state.pos.Add(state.dir),
		dir: state.dir,
	})

	// Turn right
	right := image.Point{X: -state.dir.Y, Y: state.dir.X}
	moves = append(moves, State{
		pos: state.pos.Add(right),
		dir: right,
	})

	// Turn left
	left := image.Point{X: state.dir.Y, Y: -state.dir.X}
	moves = append(moves, State{
		pos: state.pos.Add(left),
		dir: left,
	})

	return moves
}

func isDeadEnd(current RouteState, bestCost int) bool {
	if len(current.path) > 10000 {
		return true
	}
	if bestCost != -1 && current.cost > bestCost {
		return true
	}
	return false
}

func isGoal(pos image.Point, grid Grid) bool {
	return grid.cells[pos] == 'E'
}

func updateBestPaths(current RouteState, bestCost int, bestPaths map[int][]map[image.Point]struct{}) (int, map[int][]map[image.Point]struct{}) { //nolint:lll
	if bestCost == -1 || current.cost <= bestCost {
		newCost := current.cost
		if bestCost == -1 || newCost < bestCost {
			bestCost = newCost
			bestPaths[newCost] = nil
		}
		bestPaths[newCost] = append(bestPaths[newCost], current.path)
	}
	return bestCost, bestPaths
}

func validMove(grid Grid, next State) bool {
	return grid.cells[next.pos] != '#'
}

func shouldSkip(visited map[State]int, next State, newCost int) bool {
	if prevCost, exists := visited[next]; exists && prevCost < newCost {
		return true
	}
	return false
}

func copyPath(currentPath map[image.Point]struct{}, nextPos image.Point) map[image.Point]struct{} {
	newPath := make(map[image.Point]struct{}, len(currentPath)+1)
	for p := range currentPath {
		newPath[p] = struct{}{}
	}
	newPath[nextPos] = struct{}{}
	return newPath
}

func calculateCost(current RouteState, next State) int {
	newCost := current.cost + 1
	if next.dir != current.state.dir {
		newCost += 1000
	}
	return newCost
}

func collectAllPositions(bestPaths map[int][]map[image.Point]struct{}, bestCost int) int {
	allPositions := make(map[image.Point]struct{})
	for _, paths := range bestPaths[bestCost] {
		for pos := range paths {
			allPositions[pos] = struct{}{}
		}
	}
	return len(allPositions)
}

func findPath(grid Grid) (int, int) {
	visited := make(map[State]int)
	bestPaths := make(map[int][]map[image.Point]struct{})

	queue := []RouteState{{
		state: State{grid.start, image.Point{X: 1}}, // Start facing east
		path:  map[image.Point]struct{}{grid.start: {}},
		cost:  0,
	}}

	bestCost := -1

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if isDeadEnd(current, bestCost) {
			continue
		}

		if isGoal(current.state.pos, grid) {
			bestCost, bestPaths = updateBestPaths(current, bestCost, bestPaths)
			continue
		}

		for _, next := range getNeighbors(current.state) {
			if !validMove(grid, next) {
				continue
			}
			newCost := calculateCost(current, next)
			if shouldSkip(visited, next, newCost) {
				continue
			}
			visited[next] = newCost
			queue = append(queue, RouteState{
				state: next,
				path:  copyPath(current.path, next.pos),
				cost:  newCost,
			})
		}
	}

	if bestCost == -1 {
		return -1, 0
	}
	return bestCost, collectAllPositions(bestPaths, bestCost)
}

func solve(lines []string) (int, int) {
	grid := parseGrid(lines)
	return findPath(grid)
}

func main() {
	input, err := utility.ParseTextFile("input")
	if err != nil {
		log.Fatal(err)
	}
	part1, part2 := solve(input)
	log.Println(part1)
	log.Println(part2)
}
