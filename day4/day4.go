package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
)

func readInput(inputFile string) [][]rune {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var input [][]rune
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		input = append(input, []rune(sc.Text()))
	}
	return input
}

func countWords(input [][]rune, word string) int {
	if len(word) == 0 {
		return 0
	}
	wordRune := []rune(word)
	wordRuneReverse := make([]rune, len(wordRune))
	copy(wordRuneReverse, wordRune)
	slices.Reverse(wordRuneReverse)

	count := 0

	for row := 0; row < len(input); row++ {
		for col := 0; col < len(input[row]); col++ {
			if input[row][col] == wordRune[0] {
				horizontalMatch := true
				verticalMatch := true
				leftDiagMatch := true
				rightDiagMatch := true
				for i := 1; i < len(wordRune); i++ {
					// Check left to right
					if horizontalMatch {
						if col+i >= len(input[row]) {
							horizontalMatch = false
						} else if input[row][col+i] != wordRune[i] {
							horizontalMatch = false
						}
					}
					// Check top to bottom
					if verticalMatch {
						if row+i >= len(input) {
							verticalMatch = false
						} else if input[row+i][col] != wordRune[i] {
							verticalMatch = false
						}
					}
					// Check left diagonal
					if leftDiagMatch {
						if row+i >= len(input) || col-i < 0 {
							leftDiagMatch = false
						} else if input[row+i][col-i] != wordRune[i] {
							leftDiagMatch = false
						}
					}
					// Check right diagonal
					if rightDiagMatch {
						if row+i >= len(input) || col+i >= len(input[row+i]) {
							rightDiagMatch = false
						} else if input[row+i][col+i] != wordRune[i] {
							rightDiagMatch = false
						}
					}

					if !(horizontalMatch || verticalMatch || leftDiagMatch || rightDiagMatch) {
						continue
					}
				}
				if horizontalMatch {
					count += 1
				}
				if verticalMatch {
					count += 1
				}
				if leftDiagMatch {
					count += 1
				}
				if rightDiagMatch {
					count += 1
				}
			}

			if input[row][col] == wordRuneReverse[0] {
				horizontalMatch := true
				verticalMatch := true
				leftDiagMatch := true
				rightDiagMatch := true
				for i := 1; i < len(wordRuneReverse); i++ {
					// Check left to right
					if horizontalMatch {
						if col+i >= len(input[row]) {
							horizontalMatch = false
						} else if input[row][col+i] != wordRuneReverse[i] {
							horizontalMatch = false
						}
					}
					// Check top to bottom
					if verticalMatch {
						if row+i >= len(input) {
							verticalMatch = false
						} else if input[row+i][col] != wordRuneReverse[i] {
							verticalMatch = false
						}
					}
					// Check left diagonal
					if leftDiagMatch {
						if row+i >= len(input) || col-i < 0 {
							leftDiagMatch = false
						} else if input[row+i][col-i] != wordRuneReverse[i] {
							leftDiagMatch = false
						}
					}
					// Check right diagonal
					if rightDiagMatch {
						if row+i >= len(input) || col+i >= len(input[row+i]) {
							rightDiagMatch = false
						} else if input[row+i][col+i] != wordRuneReverse[i] {
							rightDiagMatch = false
						}
					}

					if !(horizontalMatch || verticalMatch || leftDiagMatch || rightDiagMatch) {
						continue
					}
				}
				if horizontalMatch {
					count += 1
				}
				if verticalMatch {
					count += 1
				}
				if leftDiagMatch {
					count += 1
				}
				if rightDiagMatch {
					count += 1
				}
			}
		}
	}

	return count
}

func countXMas(input [][]rune) int {
	ARune := rune('A')
	MRune := rune('M')
	SRune := rune('S')

	count := 0

	for row := 1; row < len(input)-1; row++ {
		for col := 1; col < len(input[row])-1; col++ {
			if input[row][col] != ARune {
				continue
			}
			lDiag := input[row-1][col-1] == MRune && input[row+1][col+1] == SRune
			lDiag = lDiag || (input[row-1][col-1] == SRune && input[row+1][col+1] == MRune)
			if !lDiag {
				continue
			}
			rDiag := input[row-1][col+1] == MRune && input[row+1][col-1] == SRune
			rDiag = rDiag || (input[row-1][col+1] == SRune && input[row+1][col-1] == MRune)
			if !rDiag {
				continue
			}
			count += 1
		}
	}

	return count
}

func main() {
	fmt.Println("Test Input")
	inputRunes := readInput("test-input.txt")
	fmt.Println("XMAS matched:", countWords(inputRunes, "XMAS"))
	fmt.Println("X-MAS matched:", countXMas(inputRunes))

	fmt.Println("Puzzle Input")
	inputRunes = readInput("input.txt")
	fmt.Println("XMAS matched:", countWords(inputRunes, "XMAS"))
	fmt.Println("X-MAS matched:", countXMas(inputRunes))
}
