package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

func readInput(inputFile string) map[string][]string {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	sc := bufio.NewScanner(file)
	m := make(map[string][]string)
	for sc.Scan() {
		pcs := strings.Split(sc.Text(), "-")
		s, ok := m[pcs[0]]
		if !ok {
			s = make([]string, 0)
		}
		s = append(s, pcs[1])
		m[pcs[0]] = s

		s, ok = m[pcs[1]]
		if !ok {
			s = make([]string, 0)
		}
		s = append(s, pcs[0])
		m[pcs[1]] = s
	}
	for key, s := range m {
		slices.Sort(s)
		m[key] = s
	}
	return m
}

func possibleLanParties(m map[string][]string) int {
	done := make(map[string]bool)
	count := 0
	for pc1, s := range m {
		for i, pc2 := range s {
			_, ok := done[pc2]
			if ok {
				continue
			}
			for j := i + 1; j < len(s); j++ {
				pc3 := s[j]
				_, ok = done[pc3]
				if ok {
					continue
				}
				pc3Conns := m[pc3]
				for _, pc := range pc3Conns {
					if pc == pc2 {
						if isChieftanPc(pc1, pc2, pc3) {
							count += 1
						}
						break
					}
					if pc > pc2 {
						break
					}
				}
			}
		}
		done[pc1] = true
	}
	return count
}

func isChieftanPc(s1, s2, s3 string) bool {
	return strings.HasPrefix(s1, "t") || strings.HasPrefix(s2, "t") || strings.HasPrefix(s3, "t")
}

func largestInterconnection(m map[string][]string) string {
	done := make(map[string]bool)
	longest := make([]string, 0)
	for pc, s := range m {
		doneInS := make(map[string]bool)
		for i, nextPc := range s {
			_, ok := done[nextPc]
			if ok {
				continue
			}
			_, ok = doneInS[nextPc]
			if ok {
				continue
			}
			doneInS[nextPc] = true
			if len(s)-i+1 <= len(longest) {
				// If remaining cannot make longest, break
				break
			}
			l := []string{nextPc}
			for j := i + 1; j < len(s); j++ {
				if len(l)+len(s)-i+1 <= len(longest) {
					// If remaining cannot make longest, break
					break
				}
				_, ok = doneInS[s[j]]
				if ok {
					continue
				}
				if aSubsetB(l, m[s[j]]) {
					l = append(l, s[j])
					doneInS[s[j]] = true
				}
			}
			l = append(l, pc)
			if len(l) > len(longest) {
				longest = l
			}
		}
		done[pc] = true
	}
	slices.Sort(longest)
	pwd := ""
	for i, pc := range longest {
		pwd += pc
		if i < len(longest)-1 {
			pwd += ","
		}
	}
	return pwd
}

func aSubsetB(a, b []string) bool {
	// Returns true is a is a subset of b
	// []string should be sorted alphabetically
	i := 0
	j := 0
	for i < len(a) && j < len(b) {
		if a[i] < b[j] {
			break
		}
		if a[i] == b[j] {
			i += 1
			j += 1
			continue
		}
		j += 1
	}
	return i == len(a)
}

func main() {
	fmt.Println("Test Input")
	m := readInput("test-input.txt")
	fmt.Println("Possible LAN Parties:", possibleLanParties(m))
	fmt.Println("Largest Interconnection:", largestInterconnection(m))

	fmt.Println("")

	fmt.Println("Test Input")
	m = readInput("input.txt")
	fmt.Println("Possible LAN Parties:", possibleLanParties(m))
	fmt.Println("Largest Interconnection:", largestInterconnection(m))
}
