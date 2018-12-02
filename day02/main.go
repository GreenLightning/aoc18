package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var ids []string
	for scanner.Scan() {
		ids = append(ids, scanner.Text())
	}

	{
		fmt.Println("--- Part One ---")
		doubles, triples := 0, 0
		for _, id := range ids {
			chars := make(map[byte]int)
			for i := 0; i < len(id); i++ {
				chars[id[i]]++
			}
			for _, count := range chars {
				if count == 2 {
					doubles++
					break
				}
			}
			for _, count := range chars {
				if count == 3 {
					triples++
					break
				}
			}
		}
		fmt.Println(doubles * triples)
	}

	{
		fmt.Println("--- Part Two ---")
		seen := make(map[string]bool)

		for _, id := range ids {
			for i := 0; i < len(id); i++ {
				truncated := id[:i] + "_" + id[i+1:]
				if seen[truncated] {
					common := strings.Replace(truncated, "_", "", 1)
					fmt.Println(common)
				}
				seen[truncated] = true
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
