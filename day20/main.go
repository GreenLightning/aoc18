package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

const ROOM_SIZE = 200
const TILE_SIZE = 2*ROOM_SIZE + 1

const (
	WALL = 0
	ROOM = 1
	DOOR = 2
)

type Position struct {
	x, y int
}

func (pos Position) Plus(other Position) Position {
	return Position{
		x: pos.x + other.x,
		y: pos.y + other.y,
	}
}

func (pos Position) Times(factor int) Position {
	return Position{
		x: factor * pos.x,
		y: factor * pos.y,
	}
}

func (pos Position) Index() int {
	return pos.y*TILE_SIZE + pos.x
}

type Destination struct {
	position Position
	distance int
}

func walk(path string, world []byte, pos Position) Position {
	if len(path) == 0 {
		return pos
	}

	index := strings.Index(path, "(")
	if index == -1 {
		index = len(path)
	}

	for i := 0; i < index; i++ {
		switch path[i] {
		case 'N':
			pos.y--
			world[pos.Index()] = DOOR
			pos.y--
			world[pos.Index()] = ROOM
		case 'S':
			pos.y++
			world[pos.Index()] = DOOR
			pos.y++
			world[pos.Index()] = ROOM
		case 'W':
			pos.x--
			world[pos.Index()] = DOOR
			pos.x--
			world[pos.Index()] = ROOM
		case 'E':
			pos.x++
			world[pos.Index()] = DOOR
			pos.x++
			world[pos.Index()] = ROOM
		default:
			panic(fmt.Sprintf("unexpected input character: %c", path[i]))
		}
	}

	path = path[index:]
	if len(path) == 0 {
		return pos
	}

	if path[0] != '(' {
		panic(fmt.Sprintf("unexpected input character: %c", path[0]))
	}
	path = path[1:]

	var options []string
	var remainder string

	start := 0
	parentheses := 0
loop:
	for i := 0; i < len(path); i++ {
		switch path[i] {
		case '(':
			parentheses++
		case '|':
			if parentheses == 0 {
				options = append(options, path[start:i])
				start = i + 1
			}
		case ')':
			if parentheses == 0 {
				options = append(options, path[start:i])
				remainder = path[i+1:]
				break loop
			}
			parentheses--
		}
	}

	empty := false

	for _, option := range options {
		endPos := walk(option, world, pos)
		if endPos == pos {
			empty = true
		} else {
			walk(remainder, world, endPos)
		}
	}

	if empty {
		return walk(remainder, world, pos)
	} else {
		return pos
	}
}

var deltas = []Position{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

func main() {
	input := readFile("input.txt")
	if !strings.HasPrefix(input, "^") || !strings.HasSuffix(input, "$") {
		panic(input)
	}
	input = input[1 : len(input)-1]

	world := make([]byte, TILE_SIZE*TILE_SIZE)
	walk(input, world, Position{x: TILE_SIZE / 2, y: TILE_SIZE / 2})

	open := make([]Destination, 0)
	reached := make([]bool, TILE_SIZE*TILE_SIZE)

	pos := Position{x: TILE_SIZE / 2, y: TILE_SIZE / 2}
	open = append(open, Destination{position: pos})
	reached[pos.Index()] = true

	longest := 0
	count := 0
	for len(open) != 0 {
		current := open[0]
		open = open[1:]

		for _, delta := range deltas {
			door := current.position.Plus(delta)
			newPos := current.position.Plus(delta.Times(2))
			if world[door.Index()] == DOOR && !reached[newPos.Index()] {
				open = append(open, Destination{
					position: newPos,
					distance: current.distance + 1,
				})
				reached[newPos.Index()] = true
			}
		}

		if current.distance > longest {
			longest = current.distance
		}
		if current.distance >= 1000 {
			count++
		}
	}

	{
		fmt.Println("--- Part One ---")
		fmt.Println(longest)
	}

	{
		fmt.Println("--- Part Two ---")
		fmt.Println(count)
	}
}

func readFile(filename string) string {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(bytes))
}
