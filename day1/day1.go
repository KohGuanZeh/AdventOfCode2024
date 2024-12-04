// https://adventofcode.com/2024/day/1
package main // Standalone Executable

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

func getSlices(inputFile string) (l1 []int, l2 []int) {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		tokens := strings.Fields(sc.Text())

		i1, err := strconv.Atoi(tokens[0])
		if err != nil {
			log.Fatal(err)
		}

		i2, err := strconv.Atoi(tokens[1])
		if err != nil {
			log.Fatal(err)
		}

		l1 = append(l1, i1)
		l2 = append(l2, i2)
	}

	return l1, l2
}

func difference(l1 []int, l2 []int) (difference int) {
	difference = 0
	for i := 0; i < len(l1); i++ {
		diff := l1[i] - l2[i]
		if diff < 0 {
			diff *= -1
		}
		difference += diff
	}
	return difference
}

func similarity(l1 []int, l2 []int) (similarity int) {
	similarity = 0

	i2 := 0
	count := 0

	for key, val := range l1 {
		if key > 0 && val == l1[key-1] {
			similarity += count * val
			continue
		}

		count = 0

		for i := i2; i < len(l2); i++ {
			if l2[i] == val {
				count += 1
				continue
			}
			if l2[i] > val {
				i2 = i
				break
			}
		}

		similarity += count * val
	}

	return similarity
}

// Program starts here
func main() {
	l1, l2 := getSlices("test-input.txt")

	slices.Sort(l1)
	slices.Sort(l2)

	fmt.Println("Test Input")
	fmt.Println("Total Distance:", difference(l1, l2))
	fmt.Println("Similarity Score:", similarity(l1, l2))

	fmt.Println("")

	l1, l2 = getSlices("input.txt")

	slices.Sort(l1)
	slices.Sort(l2)

	fmt.Println("Puzzle Input")
	fmt.Println("Total Distance:", difference(l1, l2))
	fmt.Println("Similarity Score:", similarity(l1, l2))
}
