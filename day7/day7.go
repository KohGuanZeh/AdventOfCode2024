package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func readInput(inputFile string) [][]int {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	numSlice := make([][]int, 0)
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		tokens := strings.Split(sc.Text(), ": ")
		nums := make([]int, 0)
		total, err := strconv.Atoi(tokens[0])
		if err != nil {
			log.Fatal(err)
		}
		nums = append(nums, total)
		tokens = strings.Split(tokens[1], " ")
		for _, numStr := range tokens {
			num, err := strconv.Atoi(numStr)
			if err != nil {
				log.Fatal(err)
			}
			nums = append(nums, num)
		}
		numSlice = append(numSlice, nums)
	}
	return numSlice
}

func sumCorrectEqn(numSlice [][]int, concat bool) int {
	sum := 0
	for _, nums := range numSlice {
		target := nums[0]
		resQueue := []int{nums[1]}
		for i := 2; i < len(nums); i++ {
			if len(resQueue) == 0 {
				break
			}
			newQueue := make([]int, 0)
			for _, value := range resQueue {
				// Sum
				nextValue := value + nums[i]
				if nextValue <= target {
					newQueue = append(newQueue, nextValue)
				}
				// Multiply
				nextValue = value * nums[i]
				if nextValue <= target {
					newQueue = append(newQueue, nextValue)
				}
				// Concatenation
				if !concat {
					continue
				}
				nextValue = concatNums(value, nums[i])
				if nextValue <= target {
					newQueue = append(newQueue, nextValue)
				}
			}
			resQueue = newQueue
		}
		for _, value := range resQueue {
			if value == target {
				sum += target
				break
			}
		}
	}
	return sum
}

func concatNums(leftNum int, rightNum int) int {
	i := 10
	for i < rightNum {
		i *= 10
	}
	return leftNum*i + rightNum
}

func main() {
	fmt.Println("Test Input")
	numSlice := readInput("test-input.txt")
	fmt.Println("Sum of Correct Equations:", sumCorrectEqn(numSlice, false))
	fmt.Println("Sum of Correct Equations (With Concatentation):", sumCorrectEqn(numSlice, true))

	fmt.Println()

	fmt.Println("Puzzle Input")
	numSlice = readInput("input.txt")
	fmt.Println("Sum of Correct Equations:", sumCorrectEqn(numSlice, false))
	fmt.Println("Sum of Correct Equations (With Concatentation):", sumCorrectEqn(numSlice, true))
}
