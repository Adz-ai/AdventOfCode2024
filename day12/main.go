package main

import (
	"aoc2024/utility"
	"log"
)

// Point represents a position in the garden
type Point struct {
	row, col int
}

// Region represents a connected region of the same plant type
type Region struct {
	points    map[Point]bool
	plantType rune
}

// NewRegion creates a new Region with the given plant type
func NewRegion(plantType rune) *Region {
	return &Region{
		points:    make(map[Point]bool),
		plantType: plantType,
	}
}

// getBounds finds the bounds of the region in the garden
func (r *Region) getBounds(garden [][]rune) (minRow, maxRow, minCol, maxCol int) {
	rows, cols := len(garden), len(garden[0])
	minRow, maxRow = rows, -1
	minCol, maxCol = cols, -1

	for point := range r.points {
		if point.row < minRow {
			minRow = point.row
		}
		if point.row > maxRow {
			maxRow = point.row
		}
		if point.col < minCol {
			minCol = point.col
		}
		if point.col > maxCol {
			maxCol = point.col
		}
	}
	return
}

// isValidTransition checks if a point and its below neighbor form a valid region transition
func (r *Region) isValidTransition(current, below Point, rows int) bool {
	if below.row >= rows {
		return false
	}

	hasCurrentPoint := r.points[current]
	hasBelowPoint := r.points[below]

	return !hasCurrentPoint && hasBelowPoint
}

// countSidesFromOrientation counts sides from current orientation
func (r *Region) countSidesFromOrientation(garden [][]rune, minRow, maxRow, minCol, maxCol int) int {
	rows := len(garden)
	sides := 0

	for row := minRow - 1; row <= maxRow; row++ {
		inRegion := false
		for col := minCol - 1; col <= maxCol; col++ {
			current := Point{row, col}
			below := Point{row + 1, col}

			if r.isValidTransition(current, below, rows) {
				if !inRegion {
					sides++
					inRegion = true
				}
			} else {
				inRegion = false
			}
		}
	}
	return sides
}

// rotateAndNormalize rotates the region and normalizes coordinates
func (r *Region) rotateAndNormalize(garden [][]rune) {
	// Rotate 90 degrees clockwise
	newPoints := make(map[Point]bool)
	for point := range r.points {
		newPoint := Point{
			row: point.col,
			col: -point.row,
		}
		newPoints[newPoint] = true
	}
	r.points = newPoints

	// Normalize coordinates
	minRow, _, minCol, _ := r.getBounds(garden)
	if minRow < 0 || minCol < 0 {
		normalized := make(map[Point]bool)
		for point := range r.points {
			normalized[Point{
				row: point.row - minRow,
				col: point.col - minCol,
			}] = true
		}
		r.points = normalized
	}
}

// countSides counts the number of distinct sides of the region
func (r *Region) countSides(garden [][]rune) int {
	totalSides := 0

	for rotation := 0; rotation < 4; rotation++ {
		minRow, maxRow, minCol, maxCol := r.getBounds(garden)
		totalSides += r.countSidesFromOrientation(garden, minRow, maxRow, minCol, maxCol)
		r.rotateAndNormalize(garden)
	}

	return totalSides
}

// calculatePerimeter counts the number of edges that don't connect to the same plant type
func (r *Region) calculatePerimeter(garden [][]rune) int {
	perimeter := 0
	rows, cols := len(garden), len(garden[0])

	for point := range r.points {
		// Each point contributes 4 to the perimeter initially
		edges := 4

		// Subtract 1 for each adjacent point of the same type
		for _, neighbor := range getNeighbors(point, rows, cols) {
			if garden[neighbor.row][neighbor.col] == r.plantType {
				if r.points[neighbor] {
					edges--
				}
			}
		}
		perimeter += edges
	}

	return perimeter
}

// getNeighbors returns valid adjacent points
func getNeighbors(p Point, rows, cols int) []Point {
	neighbors := []Point{
		{p.row - 1, p.col}, // up
		{p.row + 1, p.col}, // down
		{p.row, p.col - 1}, // left
		{p.row, p.col + 1}, // right
	}

	valid := make([]Point, 0, 4)
	for _, n := range neighbors {
		if n.row >= 0 && n.row < rows && n.col >= 0 && n.col < cols {
			valid = append(valid, n)
		}
	}
	return valid
}

// findRegion performs a flood fill to find all connected points of the same plant type
func findRegion(garden [][]rune, start Point, visited map[Point]bool) *Region {
	if visited[start] {
		return nil
	}

	plantType := garden[start.row][start.col]
	region := NewRegion(plantType)
	rows, cols := len(garden), len(garden[0])

	// Stack-based flood fill
	stack := []Point{start}
	visited[start] = true
	region.points[start] = true

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		for _, neighbor := range getNeighbors(current, rows, cols) {
			if !visited[neighbor] && garden[neighbor.row][neighbor.col] == plantType {
				visited[neighbor] = true
				region.points[neighbor] = true
				stack = append(stack, neighbor)
			}
		}
	}

	return region
}

func part1(lines []string) int {
	garden := make([][]rune, len(lines))
	for i, line := range lines {
		garden[i] = []rune(line)
	}

	rows, cols := len(garden), len(garden[0])
	visited := make(map[Point]bool)
	totalPrice := 0

	// Find all regions
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			point := Point{row, col}
			if region := findRegion(garden, point, visited); region != nil {
				area := len(region.points)
				perimeter := region.calculatePerimeter(garden)
				price := area * perimeter
				totalPrice += price
			}
		}
	}

	return totalPrice
}

func part2(lines []string) int {
	garden := make([][]rune, len(lines))
	for i, line := range lines {
		garden[i] = []rune(line)
	}

	rows, cols := len(garden), len(garden[0])
	visited := make(map[Point]bool)
	totalPrice := 0

	// Find all regions
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			point := Point{row, col}
			if region := findRegion(garden, point, visited); region != nil {
				area := len(region.points)
				sides := region.countSides(garden)
				price := area * sides
				totalPrice += price
			}
		}
	}

	return totalPrice
}

func main() {
	lines, err := utility.ParseTextFile("input")
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println(part1(lines))
	log.Println(part2(lines))
}
