package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"sync"

	"github.com/sergiovaneg/GoStudy/utils"
)

func processPattern(line string) []bool {
	result := make([]bool, len(line))

	for idx, char := range line {
		if char == 'R' {
			result[idx] = true
		}
	}

	return result
}

func processNode(line string) [3][3]byte {
	matches := regexp.MustCompile("([A-Z|0-9]{3})").FindAllString(line, 3)

	var result [3][3]byte
	for idx, match := range matches {
		result[idx] = [3]byte{match[0], match[1], match[2]}
	}

	return result
}

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var wg sync.WaitGroup
	var mu sync.Mutex

	n, _ := utils.LineCounter(file)
	n -= 2

	scanner.Scan()
	pattern := processPattern(scanner.Text())
	scanner.Scan()

	nodes := make(map[[3]byte][2][3]byte, n)
	positions := make([][3]byte, 0)

	// Get node information from file
	wg.Add(n)
	for scanner.Scan() {
		go func(line string) {
			defer wg.Done()
			node := processNode(line)

			mu.Lock()
			nodes[node[0]] = [2][3]byte{node[1], node[2]}
			if node[0][2] == 'A' {
				positions = append(positions, node[0])
			}
			mu.Unlock()
		}(scanner.Text())
	}
	wg.Wait()

	cycleIdx, posLen := 0, len(positions)
	periods, steps := make([]int, posLen), 0

	for cycleIdx < posLen {
		for _, goRight := range pattern {
			steps++
			for posIdx := 0; posIdx < len(positions); {
				pos := positions[posIdx]
				if goRight {
					positions[posIdx] = nodes[pos][1]
				} else {
					positions[posIdx] = nodes[pos][0]
				}

				if positions[posIdx][2] == 'Z' {
					periods[cycleIdx] = steps
					cycleIdx++
					positions = slices.Delete(positions, posIdx, posIdx+1)
				} else {
					posIdx++
				}
			}
		}
	}

	fmt.Println(utils.LCM(periods[0], periods[1], periods[2:]...))
}
