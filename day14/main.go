package main

import (
	"aoc2024/utility"
	"log"
	"math"
	"strconv"
	"strings"
)

// Vector represents a 2D coordinate or velocity vector
type Vector struct {
	x, y int
}

// Robot represents a robot with position and velocity vectors
type Robot struct {
	pos Vector
	vel Vector
}

// move updates the robot's position based on its velocity, handling wraparound
// within the given width and height bounds
func (r *Robot) move(width, height int) {
	r.pos.x = (r.pos.x + r.vel.x + width) % width
	r.pos.y = (r.pos.y + r.vel.y + height) % height
}

// parseRobot converts an input line in the format "p=x,y v=dx,dy" into a Robot struct
// Example input: "p=0,4 v=3,-3"
func parseRobot(line string) Robot {
	parts := strings.Split(line, " ")
	pos := strings.Split(parts[0][2:], ",") // Skip "p="
	vel := strings.Split(parts[1][2:], ",") // Skip "v="

	x, _ := strconv.Atoi(pos[0])
	y, _ := strconv.Atoi(pos[1])
	dx, _ := strconv.Atoi(vel[0])
	dy, _ := strconv.Atoi(vel[1])

	return Robot{Vector{x, y}, Vector{dx, dy}}
}

// parseRobots converts multiple input lines into a slice of Robots,
// skipping empty lines
func parseRobots(input []string) []Robot {
	robots := make([]Robot, 0, len(input))
	for _, line := range input {
		if line != "" {
			robots = append(robots, parseRobot(line))
		}
	}
	return robots
}

// simulateRobots moves all robots for a specified number of steps
// within the given width and height bounds
func simulateRobots(robots []Robot, steps, width, height int) {
	for t := 0; t < steps; t++ {
		for i := range robots {
			robots[i].move(width, height)
		}
	}
}

// getQuadrant determines which quadrant a position falls into based on midpoints.
// Returns row (0/1), column (0/1), and whether the position is valid (not on axes)
func getQuadrant(pos Vector, midX, midY int) (int, int, bool) {
	if pos.x == midX || pos.y == midY {
		return 0, 0, false
	}

	row := 0
	if pos.y > midY {
		row = 1
	}
	col := 0
	if pos.x >= midX {
		col = 1
	}

	return row, col, true
}

// countQuadrants counts robots in each quadrant and returns their product.
// Robots on axes are ignored
func countQuadrants(robots []Robot, width, height int) int {
	quadrants := [2][2]int{}
	midX := width / 2
	midY := height / 2

	for _, robot := range robots {
		row, col, valid := getQuadrant(robot.pos, midX, midY)
		if valid {
			quadrants[row][col]++
		}
	}

	product := 1
	for _, row := range quadrants {
		for _, count := range row {
			product *= count
		}
	}
	return product
}

// getPositionCoordinates extracts x and y coordinates from robots into separate slices
func getPositionCoordinates(robots []Robot) ([]float64, []float64) {
	xCoords := make([]float64, len(robots))
	yCoords := make([]float64, len(robots))

	for i, robot := range robots {
		xCoords[i] = float64(robot.pos.x)
		yCoords[i] = float64(robot.pos.y)
	}

	return xCoords, yCoords
}

// calculateMean computes the arithmetic mean of a slice of numbers
func calculateMean(nums []float64) float64 {
	var sum float64
	for _, num := range nums {
		sum += num
	}
	return sum / float64(len(nums))
}

// standardDeviation calculates the standard deviation of a slice of numbers.
// Returns 0 for empty slices
func standardDeviation(nums []float64) float64 {
	if len(nums) == 0 {
		return 0
	}

	mean := calculateMean(nums)
	var variance float64
	for _, num := range nums {
		diff := num - mean
		variance += diff * diff
	}
	variance /= float64(len(nums))
	return math.Sqrt(variance)
}

// isChristmasTreePattern determines if robots have formed the Christmas tree pattern
// by checking if their standard deviations are below a threshold
func isChristmasTreePattern(robots []Robot) bool {
	xCoords, yCoords := getPositionCoordinates(robots)
	xStdDev := standardDeviation(xCoords)
	yStdDev := standardDeviation(yCoords)

	const maxStdDev = 20.0
	return xStdDev < maxStdDev && yStdDev < maxStdDev
}

// simulateUntilPattern continues moving robots until they form the Christmas tree pattern.
// Returns the number of steps taken
func simulateUntilPattern(robots []Robot, width, height int) int {
	for time := 1; ; time++ {
		for i := range robots {
			robots[i].move(width, height)
		}

		if isChristmasTreePattern(robots) {
			return time
		}
	}
}

// part1 solves the first part of the puzzle:
// Simulates robot movement for 100 steps and calculates the product of robots in each quadrant
func part1(input []string) int {
	const (
		width  = 101
		height = 103
		steps  = 100
	)

	robots := parseRobots(input)
	simulateRobots(robots, steps, width, height)
	return countQuadrants(robots, width, height)
}

// part2 solves the second part of the puzzle:
// Finds how many steps it takes for robots to form a Christmas tree pattern
func part2(input []string) int {
	const (
		width  = 101
		height = 103
	)

	robots := parseRobots(input)
	return simulateUntilPattern(robots, width, height)
}

func main() {
	input, err := utility.ParseTextFile("input")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(part1(input))
	log.Println(part2(input))
}
