package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

const (
	// These types must match, so that a tool has the same index as the region
	// it cannot be used in, for easy comparison.

	// region types
	ROCKY  = 0
	WET    = 1
	NARROW = 2

	// tool types
	NEITHER       = 0
	TORCH         = 1
	CLIMBING_GEAR = 2
)

var deltas = []Position{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

type Position struct {
	x, y int
}

func (p Position) Plus(o Position) Position {
	return Position{
		x: p.x + o.x,
		y: p.y + o.y,
	}
}

func (p Position) Less(o Position) bool {
	if p.y != o.y {
		return p.y < o.y
	}
	return p.x < o.x
}

type Destination struct {
	position Position
	tool     int
	distance int
}

func (d Destination) Less(o Destination) bool {
	if d.distance != o.distance {
		return d.distance < o.distance
	}
	return d.position.Less(o.position)
}

func main() {
	lines := readLines("input.txt")

	var depth int
	var target Position
	{
		depthResult := regexp.MustCompile(`depth: (\d+)`).FindStringSubmatch(lines[0])
		depth = toInt(depthResult[1])
		targetResult := regexp.MustCompile(`target: (\d+),(\d+)`).FindStringSubmatch(lines[1])
		target.x = toInt(targetResult[1])
		target.y = toInt(targetResult[2])
	}

	width, height := 4000, 4000

	erosionLevel := make([][]int, height)
	for y := range erosionLevel {
		erosionLevel[y] = make([]int, width)
	}

	erosionLevel[0][0] = depth % 20183
	for x := 1; x < width; x++ {
		geologicalIndex := 16807 * x
		erosionLevel[0][x] = (geologicalIndex + depth) % 20183
	}
	for y := 1; y < height; y++ {
		geologicalIndex := 48271 * y
		erosionLevel[y][0] = (geologicalIndex + depth) % 20183
	}
	for y := 1; y < height; y++ {
		for x := 1; x < width; x++ {
			geologicalIndex := erosionLevel[y-1][x] * erosionLevel[y][x-1]
			erosionLevel[y][x] = (geologicalIndex + depth) % 20183
		}
	}
	erosionLevel[target.y][target.x] = depth % 20183

	world := make([][]int, height)
	for y := range world {
		world[y] = make([]int, width)
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			world[y][x] = erosionLevel[y][x] % 3
		}
	}

	{
		fmt.Println("--- Part One ---")

		risk := 0
		for y := 0; y <= target.y; y++ {
			for x := 0; x <= target.x; x++ {
				risk += world[y][x]
			}
		}

		fmt.Println(risk)
	}

	{
		fmt.Println("--- Part Two ---")

		openByDistance := make(map[int][]Destination)
		nextDistance := 0

		reached := make([][][]bool, height)
		for y := range reached {
			reached[y] = make([][]bool, width)
			for x := range reached[y] {
				reached[y][x] = make([]bool, 3)
			}
		}

		start := Destination{
			position: Position{0, 0},
			tool:     TORCH,
			distance: 0,
		}
		openByDistance[start.distance] = append(openByDistance[start.distance], start)

		for len(openByDistance) != 0 {
			open := openByDistance[nextDistance]
			if len(open) == 0 {
				delete(openByDistance, nextDistance)
				nextDistance++
				continue
			}
			current := open[len(open)-1]
			openByDistance[nextDistance] = open[:len(open)-1]

			if reached[current.position.y][current.position.x][current.tool] {
				continue
			}
			reached[current.position.y][current.position.x][current.tool] = true

			for tool := 0; tool < 3; tool++ {
				if tool != current.tool && tool != world[current.position.y][current.position.x] {
					next := Destination{
						position: current.position,
						tool:     tool,
						distance: current.distance + 7,
					}
					if !reached[next.position.y][next.position.x][next.tool] {
						openByDistance[next.distance] = append(openByDistance[next.distance], next)
					}
				}
			}

			for _, delta := range deltas {
				p := current.position.Plus(delta)
				if p.x >= 0 && p.y >= 0 && current.tool != world[p.y][p.x] {
					next := Destination{
						position: p,
						tool:     current.tool,
						distance: current.distance + 1,
					}
					if !reached[next.position.y][next.position.x][next.tool] {
						openByDistance[next.distance] = append(openByDistance[next.distance], next)
					}
				}
			}

			if current.position == target && current.tool == TORCH {
				fmt.Println(current.distance)
				break
			}
		}
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
