package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

type Node struct {
	left  int
	right int
}

func play(playerCount, lastMarble int) int {
	scores := make([]int, playerCount)
	nodes := make([]Node, lastMarble+1)

	current := 0
	nodes[current].left = current
	nodes[current].right = current

	player := 0
	for marble := 1; marble <= lastMarble; marble++ {
		if marble%23 == 0 {
			remove := current
			for i := 0; i < 7; i++ {
				remove = nodes[remove].left
			}
			scores[player] += marble + remove
			left, right := nodes[remove].left, nodes[remove].right
			nodes[left].right = right
			nodes[right].left = left
			current = nodes[remove].right
		} else {
			insert := marble
			left := nodes[current].right
			right := nodes[left].right
			nodes[insert].left = left
			nodes[insert].right = right
			nodes[left].right = insert
			nodes[right].left = insert
			current = insert
		}
		player = (player + 1) % playerCount
	}

	highScore := 0
	for _, score := range scores {
		if score > highScore {
			highScore = score
		}
	}
	return highScore
}

func main() {
	input := readFile("input.txt")

	result := regexp.MustCompile(`(\d+) players; last marble is worth (\d+) points`).FindStringSubmatch(input)
	if result == nil {
		panic(input)
	}

	playerCount, lastMarble := toInt(result[1]), toInt(result[2])

	{
		fmt.Println("--- Part One ---")
		fmt.Println(play(playerCount, lastMarble))
	}

	{
		fmt.Println("--- Part Two ---")
		fmt.Println(play(playerCount, 100*lastMarble))
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
