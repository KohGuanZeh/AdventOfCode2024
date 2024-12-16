package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"os"
)

const UP byte = 0
const DOWN byte = 1
const LEFT byte = 2
const RIGHT byte = 3

const FREE = rune('.')
const WALL = rune('#')
const START = rune('S')
const END = rune('E')

type MoveState struct {
	row int
	col int
	dir byte
}

type Move struct {
	state MoveState
	prev  *Move
	cost  int
}

// For Priority Queue Implementation
type PriorityQueue []*Move

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].cost < pq[j].cost
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Move))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

func readInput(inputFile string) [][]rune {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	inputMap := make([][]rune, 0)
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		inputMap = append(inputMap, []rune(sc.Text()))
	}
	return inputMap
}

// Returns shortest cost and number of best seats
func shortestPath(inputMap [][]rune) (int, int) {
	bestCost := -1
	var startState MoveState
	startFound := false
	for r := range inputMap {
		for c := range inputMap[r] {
			if inputMap[r][c] == START {
				startFound = true
				startState = MoveState{row: r, col: c, dir: RIGHT}
				break
			}
		}
		if startFound {
			break
		}
	}

	// Keep track of previous MoveStates to best path
	visitedCosts := make(map[MoveState]int)
	visitedChain := make(map[MoveState][]Move)
	endStates := make([]MoveState, 0)
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	heap.Push(&pq, &Move{state: startState, cost: 0, prev: nil})
	for {
		if len(pq) == 0 {
			log.Fatalln("Error. Did not make it to end of maze.")
		}
		move := *heap.Pop(&pq).(*Move)
		ms := move.state

		cost, ok := visitedCosts[ms]
		if ok {
			if cost == move.cost {
				prevMoves := visitedChain[ms]
				prevMoves = append(prevMoves, *move.prev)
				visitedChain[ms] = prevMoves
			}
			// If processed before, skip
			continue
		}

		visitedCosts[ms] = move.cost
		if move.prev != nil {
			visitedChain[ms] = []Move{*move.prev}
		} else {
			visitedChain[ms] = []Move{}
		}

		if bestCost > -1 && move.cost > bestCost {
			// No more best routes recorded
			break
		}

		if inputMap[ms.row][ms.col] == END {
			// Reached end at lowest cost
			bestCost = move.cost
			endStates = append(endStates, ms)
			continue
		}

		if inputMap[ms.row][ms.col] == WALL {
			// If hit wall, do not continue this iteration
			continue
		}

		var forward MoveState
		var rotClockwise MoveState
		var rotCounterClockwise MoveState
		switch move.state.dir {
		case UP:
			forward = MoveState{row: ms.row - 1, col: ms.col, dir: UP}
			rotClockwise = MoveState{row: ms.row, col: ms.col, dir: RIGHT}
			rotCounterClockwise = MoveState{row: ms.row, col: ms.col, dir: LEFT}
			break
		case DOWN:
			forward = MoveState{row: ms.row + 1, col: ms.col, dir: DOWN}
			rotClockwise = MoveState{row: ms.row, col: ms.col, dir: LEFT}
			rotCounterClockwise = MoveState{row: ms.row, col: ms.col, dir: RIGHT}
			break
		case RIGHT:
			forward = MoveState{row: ms.row, col: ms.col + 1, dir: RIGHT}
			rotClockwise = MoveState{row: ms.row, col: ms.col, dir: DOWN}
			rotCounterClockwise = MoveState{row: ms.row, col: ms.col, dir: UP}
			break
		case LEFT:
			forward = MoveState{row: ms.row, col: ms.col - 1, dir: LEFT}
			rotClockwise = MoveState{row: ms.row, col: ms.col, dir: UP}
			rotCounterClockwise = MoveState{row: ms.row, col: ms.col, dir: DOWN}
			break
		default:
			log.Fatalln("Invalid direction given.")
			break
		}
		heap.Push(&pq, &Move{state: forward, cost: move.cost + 1, prev: &move})
		heap.Push(&pq, &Move{state: rotClockwise, cost: move.cost + 1000, prev: &move})
		heap.Push(&pq, &Move{state: rotCounterClockwise, cost: move.cost + 1000, prev: &move})
	}

	bestSeats := 0
	seatMap := make(map[MoveState]bool)
	for {
		if len(endStates) == 0 {
			break
		}
		prevStates := make([]MoveState, 0)
		for _, ms := range endStates {
			prevMoves := visitedChain[ms]
			for _, move := range prevMoves {
				prevStates = append(prevStates, move.state)
			}
			seat := MoveState{row: ms.row, col: ms.col, dir: 255}
			_, ok := seatMap[seat]
			if ok {
				continue
			}
			// For debugging best seats
			inputMap[seat.row][seat.col] = rune('O')
			seatMap[seat] = true
			bestSeats += 1
		}
		endStates = prevStates
	}

	// Debug best seats
	// for _, r := range inputMap {
	// fmt.Println(string(r))
	// }

	return bestCost, bestSeats
}

func main() {
	fmt.Println("Test Input")
	inputMap := readInput("test-input.txt")
	bestCost, bestSeats := shortestPath(inputMap)
	fmt.Println("Shortest Cost to End:", bestCost)
	fmt.Println("Best Seats:", bestSeats)

	fmt.Println("")

	fmt.Println("Puzzle Input")
	inputMap = readInput("input.txt")
	bestCost, bestSeats = shortestPath(inputMap)
	fmt.Println("Shortest Cost to End:", bestCost)
	fmt.Println("Best Seats:", bestSeats)
}
