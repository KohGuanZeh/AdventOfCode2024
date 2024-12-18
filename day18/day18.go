package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const FREE = rune('.')
const WALL = rune('#')

type Position struct {
	row int
	col int
}

type PositionState struct {
	curr Position
	prev *PositionState
}

func readInput(inputFile string) []Position {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	sc := bufio.NewScanner(file)
	positions := make([]Position, 0)
	for sc.Scan() {
		tokens := strings.Split(sc.Text(), ",")
		row, err := strconv.Atoi(tokens[0])
		if err != nil {
			log.Fatalln(err)
		}
		col, err := strconv.Atoi(tokens[1])
		if err != nil {
			log.Fatalln(err)
		}
		positions = append(positions, Position{row: row, col: col})
	}
	return positions
}

func traverse(size int, positions []Position) int {
	mem := make([][]rune, size)
	for i := range len(mem) {
		mem[i] = make([]rune, size)
		for j := range len(mem[i]) {
			mem[i][j] = FREE
		}
	}

	for _, pos := range positions {
		mem[pos.row][pos.col] = WALL
	}

	maxCoord := size - 1

	queue := []Position{{row: 0, col: 0}}
	steps := 0
	for {
		newQueue := make([]Position, 0)
		if len(queue) == 0 {
			break
		}
		for _, pos := range queue {

			if pos.row == maxCoord && pos.col == maxCoord {
				return steps
			} else if pos.row < 0 || pos.row > maxCoord || pos.col < 0 || pos.col > maxCoord {
				continue
			}

			if mem[pos.row][pos.col] == WALL {
				continue
			}

			mem[pos.row][pos.col] = WALL
			newQueue = append(newQueue, Position{row: pos.row - 1, col: pos.col})
			newQueue = append(newQueue, Position{row: pos.row + 1, col: pos.col})
			newQueue = append(newQueue, Position{row: pos.row, col: pos.col - 1})
			newQueue = append(newQueue, Position{row: pos.row, col: pos.col + 1})
		}
		steps += 1
		queue = newQueue
	}
	return -1
}

func shortestPath(mem [][]rune) map[Position]bool {
	maxCoord := len(mem) - 1
	donePos := make(map[Position]bool)
	pathMap := make(map[Position]bool)
	queue := []PositionState{{curr: Position{row: 0, col: 0}, prev: nil}}
	for {
		if len(queue) == 0 {
			break
		}
		newQueue := make([]PositionState, 0)
		for _, state := range queue {
			pos := state.curr
			if pos.row == maxCoord && pos.col == maxCoord {
				for state.prev != nil {
					pathMap[state.curr] = true
					state = *state.prev
				}
				break
			}

			_, ok := donePos[pos]
			if ok {
				continue
			}
			donePos[pos] = true

			if pos.col < 0 || pos.row < 0 || pos.col > maxCoord || pos.row > maxCoord {
				continue
			}

			if mem[pos.row][pos.col] == WALL {
				continue
			}

			up := Position{row: pos.row - 1, col: pos.col}
			down := Position{row: pos.row + 1, col: pos.col}
			left := Position{row: pos.row, col: pos.col - 1}
			right := Position{row: pos.row, col: pos.col + 1}
			newQueue = append(newQueue, PositionState{curr: up, prev: &state})
			newQueue = append(newQueue, PositionState{curr: down, prev: &state})
			newQueue = append(newQueue, PositionState{curr: left, prev: &state})
			newQueue = append(newQueue, PositionState{curr: right, prev: &state})
		}
		queue = newQueue
	}
	return pathMap
}

func minToBlock(size int, positions []Position) (int, int) {
	mem := make([][]rune, size)
	for i := range len(mem) {
		mem[i] = make([]rune, size)
		for j := range len(mem[i]) {
			mem[i][j] = FREE
		}
	}

	posIndex := 0
	for {
		path := shortestPath(mem)
		// fmt.Println("Min Steps after", posIndex, ":", len(path))
		if len(path) == 0 {
			break
		}
		for i := posIndex; i < len(positions); i++ {
			pos := positions[i]
			mem[pos.row][pos.col] = WALL
			_, ok := path[positions[i]]
			if ok {
				posIndex = i + 1
				break
			}
			if i == len(positions)-1 {
				return -1, -1
			}
		}
	}
	posIndex -= 1
	return positions[posIndex].row, positions[posIndex].col
}

func main() {
	fmt.Println("Test Input")
	positions := readInput("test-input.txt")
	fmt.Println("Min Steps Needed:", traverse(7, positions[:12]))
	r, c := minToBlock(7, positions)
	fmt.Println("Min Bytes to Block:", r, ",", c)

	fmt.Println("")

	fmt.Println("Puzzle Input")
	positions = readInput("input.txt")
	fmt.Println("Min Steps Needed:", traverse(71, positions[:1024]))
	r, c = minToBlock(71, positions)
	fmt.Println("Min Bytes to Block:", r, ",", c)
}
