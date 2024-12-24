package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

const XOR = "XOR"
const OR = "OR"
const AND = "AND"

func readInput(inputFile string) (map[string]bool, map[string][3]string) {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	sc := bufio.NewScanner(file)
	valMap := make(map[string]bool)
	for sc.Scan() {
		if len(sc.Text()) == 0 {
			break
		}
		tokens := strings.Split(sc.Text(), ": ")
		key := tokens[0]
		val := tokens[1] == "1"
		valMap[key] = val
	}

	exprMap := make(map[string][3]string)
	for sc.Scan() {
		tokens := strings.Split(sc.Text(), " -> ")
		key := tokens[1]
		tokens = strings.Split(tokens[0], " ")
		exprMap[key] = [3]string{tokens[0], tokens[1], tokens[2]}
	}
	return valMap, exprMap
}

func zValue(valMap map[string]bool, exprMap map[string][3]string) int {
	i := 0
	result := 0
	for {
		zKey := getIthKey("z", i)
		_, ok := exprMap[zKey]
		if !ok {
			// No more target z values
			break
		}
		v := evalExpr(zKey, valMap, exprMap)
		if v {
			result += 1 << i
		}
		i += 1
	}
	return result
}

func evalExpr(key string, valMap map[string]bool, exprMap map[string][3]string) bool {
	v, ok := valMap[key]
	if ok {
		return v
	}
	expr := exprMap[key]
	x1 := evalExpr(expr[0], valMap, exprMap)
	x2 := evalExpr(expr[2], valMap, exprMap)
	switch expr[1] {
	case XOR:
		v = x1 != x2
		break
	case OR:
		v = x1 || x2
	case AND:
		v = x1 && x2
	default:
		log.Fatalln("Error. Invalid expression detected.")
		break
	}
	valMap[key] = v
	return v
}

func getIthKey(prefix string, i int) string {
	key := prefix
	if i < 10 {
		key += "0"
	}
	return key + strconv.Itoa(i)
}

func findWrongGates(exprMap map[string][3]string) {
	i := 0
	for {
		zKey := getIthKey("z", i+1)
		_, ok := exprMap[zKey]
		if !ok {
			// Is last key
			break
		}

		zKey = getIthKey("z", i)
		expr, ok := exprMap[zKey]
		if expr[1] != XOR {
			fmt.Println(fmtExpr(zKey, expr))
		}
		i += 1
	}

	fmt.Println()

	for key, expr := range exprMap {
		if strings.HasPrefix(key, "z") {
			continue
		}
		if strings.HasPrefix(expr[0], "x") || strings.HasPrefix(expr[0], "y") {
			continue
		}
		if expr[1] == XOR {
			fmt.Println(fmtExpr(key, expr))
			fmt.Println(fmtExpr(expr[0], exprMap[expr[0]]))
			fmt.Println(fmtExpr(expr[2], exprMap[expr[2]]))
			fmt.Println()
		}
	}
}

func fmtExpr(key string, expr [3]string) string {
	return key + ": " + expr[0] + " " + expr[1] + " " + expr[2]
}

func updateMap(exprMap map[string][3]string, keysToSwap [][2]string) {
	for _, keys := range keysToSwap {
		exprMap[keys[0]], exprMap[keys[1]] = exprMap[keys[1]], exprMap[keys[0]]
	}
}

func logExprMapOperations(outFile string, exprMap map[string][3]string) {
	f, err := os.Create(outFile)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	for i := 0; i >= 0; i++ {
		zKey := getIthKey("z", i+1)
		_, ok := exprMap[zKey]
		if !ok {
			zKey := getIthKey("z", i)
			expr := exprMap[zKey]
			f.WriteString(fmtExpr(zKey, expr) + "\n")
			k1, k2 := expr[0], expr[2]
			expr1, expr2 := exprMap[k1], exprMap[k2]
			f.WriteString(fmtExpr(k1, expr1) + "\n")
			f.WriteString(fmtExpr(k2, expr2) + "\n")
			break
		}

		zKey = getIthKey("z", i)
		expr := exprMap[zKey]
		f.WriteString(fmtExpr(zKey, expr) + "\n")
		if i == 0 {
			f.WriteString("\n\n")
			continue
		}

		f.WriteString("\n")

		k1, k2 := expr[0], expr[2]
		expr1, expr2 := exprMap[k1], exprMap[k2]
		f.WriteString(fmtExpr(k1, expr1) + "\n")
		if expr1[1] == OR {
			k3, k4 := expr1[0], expr1[2]
			expr3, expr4 := exprMap[k3], exprMap[k4]
			f.WriteString(fmtExpr(k3, expr3) + "\n")
			f.WriteString(fmtExpr(k4, expr4) + "\n")
		}

		f.WriteString("\n")

		f.WriteString(fmtExpr(k2, expr2) + "\n")
		if expr2[1] == OR {
			k3, k4 := expr2[0], expr2[2]
			expr3, expr4 := exprMap[k3], exprMap[k4]
			f.WriteString(fmtExpr(k3, expr3) + "\n")
			f.WriteString(fmtExpr(k4, expr4) + "\n")
		}
		f.WriteString("\n\n")
	}
}

func main() {
	fmt.Println("Test Input")
	valMap, exprMap := readInput("test-input.txt")
	fmt.Println("Z value:", zValue(valMap, exprMap))

	fmt.Println("")

	fmt.Println("Puzzle Input")
	valMap, exprMap = readInput("input.txt")
	fmt.Println("Z value:", zValue(valMap, exprMap))
	findWrongGates(exprMap)
	updateMap(exprMap, [][2]string{{"hmt", "z18"}, {"bfq", "z27"}, {"hkh", "z31"}})
	logExprMapOperations("out.txt", exprMap)

	s := []string{"bng", "fjp", "bfq", "hkh", "hmt", "z18", "z27", "z31"}
	slices.Sort(s)
	fmt.Print("Wrong Keys: ")
	for i, gate := range s {
		fmt.Print(gate)
		if i < len(s)-1 {
			fmt.Print(",")
		}
	}
}
