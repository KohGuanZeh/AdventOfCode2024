package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Pair struct {
	x int
	y int
}

type Robot struct {
	pos Pair
	vel Pair
}

const TESTFOLDER = "test-render"
const PUZZLEFOLDER = "render"

func readInput(inputFile string) []Robot {
	file, err := os.Open(inputFile)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	robots := make([]Robot, 0)
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		tokens := strings.Split(sc.Text(), " ")
		pos := strings.Split(strings.Split(tokens[0], "=")[1], ",")
		posPair := Pair{x: atoi(pos[0]), y: atoi(pos[1])}
		vel := strings.Split(strings.Split(tokens[1], "=")[1], ",")
		velPair := Pair{x: atoi(vel[0]), y: atoi(vel[1])}
		robot := Robot{pos: posPair, vel: velPair}
		robots = append(robots, robot)
	}
	return robots
}

func safetyFactor(width, height, simulations int, robots []Robot) int {
	q1, q2, q3, q4 := 0, 0, 0, 0 // Top left, top right, bottom right, bottom left.
	// Assumption that width and height are always odd.
	middleW := width / 2
	middleH := height / 2

	for _, robot := range robots {
		newX := robot.pos.x + simulations*robot.vel.x
		newX = mod(newX, width)
		newY := robot.pos.y + simulations*robot.vel.y
		newY = mod(newY, height)
		if newX == middleW || newY == middleH {
			// fmt.Println(newX, ",", newY)
			continue
		}
		if newX < middleW {
			if newY < middleH {
				q1 += 1
			} else {
				q3 += 1
			}
		} else {
			if newY < middleH {
				q2 += 1
			} else {
				q4 += 1
			}
		}
	}
	fmt.Println("Q1:", q1, "Q2:", q2, "Q3:", q3, "Q4:", q4)
	return q1 * q2 * q3 * q4
}

// Returns x mod y
func mod(x, y int) int {
	x %= y
	if x < 0 {
		x += y
	}
	return x
}

func atoi(numString string) int {
	v, err := strconv.Atoi(numString)
	if err != nil {
		log.Fatal(err)
	}
	return v
}

func simultate(robots []Robot, width, height, simulations int, saveFolder string) {
	xPos := make([]int, len(robots))
	yPos := make([]int, len(robots))
	for i := range robots {
		xPos[i] = robots[i].pos.x
		yPos[i] = robots[i].pos.y
	}
	newPositions := make([][]bool, height)
	for y := range height {
		rowDisplay := make([]bool, width)
		newPositions[y] = rowDisplay
	}
	for s := range simulations {
		for y := range height {
			for x := range width {
				newPositions[y][x] = false
			}
		}
		for i, robot := range robots {
			xPos[i] += robot.vel.x
			xPos[i] = mod(xPos[i], width)
			yPos[i] += robot.vel.y
			yPos[i] = mod(yPos[i], height)
			newPositions[yPos[i]][xPos[i]] = true
		}

		img := image.NewGray(image.Rect(0, 0, width, height))
		for y := range newPositions {
			for x := range newPositions[y] {
				if newPositions[y][x] {
					img.SetGray(x, y, color.Gray{Y: 255})
				}
			}
		}

		fileName := strconv.Itoa(s) + ".png"
		file, err := os.Create(filepath.Join(saveFolder, fileName))
		if err != nil {
			log.Fatalln(err)
		}
		png.Encode(file, img)
		file.Close()
	}
}

func main() {
	fmt.Println("Test Input")
	robots := readInput("test-input.txt")
	fmt.Println("Safety Factor:", safetyFactor(11, 7, 100, robots))
	saveFolder := filepath.Join(".", TESTFOLDER)
	err := os.RemoveAll(saveFolder)
	if err != nil {
		log.Println(err)
	}
	err = os.MkdirAll(saveFolder, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
	simultate(robots, 11, 7, 10000, saveFolder)

	fmt.Println("")

	fmt.Println("Puzzle Input")
	robots = readInput("input.txt")
	fmt.Println("Safety Factor:", safetyFactor(101, 103, 100, robots))
	saveFolder = filepath.Join(".", PUZZLEFOLDER)
	err = os.RemoveAll(saveFolder)
	if err != nil {
		log.Println(err)
	}
	err = os.MkdirAll(saveFolder, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
	simultate(robots, 101, 103, 10000, saveFolder)
}
