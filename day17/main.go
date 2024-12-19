package main

import (
	"aoc2024/utility"
	"fmt"
	"log"
	"math"
	"slices"
	"strings"
)

func parseInput(input []string) ([]uint64, uint64, error) {
	var seed uint64
	n, err := fmt.Sscanf(input[0], "Register A: %d", &seed)
	if err != nil {
		return nil, 0, fmt.Errorf("error parsing seed: %w", err)
	}
	if n != 1 {
		return nil, 0, fmt.Errorf("could not parse seed from line: %q", input[0])
	}

	programStr := strings.TrimPrefix(input[4], "Program: ")
	parts := strings.Split(programStr, ",")
	program := make([]uint64, len(parts))
	for i, part := range parts {
		part = strings.TrimSpace(part)
		var num uint64
		n, err := fmt.Sscanf(part, "%d", &num)
		if err != nil {
			return nil, 0, fmt.Errorf("error parsing program part %q: %w", part, err)
		}
		if n != 1 {
			return nil, 0, fmt.Errorf("could not parse number from part: %q", part)
		}
		program[i] = num
	}

	return program, seed, nil
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

func findLowestSelfReplicatingSeed(program []uint64) uint64 {
	var seed uint64
	for itr := len(program) - 1; itr >= 0; itr-- {
		seed <<= 3
		for !slices.Equal(executeProgram(program, seed), program[itr:]) {
			seed++
		}
	}
	return seed
}

func executeProgram(program []uint64, seed uint64) []uint64 {
	var res []uint64
	registers := map[rune]uint64{
		'A': seed,
		'B': 0,
		'C': 0,
	}

	for ip := 0; ip < len(program)-1; {
		operator := program[ip]
		operand := program[ip+1]

		newIP, jump, newRes := executeInstruction(operator, operand, registers, ip, res)
		res = newRes
		if jump {
			ip = newIP
		} else {
			ip += 2
		}
	}
	return res
}

func executeInstruction(operator, operand uint64, registers map[rune]uint64, ip int, res []uint64) (int, bool, []uint64) { //nolint:lll
	switch operator {
	case 0:
		return opADV(registers, operand, ip, res)
	case 1:
		return opBXL(registers, operand, ip, res)
	case 2:
		return opBST(registers, operand, ip, res)
	case 3:
		return opJNZ(registers, operand, ip, res)
	case 4:
		return opBXC(registers, ip, res)
	case 5:
		return opOUT(registers, operand, ip, res)
	case 6:
		return opBDV(registers, operand, ip, res)
	case 7:
		return opCDV(registers, operand, ip, res)
	default:
		log.Fatalf("Invalid operator: %d", operator)
		return ip, false, res
	}
}

func opADV(registers map[rune]uint64, operand uint64, ip int, res []uint64) (int, bool, []uint64) {
	registers['A'] >>= getValue(operand, registers)
	return ip, false, res
}

func opBXL(registers map[rune]uint64, operand uint64, ip int, res []uint64) (int, bool, []uint64) {
	registers['B'] ^= operand
	return ip, false, res
}

func opBST(registers map[rune]uint64, operand uint64, ip int, res []uint64) (int, bool, []uint64) {
	registers['B'] = getValue(operand, registers) & 7
	return ip, false, res
}

func opJNZ(registers map[rune]uint64, operand uint64, ip int, res []uint64) (int, bool, []uint64) {
	if registers['A'] != 0 {
		if operand > math.MaxInt64 {
			log.Fatal("operand too large to convert to int")
		}
		ip = int(operand)
		return ip, true, res
	}
	return ip, false, res
}

func opBXC(registers map[rune]uint64, ip int, res []uint64) (int, bool, []uint64) {
	registers['B'] ^= registers['C']
	return ip, false, res
}

func opOUT(registers map[rune]uint64, operand uint64, ip int, res []uint64) (int, bool, []uint64) {
	val := getValue(operand, registers) & 7
	res = append(res, val)
	return ip, false, res
}

func opBDV(registers map[rune]uint64, operand uint64, ip int, res []uint64) (int, bool, []uint64) {
	registers['B'] = registers['A'] >> getValue(operand, registers)
	return ip, false, res
}

func opCDV(registers map[rune]uint64, operand uint64, ip int, res []uint64) (int, bool, []uint64) {
	registers['C'] = registers['A'] >> getValue(operand, registers)
	return ip, false, res
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
		log.Fatalf("Invalid combo operand: %d", comboOperand)
		return 0 // unreachable
	}
}

func main() {
	input, err := utility.ParseTextFile("input")
	if err != nil {
		log.Fatal(err)
	}
	program, seed, err := parseInput(input)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Part 1:", partOne(program, seed))
	log.Println("Part 2:", findLowestSelfReplicatingSeed(program))
}
