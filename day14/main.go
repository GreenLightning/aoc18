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
	inputBytes := make([]byte, len(inputText))
	for i, digit := range inputText {
		inputBytes[i] = byte(digit - '0')
	}

	recipes := make([]byte, 2, 30000000)
	recipes[0] = 3
	recipes[1] = 7
	first, second := 0, 1
	inputFound, beforeInput := false, 0

	addRecipe := func(score byte) {
		recipes = append(recipes, score)
		offset := len(recipes) - len(inputBytes)
		if !inputFound && offset >= 0 {
			inputFound = true
			beforeInput = offset
			for i, value := range inputBytes {
				if value != recipes[offset+i] {
					inputFound = false
					break
				}
			}
		}
	}

	for len(recipes) < input+10 || !inputFound {
		sum := recipes[first] + recipes[second]
		if sum >= 10 {
			addRecipe(sum / 10)
			sum -= 10
		}
		addRecipe(sum)
		first = (first + int(recipes[first]) + 1)
		for first >= len(recipes) {
			first -= len(recipes)
		}
		second = (second + int(recipes[second]) + 1)
		for second >= len(recipes) {
			second -= len(recipes)
		}
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
