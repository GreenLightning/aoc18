package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Node struct {
	children []Node
	metadata []int
}

func readNode(input []int) (Node, []int) {
	var node Node
	node.children = make([]Node, input[0])
	node.metadata = make([]int, input[1])
	input = input[2:]
	for i := range node.children {
		node.children[i], input = readNode(input)
	}
	copy(node.metadata, input)
	input = input[len(node.metadata):]
	return node, input
}

func (node Node) sumMeta() int {
	result := 0
	for _, child := range node.children {
		result += child.sumMeta()
	}
	for _, meta := range node.metadata {
		result += meta
	}
	return result
}

func (node Node) value() int {
	if len(node.children) == 0 {
		return node.sumMeta()
	}
	result := 0
	for _, meta := range node.metadata {
		index := meta - 1
		if index >= 0 && index < len(node.children) {
			result += node.children[index].value()
		}
	}
	return result
}

func main() {
	input := readNumbers("input.txt")

	root, remainder := readNode(input)
	if len(remainder) != 0 {
		panic("input too big")
	}

	{
		fmt.Println("--- Part One ---")
		fmt.Println(root.sumMeta())
	}

	{
		fmt.Println("--- Part Two ---")
		fmt.Println(root.value())
	}
}

func readNumbers(filename string) []int {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)

	var numbers []int
	for scanner.Scan() {
		numbers = append(numbers, toInt(scanner.Text()))
	}
	return numbers
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
