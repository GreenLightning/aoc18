package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	x, y int
}

const N = 500

func main() {
	dimFlag := flag.Bool("dim", false, "print dimensions of input")
	visFlag := flag.String("vis", "", "print a visualisation to the specified filename")
	flag.Parse()

	lines := readLines("input.txt")

	var points []Point
	for _, line := range lines {
		split := strings.Split(line, ", ")
		var point Point
		point.x = toInt(split[0])
		point.y = toInt(split[1])
		points = append(points, point)
	}

	if *dimFlag {
		minX, maxX := math.MaxInt32, math.MinInt32
		minY, maxY := math.MaxInt32, math.MinInt32

		for _, p := range points {
			if p.x < minX {
				minX = p.x
			} else if p.x > maxX {
				maxX = p.x
			}
			if p.y < minY {
				minY = p.y
			} else if p.y > maxY {
				maxY = p.y
			}
		}

		fmt.Printf("Input dimensions: (%d, %d) - (%d, %d)\n", minX, minY, maxX, maxY)
		fmt.Printf("Test area: (0, 0) - (%d, %d)\n", N-1, N-1)
	}

	{
		fmt.Println("--- Part One ---")

		var field [N][N]int

		for x := 0; x < N; x++ {
			for y := 0; y < N; y++ {
				bestIndex := -1
				bestDistance := math.MaxInt32
				for i, p := range points {
					dist := abs(x-p.x) + abs(y-p.y)
					if dist < bestDistance {
						bestIndex = i
						bestDistance = dist
					} else if dist == bestDistance {
						bestIndex = -1
					}
				}
				field[x][y] = bestIndex
			}
		}

		if *visFlag != "" {
			file, err := os.Create(*visFlag)
			if err != nil {
				panic(err)
			}
			defer file.Close()

			writer := bufio.NewWriter(file)
			defer writer.Flush()

			for y := 0; y < N; y++ {
				for x := 0; x < N; x++ {
					found := false
					for _, p := range points {
						if p.x == x && p.y == y {
							found = true
							break
						}
					}

					value := field[x][y]
					if found {
						writer.WriteRune('#')
					} else if value == -1 {
						writer.WriteRune('.')
					} else if value >= 26 {
						writer.WriteRune(rune('A' + value - 26))
					} else {
						writer.WriteRune(rune('a' + value))
					}
				}
				writer.WriteRune('\n')
			}
		}

		best := 0
	search:
		for i := range points {
			for j := 0; j < N; j++ {
				if field[j][0] == i || field[j][N-1] == i || field[0][j] == i || field[N-1][j] == i {
					continue search
				}
			}

			count := 0
			for x := 0; x < N; x++ {
				for y := 0; y < N; y++ {
					if field[x][y] == i {
						count++
					}
				}
			}

			if count > best {
				best = count
			}
		}

		fmt.Println(best)
	}

	{
		fmt.Println("--- Part Two ---")
		count := 0
		for x := 0; x < N; x++ {
			for y := 0; y < N; y++ {
				dist := 0
				for _, p := range points {
					dist += abs(x-p.x) + abs(y-p.y)
				}
				if dist < 10000 {
					count++
				}
			}
		}
		fmt.Println(count)
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
