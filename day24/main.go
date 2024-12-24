package main

import (
	"aoc2024/utility"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
)

// Gate represents a logic gate with its operation and connections
type Gate struct {
	op          string
	input1      string
	input2      string
	output      string
	hasComputed bool
}

// Circuit represents the entire circuit system
type Circuit struct {
	gates      map[string]*Gate    // Maps output wire to gate
	wireValues map[string]int      // Current values of wires (using ints instead of bools)
	computed   map[string]struct{} // Tracks which wires have been computed
}

func NewCircuit() *Circuit {
	return &Circuit{
		gates:      make(map[string]*Gate),
		wireValues: make(map[string]int),
		computed:   make(map[string]struct{}),
	}
}

// ParseInput parses the input lines and initializes the circuit
func (c *Circuit) ParseInput(lines []string) error {
	parsingInitialValues := true

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			parsingInitialValues = false
			continue
		}

		if parsingInitialValues {
			parts := strings.Split(line, ":")
			if len(parts) != 2 {
				return fmt.Errorf("invalid initial value format: %s", line)
			}
			wire := strings.TrimSpace(parts[0])
			value, err := strconv.Atoi(strings.TrimSpace(parts[1]))
			if err != nil {
				return fmt.Errorf("invalid wire value: %v", err)
			}
			c.wireValues[wire] = value
			c.computed[wire] = struct{}{}
		} else {
			if err := c.parseGate(line); err != nil {
				return fmt.Errorf("error parsing gate: %v", err)
			}
		}
	}
	return nil
}

// parseGate parses a single gate definition line
func (c *Circuit) parseGate(line string) error {
	parts := strings.Split(line, "->")
	if len(parts) != 2 {
		return fmt.Errorf("invalid gate format: %s", line)
	}

	inputs := strings.Fields(strings.TrimSpace(parts[0]))
	output := strings.TrimSpace(parts[1])

	if len(inputs) != 3 {
		return fmt.Errorf("invalid number of inputs: %s", line)
	}

	gate := &Gate{
		op:          inputs[1],
		input1:      inputs[0],
		input2:      inputs[2],
		output:      output,
		hasComputed: false,
	}

	c.gates[output] = gate
	return nil
}

// computeGate evaluates a single gate's output using integer operations
func (c *Circuit) computeGate(gate *Gate) (bool, error) {
	_, input1Ready := c.computed[gate.input1]
	_, input2Ready := c.computed[gate.input2]

	if !input1Ready || !input2Ready {
		return false, nil
	}

	input1Val := c.wireValues[gate.input1]
	input2Val := c.wireValues[gate.input2]

	var result int
	switch gate.op {
	case "AND":
		result = input1Val & input2Val
	case "OR":
		result = input1Val | input2Val
	case "XOR":
		result = input1Val ^ input2Val
	default:
		return false, fmt.Errorf("unknown gate operation: %s", gate.op)
	}

	c.wireValues[gate.output] = result
	c.computed[gate.output] = struct{}{}
	gate.hasComputed = true
	return true, nil
}

// Simulate runs the circuit simulation until all gates have computed their outputs
func (c *Circuit) Simulate() error {
	for {
		progress := false
		for _, gate := range c.gates {
			if gate.hasComputed {
				continue
			}

			computed, err := c.computeGate(gate)
			if err != nil {
				return err
			}
			if computed {
				progress = true
			}
		}

		if !progress {
			allComputed := true
			for _, gate := range c.gates {
				if !gate.hasComputed {
					allComputed = false
					break
				}
			}
			if allComputed {
				return nil
			}
			return fmt.Errorf("circuit simulation stuck - possible circular dependency")
		}
	}
}

// getBinaryNumber gets a binary number from wires starting with prefix
func (c *Circuit) getBinaryNumber(prefix string) int {
	var wires []string
	for wire := range c.wireValues {
		if strings.HasPrefix(wire, prefix) {
			wires = append(wires, wire)
		}
	}

	// Sort wires in descending order
	sort.Slice(wires, func(i, j int) bool {
		return wires[i] > wires[j]
	})

	var bits strings.Builder
	for _, wire := range wires {
		bits.WriteString(strconv.Itoa(c.wireValues[wire]))
	}

	result, _ := strconv.ParseInt(bits.String(), 2, 64)
	return int(result)
}

// validateAddition checks if z = x + y for current wire values
func (c *Circuit) validateAddition() bool {
	x := c.getBinaryNumber("x")
	y := c.getBinaryNumber("y")
	z := c.getBinaryNumber("z")
	return x+y == z
}

// trySwapGates attempts to swap output wires of two gates
func (c *Circuit) trySwapGates(g1, g2 *Gate) {
	// Save original outputs
	g1.output, g2.output = g2.output, g1.output
}

// findGateOutput finds a gate's output wire given inputs and operation
func (c *Circuit) findGateOutput(input1, input2, operation string) string {
	for _, gate := range c.gates {
		if gate.op == operation &&
			((gate.input1 == input1 && gate.input2 == input2) ||
				(gate.input1 == input2 && gate.input2 == input1)) {
			return gate.output
		}
	}
	return ""
}

// findSwappedGates analyzes the circuit as a binary adder to find swapped gates
func (c *Circuit) findSwappedGates() ([]string, error) {
	var swappedWires []string
	var carryWire string

	// Track which wires we've already swapped
	swapped := make(map[string]bool)

	for i := 0; i < 45; i++ {
		bit := fmt.Sprintf("%02d", i)
		x := "x" + bit
		y := "y" + bit

		// Find initial XOR and AND gates (half adder components)
		xorWire := c.findGateOutput(x, y, "XOR")
		andWire := c.findGateOutput(x, y, "AND")

		if carryWire != "" {
			// Full adder
			carryAndWire := c.findGateOutput(carryWire, xorWire, "AND")
			sumWire := c.findGateOutput(carryWire, xorWire, "XOR")

			// If we can't find the expected gate, the XOR and AND might be swapped
			if carryAndWire == "" {
				if !swapped[xorWire] && !swapped[andWire] {
					swappedWires = append(swappedWires, xorWire, andWire)
					swapped[xorWire] = true
					swapped[andWire] = true
				}
				// Update our understanding for further checks
				xorWire, andWire = andWire, xorWire
				carryAndWire = c.findGateOutput(carryWire, xorWire, "AND")
			}

			// Check if sum wire is involved in a swap
			if strings.HasPrefix(sumWire, "z") && sumWire != fmt.Sprintf("z%02d", i) {
				// Find the correct z-wire for this position
				correctZWire := fmt.Sprintf("z%02d", i)
				if !swapped[sumWire] && !swapped[correctZWire] {
					swappedWires = append(swappedWires, sumWire, correctZWire)
					swapped[sumWire] = true
					swapped[correctZWire] = true
				}
			}

			// Find next carry
			nextCarry := c.findGateOutput(carryAndWire, andWire, "OR")

			// Check if carry is swapped with a z-wire
			if strings.HasPrefix(nextCarry, "z") && nextCarry != "z45" {
				if !swapped[nextCarry] && !swapped[sumWire] {
					swappedWires = append(swappedWires, nextCarry, sumWire)
					swapped[nextCarry] = true
					swapped[sumWire] = true
				}
			}

			carryWire = nextCarry
		} else {
			// First bit (half adder)
			carryWire = andWire
		}
	}

	// Filter and sort unique wires
	uniqueWires := make([]string, 0, 8)
	seen := make(map[string]bool)
	for i := 0; i < len(swappedWires); i += 2 {
		wire1, wire2 := swappedWires[i], swappedWires[i+1]
		if !seen[wire1] && !seen[wire2] {
			uniqueWires = append(uniqueWires, wire1, wire2)
			seen[wire1] = true
			seen[wire2] = true
			if len(uniqueWires) == 8 {
				break
			}
		}
	}

	if len(uniqueWires) != 8 {
		return nil, fmt.Errorf("expected 8 swapped wires, found %d", len(uniqueWires))
	}

	sort.Strings(uniqueWires)
	return uniqueWires, nil
}

func solve(input []string) (int, string) {
	circuit := NewCircuit()

	if err := circuit.ParseInput(input); err != nil {
		log.Fatal(err)
	}

	if err := circuit.Simulate(); err != nil {
		log.Fatal(err)
	}
	part1 := circuit.getBinaryNumber("z")

	swappedWires, err := circuit.findSwappedGates()
	if err != nil {
		log.Fatal(err)
	}

	part2 := strings.Join(swappedWires, "")
	return part1, part2
}

func main() {
	input, err := utility.ParseTextFile("input")
	if err != nil {
		log.Fatal(err)
	}

	part1, part2 := solve(input)

	log.Printf("Part 1: %d\n", part1)
	log.Printf("Part 2: %s\n", part2)
}
