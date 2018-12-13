package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const (
	UP = iota
	RIGHT
	DOWN
	LEFT
)

const (
	TURN_LEFT = iota
	TURN_STRAIGHT
	TURN_RIGHT

	TURN_COUNT
)

type Cart struct {
	x, y int
	dir  int
	turn int
}

func main() {
	world := readLines("input.txt")

	var carts []*Cart
	for y, line := range world {
		for x, tile := range line {
			var dir int

			switch tile {
			case '^':
				dir = UP
			case 'v':
				dir = DOWN
			case '<':
				dir = LEFT
			case '>':
				dir = RIGHT
			default:
				continue
			}

			carts = append(carts, &Cart{
				x:   x,
				y:   y,
				dir: dir,
			})
		}
	}

	hadFirstCrash := false
	for {
		// Carts move based on their position, from top to bottom, then from
		// left to right, so sort the carts in the order they move.
		sort.Slice(carts, func(i, j int) bool {
			if carts[i].y != carts[j].y {
				return carts[i].y < carts[j].y
			}
			return carts[i].x < carts[j].x
		})

		for cartIndex := 0; cartIndex < len(carts); cartIndex++ {
			cart := carts[cartIndex]

			// Update position.
			switch cart.dir {
			case UP:
				cart.y--
			case DOWN:
				cart.y++
			case LEFT:
				cart.x--
			case RIGHT:
				cart.x++
			}

			// Check for turns or intersections at the new position.
			switch world[cart.y][cart.x] {
			case '/':
				switch cart.dir {
				case UP:
					cart.dir = RIGHT
				case DOWN:
					cart.dir = LEFT
				case LEFT:
					cart.dir = DOWN
				case RIGHT:
					cart.dir = UP
				}
			case '\\':
				switch cart.dir {
				case UP:
					cart.dir = LEFT
				case DOWN:
					cart.dir = RIGHT
				case LEFT:
					cart.dir = UP
				case RIGHT:
					cart.dir = DOWN
				}
			case '+':
				// Handle intersection.
				switch cart.turn {
				case TURN_LEFT:
					cart.dir = (cart.dir + 3) % 4
				case TURN_STRAIGHT:
					cart.dir = cart.dir
				case TURN_RIGHT:
					cart.dir = (cart.dir + 1) % 4
				}
				cart.turn = (cart.turn + 1) % TURN_COUNT
			}

			// Check for collisions.
			for otherIndex, other := range carts {
				if cartIndex != otherIndex && cart.y == other.y && cart.x == other.x {
					if !hadFirstCrash {
						hadFirstCrash = true
						fmt.Println("--- Part One ---")
						fmt.Printf("%d,%d\n", cart.x, cart.y)
					}

					// Remove other cart.
					for i := otherIndex; i < len(carts)-1; i++ {
						carts[i] = carts[i+1]
					}
					carts = carts[:len(carts)-1]

					// If the other cart was before the current cart in the
					// list, the current cart has been moved down by one.
					if cartIndex > otherIndex {
						cartIndex--
					}

					// Remove current cart.
					for i := cartIndex; i < len(carts)-1; i++ {
						carts[i] = carts[i+1]
					}
					carts = carts[:len(carts)-1]

					// Move index to the "before the next cart" position,
					// so that we do not skip a cart in the cart loop.
					cartIndex--

					// There can only be one collision, as we check after each
					// cart moves.
					break
				}
			}
		}

		if len(carts) == 1 {
			cart := carts[0]
			fmt.Println("--- Part Two ---")
			fmt.Printf("%d,%d\n", cart.x, cart.y)
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
