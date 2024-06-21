package main

import (
	"bufio"
	"log"
	"os"

	"github.com/sergiovaneg/GoStudy/utils"
)

const SIZE = 100
const STEPS = 100

type Grid [][]int

func parseGrid(lines []string) Grid {
	g := make(Grid, SIZE+2)
	g[0], g[SIZE+1] = make([]int, SIZE+2), make([]int, SIZE+2)

	for i, line := range lines {
		g[i+1] = make([]int, SIZE+2)
		for j, char := range line {
			if char == '#' {
				g[i+1][j+1] = 1
			}
		}
	}

	return g
}

func (g Grid) toggleCorners() {
	g[1][1] = 1
	g[1][SIZE] = 1
	g[SIZE][1] = 1
	g[SIZE][SIZE] = 1
}

func (count Grid) getNeighbourCount(g Grid) {
	for i := range SIZE {
		for j := range SIZE {
			count[i][j] = g[i][j] + g[i][j+1] + g[i][j+2]
			count[i][j] += g[i+2][j] + g[i+2][j+1] + g[i+2][j+2]
			count[i][j] += g[i+1][j] + g[i+1][j+2]
		}
	}
}

func (g Grid) updateState(count Grid) {
	for i, row := range g[1 : SIZE+1] {
		for j, val := range row[1 : SIZE+1] {
			n := count[i][j]
			if val == 1 {
				if !(n == 2 || n == 3) {
					g[i+1][j+1] = 0
				}
			} else if n == 3 {
				g[i+1][j+1] = 1
			}
		}
	}
}

func (g Grid) animate(steps int) {
	count := make(Grid, SIZE)
	for i := range count {
		count[i] = make([]int, SIZE)
	}

	for range steps {
		g.toggleCorners() // Comment for part 1
		count.getNeighbourCount(g)
		g.updateState(count)
	}
	g.toggleCorners() // Comment for part 1
}

func (g Grid) countOn() int {
	cnt := 0

	for _, row := range g[1 : SIZE+1] {
		for _, val := range row[1 : SIZE+1] {
			cnt += val
		}
	}

	return cnt
}

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	n, err := utils.LineCounter(file)
	if err != nil {
		log.Fatal(err)
	}

	lines := make([]string, 0, n)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	g := parseGrid(lines)
	g.animate(STEPS)
	println(g.countOn())
}
