package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
)

type Star struct {
	x, y, dx, dy int
}

func getDimensions(stars []Star) (minX, maxX, minY, maxY int) {
	minX, maxX = math.MaxInt32, math.MinInt32
	minY, maxY = math.MaxInt32, math.MinInt32
	for _, star := range stars {
		if star.x < minX {
			minX = star.x
		} else if star.x > maxX {
			maxX = star.x
		}
		if star.y < minY {
			minY = star.y
		} else if star.y > maxY {
			maxY = star.y
		}
	}
	return
}

func print(stars []Star) {
	minX, maxX, minY, maxY := getDimensions(stars)

	writer := bufio.NewWriter(os.Stdin)
	defer writer.Flush()

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			found := false
			for _, star := range stars {
				if star.y == y && star.x == x {
					found = true
					break
				}
			}
			if found {
				writer.WriteRune('#')
			} else {
				writer.WriteRune('.')
			}
		}
		writer.WriteRune('\n')
	}
}

func main() {
	lines := readLines("input.txt")

	var stars []Star
	regex := regexp.MustCompile(`position=<\s*(-?\d+),\s*(-?\d+)> velocity=<\s*(-?\d+),\s*(-?\d+)>`)
	for _, line := range lines {
		result := regex.FindStringSubmatch(line)
		if result == nil {
			panic(line)
		}
		stars = append(stars, Star{
			x:  toInt(result[1]),
			y:  toInt(result[2]),
			dx: toInt(result[3]),
			dy: toInt(result[4]),
		})
	}

	seconds := 0
	for {
		minX, maxX, minY, maxY := getDimensions(stars)
		before := (maxX - minX) * (maxY - minY)

		// Simulate.
		for i, star := range stars {
			stars[i].x += star.dx
			stars[i].y += star.dy
		}

		minX, maxX, minY, maxY = getDimensions(stars)
		after := (maxX - minX) * (maxY - minY)

		// The stars are the most dense when they form the message.
		// If we are expanding, we must have just passed the message,
		// so we reverse the last step and then stop searching.
		if after > before {
			for i, star := range stars {
				stars[i].x -= star.dx
				stars[i].y -= star.dy
			}
			break
		}

		seconds++
	}

	{
		fmt.Println("--- Part One ---")
		print(stars)
	}

	{
		fmt.Println("--- Part Two ---")
		fmt.Println(seconds)
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
