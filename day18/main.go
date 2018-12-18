package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

const N = 50
const WIDTH = N + 2

const (
	OPEN       = 1
	TREES      = 2
	LUMBERYARD = 3
)

func step(old []byte, new []byte) {
	for y := 1; y <= N; y++ {
		for x := 1; x <= N; x++ {
			var counts [4]int
			// Do not count current tile.
			counts[old[y*WIDTH+x]]--
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					counts[old[(y+dy)*WIDTH+(x+dx)]]++
				}
			}

			state := old[y*WIDTH+x]
			switch state {
			case OPEN:
				if counts[TREES] >= 3 {
					state = TREES
				}
			case TREES:
				if counts[LUMBERYARD] >= 3 {
					state = LUMBERYARD
				}
			case LUMBERYARD:
				if counts[LUMBERYARD] == 0 || counts[TREES] == 0 {
					state = OPEN
				}
			}
			new[y*WIDTH+x] = state
		}
	}
}

func value(world []byte) int {
	var counts [4]int
	for y := 1; y <= N; y++ {
		for x := 1; x <= N; x++ {
			counts[world[y*WIDTH+x]]++
		}
	}
	return counts[TREES] * counts[LUMBERYARD]
}

func hash(world []byte) int {
	hash := 181
	for y := 1; y <= N; y++ {
		for x := 1; x <= N; x++ {
			hash = (hash << 1) + hash + int(world[y*WIDTH+x])
		}
	}
	return hash
}

type Node struct {
	next  *Node
	world []byte
}

type NodeList []*Node

func main() {
	lines := readLines("input.txt")

	world := make([]byte, (N+2)*(N+2))

	for y := 1; y <= N; y++ {
		for x := 1; x <= N; x++ {
			switch lines[y-1][x-1] {
			case '.':
				world[y*WIDTH+x] = OPEN
			case '|':
				world[y*WIDTH+x] = TREES
			case '#':
				world[y*WIDTH+x] = LUMBERYARD
			}
		}
	}

	iterations := 0

	{
		fmt.Println("--- Part One ---")

		buffer := make([]byte, (N+2)*(N+2))
		for ; iterations < 10; iterations++ {
			step(world, buffer)
			world, buffer = buffer, world
		}

		fmt.Println(value(world))
	}

	{
		fmt.Println("--- Part Two ---")

		maxIterations := 1000000000

		listFromHash := make(map[int]NodeList)
		previousNode := new(Node)
		previousNode.world = world
		listFromHash[hash(world)] = NodeList{previousNode}

		var looper *Node

	search:
		for ; iterations < maxIterations; iterations++ {
			buffer := make([]byte, (N+2)*(N+2))

			step(world, buffer)

			hashValue := hash(buffer)
			list := listFromHash[hashValue]

			if list != nil {
				for _, node := range list {
					if bytes.Equal(buffer, node.world) {
						previousNode.next = node
						looper = node
						iterations++
						break search
					}
				}
			}

			currentNode := new(Node)
			currentNode.world = buffer
			previousNode.next = currentNode

			list = append(list, currentNode)
			listFromHash[hashValue] = list

			previousNode = currentNode
			world = buffer
		}

		if looper != nil {
			for ; iterations < maxIterations; iterations++ {
				looper = looper.next
			}
			world = looper.world
		}

		fmt.Println(value(world))
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
