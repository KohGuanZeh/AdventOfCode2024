package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type RunePosition struct {
	val rune
	row int
	col int
}

const FREE = rune('.')
const BOX = rune('O')
const WALL = rune('#')
const FISH = rune('@')

const LEFTBOX = rune('[')
const RIGHTBOX = rune(']')

const UP = rune('^')
const DOWN = rune('v')
const LEFT = rune('<')
const RIGHT = rune('>')

func readInput(inputFile string) ([][]rune, []rune) {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	inputMap := make([][]rune, 0)
	moves := make([]rune, 0)
	finishedMap := false
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		text := sc.Text()
		if len(text) == 0 {
			finishedMap = true
			continue
		}
		if finishedMap {
			moves = append(moves, []rune(sc.Text())...)
			continue
		}
		row := []rune(sc.Text())
		inputMap = append(inputMap, row)
	}
	return inputMap, moves
}

func moveBoxes(inputMap [][]rune, moves []rune) [][]rune {
	mapCopy := make([][]rune, len(inputMap))
	for r := range mapCopy {
		mapCopy[r] = make([]rune, len(inputMap[r]))
		for c := range mapCopy[r] {
			mapCopy[r][c] = inputMap[r][c]
		}
	}
	rPos := -1
	cPos := -1
	for r := range mapCopy {
		for c := range mapCopy {
			if mapCopy[r][c] == FISH {
				rPos = r
				cPos = c
				break
			}
		}
		if rPos > -1 || cPos > -1 {
			break
		}
	}

	for i := 0; i < len(moves); i++ {
		moveBy := 1
		r := moves[i]
		for i+1 < len(moves) && moves[i+1] == r {
			i += 1
			moveBy += 1
		}
		boxes := 0
		mapCopy[rPos][cPos] = FREE
		switch r {
		case UP:
			for i := 0; i < moveBy; i++ {
				// No need to check bounds since surrounded by walls
				if mapCopy[rPos-1][cPos] == WALL {
					break
				}
				rPos -= 1
				if mapCopy[rPos][cPos] == BOX {
					boxes += 1
					mapCopy[rPos][cPos] = FREE
					i -= 1
				}
			}
			for range boxes {
				mapCopy[rPos][cPos] = BOX
				rPos += 1
			}
			break
		case DOWN:
			for i := 0; i < moveBy; i++ {
				// No need to check bounds since surrounded by walls
				if mapCopy[rPos+1][cPos] == WALL {
					break
				}
				rPos += 1
				if mapCopy[rPos][cPos] == BOX {
					boxes += 1
					mapCopy[rPos][cPos] = FREE
					i -= 1
				}
			}
			for range boxes {
				mapCopy[rPos][cPos] = BOX
				rPos -= 1
			}
			break
		case LEFT:
			for i := 0; i < moveBy; i++ {
				// No need to check bounds since surrounded by walls
				if mapCopy[rPos][cPos-1] == WALL {
					break
				}
				cPos -= 1
				if mapCopy[rPos][cPos] == BOX {
					boxes += 1
					mapCopy[rPos][cPos] = FREE
					i -= 1
				}
			}
			for range boxes {
				mapCopy[rPos][cPos] = BOX
				cPos += 1
			}
			break
		case RIGHT:
			for i := 0; i < moveBy; i++ {
				// No need to check bounds since surrounded by walls
				if mapCopy[rPos][cPos+1] == WALL {
					break
				}
				cPos += 1
				if mapCopy[rPos][cPos] == BOX {
					boxes += 1
					mapCopy[rPos][cPos] = FREE
					i -= 1
				}
			}
			for range boxes {
				mapCopy[rPos][cPos] = BOX
				cPos -= 1
			}
			break
		default:
			log.Fatal("Wrong rune used for movement.")
			break
		}
		mapCopy[rPos][cPos] = FISH

		// fmt.Println("Move:", string(r), ", by", moveBy, "times")
		// fmt.Println("Result Map:")
		// for r := range mapCopy {
		// fmt.Println(string(mapCopy[r]))
		// }
	}
	return mapCopy
}

func moveWideBoxes(inputMap [][]rune, moves []rune) [][]rune {
	// Convert to wide map
	rPos := -1
	cPos := -1
	mapCopy := make([][]rune, len(inputMap))
	for r := range mapCopy {
		mapCopy[r] = make([]rune, len(inputMap[r])*2)
		for c := range inputMap[r] {
			switch inputMap[r][c] {
			case FREE:
				mapCopy[r][c*2] = FREE
				mapCopy[r][c*2+1] = FREE
				break
			case WALL:
				mapCopy[r][c*2] = WALL
				mapCopy[r][c*2+1] = WALL
				break
			case FISH:
				mapCopy[r][c*2] = FISH
				mapCopy[r][c*2+1] = FREE
				rPos = r
				cPos = c * 2
				break
			case BOX:
				mapCopy[r][c*2] = LEFTBOX
				mapCopy[r][c*2+1] = RIGHTBOX
				break
			}
		}
	}

	for i := 0; i < len(moves); i++ {
		moveBy := 1
		r := moves[i]
		for i+1 < len(moves) && moves[i+1] == r {
			i += 1
			moveBy += 1
		}
		mapCopy[rPos][cPos] = FREE
		switch r {
		case UP:
			runePosMap := make(map[RunePosition]int)
			fishRunePos := RunePosition{val: FISH, row: rPos, col: cPos}
			runePosMap[fishRunePos] = 0
			for i := 0; i < moveBy; i++ {
				hitWall := false
				runePosToAdd := make(map[RunePosition]int)
				for runePos, moved := range runePosMap {
					r, c := runePos.row+moved-1, runePos.col
					switch mapCopy[r][c] {
					case WALL:
						hitWall = true
						break
					case LEFTBOX:
						mapCopy[r][c] = FREE
						mapCopy[r][c+1] = FREE
						r1Pos := RunePosition{val: LEFTBOX, row: r, col: c}
						r2Pos := RunePosition{val: RIGHTBOX, row: r, col: c + 1}
						runePosToAdd[r1Pos] = 0
						runePosToAdd[r2Pos] = 0
						break
					case RIGHTBOX:
						mapCopy[r][c] = FREE
						mapCopy[r][c-1] = FREE
						r1Pos := RunePosition{val: LEFTBOX, row: r, col: c - 1}
						r2Pos := RunePosition{val: RIGHTBOX, row: r, col: c}
						runePosToAdd[r1Pos] = 0
						runePosToAdd[r2Pos] = 0
						break
					}
					if hitWall {
						break
					}
				}
				if len(runePosToAdd) > 0 {
					for runePos, moved := range runePosToAdd {
						runePosMap[runePos] = moved
					}
					i -= 1
					if hitWall {
						break
					}
					continue
				}
				if hitWall {
					break
				}
				for runePos, moved := range runePosMap {
					runePosMap[runePos] = moved - 1
				}
			}
			for runePos, moved := range runePosMap {
				mapCopy[runePos.row+moved][runePos.col] = runePos.val
				if runePos.val == FISH {
					rPos = runePos.row + moved
				}
			}
			break
		case DOWN:
			runePosMap := make(map[RunePosition]int)
			fishRunePos := RunePosition{val: FISH, row: rPos, col: cPos}
			runePosMap[fishRunePos] = 0
			for i := 0; i < moveBy; i++ {
				hitWall := false
				runePosToAdd := make(map[RunePosition]int)
				for runePos, moved := range runePosMap {
					r, c := runePos.row+moved+1, runePos.col
					switch mapCopy[r][c] {
					case WALL:
						hitWall = true
						break
					case LEFTBOX:
						mapCopy[r][c] = FREE
						mapCopy[r][c+1] = FREE
						r1Pos := RunePosition{val: LEFTBOX, row: r, col: c}
						r2Pos := RunePosition{val: RIGHTBOX, row: r, col: c + 1}
						runePosToAdd[r1Pos] = 0
						runePosToAdd[r2Pos] = 0
						break
					case RIGHTBOX:
						mapCopy[r][c] = FREE
						mapCopy[r][c-1] = FREE
						r1Pos := RunePosition{val: LEFTBOX, row: r, col: c - 1}
						r2Pos := RunePosition{val: RIGHTBOX, row: r, col: c}
						runePosToAdd[r1Pos] = 0
						runePosToAdd[r2Pos] = 0
						break
					}
					if hitWall {
						break
					}
				}
				if len(runePosToAdd) > 0 {
					for runePos, moved := range runePosToAdd {
						runePosMap[runePos] = moved
					}
					i -= 1
					if hitWall {
						break
					}
					continue
				}
				if hitWall {
					break
				}
				for runePos, moved := range runePosMap {
					runePosMap[runePos] = moved + 1
				}
			}
			for runePos, moved := range runePosMap {
				mapCopy[runePos.row+moved][runePos.col] = runePos.val
				if runePos.val == FISH {
					rPos = runePos.row + moved
				}
			}
			break
		case LEFT:
			boxes := 0
			for i := 0; i < moveBy; i++ {
				// No need to check bounds since surrounded by walls
				if mapCopy[rPos][cPos-1] == WALL {
					break
				}
				cPos -= 1
				if mapCopy[rPos][cPos] == RIGHTBOX {
					boxes += 1
					mapCopy[rPos][cPos] = FREE
					mapCopy[rPos][cPos-1] = FREE
					cPos -= 1
					i -= 1
				}
			}
			for range boxes {
				mapCopy[rPos][cPos] = LEFTBOX
				mapCopy[rPos][cPos+1] = RIGHTBOX
				cPos += 2
			}
			mapCopy[rPos][cPos] = FISH
			break
		case RIGHT:
			boxes := 0
			for i := 0; i < moveBy; i++ {
				// No need to check bounds since surrounded by walls
				if mapCopy[rPos][cPos+1] == WALL {
					break
				}
				cPos += 1
				if mapCopy[rPos][cPos] == LEFTBOX {
					boxes += 1
					mapCopy[rPos][cPos] = FREE
					mapCopy[rPos][cPos+1] = FREE
					cPos += 1
					i -= 1
				}
			}
			for range boxes {
				mapCopy[rPos][cPos] = RIGHTBOX
				mapCopy[rPos][cPos-1] = LEFTBOX
				cPos -= 2
			}
			mapCopy[rPos][cPos] = FISH
			break
		default:
			log.Fatal("Wrong rune used for movement.")
			break
		}

		// fmt.Println("Move:", string(r), ", by", moveBy, "times")
		// fmt.Println("Result Map:")
		// for r := range mapCopy {
		// fmt.Println(string(mapCopy[r]))
		// }
	}
	return mapCopy
}

func getGpsCoordsSum(resultMap [][]rune) int {
	gpsCoordSum := 0
	fmt.Println("Result Map:")
	for r := range resultMap {
		fmt.Println(string(resultMap[r]))
		for c := range resultMap[r] {
			if resultMap[r][c] == BOX || resultMap[r][c] == LEFTBOX {
				gpsCoordSum += r*100 + c
			}
		}
	}
	return gpsCoordSum
}

func main() {
	fmt.Println("Small Test Input")
	inputMap, moves := readInput("small-test-input.txt")
	resultMap := moveBoxes(inputMap, moves)
	fmt.Println("GPS Coordinates Sum:", getGpsCoordsSum(resultMap))
	resultMap = moveWideBoxes(inputMap, moves)
	fmt.Println("GPS Coordinates Sum (Wide):", getGpsCoordsSum(resultMap))

	fmt.Println("")

	fmt.Println("Test Input")
	inputMap, moves = readInput("test-input.txt")
	resultMap = moveBoxes(inputMap, moves)
	fmt.Println("GPS Coordinates Sum:", getGpsCoordsSum(resultMap))
	resultMap = moveWideBoxes(inputMap, moves)
	fmt.Println("GPS Coordinates Sum (Wide):", getGpsCoordsSum(resultMap))

	fmt.Println("")

	fmt.Println("Puzzle Input")
	inputMap, moves = readInput("input.txt")
	resultMap = moveBoxes(inputMap, moves)
	fmt.Println("GPS Coordinates Sum:", getGpsCoordsSum(resultMap))
	resultMap = moveWideBoxes(inputMap, moves)
	fmt.Println("GPS Coordinates Sum (Wide):", getGpsCoordsSum(resultMap))
}
