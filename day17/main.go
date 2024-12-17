package main

import (
	"aoc2024/utility"
	"fmt"
	"log"
	"slices"
	"strings"
)

func parseInput(input []string) ([]uint64, uint64) {
	// Parse registers
	var seed uint64
	fmt.Sscanf(input[0], "Register A: %d", &seed)

	// Parse program
	programStr := strings.TrimPrefix(input[4], "Program: ")
	parts := strings.Split(programStr, ",")
	program := make([]uint64, len(parts))
	for i, part := range parts {
		var num uint64
		fmt.Sscanf(strings.TrimSpace(part), "%d", &num)
		program[i] = num
	}
	return program, seed
}

func partOne(program []uint64, seed uint64) string {
	res := executeProgram(program, seed)
	// Convert output to comma-separated string
	strNums := make([]string, len(res))
	for i, num := range res {
		strNums[i] = fmt.Sprintf("%d", num)
	}
	return strings.Join(strNums, ",")
}

func executeProgram(program []uint64, seed uint64) []uint64 {
	var res []uint64
	registers := map[rune]uint64{
		'A': seed,
		'B': 0,
		'C': 0,
	}

	for ip := 0; ip < len(program)-1; {
		operand := program[ip+1]
		switch operator := program[ip]; operator {
		case 0: // adv
			registers['A'] >>= getValue(operand, registers)
		case 1: // bxl
			registers['B'] ^= operand
		case 2: // bst
			registers['B'] = getValue(operand, registers) & 7
		case 3: // jnz
			if registers['A'] != 0 {
				ip = int(operand)
				continue
			}
		case 4: // bxc
			registers['B'] ^= registers['C']
		case 5: // out
			val := getValue(operand, registers) & 7
			res = append(res, val)
		case 6: // bdv
			registers['B'] = registers['A'] >> getValue(operand, registers)
		case 7: // cdv
			registers['C'] = registers['A'] >> getValue(operand, registers)
		}
		ip += 2
	}
	return res
}

func getValue(comboOperand uint64, registers map[rune]uint64) uint64 {
	if comboOperand <= 3 {
		return comboOperand
	}
	switch comboOperand {
	case 4:
		return registers['A']
	case 5:
		return registers['B']
	case 6:
		return registers['C']
	default:
		panic("Invalid combo operand!")
	}
}

func main() {
	input, err := utility.ParseTextFile("input")
	if err != nil {
		log.Fatal(err)
	}

	program, seed := parseInput(input)
	fmt.Println("Part 1:", partOne(program, seed))

}
