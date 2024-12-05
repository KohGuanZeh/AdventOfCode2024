package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func pageOrderRules(orderFile string) map[int]map[int]byte {
	ruleMap := make(map[int]map[int]byte)

	file, err := os.Open(orderFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		tokens := strings.Split(sc.Text(), "|")
		pageA, err := strconv.Atoi(tokens[0])
		if err != nil {
			log.Fatal(err)
		}
		pageB, err := strconv.Atoi(tokens[1])
		if err != nil {
			log.Fatal(err)
		}
		m, ok := ruleMap[pageA]
		if !ok {
			ruleMap[pageA] = map[int]byte{}
			m = ruleMap[pageA]
		}
		m[pageB] = 0

	}
	return ruleMap
}

func getPageUpdates(updateFile string) [][]int {
	updates := [][]int{}

	file, err := os.Open(updateFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		tokens := strings.Split(sc.Text(), ",")
		tokensInt := make([]int, len(tokens))
		for i, val := range tokens {
			tokensInt[i], err = strconv.Atoi(val)
			if err != nil {
				log.Fatal(err)
			}
		}
		updates = append(updates, tokensInt)
	}

	return updates
}

func sumCorrectOrder(updates [][]int, ruleMap map[int]map[int]byte) (int, int) {
	correctSum := 0
	fixSum := 0
	for _, s := range updates {
		isCorrect := true
		for i := 0; i < len(s)-1; i++ {
			swapped := false
			ti := i
			for j := i + 1; j < len(s); j++ {
				_, ok := ruleMap[s[ti]][s[j]]
				if ok {
					continue
				}
				_, ok = ruleMap[s[j]][s[ti]]
				if ok {
					// Swap incorrect value placement
					swapped = true
					tVal := s[ti]
					s[ti] = s[j]
					if ti+1 == j {
						s[j] = tVal
						continue
					}
					ti += 1
					s[j] = s[ti]
					s[ti] = tVal
				}
			}
			if swapped {
				isCorrect = false
				i -= 1
			}
		}
		if isCorrect {
			correctSum += s[len(s)/2]
		} else {
			fixSum += s[len(s)/2]
		}
	}
	return correctSum, fixSum
}

func main() {
	fmt.Println("Test Input")
	ruleMap := pageOrderRules("test-input-1.txt")
	updates := getPageUpdates("test-input-2.txt")
	correct, fix := sumCorrectOrder(updates, ruleMap)
	fmt.Println("Sum of Correct Order:", correct)
	fmt.Println("Fix Correct Order:", fix)

	fmt.Println("Puzzle Input")
	ruleMap = pageOrderRules("input-1.txt")
	updates = getPageUpdates("input-2.txt")
	correct, fix = sumCorrectOrder(updates, ruleMap)
	fmt.Println("Sum of Correct Order:", correct)
	fmt.Println("Sum of Correct Order:", fix)

}
