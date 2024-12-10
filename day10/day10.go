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

func readInput(inputFile string) [][]int {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	topoMap := make([][]int, 0)
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		r := []rune(sc.Text())
		atoi := make([]int, len(r))
		for i, v := range r {
			intVal := int(v - '0')
			atoi[i] = intVal
		}
		topoMap = append(topoMap, atoi)
	}
	return topoMap
}

func scoreTrailheads(topoMap [][]int) (int, int) {
	uniqueNines := make(map[Position]map[Position]int)
	pathsToNine := make(map[Position]int)
	score := 0
	pathScore := 0
	for i := range topoMap {
		for j := range topoMap[i] {
			currPos := Position{row: i, col: j}
			trails := len(findPathToNine(topoMap, uniqueNines, pathsToNine, currPos))
			if topoMap[i][j] == 0 {
				score += trails
				pathScore += pathsToNine[currPos]
			}
		}
	}
	return score, pathScore
}

// Returns the number of paths to nine.
func findPathToNine(topoMap [][]int, uniqueNines map[Position]map[Position]int, pathsToNine map[Position]int, currPos Position) map[Position]int {
	_, ok := uniqueNines[currPos]
	if ok {
		return uniqueNines[currPos]
	}

	pathsToNine[currPos] = 0
	currP := make(map[Position]int)
	currV := topoMap[currPos.row][currPos.col]

	if currV == 9 {
		pathsToNine[currPos] = 1
		currP[currPos] = 0
		uniqueNines[currPos] = currP
		return uniqueNines[currPos]
	}

	// Check left
	nextPos := Position{row: currPos.row, col: currPos.col - 1}
	if nextPos.col > -1 && topoMap[nextPos.row][nextPos.col]-currV == 1 {
		m := findPathToNine(topoMap, uniqueNines, pathsToNine, nextPos)
		if len(m) > 0 {
			pathsToNine[currPos] += pathsToNine[nextPos]
			for k := range m {
				currP[k] = 0
			}
		}
	}
	// Check up
	nextPos = Position{row: currPos.row - 1, col: currPos.col}
	if nextPos.row > -1 && topoMap[nextPos.row][nextPos.col]-currV == 1 {
		m := findPathToNine(topoMap, uniqueNines, pathsToNine, nextPos)
		if len(m) > 0 {
			pathsToNine[currPos] += pathsToNine[nextPos]
			for k := range m {
				currP[k] = 0
			}
		}
	}
	// Check right
	nextPos = Position{row: currPos.row, col: currPos.col + 1}
	if nextPos.col < len(topoMap[nextPos.row]) && topoMap[nextPos.row][nextPos.col]-currV == 1 {
		m := findPathToNine(topoMap, uniqueNines, pathsToNine, nextPos)
		if len(m) > 0 {
			pathsToNine[currPos] += pathsToNine[nextPos]
			for k := range m {
				currP[k] = 0
			}
		}
	}
	// Check down
	nextPos = Position{row: currPos.row + 1, col: currPos.col}
	if nextPos.row < len(topoMap) && topoMap[nextPos.row][nextPos.col]-currV == 1 {
		m := findPathToNine(topoMap, uniqueNines, pathsToNine, nextPos)
		if len(m) > 0 {
			pathsToNine[currPos] += pathsToNine[nextPos]
			for k := range m {
				currP[k] = 0
			}
		}
	}

	uniqueNines[currPos] = currP
	return uniqueNines[currPos]
}

func main() {
	fmt.Println("Test Input")
	topoMap := readInput("test-input.txt")
	score, pathScore := scoreTrailheads(topoMap)
	fmt.Println("Sum of Score:", score)
	fmt.Println("Sum of Path Score:", pathScore)

	fmt.Println("")

	fmt.Println("Puzzle Input")
	topoMap = readInput("input.txt")
	score, pathScore = scoreTrailheads(topoMap)
	fmt.Println("Sum of Score:", score)
	fmt.Println("Sum of Path Score:", pathScore)
}
