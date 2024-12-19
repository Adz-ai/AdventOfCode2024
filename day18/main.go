package main

import (
	"aoc2024/utility"
	"fmt"
	"log"
	"strconv"
	"strings"
)

const size = 71

// Point represents a coordinate in the memory space
type Point struct {
	x, y int
}

type queueItem struct {
	pt    Point
	steps int
}

// parseInput reads coordinates from string slice and returns a slice of Points
func parseInput(input []string) ([]Point, error) {
	points := make([]Point, 0, len(input))
	for _, line := range input {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		coords := strings.Split(line, ",")
		if len(coords) != 2 {
			return nil, fmt.Errorf("invalid input format: %s", line)
		}

		x, err := strconv.Atoi(coords[0])
		if err != nil {
			return nil, fmt.Errorf("invalid x coordinate: %s", coords[0])
		}

		y, err := strconv.Atoi(coords[1])
		if err != nil {
			return nil, fmt.Errorf("invalid y coordinate: %s", coords[1])
		}

		points = append(points, Point{x, y})
	}
	return points, nil
}

// createGrid creates a grid representation with the first n corrupted points
func createGrid(points []Point, gridSize, n int) [][]bool {
	grid := make([][]bool, gridSize)
	for i := range grid {
		grid[i] = make([]bool, gridSize)
	}

	numPoints := n
	if n > len(points) {
		numPoints = len(points)
	}

	for i := 0; i < numPoints; i++ {
		p := points[i]
		if p.x < gridSize && p.y < gridSize {
			grid[p.y][p.x] = true
		}
	}

	return grid
}

// getNeighbors returns all valid neighbors for a point within the grid bounds
func getNeighbors(p Point, sizeOfGrid int) []Point {
	dirs := []Point{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}
	neighbors := make([]Point, 0, 4)
	for _, d := range dirs {
		nx, ny := p.x+d.x, p.y+d.y
		if nx >= 0 && nx < sizeOfGrid && ny >= 0 && ny < sizeOfGrid {
			neighbors = append(neighbors, Point{nx, ny})
		}
	}
	return neighbors
}

// canVisit checks if a given point can be visited (not visited before and not corrupted)
func canVisit(grid [][]bool, visited [][]bool, p Point) bool {
	return !visited[p.y][p.x] && !grid[p.y][p.x]
}

// bfsGeneric performs a BFS and returns whether the end was found and the steps taken
// to reach it (or -1 if not found).
func bfsGeneric(grid [][]bool, start, end Point) (bool, int) {
	sizeOfGrid := len(grid)
	visited := make([][]bool, sizeOfGrid)
	for i := range visited {
		visited[i] = make([]bool, sizeOfGrid)
	}

	queue := []queueItem{{start, 0}}
	visited[start.y][start.x] = true

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current.pt == end {
			return true, current.steps
		}

		for _, next := range getNeighbors(current.pt, sizeOfGrid) {
			if canVisit(grid, visited, next) {
				visited[next.y][next.x] = true
				queue = append(queue, queueItem{next, current.steps + 1})
			}
		}
	}
	return false, -1
}

// bfsCheckPath returns true if a path exists from start to end, false otherwise.
func bfsCheckPath(grid [][]bool, start, end Point) bool {
	found, _ := bfsGeneric(grid, start, end)
	return found
}

// bfsShortestPath returns the number of steps in the shortest path from start to end, or -1 if none.
func bfsShortestPath(grid [][]bool, start, end Point) int {
	found, steps := bfsGeneric(grid, start, end)
	if found {
		return steps
	}
	return -1
}

// findShortestPath finds the shortest path from start to end avoiding corrupted memory
func findShortestPath(grid [][]bool, start, end Point) int {
	return bfsShortestPath(grid, start, end)
}

// hasPath checks if there is a path from start to end
func hasPath(grid [][]bool, start, end Point) bool {
	return bfsCheckPath(grid, start, end)
}

// findBlockingByte finds the first byte that blocks all paths to the exit
func findBlockingByte(points []Point, gridSize int) Point {
	start := Point{0, 0}
	end := Point{gridSize - 1, gridSize - 1}

	for i := 0; i < len(points); i++ {
		grid := createGrid(points, gridSize, i+1)
		if !hasPath(grid, start, end) {
			return points[i]
		}
	}

	return Point{-1, -1} // No blocking byte found
}

func part1(input []string) int {
	points, err := parseInput(input)
	if err != nil {
		log.Fatal(err)
	}

	n := 1024 // first kilobyte of bytes
	grid := createGrid(points, size, n)
	start := Point{0, 0}
	end := Point{size - 1, size - 1}

	steps := findShortestPath(grid, start, end)
	if steps == -1 {
		log.Println("No path found!")
	}
	return steps
}

func part2(input []string) string {
	points, err := parseInput(input)
	if err != nil {
		log.Fatal(err)
	}

	blockingByte := findBlockingByte(points, size)
	if blockingByte.x == -1 {
		log.Println("No blocking byte found!")
	}
	return fmt.Sprintf("%d,%d\n", blockingByte.x, blockingByte.y)
}

func main() {
	input, err := utility.ParseTextFile("input")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(part1(input))
	log.Println(part2(input))
}
