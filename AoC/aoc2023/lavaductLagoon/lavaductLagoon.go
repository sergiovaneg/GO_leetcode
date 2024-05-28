package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/sergiovaneg/GoStudy/utils"
)

type Instruction struct {
	direction [2]int
	length    int
	colour    string
}

func parseInstruction(line string) Instruction {
	fields := strings.Split(line, " ")

	var direction [2]int
	switch fields[0] {
	case "U":
		direction = [2]int{-1, 0}
	case "D":
		direction = [2]int{1, 0}
	case "L":
		direction = [2]int{0, -1}
	case "R":
		direction = [2]int{0, 1}
	}

	length, _ := strconv.Atoi(fields[1])

	idx := strings.Index(fields[2], "#")
	hexStr := fields[2][idx+1 : idx+7]

	return Instruction{
		direction: direction,
		length:    length,
		colour:    hexStr,
	}
}

func getCapacity(instructionSet []Instruction) int {
	// Using Shoelace formula and Pick's theorem
	x := [2]int{0, 0}
	var perimeter, area int

	for _, inst := range instructionSet {
		perimeter += inst.length

		y := [2]int{
			x[0] + inst.length*inst.direction[0],
			x[1] + inst.length*inst.direction[1],
		}
		area += x[0]*y[1] - x[1]*y[0]
		x = y
	}
	if area < 0 {
		area = -area
	}
	result := perimeter>>1 + area>>1 + 1

	return result
}

func patchInstruction(instruction *Instruction) {
	switch instruction.colour[5] {
	case '0':
		instruction.direction = [2]int{0, 1}
	case '1':
		instruction.direction = [2]int{1, 0}
	case '2':
		instruction.direction = [2]int{0, -1}
	case '3':
		instruction.direction = [2]int{-1, 0}
	}

	newLength, _ := strconv.ParseInt(instruction.colour[:5], 16, 0)
	instruction.length = int(newLength)
}

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	n, err := utils.LineCounter(file)
	if err != nil {
		log.Fatal(err)
	}

	instructionSet := make([]Instruction, 0, n)
	for scanner.Scan() {
		instructionSet = append(instructionSet, parseInstruction(scanner.Text()))
	}

	fmt.Println(getCapacity(instructionSet))
	for idx := range instructionSet {
		patchInstruction(&instructionSet[idx])
	}
	fmt.Println(getCapacity(instructionSet))
}
