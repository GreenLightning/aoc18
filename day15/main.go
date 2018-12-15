package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

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

type Unit struct {
	character byte
	position  Position
	hitPoints int
}

func main() {
	input := readLines("input.txt")

	height, width := len(input), len(input[0])

	deltas := []Position{{0, -1}, {-1, 0}, {1, 0}, {0, 1}} // in reading order

	reached := make([][]bool, height)
	for y := range reached {
		reached[y] = make([]bool, width)
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
			// fmt.Println("After round:", rounds)
			// for y, line := range world {
			// 	fmt.Print(string(line))
			// 	for x := range line {
			// 		for _, unit := range units {
			// 			if unit.position.y == y && unit.position.x == x {
			// 				fmt.Printf(" %c(%d)", unit.character, unit.hitPoints)
			// 			}
			// 		}
			// 	}
			// 	fmt.Println()
			// }

			// Units take their turns in reading order.
			sort.Slice(units, func(i, j int) bool {
				return units[i].position.Less(units[j].position)
			})

			for unitIndex := 0; unitIndex < len(units); unitIndex++ {
				unit := units[unitIndex]

				// Move!
				{
					var targets []*Unit
					for otherIndex, other := range units {
						if otherIndex != unitIndex && other.character != unit.character {
							targets = append(targets, other)
						}
					}

					if len(targets) == 0 {
						break combat
					}

					var inRange []Destination
					var currentlyInRange bool
				rangeSearch:
					for _, target := range targets {
						for _, delta := range deltas {
							p := target.position.plus(delta)
							if p == unit.position {
								currentlyInRange = true
								break rangeSearch
							}
							if world[p.y][p.x] == '.' {
								inRange = append(inRange, Destination{
									position: p,
									distance: math.MaxInt32,
								})
							}
						}
					}

					if !currentlyInRange && len(inRange) == 0 {
						continue // with the next unit
					}

					if !currentlyInRange {
						var open []Destination

						for _, delta := range deltas {
							p := unit.position.plus(delta)
							if world[p.y][p.x] == '.' {
								open = append(open, Destination{
									position: p,
									distance: 1,
								})
							}
						}

						for y := 0; y < height; y++ {
							for x := 0; x < width; x++ {
								reached[y][x] = false
							}
						}

						for len(open) != 0 {
							current := open[0]
							open = open[1:]

							if reached[current.position.y][current.position.x] {
								continue
							}

							reached[current.position.y][current.position.x] = true

							for _, delta := range deltas {
								p := current.position.plus(delta)
								if world[p.y][p.x] == '.' {
									open = append(open, Destination{
										position: p,
										distance: current.distance + 1,
									})
								}
							}

							sort.Slice(open, func(i, j int) bool {
								return open[i].distance < open[j].distance
							})

							for i, dest := range inRange {
								if current.position == dest.position && current.distance < dest.distance {
									inRange[i].distance = current.distance
								}
							}
						}

						sort.Slice(inRange, func(i, j int) bool {
							if inRange[i].distance != inRange[j].distance {
								return inRange[i].distance < inRange[j].distance
							}
							return inRange[i].position.Less(inRange[j].position)
						})

						destination := inRange[0]
						if destination.distance == math.MaxInt32 {
							continue // with the next unit
						}

						for _, delta := range deltas {
							step := unit.position.plus(delta)
							dest := destination
							dest.distance = math.MaxInt32

							if world[step.y][step.x] != '.' {
								continue
							}

							open = append(open, Destination{
								position: step,
								distance: 1,
							})

							for y := 0; y < height; y++ {
								for x := 0; x < width; x++ {
									reached[y][x] = false
								}
							}

							for len(open) != 0 {
								current := open[0]
								open = open[1:]

								if reached[current.position.y][current.position.x] {
									continue
								}

								reached[current.position.y][current.position.x] = true

								for _, delta := range deltas {
									p := current.position.plus(delta)
									if world[p.y][p.x] == '.' {
										open = append(open, Destination{
											position: p,
											distance: current.distance + 1,
										})
									}
								}

								sort.Slice(open, func(i, j int) bool {
									return open[i].distance < open[j].distance
								})

								if current.position == dest.position && current.distance < dest.distance {
									dest.distance = current.distance
								}
							}

							if dest.distance == destination.distance {
								world[unit.position.y][unit.position.x] = '.'
								unit.position = step
								world[unit.position.y][unit.position.x] = unit.character
								break
							}
						}
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
