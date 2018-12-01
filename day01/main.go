package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var changes []int
	for scanner.Scan() {
		line := scanner.Text()
		changes = append(changes, toInt(line))
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

func toInt(s string) int {
	result, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return result
}
