package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
)

type Step struct {
	remaining    int
	dependencies []string
}

type Steps map[string]*Step

func (steps Steps) get(id string) *Step {
	result := steps[id]
	if result == nil {
		result = new(Step)
		result.remaining = 60 + int(id[0]-'A'+1)
		steps[id] = result
	}
	return result
}

func (steps Steps) copy() Steps {
	result := make(Steps)
	for id, step := range steps {
		cpy := new(Step)
		cpy.remaining = step.remaining
		cpy.dependencies = make([]string, len(step.dependencies))
		copy(cpy.dependencies, step.dependencies)
		result[id] = cpy
	}
	return result
}

func (steps Steps) getReady() []string {
	var ready []string
	for id, step := range steps {
		if len(step.dependencies) == 0 {
			ready = append(ready, id)
		}
	}
	sort.Strings(ready)
	return ready
}

func (steps Steps) remove(id string) {
	delete(steps, id)
	for _, step := range steps {
		for i, dep := range step.dependencies {
			if dep == id {
				end := len(step.dependencies) - 1
				step.dependencies[i] = step.dependencies[end]
				step.dependencies = step.dependencies[:end]
				break
			}
		}
	}
}

func main() {
	lines := readLines("input.txt")

	input := make(Steps)

	{
		regex := regexp.MustCompile(`Step (\w+) must be finished before step (\w+) can begin.`)
		for _, line := range lines {
			result := regex.FindStringSubmatch(line)
			if result == nil {
				panic(line)
			}
			input.get(result[1])
			second := input.get(result[2])
			second.dependencies = append(second.dependencies, result[1])
		}
	}

	{
		fmt.Println("--- Part One ---")
		steps := input.copy()

		var result string
		for len(steps) != 0 {
			ready := steps.getReady()
			steps.remove(ready[0])
			result = result + ready[0]
		}

		fmt.Println(result)
	}

	{
		fmt.Println("--- Part Two ---")
		steps := input.copy()
		workers := 5

		time := 0
		for len(steps) != 0 {
			time++
			ready := steps.getReady()
			for w := 0; w < workers && w < len(ready); w++ {
				step := steps[ready[w]]
				step.remaining--
				if step.remaining == 0 {
					steps.remove(ready[w])
				}
			}
		}

		fmt.Println(time)
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

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
