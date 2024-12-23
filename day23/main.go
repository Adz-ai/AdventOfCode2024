package main

import (
	"aoc2024/utility"
	"log"
	"sort"
	"strings"
)

// Graph represents connections between computers
type Graph map[string]map[string]bool

// findMaxClique returns the largest set of fully connected computers
func (g Graph) findMaxClique() []string {
	var maxClique []string
	computers := make([]string, 0)
	for computer := range g {
		computers = append(computers, computer)
	}

	// Try all possible combinations of computers, starting from the largest possible size
	for size := len(computers); size > 0; size-- {
		// If we already found a clique, we don't need to check smaller sizes
		if len(maxClique) > 0 {
			break
		}

		var checkClique func([]string, int, int)
		checkClique = func(current []string, start int, remaining int) {
			// If we already found a clique, stop searching
			if len(maxClique) > 0 {
				return
			}

			// If we have enough computers in our current set
			if remaining == 0 {
				if g.areConnected(current) {
					maxClique = make([]string, len(current))
					copy(maxClique, current)
				}
				return
			}

			// If we can't possibly get enough computers, stop this branch
			if len(computers)-start < remaining {
				return
			}

			// Try adding each remaining computer
			for i := start; i < len(computers); i++ {
				// Check if this computer could be added to our current set
				canAdd := true
				for _, c := range current {
					if !g[c][computers[i]] {
						canAdd = false
						break
					}
				}
				if canAdd {
					current = append(current, computers[i])
					checkClique(current, i+1, remaining-1)
					current = current[:len(current)-1]
				}
			}
		}

		checkClique(make([]string, 0), 0, size)
	}

	return maxClique
}

// addConnection adds a bidirectional connection between two computers
func (g Graph) addConnection(c1, c2 string) {
	if _, exists := g[c1]; !exists {
		g[c1] = make(map[string]bool)
	}
	if _, exists := g[c2]; !exists {
		g[c2] = make(map[string]bool)
	}
	g[c1][c2] = true
	g[c2][c1] = true
}

// areConnected checks if all computers in the set are interconnected
func (g Graph) areConnected(computers []string) bool {
	for i := 0; i < len(computers); i++ {
		for j := i + 1; j < len(computers); j++ {
			if !g[computers[i]][computers[j]] {
				return false
			}
		}
	}
	return true
}

func part1(input []string) int {
	graph := make(Graph)

	for _, line := range input {
		connection := strings.Split(line, "-")
		if len(connection) != 2 {
			continue
		}
		graph.addConnection(connection[0], connection[1])
	}

	computers := make([]string, 0)
	for computer := range graph {
		computers = append(computers, computer)
	}

	count := 0
	for i := 0; i < len(computers); i++ {
		for j := i + 1; j < len(computers); j++ {
			for k := j + 1; k < len(computers); k++ {
				set := []string{computers[i], computers[j], computers[k]}
				if graph.areConnected(set) {
					hasT := false
					for _, comp := range set {
						if strings.HasPrefix(comp, "t") {
							hasT = true
							break
						}
					}
					if hasT {
						count++
					}
				}
			}
		}
	}

	return count
}

func part2(input []string) string {
	graph := make(Graph)
	for _, line := range input {
		connection := strings.Split(line, "-")
		if len(connection) != 2 {
			continue
		}
		graph.addConnection(connection[0], connection[1])
	}

	maxClique := graph.findMaxClique()

	sort.Strings(maxClique)

	return strings.Join(maxClique, ",")
}

func main() {
	input, err := utility.ParseTextFile("input")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(part1(input))
	log.Println(part2(input))
}
