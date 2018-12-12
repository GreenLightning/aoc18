package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type State struct {
	data   string
	offset int
}

func (state *State) simulate(rules map[string]string) {
	padded := fmt.Sprintf("....%s....", state.data)
	var result bytes.Buffer
	for i := 0; i < len(padded)-5; i++ {
		result.WriteString(rules[padded[i:i+5]])
	}
	state.data = result.String()
	state.offset -= 2 // because of the padding we extended data by 2 entries on each side
	leftTrim := strings.Index(state.data, "#")
	rightTrim := strings.LastIndex(state.data, "#")
	state.data = state.data[leftTrim : rightTrim+1]
	state.offset += leftTrim
}

func (state *State) sum() int {
	result := 0
	for i, value := range state.data {
		if value == '#' {
			result += i + state.offset
		}
	}
	return result
}

func main() {
	lines := readLines("input.txt")

	var initial State
	rules := make(map[string]string)

	result := regexp.MustCompile(`initial state: ([.#]*)`).FindStringSubmatch(lines[0])
	if result == nil {
		panic(lines[0])
	}

	initial.data = result[1]
	initial.offset = 0

	regex := regexp.MustCompile(`([.#]{5}) => ([.#])`)
	for index := 1; index < len(lines); index++ {
		line := lines[index]
		if line == "" {
			continue
		}

		result := regex.FindStringSubmatch(line)
		if result == nil {
			panic(line)
		}

		rules[result[1]] = result[2]
	}

	{
		fmt.Println("--- Part One ---")
		state := initial
		for generation := 0; generation < 20; generation++ {
			state.simulate(rules)
		}
		fmt.Println(state.sum())
	}

	{
		fmt.Println("--- Part Two ---")
		state := initial
		generation, generations := 0, 50000000000
		for {
			previous := state
			state.simulate(rules)
			generation++
			if state.data == previous.data {
				state.offset += (generations - generation) * (state.offset - previous.offset)
				break
			}
		}
		fmt.Println(state.sum())
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
