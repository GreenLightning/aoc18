package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	x, y, z, w int
}

func (p Point) Distance(o Point) int {
	return abs(p.x-o.x) + abs(p.y-o.y) + abs(p.z-o.z) + abs(p.w-o.w)
}

func main() {
	lines := readLines("input.txt")

	var points []Point
	for _, line := range lines {
		parts := strings.Split(line, ",")
		var point Point
		point.x = toInt(parts[0])
		point.y = toInt(parts[1])
		point.z = toInt(parts[2])
		point.w = toInt(parts[3])
		points = append(points, point)
	}

	var constellations [][]Point
	for _, point := range points {
		var matching []int
		for index, constellation := range constellations {
			for _, other := range constellation {
				if point.Distance(other) <= 3 {
					matching = append(matching, index)
					break
				}
			}
		}
		if len(matching) == 0 {
			constellations = append(constellations, []Point{point})
		} else if len(matching) == 1 {
			index := matching[0]
			constellations[index] = append(constellations[index], point)
		} else {
			var combined []Point
			combined = append(combined, point)
			for _, index := range matching {
				combined = append(combined, constellations[index]...)
			}
			for i := len(matching) - 1; i >= 0; i-- {
				index := matching[i]
				constellations[index] = constellations[len(constellations)-1]
				constellations = constellations[:len(constellations)-1]
			}
			constellations = append(constellations, combined)
		}
	}

	{
		fmt.Println("--- Part One ---")
		fmt.Println(len(constellations))
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
