package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

var deltas = []Position{{0, -1}, {-1, 0}, {1, 0}, {0, 1}} // in reading order

type Position struct {
	x, y int
}

func (p Position) plus(o Position) Position {
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
	distance int
}

func (d Destination) Less(o Destination) bool {
	if d.distance != o.distance {
		return d.distance < o.distance
	}
	return d.position.Less(o.position)
}

type SearchRequest struct {
	start     Position
	positions []Position
	world     [][]byte
}

type SearchResponse struct {
	start       Position
	destination Destination
}

func searcher(width, height int, requests <-chan *SearchRequest, responses chan<- *SearchResponse) {
	open := make([]Destination, 0, width*height)

	reached := make([][]bool, height)
	for y := range reached {
		reached[y] = make([]bool, width)
	}

	for request := range requests {
		open = open[:1]
		open[0].position = request.start
		open[0].distance = 1

		destinations := make([]Destination, len(request.positions))
		for i, position := range request.positions {
			destinations[i].position = position
			destinations[i].distance = math.MaxInt32
		}

		for len(open) != 0 {
			var currentIndex int
			var current Destination
			current.distance = math.MaxInt32

			for index, destination := range open {
				if destination.distance <= current.distance {
					currentIndex = index
					current = destination
				}
			}

			open[currentIndex] = open[len(open)-1]
			open = open[:len(open)-1]

			if reached[current.position.y][current.position.x] {
				continue
			}

			reached[current.position.y][current.position.x] = true

			for _, delta := range deltas {
				p := current.position.plus(delta)
				if request.world[p.y][p.x] == '.' {
					open = append(open, Destination{
						position: p,
						distance: current.distance + 1,
					})
				}
			}

			for i, destination := range destinations {
				if current.position == destination.position && current.distance < destination.distance {
					destinations[i].distance = current.distance
				}
			}
		}

		var best Destination
		best.distance = math.MaxInt32

		for _, destination := range destinations {
			if destination.Less(best) {
				best = destination
			}
		}

		responses <- &SearchResponse{
			start:       request.start,
			destination: best,
		}

		// Reset reached buffer.
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				reached[y][x] = false
			}
		}
	}
}

type Unit struct {
	character byte
	position  Position
	hitPoints int
}

func main() {
	input := readLines("input.txt")

	height, width := len(input), len(input[0])

	requests := make(chan *SearchRequest, 4)
	responses := make(chan *SearchResponse, 4)

	for i := 0; i < 4; i++ {
		go searcher(width, height, requests, responses)
	}

	for attackPower := 3; ; attackPower++ {
		world := make([][]byte, height)
		for y, line := range input {
			world[y] = []byte(line)
		}

		var units []*Unit
		for y, line := range world {
			for x, tile := range line {
				if tile == 'G' || tile == 'E' {
					units = append(units, &Unit{
						character: tile,
						position:  Position{x, y},
						hitPoints: 200,
					})
				}
			}
		}

		elvesBefore := 0
		for _, unit := range units {
			if unit.character == 'E' {
				elvesBefore++
			}
		}

		rounds := 0
	combat:
		for {
			// Units take their turns in reading order.
			sort.Slice(units, func(i, j int) bool {
				return units[i].position.Less(units[j].position)
			})

			for unitIndex := 0; unitIndex < len(units); unitIndex++ {
				unit := units[unitIndex]

				// Move!
				{
					var targets int
					var inRange []Position
					var currentlyInRange bool

				rangeSearch:
					for _, target := range units {
						if target.character != unit.character {
							targets++
							for _, delta := range deltas {
								p := target.position.plus(delta)
								if p == unit.position {
									currentlyInRange = true
									break rangeSearch
								}
								if world[p.y][p.x] == '.' {
									inRange = append(inRange, p)
								}
							}
						}
					}

					if targets == 0 {
						break combat
					}

					if !currentlyInRange && len(inRange) == 0 {
						continue // with the next unit
					}

					if !currentlyInRange {
						searchCount := 0

						for _, delta := range deltas {
							step := unit.position.plus(delta)

							if world[step.y][step.x] == '.' {
								requests <- &SearchRequest{
									start:     step,
									positions: inRange,
									world:     world,
								}
								searchCount++
							}
						}

						var bestStep Position
						var bestDestination Destination
						bestDestination.distance = math.MaxInt32

						for i := 0; i < searchCount; i++ {
							response := <-responses
							if response.destination.Less(bestDestination) {
								bestStep = response.start
								bestDestination = response.destination
							} else if response.destination == bestDestination && response.start.Less(bestStep) {
								bestStep = response.start
							}
						}

						if bestDestination.distance == math.MaxInt32 {
							continue // with the next unit
						}

						world[unit.position.y][unit.position.x] = '.'
						unit.position = bestStep
						world[unit.position.y][unit.position.x] = unit.character
					}
				}

				// Attack!
				{
					var target *Unit
					var targetIndex int

					for _, delta := range deltas {
						p := unit.position.plus(delta)
						for otherIndex, other := range units {
							if other.position == p && other.character != unit.character {
								if target == nil || other.hitPoints < target.hitPoints {
									target = other
									targetIndex = otherIndex
								}
							}
						}
					}

					if target == nil {
						continue // with the next unit
					}

					if unit.character == 'E' {
						target.hitPoints -= attackPower
					} else {
						target.hitPoints -= 3
					}

					if target.hitPoints <= 0 {
						world[target.position.y][target.position.x] = '.'
						for i := targetIndex; i < len(units)-1; i++ {
							units[i] = units[i+1]
						}
						units = units[:len(units)-1]
						if unitIndex > targetIndex {
							unitIndex--
						}
					}
				}
			}

			rounds++
		}

		elvesAfter := 0
		for _, unit := range units {
			if unit.character == 'E' {
				elvesAfter++
			}
		}

		sum := 0
		for _, unit := range units {
			sum += unit.hitPoints
		}
		outcome := sum * rounds

		if attackPower == 3 {
			fmt.Println("--- Part One ---")
			fmt.Println(outcome)
		}

		if elvesAfter == elvesBefore {
			fmt.Println("--- Part Two ---")
			fmt.Println(outcome)
			break
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
