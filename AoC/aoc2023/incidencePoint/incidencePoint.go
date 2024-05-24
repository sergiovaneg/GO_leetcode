package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func processLine(line string) []bool {
	result := make([]bool, len(line))

	for idx, char := range line {
		if char == '#' {
			result[idx] = true
		}
	}

	return result
}

func findSymmetryAxis(n1, n2 []uint, n int) int {
	shift := n & 0x01
	for mask := (uint(1) << (n - shift)) - 1; mask != 0; mask, shift = mask>>2, shift+2 {
		valid := true
		for idx, num := range n1 {
			if (num^(n2[idx]>>shift))&mask != 0 {
				valid = false
				break
			}
		}
		if valid {
			return shift
		}
	}

	return -1
}

func processBatch(buffer [][]bool) uint {
	// First try horizontal
	var candidate int
	m, n := len(buffer), len(buffer[0])
	numsFwd, numsBwd := make([]uint, m), make([]uint, m)

	for i, row := range buffer {
		for j, bit := range row {
			numsFwd[i] <<= 1
			if bit {
				numsFwd[i] += 0x01
				numsBwd[i] += 1 << j
			}
		}
	}

	candidate = findSymmetryAxis(numsFwd, numsBwd, n)
	if candidate != -1 {
		return uint(n - n>>1 + candidate>>1)
	}

	candidate = findSymmetryAxis(numsBwd, numsFwd, n)
	if candidate != -1 {
		return uint(n>>1 - candidate>>1)
	}

	numsFwd, numsBwd = make([]uint, n), make([]uint, n)

	for i, row := range buffer {
		for j, bit := range row {
			numsFwd[j] <<= 1
			if bit {
				numsFwd[j] += 0x01
				numsBwd[j] += 1 << i
			}
		}
	}

	candidate = findSymmetryAxis(numsFwd, numsBwd, m)
	if candidate != -1 {
		return 100 * uint(m-m>>1+candidate>>1)
	}

	candidate = findSymmetryAxis(numsBwd, numsFwd, m)
	if candidate != -1 {
		return 100 * uint(m>>1-candidate>>1)
	}

	return 0
}

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	buffer := make([][]bool, 0)

	var res uint
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			res += processBatch(buffer)
			buffer = make([][]bool, 0)
		} else {
			buffer = append(buffer, processLine(line))
		}
	}

	if len(buffer) > 0 {
		res += processBatch(buffer)
		buffer = nil
	}

	fmt.Println(res)
}
