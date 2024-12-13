package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Position struct {
	x int
	y int
}

type Machine struct {
	aMove Position
	bMove Position
	prize Position
}

type ButtonCount struct {
	a byte
	b byte
}

type ClawState struct {
	tokens int
	bCount ButtonCount
	pos    Position
}

// For Priority Queue Implementation
type PQ []*ClawState

func (pq PQ) Len() int {
	return len(pq)
}

func (pq PQ) Less(i, j int) bool {
	return pq[i].tokens < pq[j].tokens
}

func (pq PQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PQ) Push(x interface{}) {
	*pq = append(*pq, x.(*ClawState))
}

func (pq *PQ) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

func readInput(inputFile string) []Machine {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	machines := make([]Machine, 0)
	sc := bufio.NewScanner(file)
	lCount := 1
	machine := Machine{}
	for sc.Scan() {
		switch lCount {
		case 1:
			tokens := strings.Split(sc.Text(), ": ")
			tokens = strings.Split(tokens[1], ", ")
			x, err := strconv.Atoi(strings.Split(tokens[0], "+")[1])
			if err != nil {
				log.Fatal(err)
			}
			y, err := strconv.Atoi(strings.Split(tokens[1], "+")[1])
			if err != nil {
				log.Fatal(err)
			}
			machine.aMove = Position{x: x, y: y}
			break
		case 2:
			tokens := strings.Split(sc.Text(), ": ")
			tokens = strings.Split(tokens[1], ", ")
			x, err := strconv.Atoi(strings.Split(tokens[0], "+")[1])
			if err != nil {
				log.Fatal(err)
			}
			y, err := strconv.Atoi(strings.Split(tokens[1], "+")[1])
			if err != nil {
				log.Fatal(err)
			}
			machine.bMove = Position{x: x, y: y}
			break
		case 3:
			tokens := strings.Split(sc.Text(), ": ")
			tokens = strings.Split(tokens[1], ", ")
			x, err := strconv.Atoi(strings.Split(tokens[0], "=")[1])
			if err != nil {
				log.Fatal(err)
			}
			y, err := strconv.Atoi(strings.Split(tokens[1], "=")[1])
			if err != nil {
				log.Fatal(err)
			}
			machine.prize = Position{x: x, y: y}
			machines = append(machines, machine)
			machine = Machine{}
			break
		default:
			lCount = 0
			break
		}
		lCount += 1
	}
	return machines
}

func tokensNeeded(m Machine) int {
	posMap := make(map[Position]map[ButtonCount]bool)
	pq := make(PQ, 0)
	heap.Init(&pq)
	heap.Push(&pq, &ClawState{tokens: 0, bCount: ButtonCount{a: 0, b: 0}, pos: Position{x: 0, y: 0}})
	for {
		if pq.Len() == 0 {
			break
		}
		state := *heap.Pop(&pq).(*ClawState)
		_, ok := posMap[state.pos]
		if !ok {
			posMap[state.pos] = make(map[ButtonCount]bool)
		}
		_, ok = posMap[state.pos][state.bCount]
		if ok {
			continue
		}
		if state.pos.x == m.prize.x && state.pos.y == m.prize.y {
			return state.tokens
		}
		posMap[state.pos][state.bCount] = true
		if state.bCount.a > 100 || state.bCount.b > 100 || state.pos.x > m.prize.x || state.pos.y > m.prize.y {
			continue
		}
		// Press Button A
		heap.Push(&pq, &ClawState{tokens: state.tokens + 3, bCount: ButtonCount{a: state.bCount.a + 1, b: state.bCount.b}, pos: Position{x: state.pos.x + m.aMove.x, y: state.pos.y + m.aMove.y}})
		// Press Button B
		heap.Push(&pq, &ClawState{tokens: state.tokens + 1, bCount: ButtonCount{a: state.bCount.a, b: state.bCount.b + 1}, pos: Position{x: state.pos.x + m.bMove.x, y: state.pos.y + m.bMove.y}})
	}
	return 0
}

func minTokens(machines []Machine) int {
	coins := 0
	for _, m := range machines {
		coins += tokensNeeded(m)
	}
	return coins
}

func fixPrizePos(machines []Machine) []Machine {
	fixedMachines := make([]Machine, len(machines))
	for i, machine := range machines {
		fixedMachines[i] = machine
		fixedMachines[i].prize.x += 10000000000000
		fixedMachines[i].prize.y += 10000000000000
	}
	return fixedMachines
}

func fastTokensNeeded(m Machine) int {
	diffConst := (m.bMove.x * m.prize.y) - (m.bMove.y * m.prize.x)
	diffA := (m.bMove.y * -m.aMove.x) - (m.bMove.x * -m.aMove.y)
	a := diffConst / diffA
	b := (m.prize.x - a*m.aMove.x) / m.bMove.x
	// fmt.Println("A:", a, ", B:", b)
	if a*m.aMove.y+b*m.bMove.y != m.prize.y {
		return 0
	}
	return 3*a + b
}

func minTokensFast(machines []Machine) int {
	coins := 0
	for _, m := range machines {
		coins += fastTokensNeeded(m)
	}
	return coins
}

func main() {
	fmt.Println("Test Input")
	machines := readInput("test-input.txt")
	fmt.Println("Minimum Coins Needed:", minTokensFast(machines))
	fixedMachines := fixPrizePos(machines)
	fmt.Println("Minimum Coins Needed:", minTokensFast(fixedMachines))

	fmt.Println("")

	fmt.Println("Puzzle Input")
	machines = readInput("input.txt")
	fmt.Println("Minimum Coins Needed:", minTokensFast(machines))
	fixedMachines = fixPrizePos(machines)
	fmt.Println("Minimum Coins Needed:", minTokensFast(fixedMachines))
}
