package main

import (
	"bufio"
	"log"
	"os"
	"slices"

	"github.com/sergiovaneg/GoStudy/utils"
)

const targetStepsFinite = 64
const targetStepsInfinite = 26501365

type Point [2]int
type Record map[Point]uint

type Garden [][]byte

func (g Garden) isValid(p Point) bool {
	if p[0] < 0 || p[0] >= len(g) {
		return false
	}
	if p[1] < 0 || p[1] >= len(g[p[0]]) {
		return false
	}
	if g[p[0]][p[1]] == '#' {
		return false
	}

	return true
}

func (g Garden) getNeighbours(p Point) []Point {
	neighbours := make([]Point, 0, 4)
	for _, offset := range []Point{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} {
		candidate := Point{p[0] + offset[0], p[1] + offset[1]}
		if g.isValid(candidate) {
			neighbours = append(neighbours, candidate)
		}
	}
	return slices.Clip(neighbours)
}

func (record *Record) populate(garden Garden, start Point) {
	queue := []struct {
		p Point
		d uint
	}{{p: start, d: 0}}

	for len(queue) > 0 {
		d, p := queue[0].d, queue[0].p
		queue = queue[1:]

		_, ok := (*record)[p]
		if ok {
			continue
		}

		(*record)[p] = d

		for _, next := range garden.getNeighbours(p) {
			if _, ok := (*record)[next]; ok {
				continue
			}
			queue = append(queue, struct {
				p Point
				d uint
			}{p: next, d: d + 1})
		}
	}
}

func (record Record) solveSinglePlot(target uint) uint {
	var result uint

	oddity := target & 0x01
	for _, dist := range record {
		if (dist&0x01) == oddity && dist <= target {
			result++
		}
	}

	return result
}

// Requires further study
func (record Record) solveInfinitePlot(g Garden, target uint) uint {
	l := uint(len(g))
	n := (target - l>>1) / l

	nOddTiles, nEvenTiles := (n+1)*(n+1), n*n

	var nOddPoints, nEvenPoints uint
	var oddCorners, evenCorners uint
	for _, dist := range record {
		if dist&0x01 == 1 {
			nOddPoints++
			if dist > l>>1 {
				oddCorners++
			}
		} else {
			nEvenPoints++
			if dist > l>>1 {
				evenCorners++
			}
		}
	}

	return nOddTiles*nOddPoints + nEvenTiles*nEvenPoints - (n+1)*oddCorners + n*evenCorners
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

	garden := make(Garden, 0, n)
	var start Point
	for i := 0; scanner.Scan(); i++ {
		line := []byte(scanner.Text())
		if j := slices.Index(line, 'S'); j != -1 {
			start = [2]int{i, j}
			line[j] = '.'
		}
		garden = append(garden, line)
	}

	record := make(Record, n*len(garden[0]))
	record.populate(garden, start)

	println(record.solveSinglePlot(targetStepsFinite))
	println(record.solveInfinitePlot(garden, targetStepsInfinite))
}
