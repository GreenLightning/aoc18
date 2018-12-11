package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

const N = 300

func findBest(cells *[N][N]int, minSize, maxSize int) (bestX, bestY, bestSize int) {
	bestSum := 0
	for size := minSize; size <= maxSize; size++ {
		for y := 0; y <= N-size; y++ {
			for x := 0; x <= N-size; x++ {
				sum := 0
				for dy := 0; dy < size; dy++ {
					for dx := 0; dx < size; dx++ {
						sum += cells[y+dy][x+dx]
					}
				}
				if sum > bestSum {
					bestSum = sum
					bestX = x + 1
					bestY = y + 1
					bestSize = size
				}
			}
		}
	}
	return
}

func main() {
	serial := toInt(readFile("input.txt"))

	var cells [N][N]int

	for y := 0; y < N; y++ {
		for x := 0; x < N; x++ {
			rack := (x + 1) + 10
			power := (rack*(y+1) + serial) * rack
			cells[y][x] = ((power / 100) % 10) - 5
		}
	}

	{
		fmt.Println("--- Part One ---")
		x, y, _ := findBest(&cells, 3, 3)
		fmt.Printf("%d,%d\n", x, y)
	}

	{
		fmt.Println("--- Part Two ---")
		x, y, size := findBest(&cells, 1, N)
		fmt.Printf("%d,%d,%d\n", x, y, size)
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
