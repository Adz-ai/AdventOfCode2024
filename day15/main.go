package main

import (
	"aoc2024/utility"
	"errors"
	"log"
)

// Board represents the warehouse grid state
type Board struct {
	grid     [][]byte
	robotRow int
	robotCol int
}

// Direction represents a movement direction with its row and column deltas
type Direction struct {
	dr, dc int
}

var (
	// Movement direction mappings
	directions = map[byte]Direction{
		'^': {dr: -1, dc: 0},
		'v': {dr: 1, dc: 0},
		'<': {dr: 0, dc: -1},
		'>': {dr: 0, dc: 1},
	}

	ErrNoRobot       = errors.New("no robot found in input")
	ErrInvalidFormat = errors.New("invalid input format: no empty line separator found")
)

// ParseInput parses the input lines and returns the board and moves
func ParseInput(input []string) ([][]byte, []byte, error) {
	var board [][]byte
	var moves []byte
	var movesIdx int
	var found bool

	// Find separator between board and moves
	for i, line := range input {
		if len(line) == 0 {
			movesIdx = i
			found = true
			break
		}
	}

	if !found {
		return nil, nil, ErrInvalidFormat
	}

	// Parse board
	board = make([][]byte, movesIdx)
	for i := 0; i < movesIdx; i++ {
		board[i] = []byte(input[i])
	}

	// Parse moves
	for _, line := range input[movesIdx+1:] {
		moves = append(moves, []byte(line)...)
	}

	return board, moves, nil
}

// findRobot locates the robot position in the board
func findRobot(board [][]byte) (int, int, error) {
	for i, row := range board {
		for j, cell := range row {
			if cell == '@' {
				return i, j, nil
			}
		}
	}
	return 0, 0, ErrNoRobot
}

// createExpandedBoard creates a board with doubled width for part 2
func createExpandedBoard(input []string, movesIdx int) *Board {
	board := &Board{}
	for i := 0; i < movesIdx; i++ {
		var newLine []byte
		for j, ch := range input[i] {
			switch ch {
			case '#':
				newLine = append(newLine, '#', '#')
			case 'O':
				newLine = append(newLine, '[', ']')
			case '.':
				newLine = append(newLine, '.', '.')
			case '@':
				newLine = append(newLine, '@', '.')
				board.robotRow, board.robotCol = len(board.grid), 2*j
			}
		}
		board.grid = append(board.grid, newLine)
	}
	return board
}

// isValidPosition checks if a position is within board bounds
func isValidPosition(board [][]byte, r, c int) bool {
	return r >= 0 && r < len(board) && c >= 0 && c < len(board[0])
}

// canMovePart1 checks if movement is possible in part 1
func canMovePart1(board [][]byte, r, c int, dir Direction) bool {
	nextR, nextC := r+dir.dr, c+dir.dc

	if !isValidPosition(board, nextR, nextC) {
		return false
	}

	if board[nextR][nextC] == 'O' {
		return canMovePart1(board, nextR, nextC, dir)
	}

	return board[nextR][nextC] == '.'
}

// canMovePart2 checks if movement is possible in part 2
func canMovePart2(board [][]byte, r, c int, dir Direction) bool {
	nextR, nextC := r+dir.dr, c+dir.dc

	if !isValidPosition(board, nextR, nextC) {
		return false
	}

	switch board[nextR][nextC] {
	case ']':
		if !canMovePart2(board, nextR, nextC, dir) {
			return false
		}
		// Handle vertical movement for wide boxes
		if dir.dr != 0 {
			return canMovePart2(board, nextR, nextC-1, dir)
		}
		return true

	case '[':
		if !canMovePart2(board, nextR, nextC, dir) {
			return false
		}
		// Handle vertical movement for wide boxes
		if dir.dr != 0 {
			return canMovePart2(board, nextR, nextC+1, dir)
		}
		return true
	}

	return board[nextR][nextC] == '.'
}

// movePart1 executes a move in part 1
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

// movePart2 executes a move in part 2
func movePart2(board [][]byte, r, c int, dir Direction) (int, int) {
	if !canMovePart2(board, r, c, dir) {
		return r, c
	}

	nextR, nextC := r+dir.dr, c+dir.dc

	switch board[nextR][nextC] {
	case ']':
		movePart2(board, nextR, nextC, dir)
		if dir.dr != 0 {
			movePart2(board, nextR, nextC-1, dir)
		}
	case '[':
		movePart2(board, nextR, nextC, dir)
		if dir.dr != 0 {
			movePart2(board, nextR, nextC+1, dir)
		}
	}

	board[nextR][nextC] = board[r][c]
	board[r][c] = '.'

	return nextR, nextC
}

// calculateGPS calculates the sum of GPS coordinates
func calculateGPS(board [][]byte, target byte) int {
	result := 0
	for i, row := range board {
		for j, cell := range row {
			if cell == target {
				result += i*100 + j
			}
		}
	}
	return result
}

func part1(input []string) int {
	// Parse input
	board, moves, err := ParseInput(input)
	if err != nil {
		log.Fatal(err)
	}

	// Find robot position
	robotRow, robotCol, err := findRobot(board)
	if err != nil {
		log.Fatal(err)
	}

	// Process moves
	for _, m := range moves {
		dir := directions[m]
		robotRow, robotCol = movePart1(board, robotRow, robotCol, dir)
	}

	return calculateGPS(board, 'O')
}

func part2(input []string) int {
	// Find separator
	movesIdx := -1
	for i, line := range input {
		if len(line) == 0 {
			movesIdx = i
			break
		}
	}
	if movesIdx == -1 {
		log.Fatal("invalid input format")
	}

	// Create expanded board
	board := createExpandedBoard(input, movesIdx)
	if err != nil {
		log.Fatal(err)
	}

	// Parse moves
	var moves []byte
	for _, line := range input[movesIdx+1:] {
		moves = append(moves, []byte(line)...)
	}

	// Process moves
	for _, m := range moves {
		dir := directions[m]
		board.robotRow, board.robotCol = movePart2(board.grid, board.robotRow, board.robotCol, dir)
	}

	return calculateGPS(board.grid, '[')
}

func main() {
	input, err := utility.ParseTextFile("input")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(part1(input))
	log.Println(part2(input))
}
