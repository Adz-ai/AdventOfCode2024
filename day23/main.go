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
	state := newCliqueState(g)

	p := make([]bool, len(state.computers))
	x := make([]bool, len(state.computers))
	skip := make([]bool, len(state.computers))
	for i := range p {
		p[i] = true
	}

	state.bronKerbosch([]string{}, p, x, skip)
	return state.maxClique
}

// findTripleConnections finds all sets of three interconnected computers
func (g Graph) findTripleConnections(computers []string) [][]string {
	var result [][]string

	for i := 0; i < len(computers); i++ {
		for j := i + 1; j < len(computers); j++ {
			for k := j + 1; k < len(computers); k++ {
				set := []string{computers[i], computers[j], computers[k]}
				if g.areConnected(set) {
					result = append(result, set)
				}
			}
		}
	}
	return result
}

// CliqueState holds the state for clique finding
type CliqueState struct {
	computers     []string
	computerIndex map[string]int
	maxClique     []string
	currentMax    int
	graph         Graph
}

// newCliqueState initializes a new CliqueState
func newCliqueState(g Graph) *CliqueState {
	computers := make([]string, 0, len(g))
	for computer := range g {
		computers = append(computers, computer)
	}

	// Sort by degree
	sort.Slice(computers, func(i, j int) bool {
		return g.getDegree(computers[i]) > g.getDegree(computers[j])
	})

	computerIndex := make(map[string]int, len(computers))
	for i, c := range computers {
		computerIndex[c] = i
	}

	return &CliqueState{
		computers:     computers,
		computerIndex: computerIndex,
		maxClique:     make([]string, 0),
		currentMax:    0,
		graph:         g,
	}
}

// findPivot finds the pivot vertex for the Bron-Kerbosch algorithm
func (s *CliqueState) findPivot(p, x []bool) int {
	pivot := -1
	maxCount := -1

	for i := range s.computers {
		if !p[i] && !x[i] {
			continue
		}
		count := s.countConnections(i, p)
		if count > maxCount {
			maxCount = count
			pivot = i
		}
	}
	return pivot
}

// countConnections counts connections from vertex to vertices in set p
func (s *CliqueState) countConnections(vertex int, p []bool) int {
	count := 0
	for j := range s.computers {
		if p[j] && s.graph[s.computers[vertex]][s.computers[j]] {
			count++
		}
	}
	return count
}

// updateMaxClique updates the maximum clique if a larger one is found
func (s *CliqueState) updateMaxClique(r []string) {
	if len(r) > s.currentMax {
		s.maxClique = make([]string, len(r))
		copy(s.maxClique, r)
		s.currentMax = len(r)
	}
}

// bronKerboschStep performs one step of the Bron-Kerbosch algorithm
func (s *CliqueState) bronKerboschStep(r []string, p, x []bool, v int) ([]string, []bool, []bool) {
	newR := append([]string(nil), r...)
	newR = append(newR, s.computers[v])

	newP := make([]bool, len(s.computers))
	newX := make([]bool, len(s.computers))

	for i := range s.computers {
		if p[i] && s.graph[s.computers[v]][s.computers[i]] {
			newP[i] = true
		}
		if x[i] && s.graph[s.computers[v]][s.computers[i]] {
			newX[i] = true
		}
	}

	return newR, newP, newX
}

// bronKerbosch implements the Bron-Kerbosch algorithm with pivoting
func (s *CliqueState) bronKerbosch(r []string, p, x, skip []bool) {
	if allFalse(p) && allFalse(x) {
		s.updateMaxClique(r)
		return
	}

	pivot := s.findPivot(p, x)

	for v := range s.computers {
		if !p[v] || (pivot != -1 && s.graph[s.computers[pivot]][s.computers[v]]) {
			continue
		}
		if skip[v] {
			continue
		}

		newR, newP, newX := s.bronKerboschStep(r, p, x, v)
		s.bronKerbosch(newR, newP, newX, skip)

		p[v] = false
		x[v] = true
	}
}

// hasComputerStartingWith checks if any computer in the set starts with the given prefix
func hasComputerStartingWith(computers []string, prefix string) bool {
	for _, comp := range computers {
		if strings.HasPrefix(comp, prefix) {
			return true
		}
	}
	return false
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

	triples := graph.findTripleConnections(computers)

	count := 0
	for _, triple := range triples {
		if hasComputerStartingWith(triple, "t") {
			count++
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
