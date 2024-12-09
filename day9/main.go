package main

import (
	"aoc2024/utility"
	"log"
	"strconv"
)

type File struct {
	id       int
	start    int
	length   int
	original bool
}

func parseInput(input string) ([]int, []File) {
	var disk []int
	var files []File
	pos := 0
	fileID := 0

	// Parse alternating file and space lengths
	for i := 0; i < len(input); i++ {
		length, _ := strconv.Atoi(string(input[i]))
		isSpace := i%2 == 1

		if !isSpace {
			// Record file information
			files = append(files, File{
				id:       fileID,
				start:    pos,
				length:   length,
				original: true,
			})
			// Fill disk with file blocks
			for j := 0; j < length; j++ {
				disk = append(disk, fileID)
			}
			fileID++
		} else {
			// Fill disk with space blocks
			for j := 0; j < length; j++ {
				disk = append(disk, -1)
			}
		}
		pos += length
	}

	return disk, files
}

func compactBlockByBlock(disk []int) {
	for i := len(disk) - 1; i >= 0; i-- {
		if disk[i] != -1 { // If it's a file block
			// Find leftmost space
			spacePos := -1
			for j := 0; j < i; j++ {
				if disk[j] == -1 {
					spacePos = j
					break
				}
			}
			// Move block if space found
			if spacePos != -1 {
				disk[spacePos] = disk[i]
				disk[i] = -1
			}
		}
	}
}

func compactWholeFiles(disk []int, files []File) {
	// Sort files by ID in descending order
	for i := len(files) - 1; i >= 0; i-- {
		file := files[i]
		// Find leftmost space that can fit the whole file
		bestPos := findBestFitPosition(disk, file.length, file.start)
		if bestPos != -1 && bestPos < file.start {
			// Move the whole file
			moveFile(disk, file, bestPos)
		}
	}
}

func findBestFitPosition(disk []int, fileLength, currentStart int) int {
	for pos := 0; pos < currentStart; pos++ {
		if hasEnoughSpace(disk, pos, fileLength) {
			return pos
		}
	}
	return -1
}

func hasEnoughSpace(disk []int, start, length int) bool {
	if start+length > len(disk) {
		return false
	}
	for i := 0; i < length; i++ {
		if disk[start+i] != -1 {
			return false
		}
	}
	return true
}

func moveFile(disk []int, file File, newStart int) {
	// Get file ID from disk
	fileID := disk[file.start]
	// Clear old location
	for i := 0; i < file.length; i++ {
		disk[file.start+i] = -1
	}
	// Write to new location
	for i := 0; i < file.length; i++ {
		disk[newStart+i] = fileID
	}
}

func calculateChecksum(disk []int) int {
	checksum := 0
	for pos, fileID := range disk {
		if fileID != -1 {
			checksum += pos * fileID
		}
	}
	return checksum
}

func part1(input []string) int {
	disk, _ := parseInput(input[0])
	compactBlockByBlock(disk)
	return calculateChecksum(disk)
}

func part2(input []string) int {
	disk, files := parseInput(input[0])
	compactWholeFiles(disk, files)
	return calculateChecksum(disk)
}

func main() {
	input, err := utility.ParseTextFile("input")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(part1(input))
	log.Println(part2(input))
}
