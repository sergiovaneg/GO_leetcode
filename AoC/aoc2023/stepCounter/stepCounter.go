package main

import (
	"bufio"
	"log"
	"os"
	"slices"

	"github.com/sergiovaneg/GoStudy/utils"
)

const targetSteps = 6

type State struct {
	i, j  int
	steps int
}

type Record map[State]bool

type Garden [][]byte

func (garden Garden) isValid(state State) bool {
	i, j := state.i, state.j
	if i < 0 || i >= len(garden) {
		return false
	}
	if j < 0 || j >= len(garden[i]) {
		return false
	}

	if garden[i][j] == '#' {
		return false
	}

	return true
}

func (s State) genNext(garden Garden) []State {
	next := make([]State, 0, 4)
	for _, shift := range [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} {
		candidate := State{
			i:     s.i + shift[0],
			j:     s.j + shift[1],
			steps: s.steps + 1,
		}
		if garden.isValid(candidate) {
			next = append(next, candidate)
		}
	}
	return next
}

func (record *Record) countReachablePlots(garden Garden, start State) int {
	var result int

	queue := []State{start}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if (*record)[current] {
			continue
		}

		(*record)[current] = true
		if current.steps == targetSteps {
			result++
		} else {
			next := current.genNext(garden)
			queue = append(queue, next...)
		}
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

	n, err := utils.LineCounter(file)
	if err != nil {
		log.Fatal(err)
	}

	garden := make(Garden, 0, n)
	var start State
	for i := 0; scanner.Scan(); i++ {
		line := []byte(scanner.Text())
		if j := slices.Index(line, 'S'); j != -1 {
			start.i, start.j = i, j
			line[j] = '.'
		}
		garden = append(garden, line)
	}

	record := make(Record, 64*n*len(garden[0]))
	println(record.countReachablePlots(garden, start))
}
