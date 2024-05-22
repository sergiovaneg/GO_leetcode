package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/sergiovaneg/GoStudy/utils"
)

const emptyDistance = 1000000 // set to 2 for part 1

func countEmpty(flags []bool) int {
	var res int
	for _, flag := range flags {
		if flag {
			res++
		}
	}
	return res
}

func manhattanDistance(a, b [2]int) int {
	var dist int

	if a[0] > b[0] {
		dist += a[0] - b[0]
	} else {
		dist += b[0] - a[0]
	}

	if a[1] > b[1] {
		dist += a[1] - b[1]
	} else {
		dist += b[1] - a[1]
	}

	return dist
}

func getSumDistances(stars [][2]int, rowFlags, colFlags []bool) int {
	var res int

	for i, a := range stars {
		for _, b := range stars[i+1:] {
			res += manhattanDistance(a, b)

			emptyRowCnt := countEmpty(rowFlags[min(a[0], b[0]):max(a[0], b[0])])
			emptyColCnt := countEmpty(colFlags[min(a[1], b[1]):max(a[1], b[1])])

			res += (emptyDistance - 1) * (emptyRowCnt + emptyColCnt)
		}
	}

	return res
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

	stars := make([][2]int, 0)
	emptyRows := make([]bool, n)
	var emptyCols []bool

	for rowIdx := 0; scanner.Scan(); rowIdx++ {
		line := []rune(scanner.Text())

		if slices.Index(line, '#') == -1 {
			emptyRows[rowIdx] = true
		} else {
			if emptyCols == nil {
				emptyCols = make([]bool, len(line))
				for colIdx := range emptyCols {
					emptyCols[colIdx] = true
				}
			}

			for colIdx, char := range line {
				if char == '#' {
					stars = append(stars, [2]int{rowIdx, colIdx})
					if emptyCols[colIdx] {
						emptyCols[colIdx] = false
					}

				}
			}
		}
	}

	fmt.Println(getSumDistances(stars, emptyRows, emptyCols))
}
