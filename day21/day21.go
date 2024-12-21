package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
)

type Position struct {
	row int
	col int
}

type PQItem[X any] struct {
	val  X
	cost int
}

type MoveState struct {
	pos  Position
	i    int
	prev *MoveState
}

type Memo struct {
	seq   string
	layer int
}

const A = 10
const EMPTY = -1

var TEMP = Position{row: -10, col: -10}
var DPADEMPTY = Position{row: -1, col: -1}
var UP = Position{row: -1, col: 0}
var DPADA = Position{row: 0, col: 0}
var LEFT = Position{row: 0, col: -1}
var DOWN = Position{row: 1, col: 0}
var RIGHT = Position{row: 0, col: 1}

var NUMPAD = [4][3]int{{7, 8, 9}, {4, 5, 6}, {1, 2, 3}, {EMPTY, 0, A}}
var DPAD = [2][3]Position{{DPADEMPTY, UP, DPADA}, {LEFT, DOWN, RIGHT}}
var DIRECTIONS = [4]Position{UP, DOWN, LEFT, RIGHT}

var keyToKey = make(map[Position]map[Position][][]Position)

func initKeyMap(m map[Position]map[Position][][]Position) {
	m[UP] = make(map[Position][][]Position)
	m[DOWN] = make(map[Position][][]Position)
	m[LEFT] = make(map[Position][][]Position)
	m[RIGHT] = make(map[Position][][]Position)
	m[DPADA] = make(map[Position][][]Position)

	// Gives most effective way to move from one D key to another
	m[UP][UP] = [][]Position{{DPADA}}
	m[UP][DOWN] = [][]Position{{DOWN, DPADA}}
	m[UP][LEFT] = [][]Position{{DOWN, LEFT, DPADA}}
	m[UP][RIGHT] = [][]Position{{RIGHT, DOWN, DPADA}, {DOWN, RIGHT, DPADA}}
	m[UP][DPADA] = [][]Position{{RIGHT, DPADA}}

	m[DOWN][UP] = [][]Position{{UP, DPADA}}
	m[DOWN][DOWN] = [][]Position{{DPADA}}
	m[DOWN][LEFT] = [][]Position{{LEFT, DPADA}}
	m[DOWN][RIGHT] = [][]Position{{RIGHT, DPADA}}
	m[DOWN][DPADA] = [][]Position{{UP, RIGHT, DPADA}, {RIGHT, UP, DPADA}}

	m[LEFT][UP] = [][]Position{{RIGHT, UP, DPADA}}
	m[LEFT][DOWN] = [][]Position{{RIGHT, DPADA}}
	m[LEFT][LEFT] = [][]Position{{DPADA}}
	m[LEFT][RIGHT] = [][]Position{{RIGHT, RIGHT, DPADA}}
	m[LEFT][DPADA] = [][]Position{{RIGHT, RIGHT, UP, DPADA}}

	m[RIGHT][UP] = [][]Position{{UP, LEFT, DPADA}, {LEFT, UP, DPADA}}
	m[RIGHT][DOWN] = [][]Position{{LEFT, DPADA}}
	m[RIGHT][LEFT] = [][]Position{{LEFT, LEFT, DPADA}}
	m[RIGHT][RIGHT] = [][]Position{{DPADA}}
	m[RIGHT][DPADA] = [][]Position{{UP, DPADA}}

	m[DPADA][UP] = [][]Position{{LEFT, DPADA}}
	m[DPADA][DOWN] = [][]Position{{DOWN, LEFT, DPADA}, {LEFT, DOWN, DPADA}}
	m[DPADA][LEFT] = [][]Position{{DOWN, LEFT, LEFT, DPADA}}
	m[DPADA][RIGHT] = [][]Position{{DOWN, DPADA}}
	m[DPADA][DPADA] = [][]Position{{DPADA}}
}

func getNumpadVal(pos Position) int {
	if pos.row < 0 || pos.row >= len(NUMPAD) || pos.col < 0 || pos.col >= len(NUMPAD[pos.row]) {
		return EMPTY
	}
	return NUMPAD[pos.row][pos.col]
}

func getDpadVal(pos Position) Position {
	if pos.row < 0 || pos.row >= len(DPAD) || pos.col < 0 || pos.col >= len(DPAD[pos.row]) {
		return DPADEMPTY
	}
	return DPAD[pos.row][pos.col]
}

func readInput(inputFile string) [][]int {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	sc := bufio.NewScanner(file)
	keySeqs := make([][]int, 0)
	for sc.Scan() {
		r := []rune(sc.Text())
		k := make([]int, len(r))
		for i := range r {
			if r[i] == rune('A') {
				k[i] = 10
				continue
			}
			k[i] = int(r[i] - rune('0'))
		}
		keySeqs = append(keySeqs, k)
	}
	return keySeqs
}

func numpadSeqs(keySeq []int) [][]Position {
	queue := []MoveState{{pos: Position{row: 3, col: 2}, i: 0, prev: nil}}
	seqs := make([][]Position, 0)
	shouldBreak := false
	targetI := 0
	for len(queue) > 0 && !shouldBreak {
		newQueue := make([]MoveState, 0)
		for _, ms := range queue {
			if ms.i == len(keySeq) {
				shouldBreak = true
				seq := make([]Position, 0)
				for ms.prev != nil {
					prev := ms.prev
					seq = append(seq, Position{row: ms.pos.row - prev.pos.row, col: ms.pos.col - prev.pos.col})
					ms = *prev
				}
				slices.Reverse(seq)
				seqs = append(seqs, seq)
				continue
			}
			numpadVal := getNumpadVal(ms.pos)
			if numpadVal == EMPTY {
				// Skip if out of bounds
				continue
			}
			if keySeq[ms.i] == numpadVal {
				if ms.i == targetI {
					newQueue = make([]MoveState, 0)
					targetI += 1
				}
				newQueue = append(newQueue, MoveState{pos: ms.pos, i: ms.i + 1, prev: &ms})
				continue
			}

			if ms.i < targetI {
				// Throw redundant states
				continue
			}
			for _, dir := range DIRECTIONS {
				newPos := Position{row: ms.pos.row + dir.row, col: ms.pos.col + dir.col}
				newQueue = append(newQueue, MoveState{pos: newPos, i: ms.i, prev: &ms})
			}
		}
		queue = newQueue
	}
	return seqs
}

func getComplexities(keySeq []int, memoMap map[Memo]int, numDpadRobots int) int {
	v := 0
	for _, k := range keySeq {
		if k == A {
			break
		}
		v *= 10
		v += k
	}

	s := numpadSeqs(keySeq)
	minMoves := -1
	for _, seq := range s {
		m := recurseDpadMove(memoMap, seq, numDpadRobots)
		if minMoves == -1 || m < minMoves {
			minMoves = m
		}
	}
	fmt.Println("V:", v, "Moves:", minMoves)

	return v * minMoves
}

// Get shortest move from 'A' to 'A'.
func recurseDpadMove(memoMap map[Memo]int, s []Position, layer int) int {
	if layer == 0 {
		return len(s)
	}

	sSplit := splitS(s)
	sum := 0
	for _, minS := range sSplit {
		key := Memo{seq: sKey(minS), layer: layer}
		minMoves, ok := memoMap[key]
		if ok {
			sum += minMoves
			continue
		}

		res := make([][]Position, 0)
		for i := 0; i < len(minS); i++ {
			if i == 0 {
				res = keyToKey[DPADA][minS[i]]
				continue
			}
			newRes := make([][]Position, 0)
			for _, newSeq := range res {
				for _, newMoves := range keyToKey[minS[i-1]][minS[i]] {
					combined := append([]Position{}, newSeq...)
					combined = append(combined, newMoves...)
					newRes = append(newRes, combined)
				}
			}
			res = newRes
		}

		minMoves = -1
		for _, resultSeq := range res {
			v := recurseDpadMove(memoMap, resultSeq, layer-1)
			if minMoves == -1 || v < minMoves {
				minMoves = v
			}
		}
		memoMap[key] = minMoves
		sum += minMoves
	}
	return sum
}

func splitS(s []Position) [][]Position {
	sSplit := make([][]Position, 0)
	nextS := make([]Position, 0)
	for _, pos := range s {
		nextS = append(nextS, pos)
		if pos == DPADA {
			sSplit = append(sSplit, nextS)
			nextS = make([]Position, 0)
		}
	}
	return sSplit
}

func sKey(s []Position) string {
	key := ""
	for _, move := range s {
		switch move {
		case UP:
			key += "^"
			break
		case DOWN:
			key += "v"
			break
		case LEFT:
			key += "<"
			break
		case RIGHT:
			key += ">"
			break
		case DPADA:
			key += "A"
			break
		}
	}
	return key
}

func sumComplexities(keySeqs [][]int, numDpadRobots int) int {
	memoMap := make(map[Memo]int)
	sum := 0
	for i := range keySeqs {
		sum += getComplexities(keySeqs[i], memoMap, numDpadRobots)
	}
	return sum
}

func main() {
	initKeyMap(keyToKey)
	fmt.Println("Test Input")
	keySeqs := readInput("test-input.txt")
	fmt.Println("Sum Complexities (2 Dpad Robots):", sumComplexities(keySeqs, 2))
	fmt.Println("Sum Complexities (25 Dpad Robots):", sumComplexities(keySeqs, 25))

	fmt.Println("")

	fmt.Println("Puzzle Input")
	keySeqs = readInput("input.txt")
	fmt.Println("Sum Complexities (2 Dpad Robots):", sumComplexities(keySeqs, 2))
	fmt.Println("Sum Complexities (25 Dpad Robots):", sumComplexities(keySeqs, 25))
}
