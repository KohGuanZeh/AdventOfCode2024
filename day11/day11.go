package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func readInput(inputFile string) []int {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	input := make([]int, 0)
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		tokens := strings.Split(sc.Text(), " ")
		for _, token := range tokens {
			intVal, err := strconv.Atoi(token)
			if err != nil {
				log.Fatal(err)
			}
			input = append(input, intVal)
		}
	}
	return input
}

func stonesAfterBlink(stones []int, blinks int) int {
	stoneMap := make(map[int]map[int]int)
	numStones := 0
	for _, stone := range stones {
		numStones += blink(stoneMap, stone, blinks)
	}
	return numStones
}

func blink(stoneMap map[int]map[int]int, stone int, blinks int) int {
	if blinks == 0 {
		return 1
	}
	_, ok := stoneMap[stone]
	if !ok {
		stoneMap[stone] = make(map[int]int)
	} else {
		numStones, ok := stoneMap[stone][blinks]
		if ok {
			return numStones
		}
	}

	if stone == 0 {
		stoneMap[stone][1] = 1
		stoneMap[stone][blinks] = blink(stoneMap, 1, blinks-1)
		return stoneMap[stone][blinks]
	}

	digits := countDigits(stone)
	if digits&1 == 0 {
		lDigits := digits / 2
		div := intPow(10, lDigits)
		lStone := stone / div
		rStone := stone - lStone*div
		stoneMap[stone][blinks] = blink(stoneMap, lStone, blinks-1) + blink(stoneMap, rStone, blinks-1)
	} else {
		stoneMap[stone][blinks] = blink(stoneMap, stone*2024, blinks-1)
	}
	return stoneMap[stone][blinks]
}

func intPow(v int, exp int) int {
	vExp := 1
	for exp > 0 {
		vExp *= v
		exp -= 1
	}
	return vExp
}

func countDigits(v int) int {
	digits := 0
	for true {
		digits += 1
		v /= 10
		if v == 0 {
			break
		}
	}
	return digits
}

func main() {
	fmt.Println("Test Input")
	stones := readInput("test-input.txt")
	blinks := 6
	fmt.Println("Number of Stones after", blinks, "blinks:", stonesAfterBlink(stones, blinks))
	blinks = 25
	fmt.Println("Number of Stones after", blinks, "blinks:", stonesAfterBlink(stones, blinks))

	fmt.Println("")

	fmt.Println("Puzzle Input")
	stones = readInput("input.txt")
	blinks = 25
	fmt.Println("Number of Stones after", blinks, "blinks:", stonesAfterBlink(stones, blinks))
	blinks = 75
	fmt.Println("Number of Stones after", blinks, "blinks:", stonesAfterBlink(stones, blinks))
}
