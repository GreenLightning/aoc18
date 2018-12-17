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
	a, b, c, d int
}

type Instruction struct {
	op, inA, inB, outC int
}

type Sample struct {
	before, after State
	instr         Instruction
}

func (s *State) getReg(n int) (int, bool) {
	switch n {
	case 0:
		return s.a, true
	case 1:
		return s.b, true
	case 2:
		return s.c, true
	case 3:
		return s.d, true
	default:
		return 0, false
	}
}

func (s *State) setReg(n int, val int) bool {
	switch n {
	case 0:
		s.a = val
		return true
	case 1:
		s.b = val
		return true
	case 2:
		s.c = val
		return true
	case 3:
		s.d = val
		return true
	default:
		return false
	}
}

func (s *State) apply(i Instruction) bool {
	switch i.op {
	case ADDR:
		a, ok := s.getReg(i.inA)
		if !ok {
			return false
		}
		b, ok := s.getReg(i.inB)
		if !ok {
			return false
		}
		return s.setReg(i.outC, a+b)
	case ADDI:
		a, ok := s.getReg(i.inA)
		if !ok {
			return false
		}
		return s.setReg(i.outC, a+i.inB)
	case MULR:
		a, ok := s.getReg(i.inA)
		if !ok {
			return false
		}
		b, ok := s.getReg(i.inB)
		if !ok {
			return false
		}
		return s.setReg(i.outC, a*b)
	case MULI:
		a, ok := s.getReg(i.inA)
		if !ok {
			return false
		}
		return s.setReg(i.outC, a*i.inB)
	case BANR:
		a, ok := s.getReg(i.inA)
		if !ok {
			return false
		}
		b, ok := s.getReg(i.inB)
		if !ok {
			return false
		}
		return s.setReg(i.outC, a&b)
	case BANI:
		a, ok := s.getReg(i.inA)
		if !ok {
			return false
		}
		return s.setReg(i.outC, a&i.inB)
	case BORR:
		a, ok := s.getReg(i.inA)
		if !ok {
			return false
		}
		b, ok := s.getReg(i.inB)
		if !ok {
			return false
		}
		return s.setReg(i.outC, a|b)
	case BORI:
		a, ok := s.getReg(i.inA)
		if !ok {
			return false
		}
		return s.setReg(i.outC, a|i.inB)
	case SETR:
		a, ok := s.getReg(i.inA)
		if !ok {
			return false
		}
		return s.setReg(i.outC, a)
	case SETI:
		return s.setReg(i.outC, i.inA)
	case GTIR:
		b, ok := s.getReg(i.inB)
		if !ok {
			return false
		}
		if i.inA > b {
			return s.setReg(i.outC, 1)
		} else {
			return s.setReg(i.outC, 0)
		}
	case GTRI:
		a, ok := s.getReg(i.inA)
		if !ok {
			return false
		}
		if a > i.inB {
			return s.setReg(i.outC, 1)
		} else {
			return s.setReg(i.outC, 0)
		}
	case GTRR:
		a, ok := s.getReg(i.inA)
		if !ok {
			return false
		}
		b, ok := s.getReg(i.inB)
		if !ok {
			return false
		}
		if a > b {
			return s.setReg(i.outC, 1)
		} else {
			return s.setReg(i.outC, 0)
		}
	case EQIR:
		b, ok := s.getReg(i.inB)
		if !ok {
			return false
		}
		if i.inA == b {
			return s.setReg(i.outC, 1)
		} else {
			return s.setReg(i.outC, 0)
		}
	case EQRI:
		a, ok := s.getReg(i.inA)
		if !ok {
			return false
		}
		if a == i.inB {
			return s.setReg(i.outC, 1)
		} else {
			return s.setReg(i.outC, 0)
		}
	case EQRR:
		a, ok := s.getReg(i.inA)
		if !ok {
			return false
		}
		b, ok := s.getReg(i.inB)
		if !ok {
			return false
		}
		if a == b {
			return s.setReg(i.outC, 1)
		} else {
			return s.setReg(i.outC, 0)
		}
	default:
		return false
	}
}

func main() {
	lines := readLines("input.txt")

	var samples []Sample
	var testProgram []Instruction

	{
		beforeRegex := regexp.MustCompile(`Before: \[(\d+), (\d+), (\d+), (\d+)\]`)
		afterRegex := regexp.MustCompile(`After:  \[(\d+), (\d+), (\d+), (\d+)\]`)
		instrRegex := regexp.MustCompile(`(\d+) (\d+) (\d+) (\d+)`)

		index := 0

		for index+2 < len(lines) {
			if lines[index] == "" {
				index++
				continue
			}

			var sample Sample

			result := beforeRegex.FindStringSubmatch(lines[index])
			if result == nil {
				break
			}

			sample.before.a = toInt(result[1])
			sample.before.b = toInt(result[2])
			sample.before.c = toInt(result[3])
			sample.before.d = toInt(result[4])

			result = instrRegex.FindStringSubmatch(lines[index+1])
			if result == nil {
				panic(lines[index+1])
			}

			sample.instr.op = toInt(result[1])
			sample.instr.inA = toInt(result[2])
			sample.instr.inB = toInt(result[3])
			sample.instr.outC = toInt(result[4])

			result = afterRegex.FindStringSubmatch(lines[index+2])
			if result == nil {
				panic(lines[index+2])
			}

			sample.after.a = toInt(result[1])
			sample.after.b = toInt(result[2])
			sample.after.c = toInt(result[3])
			sample.after.d = toInt(result[4])

			samples = append(samples, sample)
			index += 3
		}

		for index < len(lines) {
			if lines[index] == "" {
				index++
				continue
			}

			var instr Instruction

			result := instrRegex.FindStringSubmatch(lines[index])
			if result == nil {
				panic(lines[index])
			}

			instr.op = toInt(result[1])
			instr.inA = toInt(result[2])
			instr.inB = toInt(result[3])
			instr.outC = toInt(result[4])

			testProgram = append(testProgram, instr)
			index++
		}
	}

	var possibleCodes [16][16]bool
	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			possibleCodes[i][j] = true
		}
	}

	count := 0
	for _, sample := range samples {
		instrCount := 0
		instr := sample.instr
		for i := 0; i < 16; i++ {
			instr.op = i
			state := sample.before
			ok := state.apply(instr)
			if ok {
				if state == sample.after {
					instrCount++
				} else {
					possibleCodes[sample.instr.op][i] = false
				}
			}
		}
		if instrCount >= 3 {
			count++
		}
	}

	{
		fmt.Println("--- Part One ---")
		fmt.Println(count)
	}

	opmap := make(map[int]int)
	for len(opmap) != 16 {
		for i := 0; i < 16; i++ {
			count := 0
			for j := 0; j < 16; j++ {
				if possibleCodes[i][j] {
					count++
				}
			}
			if count == 1 {
				for j := 0; j < 16; j++ {
					if possibleCodes[i][j] {
						opmap[i] = j
						for k := 0; k < 16; k++ {
							possibleCodes[k][j] = false
						}
					}
				}
			}
		}
	}

	for i, instr := range testProgram {
		testProgram[i].op = opmap[instr.op]
	}

	var state State
	for _, instr := range testProgram {
		state.apply(instr)
	}

	{
		fmt.Println("--- Part Two ---")
		fmt.Println(state.a)
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
