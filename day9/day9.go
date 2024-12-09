package main

import (
	"fmt"
	"log"
	"os"
)

func readInput(inputFile string) []int {
	fileBytes, err := os.ReadFile(inputFile) // just pass the file name
	if err != nil {
		log.Fatal(err)
	}
	runes := []rune(string(fileBytes))
	atoi := make([]int, len(runes))
	for i, intRune := range runes {
		atoi[i] = int(intRune - '0')
	}
	return atoi
}

func fragmentChecksum(input []int) int {
	inputLen := len(input)
	lastFileId := inputLen / 2
	// If input length is even, last block free space.
	if inputLen&1 == 0 {
		lastFileId -= 1
		inputLen -= 1
	}
	checksum := 0
	nextFileId := 0   // Next file ID.
	pos := 0          // Current position in file disk
	i := -1           // Index from front for input traversal
	j := inputLen - 1 // Index from back
	remainingBlocks := input[j]
	for true {
		if nextFileId > lastFileId {
			break
		}
		i += 1
		v := input[i]
		// If i is odd, is free space
		if i&1 == 1 {
			for v > 0 {
				checksum += pos * lastFileId
				pos += 1
				v -= 1
				remainingBlocks -= 1
				if remainingBlocks <= 0 {
					lastFileId -= 1
					j -= 2
					remainingBlocks = input[j]
				}
			}
			continue
		}
		if nextFileId == lastFileId {
			v = remainingBlocks
		}
		for v > 0 {
			checksum += pos * nextFileId
			pos += 1
			v -= 1
		}
		nextFileId += 1
	}
	return checksum
}

func moveFileChecksum(input []int) int {
	inputLen := len(input)
	lastFileId := inputLen / 2
	// If input length is even, last block free space.
	if inputLen&1 == 0 {
		lastFileId -= 1
		inputLen -= 1
	}
	checksum := 0
	iPositions := make([]int, inputLen)
	for i := range iPositions {
		if i == 0 {
			iPositions[i] = 0
			continue
		}
		iPositions[i] = iPositions[i-1] + input[i-1]
	}
	for i := inputLen - 1; i > -1; i -= 2 {
		if i == 0 {
			fileSize := input[i]
			for fileSize > 0 {
				checksum += lastFileId * iPositions[i]
				iPositions[i] += 1
				fileSize -= 1
			}
			break
		}

		j := 1
		// If it cannot fit
		fileSize := input[i]
		for fileSize > input[j] && j < i {
			j += 2
		}
		if j > i || fileSize > input[j] {
			// Remain if cannot move.
			for fileSize > 0 {
				checksum += lastFileId * iPositions[i]
				iPositions[i] += 1
				fileSize -= 1
			}
		} else {
			input[j] -= fileSize
			for fileSize > 0 {
				checksum += lastFileId * iPositions[j]
				iPositions[j] += 1
				fileSize -= 1
			}
		}
		lastFileId -= 1
	}
	return checksum
}

func main() {
	fmt.Println("Test Input")
	inputString := readInput("test-input.txt")
	fmt.Println("Fragment Checksum:", fragmentChecksum(inputString))
	fmt.Println("Move File Checksum:", moveFileChecksum(inputString))

	fmt.Println()

	fmt.Println("Puzzle Input")
	inputString = readInput("input.txt")
	fmt.Println("Fragment Checksum:", fragmentChecksum(inputString))
	fmt.Println("Move File Checksum:", moveFileChecksum(inputString))
}
