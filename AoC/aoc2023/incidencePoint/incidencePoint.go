package main

import (
	"bufio"
	"fmt"
	"log"
	"math/bits"
	"os"
)

const expectedDiff = 1

func processLine(line string) []bool {
	result := make([]bool, len(line))

	for idx, char := range line {
		if char == '#' {
			result[idx] = true
		}
	}

	return result
}

func getDiffMatrix(masks []uint) map[[2]int]int {
	n := len(masks)
	result := make(map[[2]int]int, (n*(n-1))>>1)

	for idxA, maskA := range masks {
		for idxB, maskB := range masks[idxA+1:] {
			result[[2]int{idxA, idxA + idxB + 1}] = bits.OnesCount(maskA ^ maskB)
		}
	}

	return result
}

func getSymIdx(diffMatrix map[[2]int]int, n int) int {
	for i := 0; i < n-1; i++ {
		acc := diffMatrix[[2]int{i, i + 1}]

		for a, b := i-1, i+2; a >= 0 && b < n && acc <= expectedDiff; a, b = a-1, b+1 {
			acc += diffMatrix[[2]int{a, b}]
		}

		if acc == expectedDiff {
			return i
		}
	}

	return -1
}

func processBatch(buffer [][]bool) uint {
	// First try horizontal
	m, n := len(buffer), len(buffer[0])
	rowMasks, colMasks := make([]uint, m), make([]uint, n)

	for i, row := range buffer {
		for j, bit := range row {
			rowMasks[i] <<= 1
			colMasks[j] <<= 1
			if bit {
				rowMasks[i] += 0x01
				colMasks[j] += 0x01
			}
		}
	}

	var candidate int

	rowDiffs := getDiffMatrix(rowMasks)
	candidate = getSymIdx(rowDiffs, m)
	if candidate != -1 {
		return 100 * uint(candidate+1)
	}

	candidate = getSymIdx(getDiffMatrix(colMasks), n)
	if candidate != -1 {
		return uint(candidate + 1)
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
