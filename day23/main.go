package main

import (
	"aoc2024/utility"
	"log"
	"sort"
	"strings"
)

// Graph represents connections between computers
type Graph map[string]map[string]bool

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

// getDegree returns the number of connections a computer has
func (g Graph) getDegree(computer string) int {
	return len(g[computer])
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

// findMaxClique returns the largest set of fully connected computers
func (g Graph) findMaxClique() []string {
	// Convert map keys to slice and sort by degree (descending)
	computers := make([]string, 0, len(g))
	for computer := range g {
		computers = append(computers, computer)
	}
	sort.Slice(computers, func(i, j int) bool {
		return g.getDegree(computers[i]) > g.getDegree(computers[j])
	})

	// Initialize maximum clique tracking
	maxClique := make([]string, 0)
	currentMax := 0

	// Use bit arrays for faster set operations
	candidates := make([]bool, len(computers))
	for i := range candidates {
		candidates[i] = true
	}

	// Helper function to get computer index
	computerIndex := make(map[string]int, len(computers))
	for i, c := range computers {
		computerIndex[c] = i
	}

	var bronKerbosch func([]string, []bool, []bool, []bool)
	bronKerbosch = func(r []string, p, x, skip []bool) {
		// If no candidates and nothing to exclude
		if allFalse(p) && allFalse(x) {
			if len(r) > currentMax {
				maxClique = make([]string, len(r))
				copy(maxClique, r)
				currentMax = len(r)
			}
			return
		}

		// Find pivot
		pivot := -1
		maxCount := -1
		for i := range computers {
			if !p[i] && !x[i] {
				continue
			}
			count := 0
			for j := range computers {
				if p[j] && g[computers[i]][computers[j]] {
					count++
				}
			}
			if count > maxCount {
				maxCount = count
				pivot = i
			}
		}

		// For each vertex not connected to pivot
		for v := range computers {
			if !p[v] || (pivot != -1 && g[computers[pivot]][computers[v]]) {
				continue
			}
			if skip[v] {
				continue
			}

			// Create new sets
			newR := append([]string(nil), r...)
			newR = append(newR, computers[v])

			newP := make([]bool, len(computers))
			newX := make([]bool, len(computers))

			// Add neighbors of v that are in p to new p and x sets
			for i := range computers {
				if p[i] && g[computers[v]][computers[i]] {
					newP[i] = true
				}
				if x[i] && g[computers[v]][computers[i]] {
					newX[i] = true
				}
			}

			bronKerbosch(newR, newP, newX, skip)

			p[v] = false
			x[v] = true
		}
	}

	// Initialize starting sets
	p := make([]bool, len(computers))
	x := make([]bool, len(computers))
	skip := make([]bool, len(computers))
	for i := range p {
		p[i] = true
	}

	// Start recursion
	bronKerbosch([]string{}, p, x, skip)

	return maxClique
}

// allFalse returns true if all values in the slice are false
func allFalse(arr []bool) bool {
	for _, v := range arr {
		if v {
			return false
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
