package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
)

const N = 2000

func main() {
	lines := readLines("input.txt")

	// 0 sand
	// 1 clay
	// 10 passed
	// 11 settled
	var world [N][N]byte

	minY, maxY := N, 0

	yRange := regexp.MustCompile(`x=(\d+), y=(\d+)\.\.(\d+)`)
	xRange := regexp.MustCompile(`y=(\d+), x=(\d+)\.\.(\d+)`)

	for _, line := range lines {
		if result := yRange.FindStringSubmatch(line); result != nil {
			x := toInt(result[1])
			start, end := toInt(result[2]), toInt(result[3])
			for y := start; y <= end; y++ {
				world[y][x] = 1
			}
			if start < minY {
				minY = start
			} else if end > maxY {
				maxY = end
			}
		} else if result := xRange.FindStringSubmatch(line); result != nil {
			y := toInt(result[1])
			start, end := toInt(result[2]), toInt(result[3])
			for x := start; x <= end; x++ {
				world[y][x] = 1
			}
			if y < minY {
				minY = y
			} else if y > maxY {
				maxY = y
			}
		} else {
			panic(line)
		}
	}

	unchangedCount := 0
	for unchangedCount < 100 {
		x, y := 500, 0
		changed := false
		backing := false
		left := (rand.Int31n(2) == 0)
	particle:
		for {
			if world[y][x] != 10 {
				changed = true
			}

			world[y][x] = 10

			if y == N-1 {
				break
			}

			if world[y+1][x]&1 == 0 {
				backing = false
				left = (rand.Int31n(2) == 0)
				y++
				continue
			}

			if world[y+1][x] == 11 {
				for xx := 0; ; xx++ {
					if world[y+1][x+xx]&1 == 0 {
						changed = true
						world[y+1][x+xx] = 11
						break particle
					}
					if world[y+1][x+xx] != 11 {
						break
					}
				}
				for xx := 0; ; xx++ {
					if world[y+1][x-xx]&1 == 0 {
						changed = true
						world[y+1][x-xx] = 11
						break particle
					}
					if world[y+1][x-xx] != 11 {
						break
					}
				}
			}

			if left {
				if !backing {
					if world[y][x-1]&1 == 0 {
						x--
						continue
					}
					backing = true
				}
				if world[y][x+1]&1 == 0 {
					x++
					continue
				}
			} else {
				if !backing {
					if world[y][x+1]&1 == 0 {
						x++
						continue
					}
					backing = true
				}
				if world[y][x-1]&1 == 0 {
					x--
					continue
				}
			}

			changed = true
			world[y][x] = 11
			break
		}

		if !changed {
			unchangedCount++
		}
	}

	passedCount := 0
	settledCount := 0

	for y := 0; y < N; y++ {
		if y >= minY && y <= maxY {
			for x := 0; x < N; x++ {
				switch world[y][x] {
				case 10:
					passedCount++
				case 11:
					settledCount++
				}
			}
		}
	}

	{
		fmt.Println("--- Part One ---")
		fmt.Println(passedCount + settledCount)
	}

	{
		fmt.Println("--- Part Two ---")
		fmt.Println(settledCount)
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
