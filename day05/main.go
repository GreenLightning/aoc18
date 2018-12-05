package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func main() {
	polymer := readFile("input.txt")

	{
		fmt.Println("--- Part One ---")
		fmt.Println(len(react(polymer)))
	}

	{
		fmt.Println("--- Part Two ---")
		best := len(react(polymer))
		for i := byte(0); i < 26; i++ {
			test := polymer
			test = strings.Replace(test, string([]byte{'a' + i}), "", -1)
			test = strings.Replace(test, string([]byte{'A' + i}), "", -1)
			length := len(react(test))
			if length < best {
				best = length
			}
		}
		fmt.Println(best)
	}
}

func react(polymer string) string {
	for {
		old := polymer
		for i := byte(0); i < 26; i++ {
			polymer = strings.Replace(polymer, string([]byte{'a' + i, 'A' + i}), "", -1)
			polymer = strings.Replace(polymer, string([]byte{'A' + i, 'a' + i}), "", -1)
		}
		if len(polymer) == len(old) {
			return polymer
		}
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
