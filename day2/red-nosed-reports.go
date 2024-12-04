package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func getReports(inputFile string) (reports [][]int) {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		tokens := strings.Fields(sc.Text())

		var slice []int

		for _, val := range tokens {
			i_val, err := strconv.Atoi(val)
			if err != nil {
				log.Fatal(err)
			}
			slice = append(slice, i_val)
		}
		reports = append(reports, slice)
	}

	return reports
}

func isSafeWithDrop(report []int, canDrop bool) bool {
	maxDrop := 0
	if canDrop {
		maxDrop = 1
	}

	diffs := make([]int, len(report)-1)
	signTotal := 0
	for i := range diffs {
		diff := report[i+1] - report[i]
		diffs[i] = diff
		if diff > 0 {
			signTotal += 1
		} else if diff < 0 {
			signTotal -= 1
		}
	}

	isNegative := signTotal < 0
	if isNegative {
		signTotal *= -1
	}

	if !canDrop && (signTotal) < len(diffs) {
		return false
	}

	for i := 0; i < len(diffs); i++ {
		if maxDrop < 0 {
			return false
		}

		val := diffs[i]
		if isNegative {
			val *= -1
		}

		if val > 0 && val <= 3 {
			continue
		}

		if maxDrop == 0 {
			return false
		}

		if (i + 1) < len(diffs) {
			val1 := diffs[i+1]
			if isNegative {
				val1 *= -1
			}
			newDiff := val + val1
			if newDiff > 0 && newDiff <= 3 {
				maxDrop -= 1
				i += 1
				continue
			}
		}

		if (i - 1) > 0 {
			val1 := diffs[i-1]
			if isNegative {
				val1 *= -1
			}
			newDiff := val + val1
			if newDiff > 0 && newDiff <= 3 {
				maxDrop -= 1
				continue
			}
		}

		if i == 0 || i == len(diffs)-1 {
			maxDrop -= 1
			continue
		}

		return false
	}

	return maxDrop >= 0
}

func main() {
	fmt.Println("Test Input")
	reports := getReports("test-input.txt")
	safeReports := 0
	safeReportsDamp1 := 0
	for _, report := range reports {
		if isSafeWithDrop(report, false) {
			safeReports += 1
		}
		if isSafeWithDrop(report, true) {
			safeReportsDamp1 += 1
		}
	}
	fmt.Println("Safe Reports:", safeReports)
	fmt.Println("Safe Reports Damp 1:", safeReportsDamp1)

	fmt.Println("Puzzle Input")
	reports = getReports("input.txt")
	safeReports = 0
	safeReportsDamp1 = 0
	for _, report := range reports {
		if isSafeWithDrop(report, false) {
			safeReports += 1
		}
		if isSafeWithDrop(report, true) {
			safeReportsDamp1 += 1
		}
	}
	fmt.Println("Safe Reports:", safeReports)
	fmt.Println("Safe Reports Damp 1:", safeReportsDamp1)
}
