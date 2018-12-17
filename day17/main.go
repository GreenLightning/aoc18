package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

const N = 2000

// The world should be initialized with sand, so SAND must be 0. Things that
// block the flow, i.e. CLAY and WATER have the least significant bit set, so
// that blocking can be checked using bit tests.
const (
	SAND    = 0 // .
	CLAY    = 1 // #
	WETSAND = 2 // |
	WATER   = 3 // ~
)

type Particle struct {
	x, y int
}

func main() {
	lines := readLines("input.txt")

	var world [N][N]byte

	minY, maxY := N, 0

	yRange := regexp.MustCompile(`x=(\d+), y=(\d+)\.\.(\d+)`)
	xRange := regexp.MustCompile(`y=(\d+), x=(\d+)\.\.(\d+)`)

	for _, line := range lines {
		if result := yRange.FindStringSubmatch(line); result != nil {
			x := toInt(result[1])
			start, end := toInt(result[2]), toInt(result[3])
			for y := start; y <= end; y++ {
				world[y][x] = CLAY
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
				world[y][x] = CLAY
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

	var particles []Particle
	particles = append(particles, Particle{x: 500, y: 0})
	for len(particles) != 0 {
		particle := particles[len(particles)-1]
		particles = particles[:len(particles)-1]

		// Mark current position.
		world[particle.y][particle.x] = WETSAND

		// Fall down unitl we hit a blocking surface.
		for particle.y+1 < N && world[particle.y+1][particle.x]&1 == 0 {
			particle.y++
			world[particle.y][particle.x] = WETSAND
		}

		// Stop if we reach the bottom of the simulation area.
		if particle.y == N-1 {
			continue
		}

		left, right := 1, 1
		leftBlocked, rightBlocked := false, false

		// Spread to the left.
		for {
			if world[particle.y][particle.x-left]&1 != 0 {
				leftBlocked = true
				break
			}
			if world[particle.y+1][particle.x-left]&1 == 0 {
				break
			}
			world[particle.y][particle.x-left] = WETSAND
			left++
		}

		// Spread to the right.
		for {
			if world[particle.y][particle.x+right]&1 != 0 {
				rightBlocked = true
				break
			}
			if world[particle.y+1][particle.x+right]&1 == 0 {
				break
			}
			world[particle.y][particle.x+right] = WETSAND
			right++
		}

		// If we are enclosed, fill everything with water and repeat the
		// process one tile up from where we came from.
		if leftBlocked && rightBlocked {
			for offset := -left + 1; offset < right; offset++ {
				world[particle.y][particle.x+offset] = WATER
			}
			particles = append(particles, Particle{x: particle.x, y: particle.y - 1})
			continue
		}

		// Otherwise spawn particles on the free ends.
		if !leftBlocked {
			particles = append(particles, Particle{x: particle.x - left, y: particle.y})
		}
		if !rightBlocked {
			particles = append(particles, Particle{x: particle.x + right, y: particle.y})
		}
	}

	passedCount := 0
	settledCount := 0

	for y := 0; y < N; y++ {
		if y >= minY && y <= maxY {
			for x := 0; x < N; x++ {
				switch world[y][x] {
				case WETSAND:
					passedCount++
				case WATER:
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
