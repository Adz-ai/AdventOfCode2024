package main

import (
	"aoc2024/utility"
	"log"
	"strconv"
)

// File represents a file on the disk with its metadata
type File struct {
	id     int
	start  int
	length int
}

// Space represents a continuous free space region
type Space struct {
	start  int
	length int
}

// parseInput takes a string input representing alternating file and space lengths
// and returns slices of File and Space structs along with the total disk length
// Input format: "FSFSFS" where F is file length and S is space length
func parseInput(input string) ([]File, []Space, int) {
	var files []File
	var spaces []Space
	pos := 0
	fileID := 0
	totalLength := 0

	// Pre-calculate total length and allocate slices with capacity
	for i := 0; i < len(input); i++ {
		length, _ := strconv.Atoi(string(input[i]))
		totalLength += length
	}

	files = make([]File, 0, len(input)/2)
	spaces = make([]Space, 0, len(input)/2)

	// Parse input
	for i := 0; i < len(input); i++ {
		length, _ := strconv.Atoi(string(input[i]))
		if i%2 == 0 {
			// File
			files = append(files, File{
				id:     fileID,
				start:  pos,
				length: length,
			})
			fileID++
		} else {
			// Space
			spaces = append(spaces, Space{
				start:  pos,
				length: length,
			})
		}
		pos += length
	}

	return files, spaces, totalLength
}

// initializeDisk sets up the initial state of the disk array based on
// the provided files and spaces. -1 represents free space
func initializeDisk(disk []int, files []File, spaces []Space) {
	// Initialize all positions to -1 (free space)
	for i := range disk {
		disk[i] = -1
	}

	// Place files
	for _, file := range files {
		for i := 0; i < file.length; i++ {
			disk[file.start+i] = file.id
		}
	}
}

// findContinuousSpaces scans the disk and returns a slice of Space structs
// representing all continuous regions of free space
func findContinuousSpaces(disk []int) []Space {
	var spaces []Space
	start := -1
	length := 0

	for i, val := range disk {
		if val == -1 {
			if start == -1 {
				start = i
			}
			length++
		} else if start != -1 {
			spaces = append(spaces, Space{start: start, length: length})
			start = -1
			length = 0
		}
	}

	if start != -1 {
		spaces = append(spaces, Space{start: start, length: length})
	}

	return spaces
}

// moveFileToSpace moves a file from its current location to a new position
// on the disk, updating both the old and new locations
func moveFileToSpace(disk []int, file File, newStart int) {
	// Clear old location
	fileID := disk[file.start]
	for i := 0; i < file.length; i++ {
		disk[file.start+i] = -1
	}
	// Write to new location
	for i := 0; i < file.length; i++ {
		disk[newStart+i] = fileID
	}
}

// updateFreeSpaces updates the list of free spaces when a file cannot be moved,
// potentially splitting spaces that overlap with the unmoved file
func updateFreeSpaces(spaces []Space, file File) {
	// Update free spaces list to account for unmoved file
	for i := range spaces {
		if spaces[i].start > file.start {
			break
		}
		if spaces[i].start+spaces[i].length > file.start {
			// Split the space if necessary
			if spaces[i].start < file.start {
				newLength := file.start - spaces[i].start
				if spaces[i].length > newLength+file.length {
					// Space continues after the file
					spaces = append(spaces, Space{
						start:  file.start + file.length,
						length: spaces[i].length - newLength - file.length,
					})
				}
				spaces[i].length = newLength
			}
		}
	}
}

// calculateChecksum computes the final checksum of the disk state by multiplying
// each file block's position by its file ID and summing the results
func calculateChecksum(disk []int) int {
	checksum := 0
	for pos, fileID := range disk {
		if fileID != -1 {
			checksum += pos * fileID
		}
	}
	return checksum
}

// part1 solves the first part of the puzzle where individual blocks
// are moved from right to left to the first available space
// Returns the checksum of the final disk state
func part1(input []string) int {
	// Parse input into files and initialize disk
	files, spaces, totalLength := parseInput(input[0])
	disk := make([]int, totalLength)
	initializeDisk(disk, files, spaces)

	// Process blocks right to left, moving each block to the leftmost available space
	for i := len(disk) - 1; i >= 0; i-- {
		if disk[i] == -1 { // Skip free space
			continue
		}

		// Find leftmost free space before current position
		for j := 0; j < i; j++ {
			if disk[j] == -1 {
				// Move block
				disk[j] = disk[i]
				disk[i] = -1
				break
			}
		}
	}

	return calculateChecksum(disk)
}

// part2 solves the second part of the puzzle where entire files
// are moved from right to left to the first available space that can fit them
// Returns the checksum of the final disk state
func part2(input []string) int {
	files, spaces, totalLength := parseInput(input[0])
	disk := make([]int, totalLength)
	initializeDisk(disk, files, spaces)

	// Create a slice of continuous free spaces
	freeSpaces := findContinuousSpaces(disk)

	// Process files in descending order of ID
	for i := len(files) - 1; i >= 0; i-- {
		file := files[i]
		moved := false

		// Find the leftmost suitable space for this file
		for j, space := range freeSpaces {
			if space.length >= file.length && space.start < file.start {
				// Move file to this space
				moveFileToSpace(disk, file, space.start)

				// Update the free space
				newSpace := Space{
					start:  space.start + file.length,
					length: space.length - file.length,
				}
				if newSpace.length > 0 {
					freeSpaces[j] = newSpace
				} else {
					// Remove this space if fully used
					freeSpaces = append(freeSpaces[:j], freeSpaces[j+1:]...)
				}

				moved = true
				break
			}
		}

		if !moved {
			// If file couldn't move, update free spaces list to exclude its current position
			updateFreeSpaces(freeSpaces, file)
		}
	}

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
