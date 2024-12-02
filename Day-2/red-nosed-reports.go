package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func get_reports(input_file string) (reports [][]int) {
	file, err := os.Open(input_file)
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

func is_safe(report []int) bool {
	diffs := make([]int, len(report)-1)
	sign_total := 0
	for i := range diffs {
		diff := report[i+1] - report[i]
		diffs[i] = diff
		if diff > 0 {
			sign_total += 1
		} else if diff < 0 {
			sign_total -= 1
		}
	}

	neg_sign := sign_total < 0
	if neg_sign {
		sign_total *= -1
	}
	if (sign_total) < len(diffs) {
		return false
	}

	for _, val := range diffs {
		if neg_sign {
			val *= -1
		}
		if val == 0 || val > 3 {
			return false
		}
	}

	return true
}

func is_safe_with_drop(report []int, can_drop bool) bool {
	max_drop := 0
	if can_drop {
		max_drop = 1
	}

	diffs := make([]int, len(report)-1)
	sign_total := 0
	for i := range diffs {
		diff := report[i+1] - report[i]
		diffs[i] = diff
		if diff > 0 {
			sign_total += 1
		} else if diff < 0 {
			sign_total -= 1
		}
	}

	neg_sign := sign_total < 0
	if neg_sign {
		sign_total *= -1
	}

	if !can_drop && (sign_total) < len(diffs) {
		return false
	}

	for i := 0; i < len(diffs); i++ {
		if max_drop < 0 {
			return false
		}

		val := diffs[i]
		if neg_sign {
			val *= -1
		}

		if val > 0 && val <= 3 {
			continue
		}

		if max_drop == 0 {
			return false
		}

		if (i + 1) < len(diffs) {
			val_1 := diffs[i+1]
			if neg_sign {
				val_1 *= -1
			}
			new_diff := val + val_1
			if new_diff > 0 && new_diff <= 3 {
				max_drop -= 1
				i += 1
				continue
			}
		}

		if (i - 1) > 0 {
			val_1 := diffs[i-1]
			if neg_sign {
				val_1 *= -1
			}
			new_diff := val + val_1
			if new_diff > 0 && new_diff <= 3 {
				max_drop -= 1
				continue
			}
		}

		if i == 0 || i == len(diffs)-1 {
			max_drop -= 1
			continue
		}

		return false
	}

	return max_drop >= 0
}

func main() {
	fmt.Println("Test Input")
	reports := get_reports("test-input.txt")
	safe_reports := 0
	safe_reports_damp_1 := 0
	for _, report := range reports {
		if is_safe_with_drop(report, false) {
			safe_reports += 1
		}
		if is_safe_with_drop(report, true) {
			safe_reports_damp_1 += 1
		}
	}
	fmt.Println("Safe Reports:", safe_reports)
	fmt.Println("Safe Reports Damp 1:", safe_reports_damp_1)

	fmt.Println("Puzzle Input")
	reports = get_reports("input.txt")
	safe_reports = 0
	safe_reports_damp_1 = 0
	for _, report := range reports {
		if is_safe_with_drop(report, false) {
			safe_reports += 1
		}
		if is_safe_with_drop(report, true) {
			safe_reports_damp_1 += 1
		}
	}
	fmt.Println("Safe Reports:", safe_reports)
	fmt.Println("Safe Reports Damp 1:", safe_reports_damp_1)
}
