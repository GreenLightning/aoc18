package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func main() {
	inputText := readFile("input.txt")
	input := toInt(inputText)

	recipes := []int{3, 7}
	first, second := 0, 1
	inputFound, beforeInput := false, 0

	for len(recipes) < input+10 || !inputFound {
		new := fmt.Sprintf("%d", recipes[first]+recipes[second])
		for _, digit := range new {
			recipes = append(recipes, int(digit)-'0')
			if !inputFound && len(recipes) >= len(inputText) {
				match := true
				for i, digit := range inputText {
					if int(digit)-'0' != recipes[len(recipes)-len(inputText)+i] {
						match = false
						break
					}
				}
				if match {
					inputFound = true
					beforeInput = len(recipes) - len(inputText)
				}
			}
		}
		first = (first + recipes[first] + 1) % len(recipes)
		second = (second + recipes[second] + 1) % len(recipes)
	}

	{
		fmt.Println("--- Part One ---")
		for i := input; i < input+10; i++ {
			fmt.Print(recipes[i])
		}
		fmt.Println()
	}

	{
		fmt.Println("--- Part Two ---")
		fmt.Println(beforeInput)
	}
}

func readFile(filename string) string {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(bytes))
}

func toInt(s string) int {
	result, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return result
}
