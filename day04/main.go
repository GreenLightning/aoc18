package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
)

const MINUTES_PER_SHIFT = 60

type Shift struct {
	asleep [MINUTES_PER_SHIFT]bool
}

type Guard struct {
	shifts []*Shift
}

func main() {
	lines := readLines("input.txt")

	sort.Strings(lines)

	shiftRegex := regexp.MustCompile(`\[\d+-\d+-\d+ \d+:\d+\] Guard #(\d+) begins shift`)
	asleepRegex := regexp.MustCompile(`\[\d+-\d+-\d+ 00:(\d+)\] falls asleep`)
	awakeRegex := regexp.MustCompile(`\[\d+-\d+-\d+ 00:(\d+)\] wakes up`)

	guards := make(map[int]*Guard)
	var lastShift *Shift
	for _, line := range lines {
		if result := shiftRegex.FindStringSubmatch(line); result != nil {
			id := toInt(result[1])
			guard := guards[id]
			if guard == nil {
				guard = new(Guard)
				guards[id] = guard
			}
			lastShift = new(Shift)
			guard.shifts = append(guard.shifts, lastShift)
		} else if result := asleepRegex.FindStringSubmatch(line); result != nil {
			for i := toInt(result[1]); i < MINUTES_PER_SHIFT; i++ {
				lastShift.asleep[i] = true
			}
		} else if result := awakeRegex.FindStringSubmatch(line); result != nil {
			for i := toInt(result[1]); i < MINUTES_PER_SHIFT; i++ {
				lastShift.asleep[i] = false
			}
		} else {
			panic(line)
		}
	}

	{
		fmt.Println("--- Part One ---")
		bestMinutes := 0
		bestResult := 0
		for id, guard := range guards {
			minutes := 0
			for _, shift := range guard.shifts {
				for i := 0; i < MINUTES_PER_SHIFT; i++ {
					if shift.asleep[i] {
						minutes++
					}
				}
			}

			bestCount := 0
			bestTarget := 0
			for i := 0; i < MINUTES_PER_SHIFT; i++ {
				count := 0
				for _, shift := range guard.shifts {
					if shift.asleep[i] {
						count++
					}
				}

				if count > bestCount {
					bestCount = count
					bestTarget = i
				}
			}

			if minutes > bestMinutes {
				bestMinutes = minutes
				bestResult = id * bestTarget
			}
		}
		fmt.Println(bestResult)
	}

	{
		fmt.Println("--- Part Two ---")
		bestCount := 0
		bestResult := 0
		for id, guard := range guards {
			for i := 0; i < MINUTES_PER_SHIFT; i++ {
				count := 0
				for _, shift := range guard.shifts {
					if shift.asleep[i] {
						count++
					}
				}

				if count > bestCount {
					bestCount = count
					bestResult = id * i
				}
			}
		}
		fmt.Println(bestResult)
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
