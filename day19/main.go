package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

const (
	ADDR = iota
	ADDI
	MULR
	MULI
	BANR
	BANI
	BORR
	BORI
	SETR
	SETI
	GTIR
	GTRI
	GTRR
	EQIR
	EQRI
	EQRR
)

type State struct {
	regs [6]int
}

type Instruction struct {
	op, inA, inB, outC int
}

func (s *State) apply(i Instruction) {
	switch i.op {
	case ADDR:
		a := s.regs[i.inA]
		b := s.regs[i.inB]
		s.regs[i.outC] = a + b
	case ADDI:
		a := s.regs[i.inA]
		s.regs[i.outC] = a + i.inB
	case MULR:
		a := s.regs[i.inA]
		b := s.regs[i.inB]
		s.regs[i.outC] = a * b
	case MULI:
		a := s.regs[i.inA]
		s.regs[i.outC] = a * i.inB
	case BANR:
		a := s.regs[i.inA]
		b := s.regs[i.inB]
		s.regs[i.outC] = a & b
	case BANI:
		a := s.regs[i.inA]
		s.regs[i.outC] = a & i.inB
	case BORR:
		a := s.regs[i.inA]
		b := s.regs[i.inB]
		s.regs[i.outC] = a | b
	case BORI:
		a := s.regs[i.inA]
		s.regs[i.outC] = a | i.inB
	case SETR:
		a := s.regs[i.inA]
		s.regs[i.outC] = a
	case SETI:
		s.regs[i.outC] = i.inA
	case GTIR:
		b := s.regs[i.inB]
		if i.inA > b {
			s.regs[i.outC] = 1
		} else {
			s.regs[i.outC] = 0
		}
	case GTRI:
		a := s.regs[i.inA]
		if a > i.inB {
			s.regs[i.outC] = 1
		} else {
			s.regs[i.outC] = 0
		}
	case GTRR:
		a := s.regs[i.inA]
		b := s.regs[i.inB]
		if a > b {
			s.regs[i.outC] = 1
		} else {
			s.regs[i.outC] = 0
		}
	case EQIR:
		b := s.regs[i.inB]
		if i.inA == b {
			s.regs[i.outC] = 1
		} else {
			s.regs[i.outC] = 0
		}
	case EQRI:
		a := s.regs[i.inA]
		if a == i.inB {
			s.regs[i.outC] = 1
		} else {
			s.regs[i.outC] = 0
		}
	case EQRR:
		a := s.regs[i.inA]
		b := s.regs[i.inB]
		if a == b {
			s.regs[i.outC] = 1
		} else {
			s.regs[i.outC] = 0
		}
	}
}

func main() {
	lines := readLines("input.txt")

	var ipReg, valueReg int
	var program []Instruction

	{
		ipRegex := regexp.MustCompile(`#ip (\d+)`)
		line := lines[0]
		result := ipRegex.FindStringSubmatch(line)
		if result == nil {
			panic(line)
		}
		ipReg = toInt(result[1])
	}

	{
		instrRegex := regexp.MustCompile(`(\w+) (\d+) (\d+) (\d+)`)

		opMap := map[string]int{
			"addr": ADDR,
			"addi": ADDI,
			"mulr": MULR,
			"muli": MULI,
			"banr": BANR,
			"bani": BANI,
			"borr": BORR,
			"bori": BORI,
			"setr": SETR,
			"seti": SETI,
			"gtir": GTIR,
			"gtri": GTRI,
			"gtrr": GTRR,
			"eqir": EQIR,
			"eqri": EQRI,
			"eqrr": EQRR,
		}

		for i := 1; i < len(lines); i++ {
			line := lines[i]
			result := instrRegex.FindStringSubmatch(line)
			if result == nil {
				panic(line)
			}
			var instr Instruction
			instr.op = opMap[result[1]]
			instr.inA = toInt(result[2])
			instr.inB = toInt(result[3])
			instr.outC = toInt(result[4])
			program = append(program, instr)
		}
	}

	{
		if 17 >= len(program) {
			panic("program too short")
		}
		if program[17].op != ADDI {
			panic("expected addi at index 17")
		}
		valueReg = program[17].outC
	}

	for part := 0; part < 2; part++ {
		if part == 0 {
			fmt.Println("--- Part One ---")
		} else {
			fmt.Println("--- Part Two ---")
		}

		var state State
		state.regs[0] = part

		for state.regs[ipReg] >= 0 && state.regs[ipReg] < len(program) {
			state.apply(program[state.regs[ipReg]])
			state.regs[ipReg]++
			if state.regs[ipReg] == 1 {
				value := state.regs[valueReg]
				result := 0
				for factor := 1; factor <= value; factor++ {
					if value%factor == 0 {
						result += factor
					}
				}
				fmt.Println(result)
				break
			}
		}
	}
}

func readLines(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func toInt(s string) int {
	result, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return result
}
