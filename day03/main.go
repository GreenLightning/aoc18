package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

const N = 2000

type Claim struct {
	id            int
	left, top     int
	width, height int
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var claims []Claim
	regex := regexp.MustCompile(`#(\d+) @ (\d+),(\d+): (\d+)x(\d+)`)

	for scanner.Scan() {
		result := regex.FindStringSubmatch(scanner.Text())
		if result == nil {
			panic(fmt.Errorf("cannot parse claim: %s", scanner.Text()))
		}
		var claim Claim
		claim.id = toInt(result[1])
		claim.left, claim.top = toInt(result[2]), toInt(result[3])
		claim.width, claim.height = toInt(result[4]), toInt(result[5])
		claims = append(claims, claim)
	}

	var count [N][N]int

	for _, claim := range claims {
		for x := claim.left; x < claim.left+claim.width; x++ {
			for y := claim.top; y < claim.top+claim.height; y++ {
				count[x][y]++
			}
		}
	}

	{
		fmt.Println("--- Part One ---")
		overlap := 0
		for x := 0; x < N; x++ {
			for y := 0; y < N; y++ {
				if count[x][y] > 1 {
					overlap++
				}
			}
		}
		fmt.Println(overlap)
	}

	{
		fmt.Println("--- Part Two ---")
	search:
		for _, claim := range claims {
			for x := claim.left; x < claim.left+claim.width; x++ {
				for y := claim.top; y < claim.top+claim.height; y++ {
					if count[x][y] > 1 {
						continue search
					}
				}
			}
			fmt.Println(claim.id)
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
