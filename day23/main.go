package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
)

var Zero = Point{0, 0, 0}

type Point struct {
	x, y, z int
}

func (p Point) Distance(o Point) int {
	return abs(p.x-o.x) + abs(p.y-o.y) + abs(p.z-o.z)
}

type Nanobot struct {
	position Point
	radius   int
}

type Box struct {
	position Point
	size     int
	inRange  int
}

type PriorityQueue []*Box

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	// The box with the most nanobots in range has the highest priority.
	return pq[i].inRange > pq[j].inRange
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	box := x.(*Box)
	*pq = append(*pq, box)
}

func (pq *PriorityQueue) Pop() interface{} {
	len := len(*pq)
	box := (*pq)[len-1]
	*pq = (*pq)[:len-1]
	return box
}

func main() {
	lines := readLines("input.txt")

	var bots []Nanobot
	{
		botRegex := regexp.MustCompile(`pos=<(-?\d+),(-?\d+),(-?\d+)>, r=(-?\d+)`)
		for _, line := range lines {
			result := botRegex.FindStringSubmatch(line)
			if result == nil {
				panic(line)
			}
			bots = append(bots, Nanobot{
				position: Point{
					x: toInt(result[1]),
					y: toInt(result[2]),
					z: toInt(result[3]),
				},
				radius: toInt(result[4]),
			})
		}
	}

	{
		fmt.Println("--- Part One ---")

		var strongest Nanobot
		for _, bot := range bots {
			if bot.radius > strongest.radius {
				strongest = bot
			}
		}

		inRange := 0
		for _, bot := range bots {
			if strongest.position.Distance(bot.position) <= strongest.radius {
				inRange++
			}
		}

		fmt.Println(inRange)
	}

	{
		fmt.Println("--- Part Two ---")

		min := Point{math.MaxInt32, math.MaxInt32, math.MaxInt32}
		max := Point{math.MinInt32, math.MinInt32, math.MinInt32}

		for _, bot := range bots {
			if bot.position.x-bot.radius < min.x {
				min.x = bot.position.x - bot.radius
			}
			if bot.position.x+bot.radius > max.x {
				max.x = bot.position.x + bot.radius
			}
			if bot.position.y-bot.radius < min.y {
				min.y = bot.position.y - bot.radius
			}
			if bot.position.y+bot.radius > max.y {
				max.y = bot.position.y + bot.radius
			}
			if bot.position.z-bot.radius < min.z {
				min.z = bot.position.z - bot.radius
			}
			if bot.position.z+bot.radius > max.z {
				max.z = bot.position.z + bot.radius
			}
		}

		var initial Box
		initial.position = min
		initial.size = math.MinInt32
		initial.inRange = len(bots)

		if width := max.x - min.x + 1; width > initial.size {
			initial.size = width
		}
		if height := max.y - min.y + 1; height > initial.size {
			initial.size = height
		}
		if depth := max.z - min.z + 1; depth > initial.size {
			initial.size = depth
		}

		queue := make(PriorityQueue, 0)
		heap.Push(&queue, &initial)

		var best Point
		var bestInRange int

		for queue.Len() != 0 {
			box := heap.Pop(&queue).(*Box)

			if box.inRange < bestInRange {
				break
			}

			if box.size == 1 {
				if box.inRange > bestInRange {
					best = box.position
					bestInRange = box.inRange
				} else if box.inRange == bestInRange && box.position.Distance(Zero) < best.Distance(Zero) {
					best = box.position
				}
				continue
			}

			halfDown := box.size / 2      // rounds down for uneven numbers
			halfUp := box.size - halfDown // rounds up for uneven numbers

			for pdx := 0; pdx <= 1; pdx++ {
				for pdy := 0; pdy <= 1; pdy++ {
					for pdz := 0; pdz <= 1; pdz++ {
						var next Box
						next.position.x = box.position.x + pdx*halfDown
						next.position.y = box.position.y + pdy*halfDown
						next.position.z = box.position.z + pdz*halfDown
						next.size = halfUp
						next.inRange = 0

						for _, bot := range bots {
							var dx, dy, dz int
							if bot.position.x < next.position.x {
								dx = next.position.x - bot.position.x
							} else if bot.position.x < next.position.x+next.size {
								dx = 0
							} else {
								dx = bot.position.x - (next.position.x + next.size - 1)
							}
							if bot.position.y < next.position.y {
								dy = next.position.y - bot.position.y
							} else if bot.position.y < next.position.y+next.size {
								dy = 0
							} else {
								dy = bot.position.y - (next.position.y + next.size - 1)
							}
							if bot.position.z < next.position.z {
								dz = next.position.z - bot.position.z
							} else if bot.position.z < next.position.z+next.size {
								dz = 0
							} else {
								dz = bot.position.z - (next.position.z + next.size - 1)
							}
							if dx+dy+dz <= bot.radius {
								next.inRange++
							}
						}

						heap.Push(&queue, &next)
					}
				}
			}

		}

		fmt.Println(best.Distance(Zero))
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

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
