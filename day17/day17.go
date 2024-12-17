package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func readInput(inputFile string) (int, int, int, []int) {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	regA, regB, regC := 0, 0, 0
	prog := make([]int, 0)
	c := 0
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		c += 1
		text := sc.Text()
		if len(text) == 0 {
			continue
		}
		tokens := strings.Split(text, ": ")
		switch c {
		case 1:
			regA, err = strconv.Atoi(tokens[1])
			if err != nil {
				log.Fatalln(err)
			}
			break
		case 2:
			regB, err = strconv.Atoi(tokens[1])
			if err != nil {
				log.Fatalln(err)
			}
			break
		case 3:
			regC, err = strconv.Atoi(tokens[1])
			if err != nil {
				log.Fatalln(err)
			}
			break
		default:
			tokens = strings.Split(tokens[1], ",")
			for _, token := range tokens {
				intToken, err := strconv.Atoi(token)
				if err != nil {
					log.Fatalln(err)
				}
				prog = append(prog, intToken)
			}
		}
	}
	return regA, regB, regC, prog
}

func runProg(regA, regB, regC int, prog []int) string {
	out := ""

	for i := 0; i < len(prog); i += 2 {
		opCode := prog[i]
		switch opCode {
		case 0:
			operand := getComboOperand(prog[i+1], regA, regB, regC)
			regA = regA >> operand
			break
		case 1:
			operand := prog[i+1]
			regB = regB ^ operand
			break
		case 2:
			operand := getComboOperand(prog[i+1], regA, regB, regC)
			regB = operand & 7
			break
		case 3:
			if regA == 0 {
				continue
			}
			i = prog[i+1]
			i -= 2 // If jump is successful, do not += 2
			break
		case 4:
			// Operand is read but ignored
			regB = regB ^ regC
			break
		case 5:
			operand := getComboOperand(prog[i+1], regA, regB, regC)
			operand &= 7
			if len(out) != 0 {
				out += ","
			}
			out += strconv.Itoa(operand)
			break
		case 6:
			operand := getComboOperand(prog[i+1], regA, regB, regC)
			regB = regA >> operand
			break
		case 7:
			operand := getComboOperand(prog[i+1], regA, regB, regC)
			regC = regA >> operand
			break
		default:
			log.Println("Error. Invalid opcode in program sequence.")
			break
		}
	}
	return out
}

func getComboOperand(code, a, b, c int) int {
	if code <= 3 {
		return code
	}
	switch code {
	case 4:
		return a
	case 5:
		return b
	case 6:
		return c
	default:
		log.Fatalln("Invalid combo operand value given.")
	}
	return 0
}

func revEng(prog []int) int {
	// Given Sequence:
	// 2,4,1,1,7,5,1,5,4,3,5,5,0,3,3,0
	// 1) B = A & 7
	// 2) B = B ^ 1
	// 3) C = A >> B
	// 4) B = B ^ 5
	// 5) B = B ^ C
	// 6) Print B & 7
	// 7) A = A >> 3
	// 8) Loop from 1) if A > 0
	regAs := []int{0}
	for i := len(prog) - 1; i >= 0; i-- {
		newRegAs := make([]int, 0)
		for _, regA := range regAs {
			regA = regA << 3
			for a := range 8 {
				resA := regA + a
				regB := a ^ 1
				regC := resA >> regB
				regB = regB ^ 5 ^ regC
				if regB&7 == prog[i] {
					newRegAs = append(newRegAs, resA)
				}
			}
		}
		regAs = newRegAs
	}
	return regAs[0]
}

func revEngTest(prog []int) int {
	// Given Sequence:
	// 0,3,5,4,3,0
	// 1) A >> 3
	// 2) Print A & 7
	// 3) Loop from 1) if A > 0
	regA := 0
	for i := len(prog) - 1; i >= 0; i-- {
		for a := range 8 {
			resA := regA + a
			if resA&7 == prog[i] {
				regA = resA
				break
			}
		}
		regA = regA << 3
	}
	return regA
}

func main() {
	fmt.Println("Test Input")
	regA, regB, regC, prog := readInput("test-input.txt")
	fmt.Println("Output:", runProg(regA, regB, regC, prog))

	fmt.Println("")

	fmt.Println("Test Input 1")
	regA, regB, regC, prog = readInput("test-input-1.txt")
	fmt.Println("Output:", runProg(regA, regB, regC, prog))
	revEngInput := revEngTest(prog)
	fmt.Println("Register A Value:", revEngInput)
	fmt.Println("Output w Reversed Engineering Input:", runProg(revEngInput, regB, regC, prog))

	fmt.Println("")

	fmt.Println("Puzzle Input")
	regA, regB, regC, prog = readInput("input.txt")
	fmt.Println("Output:", runProg(regA, regB, regC, prog))
	revEngInput = revEng(prog)
	fmt.Println("Register A Value:", revEngInput)
	fmt.Println("Output w Reversed Engineering Input:", runProg(revEngInput, regB, regC, prog))
}
