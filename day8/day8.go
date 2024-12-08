package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const EMPTY = rune('.')

type Position struct {
	row int
	col int
}

func readMap(inputFile string) (map[rune][]Position, int, int) {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	antennas := make(map[rune][]Position)
	sc := bufio.NewScanner(file)
	rowPos := 0
	maxColPos := -1
	for sc.Scan() {
		runeSlice := []rune(sc.Text())
		for colPos, char := range runeSlice {
			if maxColPos < colPos {
				maxColPos = colPos
			}
			if char == EMPTY {
				continue
			}
			s, ok := antennas[char]
			if !ok {
				s = make([]Position, 0)
			}
			s = append(s, Position{row: rowPos, col: colPos})
			antennas[char] = s
		}
		rowPos += 1
	}
	return antennas, rowPos - 1, maxColPos
}

func countAntinodes(antennas map[rune][]Position, maxRow int, maxCol int) int {
	antinodePositions := make(map[Position]int)
	for _, positions := range antennas {
		for i, position1 := range positions {
			for j := i + 1; j < len(positions); j++ {
				position2 := positions[j]
				positionDiff := Position{row: position2.row - position1.row, col: position2.col - position1.col}
				antinode1Pos := Position{row: position1.row - positionDiff.row, col: position1.col - positionDiff.col}
				if positionWithinMap(antinode1Pos, maxRow, maxCol) {
					antinodePositions[antinode1Pos] = 0
				}
				antinode2Pos := Position{row: position2.row + positionDiff.row, col: position2.col + positionDiff.col}
				if positionWithinMap(antinode2Pos, maxRow, maxCol) {
					antinodePositions[antinode2Pos] = 0
				}
			}
		}
	}
	return len(antinodePositions)
}

func countAntinodesHarmonics(antennas map[rune][]Position, maxRow int, maxCol int) int {
	antinodePositions := make(map[Position]int)
	for _, positions := range antennas {
		for i, position1 := range positions {
			for j := i + 1; j < len(positions); j++ {
				position2 := positions[j]
				antinodePositions[position1] = 0
				antinodePositions[position2] = 0
				positionDiff := Position{row: position2.row - position1.row, col: position2.col - position1.col}
				antinode1Pos := Position{row: position1.row, col: position1.col}
				for true {
					antinode1Pos = Position{row: antinode1Pos.row - positionDiff.row, col: antinode1Pos.col - positionDiff.col}
					if !positionWithinMap(antinode1Pos, maxRow, maxCol) {
						break
					}
					antinodePositions[antinode1Pos] = 0
				}
				antinode2Pos := Position{row: position2.row, col: position2.col}
				for true {
					antinode2Pos = Position{row: antinode2Pos.row + positionDiff.row, col: antinode2Pos.col + positionDiff.col}
					if !positionWithinMap(antinode2Pos, maxRow, maxCol) {
						break
					}
					antinodePositions[antinode2Pos] = 0
				}
			}
		}
	}
	return len(antinodePositions)
}

func positionWithinMap(position Position, maxRow int, maxCol int) bool {
	return position.row >= 0 && position.row <= maxRow && position.col >= 0 && position.col <= maxCol
}

func main() {
	fmt.Println("Test Input")
	antennas, maxRow, maxCol := readMap("test-input.txt")
	fmt.Println("Antinodes:", countAntinodes(antennas, maxRow, maxCol))
	fmt.Println("Antinodes:", countAntinodesHarmonics(antennas, maxRow, maxCol))

	fmt.Println()

	fmt.Println("Puzzle Input")
	antennas, maxRow, maxCol = readMap("input.txt")
	fmt.Println("Antinodes:", countAntinodes(antennas, maxRow, maxCol))
	fmt.Println("Antinodes:", countAntinodesHarmonics(antennas, maxRow, maxCol))
}
