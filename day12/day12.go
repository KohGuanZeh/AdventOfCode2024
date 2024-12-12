package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Position struct {
	row int
	col int
}

const UP byte = 1
const RIGHT byte = 2
const DOWN byte = 3
const LEFT byte = 4

func readInput(inputFile string) [][]rune {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	input := make([][]rune, 0)
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		r := []rune(sc.Text())
		input = append(input, r)
	}
	return input
}

func calculateCost(land [][]rune) (int, int) {
	isDone := make([][]bool, 0)
	for i := range land {
		isDone = append(isDone, make([]bool, len(land[i])))
	}

	cost := 0
	discountedCost := 0

	for i := range land {
		for j := range land[i] {
			if isDone[i][j] {
				continue
			}
			pos := Position{row: i, col: j}
			area, fences, longFences := getStatistics(land, pos, isDone)
			cost += area * fences
			discountedCost += area * longFences
		}
	}
	return cost, discountedCost
}

// Returns area, fences and longFences
func getStatistics(land [][]rune, pos Position, isDone [][]bool) (int, int, int) {
	targetRune := land[pos.row][pos.col]
	checkQueue := []Position{pos}
	perimeterMap := make(map[Position]map[byte]bool)
	area := 0
	fences := 0
	longFences := 0
	for len(checkQueue) > 0 {
		newQueue := make([]Position, 0)
		for _, checkPos := range checkQueue {
			if isDone[checkPos.row][checkPos.col] {
				continue
			}
			perimeterMap[checkPos] = make(map[byte]bool)
			isDone[checkPos.row][checkPos.col] = true
			area += 1

			up := Position{row: checkPos.row - 1, col: checkPos.col}
			right := Position{row: checkPos.row, col: checkPos.col + 1}
			down := Position{row: checkPos.row + 1, col: checkPos.col}
			left := Position{row: checkPos.row, col: checkPos.col - 1}

			if up.row < 0 || land[up.row][up.col] != targetRune {
				fences += 1
				exists, isBridge := hasExisitingFence(perimeterMap, UP, []Position{left, right})
				if !exists {
					longFences += 1
				} else if isBridge {
					longFences -= 1
				}
				perimeterMap[checkPos][UP] = true
			} else if !isDone[up.row][up.col] {
				newQueue = append(newQueue, up)
			}

			if right.col >= len(land[right.row]) || land[right.row][right.col] != targetRune {
				fences += 1
				exists, isBridge := hasExisitingFence(perimeterMap, RIGHT, []Position{up, down})
				if !exists {
					longFences += 1
				} else if isBridge {
					longFences -= 1
				}
				perimeterMap[checkPos][RIGHT] = true
			} else if !isDone[right.row][right.col] {
				newQueue = append(newQueue, right)
			}

			if down.row >= len(land) || land[down.row][down.col] != targetRune {
				fences += 1
				exists, isBridge := hasExisitingFence(perimeterMap, DOWN, []Position{left, right})
				if !exists {
					longFences += 1
				} else if isBridge {
					longFences -= 1
				}
				perimeterMap[checkPos][DOWN] = true
			} else if !isDone[down.row][down.col] {
				newQueue = append(newQueue, down)
			}

			if left.col < 0 || land[left.row][left.col] != targetRune {
				fences += 1
				exists, isBridge := hasExisitingFence(perimeterMap, LEFT, []Position{up, down})
				if !exists {
					longFences += 1
				} else if isBridge {
					longFences -= 1
				}
				perimeterMap[checkPos][LEFT] = true
			} else if !isDone[left.row][left.col] {
				newQueue = append(newQueue, left)
			}
		}
		checkQueue = newQueue
	}
	return area, fences, longFences
}

func hasExisitingFence(perimeterMap map[Position]map[byte]bool, targetDir byte, positions []Position) (bool, bool) {
	existingFence := false
	isBridge := true
	for _, pos := range positions {
		_, ok := perimeterMap[pos]
		if ok {
			_, ok = perimeterMap[pos][targetDir]
			if ok {
				existingFence = true
				isBridge = !isBridge
			}
		}
	}
	return existingFence, isBridge && existingFence
}

func main() {
	fmt.Println("Test Input")
	land := readInput("test-input.txt")
	cost, discountedCost := calculateCost(land)
	fmt.Println("Total Cost:", cost)
	fmt.Println("Total Cost With Discount:", discountedCost)

	fmt.Println("")

	fmt.Println("Test Input 2")
	land = readInput("test-input-2.txt")
	cost, discountedCost = calculateCost(land)
	fmt.Println("Total Cost:", cost)
	fmt.Println("Total Cost With Discount:", discountedCost)

	fmt.Println("")

	fmt.Println("Puzzle Input")
	land = readInput("input.txt")
	cost, discountedCost = calculateCost(land)
	fmt.Println("Total Cost:", cost)
	fmt.Println("Total Cost With Discount:", discountedCost)

}
