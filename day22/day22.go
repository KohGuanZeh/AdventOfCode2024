package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func readInput(inputFile string) []int {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	sc := bufio.NewScanner(file)
	secrets := make([]int, 0)
	for sc.Scan() {
		v, err := strconv.Atoi(sc.Text())
		if err != nil {
			log.Fatalln(err)
		}
		secrets = append(secrets, v)
	}
	return secrets
}

func mix(a, b int) int {
	return a ^ b
}

func prune(a int) int {
	return a % 16777216
}

func getNextSecret(secret int) int {
	secret = prune(mix(secret, secret<<6))
	secret = prune(mix(secret, secret>>5))
	secret = prune(mix(secret, secret<<11))
	return secret
}

func sumIthSecret(secrets []int, i int) int {
	copy := make([]int, len(secrets))
	for j := range secrets {
		copy[j] = secrets[j]
	}
	for range i {
		for j := range copy {
			copy[j] = getNextSecret(copy[j])
		}
	}
	res := 0
	for j := range copy {
		res += copy[j]
	}
	return res
}

func getPrice(secret int) int {
	return secret - ((secret / 10) * 10)
}

func maximizeBananas(secrets []int, i int) int {
	copy := make([]int, len(secrets))
	for j := range secrets {
		copy[j] = secrets[j]
	}
	seqMap := make(map[string]map[int]int)
	prices := make([][]int, len(copy))
	changes := make([][]int, len(copy))
	for j := range i + 1 {
		for k := range copy {
			if j == 0 {
				prices[k] = make([]int, i+1)
				changes[k] = make([]int, i)
				prices[k][0] = getPrice(copy[k])
				continue
			}
			copy[k] = getNextSecret(copy[k])
			prices[k][j] = getPrice(copy[k])
			changes[k][j-1] = prices[k][j] - prices[k][j-1]
			if j > 3 {
				seqKey := getSeqKey(changes[k][j-4], changes[k][j-3], changes[k][j-2], changes[k][j-1])
				_, ok := seqMap[seqKey]
				if !ok {
					seqMap[seqKey] = make(map[int]int)
				}
				_, ok = seqMap[seqKey][k]
				if !ok {
					seqMap[seqKey][k] = prices[k][j]
				}
			}
		}
	}
	maxBananas := 0
	for _, m := range seqMap {
		val := 0
		for _, v := range m {
			val += v
		}
		if val > maxBananas {
			maxBananas = val
		}
	}
	return maxBananas
}

func getSeqKey(a, b, c, d int) string {
	return strconv.Itoa(a) + strconv.Itoa(b) + strconv.Itoa(c) + strconv.Itoa(d)
}

func main() {
	fmt.Println("Test Input")
	secrets := readInput("test-input.txt")
	fmt.Println("Sum of 2000th Secret:", sumIthSecret(secrets, 2000))
	fmt.Println("Maximize Bananas (2000 Secrets):", maximizeBananas(secrets, 2000))

	fmt.Println("")

	fmt.Println("Test Input 1")
	secrets = readInput("test-input-1.txt")
	fmt.Println("Sum of 2000th Secret:", sumIthSecret(secrets, 2000))
	fmt.Println("Maximize Bananas (2000 Secrets):", maximizeBananas(secrets, 2000))

	fmt.Println("")

	fmt.Println("Puzzle Input")
	secrets = readInput("input.txt")
	fmt.Println("Sum of 2000th Secret:", sumIthSecret(secrets, 2000))
	fmt.Println("Maximize Bananas (2000 Secrets):", maximizeBananas(secrets, 2000))
}
