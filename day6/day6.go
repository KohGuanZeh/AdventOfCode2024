package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Direction struct {
	key    string
	rowDir int
	colDir int
}

const FreeSpace = rune('.')
const Obstacle = rune('#')
const Visited = rune('X')
const Guard = rune('^')

var Left = Direction{key: "LEFT", rowDir: 0, colDir: -1}
var Right = Direction{key: "RIGHT", rowDir: 0, colDir: 1}
var Up = Direction{key: "UP", rowDir: -1, colDir: 0}
var Down = Direction{key: "DOWN", rowDir: 1, colDir: 0}

func readMap(inputFile string) [][]rune {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var inputMap [][]rune
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		inputMap = append(inputMap, []rune(sc.Text()))
	}
	return inputMap
}

func rotateClockwise(direction Direction) Direction {
	switch direction {
	case Up:
		return Right
	case Right:
		return Down
	case Down:
		return Left
	case Left:
		return Up
	}
	// Should not reach here.
	return Up
}

func countVisitedSpaces(inputMap [][]rune) (int, int) {
	rowOrigin := -1
	colOrigin := -1
	for row, s := range inputMap {
		for col, r := range s {
			if r == Guard {
				rowOrigin = row
				colOrigin = col
			}
		}
		if colOrigin > -1 || rowOrigin > -1 {
			break
		}
	}

	rowPos := rowOrigin
	colPos := colOrigin

	visited := 0
	possibleLoops := 0
	currDir := Up
	placedMap := make(map[int]map[int]byte)
	for true {
		newColPos := colPos + currDir.colDir
		newRowPos := rowPos + currDir.rowDir
		if newRowPos < 0 || newRowPos >= len(inputMap) || newColPos < 0 || newColPos >= len(inputMap[newRowPos]) {
			visited += 1
			break
		}
		switch inputMap[newRowPos][newColPos] {
		case FreeSpace:
			inputMap[rowPos][colPos] = Visited
			if registerPlacedMap(placedMap, newRowPos, newColPos) {
				inputMap[newRowPos][newColPos] = Obstacle
				if canLoop(inputMap, rowPos, colPos, currDir) {
					possibleLoops += 1
				}
				inputMap[newRowPos][newColPos] = FreeSpace
			}
			visited += 1
			colPos = newColPos
			rowPos = newRowPos
			break
		case Visited:
			inputMap[rowPos][colPos] = Visited
			if registerPlacedMap(placedMap, newRowPos, newColPos) {
				inputMap[newRowPos][newColPos] = Obstacle
				if (newRowPos != rowOrigin && newColPos != colOrigin) && canLoop(inputMap, rowPos, colPos, currDir) {
					possibleLoops += 1
				}
				inputMap[newRowPos][newColPos] = Visited
			}
			colPos = newColPos
			rowPos = newRowPos
			break
		case Obstacle:
			currDir = rotateClockwise(currDir)
			break
		}
	}

	return visited, possibleLoops
}

// Returns false if already registered
func registerPlacedMap(cMap map[int]map[int]byte, rowPos int, colPos int) bool {
	_, ok := cMap[rowPos]
	if !ok {
		cMap[rowPos] = make(map[int]byte)
	}

	_, ok = cMap[rowPos][colPos]
	if !ok {
		cMap[rowPos][colPos] = 0
	}

	return !ok
}

func canLoop(inputMap [][]rune, rowPos int, colPos int, currDir Direction) bool {
	// Check if can loop
	cMap := make(map[int]map[int]map[string]byte)
	registerLoopMap(cMap, rowPos, colPos, currDir.key)
	rotateClockwise(currDir)
	for true {
		newColPos := colPos + currDir.colDir
		newRowPos := rowPos + currDir.rowDir
		if newRowPos < 0 || newRowPos >= len(inputMap) || newColPos < 0 || newColPos >= len(inputMap[newRowPos]) {
			return false
		}
		switch inputMap[newRowPos][newColPos] {
		case Obstacle:
			currDir = rotateClockwise(currDir)
			break
		default:
			if !registerLoopMap(cMap, rowPos, colPos, currDir.key) {
				// Not new entry, means it has looped
				return true
			}
			colPos = newColPos
			rowPos = newRowPos
			break
		}
	}
	// Should not reach here
	return false
}

// Returns false if already registered
func registerLoopMap(cMap map[int]map[int]map[string]byte, rowPos int, colPos int, key string) bool {
	_, ok := cMap[rowPos]
	if !ok {
		cMap[rowPos] = make(map[int]map[string]byte)
	}

	_, ok = cMap[rowPos][colPos]
	if !ok {
		cMap[rowPos][colPos] = make(map[string]byte)
	}

	_, ok = cMap[rowPos][colPos][key]
	if !ok {
		cMap[rowPos][colPos][key] = 0
	}

	return !ok
}

func main() {
	fmt.Println("Test Input")
	inputMap := readMap("test-input.txt")
	visited, loops := countVisitedSpaces(inputMap)
	fmt.Println("Distinct Positions:", visited)
	fmt.Println("Possible Loops:", loops)

	fmt.Println("Puzzle Input")
	inputMap = readMap("input.txt")
	visited, loops = countVisitedSpaces(inputMap)
	fmt.Println("Distinct Positions:", visited)
	fmt.Println("Possible Loops:", loops)
}
