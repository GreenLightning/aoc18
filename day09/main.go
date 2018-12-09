package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

type Node struct {
	value int
	left  *Node
	right *Node
}

func play(playerCount, lastMarble int) int {
	scores := make([]int, playerCount)

	current := &Node{}
	current.value = 0
	current.left = current
	current.right = current

	player := 0
	for marble := 1; marble <= lastMarble; marble++ {
		if marble%23 == 0 {
			remove := current
			for i := 0; i < 7; i++ {
				remove = remove.left
			}
			scores[player] += marble + remove.value
			left, right := remove.left, remove.right
			left.right = right
			right.left = left
			current = remove.right
		} else {
			insert := &Node{value: marble}
			left, right := current.right, current.right.right
			insert.left = left
			insert.right = right
			left.right = insert
			right.left = insert
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
