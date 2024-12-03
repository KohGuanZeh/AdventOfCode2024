package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
)

func read_input(input_file string) string {
	file, err := os.Open(input_file)
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

		val_1, err := strconv.Atoi(match[1])
		if err != nil {
			log.Fatal(err)
		}
		val_2, err := strconv.Atoi(match[2])
		if err != nil {
			log.Fatal(err)
		}
		value += val_1 * val_2
	}
	return value
}

func main() {
	r1 := regexp.MustCompile(`mul\(([\d]+),([\d]+)\)`)
	r2 := regexp.MustCompile(`do\(\)|don't\(\)|mul\(([\d]+),([\d]+)\)`)
	fmt.Println("Test Input")
	input := read_input("test-input.txt")
	matches := r1.FindAllStringSubmatch(input, -1)
	fmt.Println("Total Value:", get_mul_sum(matches))
	input = read_input("test-input-2.txt")
	matches = r2.FindAllStringSubmatch(input, -1)
	fmt.Println("Total Value (do(), don't()):", get_mul_sum(matches))

	fmt.Println()

	fmt.Println("Puzzle Input")
	input = read_input("input.txt")
	matches = r1.FindAllStringSubmatch(input, -1)
	fmt.Println("Total Value:", get_mul_sum(matches))
	matches = r2.FindAllStringSubmatch(input, -1)
	fmt.Println("Total Value (do(), don't()):", get_mul_sum(matches))
}
