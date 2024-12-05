package main

import (
	"aoc2024/utility"
	"fmt"
	"log"
	"strings"
)

// Rule represents a directed edge in the dependency graph where 'before' must precede 'after'
type Rule struct {
	before, after int
}

// Node represents a vertex in the dependency graph for topological sorting.
// Each node tracks its incoming edges (inDegree) and outgoing connections (outNodes).
type Node struct {
	page     int
	inDegree int
	outNodes map[int]bool
}

// parseInput processes the input file and returns two slices:
// - rules: collection of page ordering rules
// - updates: groups of pages that need to be printed
// The function handles malformed input by skipping invalid entries.
func parseInput(lines []string) ([]Rule, [][]int) {
	var rules []Rule
	var updates [][]int

	i := 0
	for ; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			i++
			break
		}

		parts := strings.Split(line, "|")
		if len(parts) != 2 {
			continue
		}

		var before, after int
		if _, err := fmt.Sscanf(parts[0], "%d", &before); err != nil {
			continue
		}
		if _, err := fmt.Sscanf(parts[1], "%d", &after); err != nil {
			continue
		}
		rules = append(rules, Rule{before, after})
	}

	for ; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		var update []int
		for _, numStr := range strings.Split(line, ",") {
			var num int
			if _, err := fmt.Sscanf(numStr, "%d", &num); err != nil {
				continue
			}
			update = append(update, num)
		}
		if len(update) > 0 {
			updates = append(updates, update)
		}
	}

	return rules, updates
}

// isValidOrder checks if a sequence of pages satisfies all applicable ordering rules.
// Time complexity: O(R * N) where R is number of rules and N is length of update.
// Space complexity: O(N) for the positions map.
func isValidOrder(update []int, rules []Rule) bool {
	positions := make(map[int]int)
	for i, page := range update {
		positions[page] = i
	}

	for _, rule := range rules {
		beforePos, beforeExists := positions[rule.before]
		afterPos, afterExists := positions[rule.after]
		if beforeExists && afterExists {
			if beforePos >= afterPos {
				return false
			}
		}
	}

	return true
}

// getMiddlePage returns the middle element from a slice of integers.
// For odd-length slices, returns the exact middle.
// For even-length slices, returns the lower middle element.
func getMiddlePage(update []int) int {
	return update[len(update)/2]
}

// buildGraph constructs a directed acyclic graph (DAG) from the given pages and rules.
// The resulting graph represents the dependencies between pages where edges indicate
// ordering requirements. Only rules applicable to the given update are included.
// Time complexity: O(P + R) where P is number of pages and R is number of rules.
// Space complexity: O(P + R) for storing nodes and edges.
func buildGraph(update []int, rules []Rule) map[int]*Node {
	pageSet := make(map[int]bool)
	for _, page := range update {
		pageSet[page] = true
	}

	graph := make(map[int]*Node)
	for page := range pageSet {
		graph[page] = &Node{
			page:     page,
			inDegree: 0,
			outNodes: make(map[int]bool),
		}
	}

	for _, rule := range rules {
		if pageSet[rule.before] && pageSet[rule.after] {
			graph[rule.before].outNodes[rule.after] = true
			graph[rule.after].inDegree++
		}
	}

	return graph
}

// topologicalSort implements Kahn's algorithm to find a valid ordering of pages.
// Kahn's algorithm works by repeatedly removing nodes with no incoming edges and
// adding them to the result. This process continues until all nodes are processed
// or a cycle is detected.
//
// Algorithm steps:
// 1. Find all nodes with no incoming edges (inDegree = 0)
// 2. While such nodes exist:
//   - Remove a node and add it to result
//   - Decrease inDegree for all its neighbors
//   - Add any neighbors with new inDegree of 0 to queue
//
// Time complexity: O(V + E) where V is number of pages and E is number of rules
// Space complexity: O(V) for queue and result
func topologicalSort(pages []int, rules []Rule) []int {
	graph := buildGraph(pages, rules)
	var result []int
	var queue []int

	for page, node := range graph {
		if node.inDegree == 0 {
			queue = append(queue, page)
		}
	}

	for len(queue) > 0 {
		page := queue[0]
		queue = queue[1:]
		result = append(result, page)

		for nextPage := range graph[page].outNodes {
			graph[nextPage].inDegree--
			if graph[nextPage].inDegree == 0 {
				queue = append(queue, nextPage)
			}
		}
	}

	return result
}

func part2(input []string) int {
	rules, updates := parseInput(input)
	sum := 0

	for _, update := range updates {
		if !isValidOrder(update, rules) {
			correctOrder := topologicalSort(update, rules)
			sum += getMiddlePage(correctOrder)
		}
	}
	return sum
}

func part1(input []string) int {
	rules, updates := parseInput(input)
	sum := 0

	for _, update := range updates {
		if isValidOrder(update, rules) {
			middlePage := getMiddlePage(update)
			sum += middlePage
		}
	}

	return sum
}

func main() {
	input, err := utility.ParseTextFile("input")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(part1(input))
	log.Println(part2(input))
}
