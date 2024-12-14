package main

import (
	"aoc2024/utility"
	"log"
	"math"
	"strconv"
	"strings"
)

type Vector struct {
	x, y int
}

type Robot struct {
	pos Vector
	vel Vector
}

// Move robot one step, handling wraparound
func (r *Robot) move(width, height int) {
	r.pos.x = (r.pos.x + r.vel.x + width) % width
	r.pos.y = (r.pos.y + r.vel.y + height) % height
}

// Parse a line like "p=0,4 v=3,-3"
func parseRobot(line string) Robot {
	parts := strings.Split(line, " ")

	// Parse position
	pos := parts[0][2:] // Skip "p="
	posNums := strings.Split(pos, ",")
	x, _ := strconv.Atoi(posNums[0])
	y, _ := strconv.Atoi(posNums[1])
	position := Vector{x, y}

	// Parse velocity
	vel := parts[1][2:] // Skip "v="
	velNums := strings.Split(vel, ",")
	dx, _ := strconv.Atoi(velNums[0])
	dy, _ := strconv.Atoi(velNums[1])
	velocity := Vector{dx, dy}

	return Robot{position, velocity}
}

// Calculate standard deviation for a slice of numbers
func standardDeviation(nums []float64) float64 {
	if len(nums) == 0 {
		return 0
	}

	// Calculate mean
	var sum float64
	for _, num := range nums {
		sum += num
	}
	mean := sum / float64(len(nums))

	// Calculate variance
	var variance float64
	for _, num := range nums {
		diff := num - mean
		variance += diff * diff
	}
	variance /= float64(len(nums))

	return math.Sqrt(variance)
}

func part1(input []string) int {
	var robots []Robot
	for _, line := range input {
		if line == "" {
			continue
		}
		robots = append(robots, parseRobot(line))
	}

	width := 101
	height := 103

	// Simulate 100 seconds
	for t := 0; t < 100; t++ {
		for i := range robots {
			robots[i].move(width, height)
		}
	}

	// Count quadrants
	quadrants := [2][2]int{}
	midX := width / 2
	midY := height / 2

	for _, robot := range robots {
		if robot.pos.x == midX || robot.pos.y == midY {
			continue
		}

		row := 0
		if robot.pos.y > midY {
			row = 1
		}
		col := 0
		if robot.pos.x >= midX {
			col = 1
		}
		quadrants[row][col]++
	}

	// Calculate product
	product := 1
	for _, row := range quadrants {
		for _, count := range row {
			product *= count
		}
	}

	return product
}

func part2(input []string) int {
	var robots []Robot
	for _, line := range input {
		if line == "" {
			continue
		}
		robots = append(robots, parseRobot(line))
	}

	width := 101
	height := 103

	// Simulate until we find the pattern
	for t := 1; ; t++ {
		// Move all robots
		for i := range robots {
			robots[i].move(width, height)
		}

		// Collect x and y coordinates
		xCoords := make([]float64, len(robots))
		yCoords := make([]float64, len(robots))
		for i, robot := range robots {
			xCoords[i] = float64(robot.pos.x)
			yCoords[i] = float64(robot.pos.y)
		}

		// Check standard deviations
		// When robots form a Christmas tree pattern, they'll be clustered together
		// resulting in low standard deviations for both x and y coordinates
		xStdDev := standardDeviation(xCoords)
		yStdDev := standardDeviation(yCoords)

		// Based on Reddit insight: Christmas tree is ~33x31, and std dev is <20 for both x and y
		if xStdDev < 20 && yStdDev < 20 {
			return t
		}
	}
}

func main() {
	input, err := utility.ParseTextFile("input")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(part1(input))
	log.Println(part2(input))
}
