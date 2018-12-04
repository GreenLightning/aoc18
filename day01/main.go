package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	lines := readLines("input.txt")

	changes := make([]int, len(lines))
	for index, line := range lines {
		changes[index] = toInt(line)
	}

	{
		fmt.Println("--- Part One ---")
		frequency := 0
		for _, change := range changes {
			frequency += change
		}
		fmt.Println(frequency)
	}

	{
		fmt.Println("--- Part Two ---")
		frequency := 0
		seen := make(map[int]bool)
		seen[0] = true

		for {
			for _, change := range changes {
				frequency += change
				if seen[frequency] {
					fmt.Println(frequency)
					return
				}
				seen[frequency] = true
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
