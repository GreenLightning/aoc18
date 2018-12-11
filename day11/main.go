package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

const N = 300

func findBest(cells *[N + 1][N + 1]int, minSize, maxSize int) (bestX, bestY, bestSize int) {
	bestSum := 0
	for size := minSize; size <= maxSize; size++ {
		offset := size - 1
		for y := 1; y <= N-size+1; y++ {
			for x := 1; x <= N-size+1; x++ {
				sum := cells[y+offset][x+offset] + cells[y-1][x-1] - cells[y+offset][x-1] - cells[y-1][x+offset]
				if sum > bestSum {
					bestSum = sum
					bestX = x
					bestY = y
					bestSize = size
				}
			}
		}
	}
	return
}

func main() {
	serial := toInt(readFile("input.txt"))

	var cells [N + 1][N + 1]int

	// Calculate power levels.
	for y := 1; y <= N; y++ {
		for x := 1; x <= N; x++ {
			rack := x + 10
			power := (rack*y + serial) * rack
			cells[y][x] = ((power / 100) % 10) - 5
		}
	}

	// Convert to prefix sum.
	{
		for y := 0; y < N+1; y++ {
			sum := 0
			for x := 0; x < N+1; x++ {
				sum += cells[y][x]
				cells[y][x] = sum
			}
		}

		for x := 0; x < N+1; x++ {
			sum := 0
			for y := 0; y < N+1; y++ {
				sum += cells[y][x]
				cells[y][x] = sum
			}
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
