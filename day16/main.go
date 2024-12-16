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
	// Define possible moves based on current direction
	var moves []State

	// Always try moving forward
	moves = append(moves, State{
		pos: state.pos.Add(state.dir),
		dir: state.dir,
	})

	// Turn right (90 degrees clockwise)
	right := image.Point{X: -state.dir.Y, Y: state.dir.X}
	moves = append(moves, State{
		pos: state.pos.Add(right),
		dir: right,
	})

	// Turn left (90 degrees counterclockwise)
	left := image.Point{X: state.dir.Y, Y: -state.dir.X}
	moves = append(moves, State{
		pos: state.pos.Add(left),
		dir: left,
	})

	return moves
}

func findPath(grid Grid) (int, int) {
	visited := make(map[State]int)
	queue := []RouteState{{
		state: State{grid.start, image.Point{X: 1}}, // Start facing east
		path:  map[image.Point]struct{}{grid.start: {}},
		cost:  0,
	}}

	bestCost := -1
	bestPaths := make(map[int][]map[image.Point]struct{})

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// Skip if path is too long
		if len(current.path) > 10000 {
			continue
		}

		// Skip if we've already found better paths
		if bestCost != -1 && current.cost > bestCost {
			continue
		}

		// Check if we've reached the end
		if grid.cells[current.state.pos] == 'E' {
			if bestCost == -1 || current.cost <= bestCost {
				bestCost = current.cost
				bestPaths[bestCost] = append(bestPaths[bestCost], current.path)
			}
			continue
		}

		// Try each possible move
		for _, next := range getNeighbors(current.state) {
			if grid.cells[next.pos] == '#' {
				continue
			}

			// Calculate new cost
			newCost := current.cost + 1
			if next.dir != current.state.dir {
				newCost += 1000
			}

			// Skip if we've found a better path to this state
			if prevCost, exists := visited[next]; exists && prevCost < newCost {
				continue
			}
			visited[next] = newCost

			// Create new path
			newPath := make(map[image.Point]struct{})
			for pos := range current.path {
				newPath[pos] = struct{}{}
			}
			newPath[next.pos] = struct{}{}

			// Add to queue
			queue = append(queue, RouteState{
				state: next,
				path:  newPath,
				cost:  newCost,
			})
		}
	}

	// Count unique positions in all best paths
	allPositions := make(map[image.Point]struct{})
	for _, paths := range bestPaths[bestCost] {
		for pos := range paths {
			allPositions[pos] = struct{}{}
		}
	}

	return bestCost, len(allPositions)
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
