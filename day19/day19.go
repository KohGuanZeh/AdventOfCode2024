package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func readInput(inputFile string) ([][]rune, [][]rune) {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	sc := bufio.NewScanner(file)
	patterns := make([][]rune, 0)

	sc.Scan()
	tokens := strings.Split(sc.Text(), ", ")
	towels := make([][]rune, len(tokens))
	for i, token := range tokens {
		towels[i] = []rune(token)
	}

	sc.Scan()

	for sc.Scan() {
		patterns = append(patterns, []rune(sc.Text()))
	}
	return towels, patterns
}

func possiblePatterns(towels, patterns [][]rune) int {
	matches := 0
	for _, pattern := range patterns {
		if canMatch(towels, pattern) {
			matches += 1
		}
	}
	return matches
}

func canMatch(towels [][]rune, pattern []rune) bool {
	if len(pattern) == 0 {
		return true
	}
	for _, towel := range towels {
		matchPattern := true
		if len(towel) > len(pattern) {
			continue
		}
		for j, r := range towel {
			if pattern[j] != r {
				matchPattern = false
				break
			}
		}
		if matchPattern {
			if canMatch(towels, pattern[len(towel):]) {
				return true
			}
		}
	}
	return false
}

func patternsSum(towels, patterns [][]rune) int {
	sum := 0
	pSumMap := make(map[string]int)
	for _, pattern := range patterns {
		sum += patternSum(towels, pattern, pSumMap)
	}
	return sum
}

func patternSum(towels [][]rune, pattern []rune, pSumMap map[string]int) int {
	if len(pattern) == 0 {
		return 1
	}
	numWays, ok := pSumMap[string(pattern)]
	if ok {
		return numWays
	}
	numWays = 0
	for _, towel := range towels {
		matchPattern := true
		if len(towel) > len(pattern) {
			continue
		}
		for j, r := range towel {
			if pattern[j] != r {
				matchPattern = false
				break
			}
		}
		if matchPattern {
			numWays += patternSum(towels, pattern[len(towel):], pSumMap)
		}
	}
	pSumMap[string(pattern)] = numWays
	return numWays
}

func main() {
	fmt.Println("Test Input")
	towels, patterns := readInput("test-input.txt")
	fmt.Println("Possible Patterns:", possiblePatterns(towels, patterns))
	fmt.Println("Possible Patterns Sum:", patternsSum(towels, patterns))

	fmt.Println("")

	fmt.Println("Puzzle Input")
	towels, patterns = readInput("input.txt")
	fmt.Println("Possible Patterns:", possiblePatterns(towels, patterns))
	fmt.Println("Possible Patterns Sum:", patternsSum(towels, patterns))
}
