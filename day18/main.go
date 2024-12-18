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

// Queue for BFS implementation
type Queue struct {
	items []Point
}

func (q *Queue) push(item Point) {
	q.items = append(q.items, item)
}

func (q *Queue) pop() Point {
	item := q.items[0]
	q.items = q.items[1:]
	return item
}

func (q *Queue) isEmpty() bool {
	return len(q.items) == 0
}

// parseInput reads coordinates from string slice and returns a slice of Points
func parseInput(input []string) ([]Point, error) {
	var points []Point

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
func createGrid(points []Point, size, n int) [][]bool {
	grid := make([][]bool, size)
	for i := range grid {
		grid[i] = make([]bool, size)
	}

	// Mark corrupted points (true = corrupted)
	numPoints := n
	if n > len(points) {
		numPoints = len(points)
	}

	for i := 0; i < numPoints; i++ {
		p := points[i]
		if p.x < size && p.y < size {
			grid[p.y][p.x] = true
		}
	}

	return grid
}

// findShortestPath finds the shortest path from start to end avoiding corrupted memory
func findShortestPath(grid [][]bool, start, end Point) int {
	sizeOfGrid := len(grid)
	visited := make([][]bool, sizeOfGrid)
	for i := range visited {
		visited[i] = make([]bool, sizeOfGrid)
	}

	// Directions: up, right, down, left
	dirs := []Point{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

	queue := Queue{}
	queue.push(start)
	visited[start.y][start.x] = true

	steps := 0
	nodesInCurrentLevel := 1
	nodesInNextLevel := 0

	for !queue.isEmpty() {
		current := queue.pop()
		nodesInCurrentLevel--

		if current == end {
			return steps
		}

		// Try all four directions
		for _, dir := range dirs {
			next := Point{current.x + dir.x, current.y + dir.y}

			// Check if next point is valid
			if next.x < 0 || next.x >= sizeOfGrid || next.y < 0 || next.y >= sizeOfGrid {
				continue
			}

			// Check if next point is unvisited and not corrupted
			if !visited[next.y][next.x] && !grid[next.y][next.x] {
				queue.push(next)
				visited[next.y][next.x] = true
				nodesInNextLevel++
			}
		}

		// Level complete, move to next level
		if nodesInCurrentLevel == 0 {
			nodesInCurrentLevel = nodesInNextLevel
			nodesInNextLevel = 0
			steps++
		}
	}

	return -1 // No path found
}

// hasPath checks if there is a path from start to end
func hasPath(grid [][]bool, start, end Point) bool {
	sizeOfGrid := len(grid)
	visited := make([][]bool, sizeOfGrid)
	for i := range visited {
		visited[i] = make([]bool, sizeOfGrid)
	}

	// Directions: up, right, down, left
	dirs := []Point{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

	queue := Queue{}
	queue.push(start)
	visited[start.y][start.x] = true

	for !queue.isEmpty() {
		current := queue.pop()

		if current == end {
			return true
		}

		// Try all four directions
		for _, dir := range dirs {
			next := Point{current.x + dir.x, current.y + dir.y}

			// Check if next point is valid
			if next.x < 0 || next.x >= sizeOfGrid || next.y < 0 || next.y >= sizeOfGrid {
				continue
			}

			// Check if next point is unvisited and not corrupted
			if !visited[next.y][next.x] && !grid[next.y][next.x] {
				queue.push(next)
				visited[next.y][next.x] = true
			}
		}
	}

	return false
}

// findBlockingByte finds the first byte that blocks all paths to the exit
func findBlockingByte(points []Point, size int) Point {
	start := Point{0, 0}
	end := Point{size - 1, size - 1}

	for i := 0; i < len(points); i++ {
		grid := createGrid(points, size, i+1)
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
		fmt.Println("No path found!")
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
		fmt.Println("No blocking byte found!")
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
