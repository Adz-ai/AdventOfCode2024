package main

import (
	"aoc2024/utility"
	"log"
)

// Direction represents a movement direction with its row and column deltas
type Direction struct {
	dr, dc int
}

var (
	directions = map[byte]Direction{
		'^': {dr: -1, dc: 0},
		'v': {dr: 1, dc: 0},
		'<': {dr: 0, dc: -1},
		'>': {dr: 0, dc: 1},
	}
)

// Part 1 functions
func canMovePart1(board [][]byte, r, c int, dir Direction) bool {
	nextR, nextC := r+dir.dr, c+dir.dc

	// Check bounds
	if nextR < 0 || nextR >= len(board) || nextC < 0 || nextC >= len(board[0]) {
		return false
	}

	if board[nextR][nextC] == 'O' {
		return canMovePart1(board, nextR, nextC, dir)
	}

	return board[nextR][nextC] == '.'
}

func movePart1(board [][]byte, r, c int, dir Direction) (int, int) {
	if !canMovePart1(board, r, c, dir) {
		return r, c
	}

	nextR, nextC := r+dir.dr, c+dir.dc

	if board[nextR][nextC] == 'O' {
		movePart1(board, nextR, nextC, dir)
	}

	board[nextR][nextC] = board[r][c]
	board[r][c] = '.'

	return nextR, nextC
}

// Part 2 functions
func canMovePart2(board [][]byte, r, c int, dir Direction) bool {
	nextR, nextC := r+dir.dr, c+dir.dc

	if nextR < 0 || nextR >= len(board) || nextC < 0 || nextC >= len(board[0]) {
		return false
	}

	if board[nextR][nextC] == ']' {
		if !canMovePart2(board, nextR, nextC, dir) {
			return false
		}
		if dir.dr != 0 { // vertical movement
			return canMovePart2(board, nextR, nextC-1, dir)
		}
		return true
	}

	if board[nextR][nextC] == '[' {
		if !canMovePart2(board, nextR, nextC, dir) {
			return false
		}
		if dir.dr != 0 { // vertical movement
			return canMovePart2(board, nextR, nextC+1, dir)
		}
		return true
	}

	return board[nextR][nextC] == '.'
}

func movePart2(board [][]byte, r, c int, dir Direction) (int, int) {
	if !canMovePart2(board, r, c, dir) {
		return r, c
	}

	nextR, nextC := r+dir.dr, c+dir.dc

	if board[nextR][nextC] == ']' {
		movePart2(board, nextR, nextC, dir)
		if dir.dr != 0 { // vertical movement
			movePart2(board, nextR, nextC-1, dir)
		}
	}

	if board[nextR][nextC] == '[' {
		movePart2(board, nextR, nextC, dir)
		if dir.dr != 0 { // vertical movement
			movePart2(board, nextR, nextC+1, dir)
		}
	}

	board[nextR][nextC] = board[r][c]
	board[r][c] = '.'

	return nextR, nextC
}

func part1(input []string) int {
	var board [][]byte
	var robotRow, robotCol int
	var movesIdx int

	// Create the board
	for i, line := range input {
		if len(line) == 0 {
			movesIdx = i
			break
		}
		boardLine := []byte(line)
		board = append(board, boardLine)

		// Find robot position
		if pos := indexOf(boardLine, '@'); pos != -1 {
			robotRow, robotCol = len(board)-1, pos
		}
	}

	// Parse moves
	var moves []byte
	for _, line := range input[movesIdx+1:] {
		moves = append(moves, []byte(line)...)
	}

	// Process all moves
	for _, m := range moves {
		dir := directions[m]
		robotRow, robotCol = movePart1(board, robotRow, robotCol, dir)
	}

	// Calculate final GPS coordinates
	result := 0
	for i, row := range board {
		for j, cell := range row {
			if cell == 'O' {
				result += i*100 + j
			}
		}
	}

	return result
}

func part2(input []string) int {
	var board [][]byte
	var robotRow, robotCol int
	var movesIdx int

	// Create the expanded board
	for i, line := range input {
		if len(line) == 0 {
			movesIdx = i
			break
		}

		var newLine []byte
		for j, ch := range line {
			switch ch {
			case '#':
				newLine = append(newLine, '#', '#')
			case 'O':
				newLine = append(newLine, '[', ']')
			case '.':
				newLine = append(newLine, '.', '.')
			case '@':
				newLine = append(newLine, '@', '.')
				robotRow, robotCol = len(board), 2*j
			}
		}
		board = append(board, newLine)
	}

	// Parse moves
	var moves []byte
	for _, line := range input[movesIdx+1:] {
		moves = append(moves, []byte(line)...)
	}

	// Process all moves
	for _, m := range moves {
		dir := directions[m]
		robotRow, robotCol = movePart2(board, robotRow, robotCol, dir)
	}

	// Calculate final GPS coordinates
	result := 0
	for i, row := range board {
		for j, cell := range row {
			if cell == '[' {
				result += i*100 + j
			}
		}
	}

	return result
}

// Helper function to find index of byte in slice
func indexOf(slice []byte, target byte) int {
	for i, b := range slice {
		if b == target {
			return i
		}
	}
	return -1
}

func main() {
	input, err := utility.ParseTextFile("input")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(part1(input))
	log.Println(part2(input))
}
