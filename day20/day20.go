package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const FREE = rune('.')
const WALL = rune('#')
const START = rune('S')
const END = rune('E')

type Position struct {
	row int
	col int
}

func readInput(inputFile string) [][]rune {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	sc := bufio.NewScanner(file)
	inputMap := make([][]rune, 0)
	for sc.Scan() {
		inputMap = append(inputMap, []rune(sc.Text()))
	}
	return inputMap
}

func findCheats(inputMap [][]rune, timeSaved int, cheatTime int) int {
	startPos := Position{row: -1, col: -1}
	endPos := Position{row: -1, col: -1}
	for i := range inputMap {
		for j := range inputMap[i] {
			if inputMap[i][j] == START {
				startPos = Position{row: i, col: j}
			}
			if inputMap[i][j] == END {
				endPos = Position{row: i, col: j}
			}
		}
		if (endPos.row > -1 || endPos.col > -1) && (startPos.row > -1 || startPos.col > -1) {
			break
		}
	}

	timeTaken := 0
	timeMap := make(map[Position]int)
	queue := []Position{endPos}
	for {
		if len(queue) == 0 {
			break
		}
		newQueue := make([]Position, 0)
		for _, pos := range queue {
			if pos.row < 0 || pos.row >= len(inputMap) || pos.col < 0 || pos.col >= len(inputMap[pos.row]) {
				continue
			}

			_, ok := timeMap[pos]
			if ok {
				continue
			}

			if inputMap[pos.row][pos.col] == WALL {
				continue
			}

			timeMap[pos] = timeTaken
			newQueue = append(newQueue, Position{row: pos.row - 1, col: pos.col})
			newQueue = append(newQueue, Position{row: pos.row + 1, col: pos.col})
			newQueue = append(newQueue, Position{row: pos.row, col: pos.col - 1})
			newQueue = append(newQueue, Position{row: pos.row, col: pos.col + 1})
		}
		queue = newQueue
		timeTaken += 1
	}
	targetTime := timeMap[startPos] - timeSaved
	return navWithCheat(inputMap, timeMap, startPos, targetTime, cheatTime)
}

func navWithCheat(inputMap [][]rune, timeMap map[Position]int, startPos Position, targetTime int, cheatTime int) int {
	type CheatKey struct {
		start Position
		end   Position
	}
	type State struct {
		pos       Position
		cheatPos  Position
		cheatTime int
	}
	DEFCHEATPOS := Position{row: -1, col: -1}
	directions := []Position{{row: -1, col: 0}, {row: 1, col: 0}, {row: 0, col: -1}, {row: 0, col: 1}}
	queue := []State{{pos: startPos, cheatPos: DEFCHEATPOS, cheatTime: cheatTime}}
	done := make(map[CheatKey]bool)
	cheatMap := make(map[CheatKey]bool)
	for len(queue) > 0 && targetTime >= 0 {
		new := make([]State, 0)
		for _, state := range queue {
			pos := state.pos

			key := CheatKey{start: state.cheatPos, end: pos}
			_, ok := done[key]
			if ok {
				continue
			}
			done[key] = true

			if state.cheatTime < cheatTime {
				time, ok := timeMap[pos]
				if ok && time <= targetTime {
					cheatMap[key] = true
				}
				if state.cheatTime <= 0 {
					continue
				}
				for _, dir := range directions {
					nextPos := Position{row: pos.row + dir.row, col: pos.col + dir.col}
					if nextPos.row < 0 || nextPos.col < 0 || nextPos.row >= len(inputMap) || nextPos.col >= len(inputMap[nextPos.row]) {
						continue
					}

					new = append(new, State{pos: nextPos, cheatPos: state.cheatPos, cheatTime: state.cheatTime - 1})
				}
				continue
			}

			for _, dir := range directions {
				nextPos := Position{row: pos.row + dir.row, col: pos.col + dir.col}
				if nextPos.row < 0 || nextPos.col < 0 || nextPos.row >= len(inputMap) || nextPos.col >= len(inputMap[nextPos.row]) {
					continue
				}

				if inputMap[nextPos.row][nextPos.col] == WALL {
					// If next one is a wall, activate cheat
					new = append(new, State{pos: nextPos, cheatPos: pos, cheatTime: cheatTime - 1})
				} else {
					new = append(new, State{pos: nextPos, cheatPos: DEFCHEATPOS, cheatTime: state.cheatTime})
					new = append(new, State{pos: nextPos, cheatPos: pos, cheatTime: state.cheatTime - 1})
				}
			}
		}
		targetTime -= 1
		queue = new
	}
	// for key := range cheatMap {
	// fmt.Println(key)
	// }
	return len(cheatMap)
}

func main() {
	fmt.Println("Test Input")
	inputMap := readInput("test-input.txt")
	timeSaved, cheatTime := 10, 2
	fmt.Println("Cheats saving >=", timeSaved, "ps (", cheatTime, "ps time limit ):", findCheats(inputMap, timeSaved, cheatTime))
	timeSaved, cheatTime = 68, 20
	fmt.Println("Cheats saving >=", timeSaved, "ps (", cheatTime, "ps time limit ):", findCheats(inputMap, timeSaved, cheatTime))

	fmt.Println("")

	fmt.Println("Puzzle Input")
	inputMap = readInput("input.txt")
	timeSaved, cheatTime = 100, 2
	fmt.Println("Cheats saving >=", timeSaved, "ps (", cheatTime, "ps time limit ):", findCheats(inputMap, timeSaved, cheatTime))
	timeSaved, cheatTime = 100, 20
	fmt.Println("Cheats saving >=", timeSaved, "ps (", cheatTime, "ps time limit ):", findCheats(inputMap, timeSaved, cheatTime))
}
