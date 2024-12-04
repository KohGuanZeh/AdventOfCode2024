package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
)

func readInput(inputFile string) string {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	return string(bytes)
}

func get_mul_sum(matches [][]string) int {
	value := 0
	multiply := true
	for _, match := range matches {
		switch match[0] {
		case "do()":
			multiply = true
			continue
		case "don't()":
			multiply = false
			continue
		}

		if !multiply {
			continue
		}

		value1, err := strconv.Atoi(match[1])
		if err != nil {
			log.Fatal(err)
		}
		value2, err := strconv.Atoi(match[2])
		if err != nil {
			log.Fatal(err)
		}
		value += value1 * value2
	}
	return value
}

func main() {
	r1 := regexp.MustCompile(`mul\(([\d]+),([\d]+)\)`)
	r2 := regexp.MustCompile(`do\(\)|don't\(\)|mul\(([\d]+),([\d]+)\)`)
	fmt.Println("Test Input")
	input := readInput("test-input.txt")
	matches := r1.FindAllStringSubmatch(input, -1)
	fmt.Println("Total Value:", get_mul_sum(matches))
	input = readInput("test-input-2.txt")
	matches = r2.FindAllStringSubmatch(input, -1)
	fmt.Println("Total Value (do(), don't()):", get_mul_sum(matches))

	fmt.Println()

	fmt.Println("Puzzle Input")
	input = readInput("input.txt")
	matches = r1.FindAllStringSubmatch(input, -1)
	fmt.Println("Total Value:", get_mul_sum(matches))
	matches = r2.FindAllStringSubmatch(input, -1)
	fmt.Println("Total Value (do(), don't()):", get_mul_sum(matches))
}
