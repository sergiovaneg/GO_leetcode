package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"

	"github.com/sergiovaneg/GoStudy/utils"
)

const minSteps, maxSteps = 4, 10

type CellState struct {
	position  [2]int
	direction [2]int
	steps     int
}

type Cell struct {
	heatLoss int
	cs       CellState
}

type CellMap [][]int

type CellQueue []Cell

func (cq CellQueue) Len() int { return len(cq) }

func (cq CellQueue) Less(i, j int) bool {
	if cq[i].heatLoss != cq[j].heatLoss {
		return cq[i].heatLoss < cq[j].heatLoss
	}
	return cq[i].cs.steps > cq[j].cs.steps
}

func (cq CellQueue) Swap(i, j int) {
	cq[i], cq[j] = cq[j], cq[i]
}

func (cq *CellQueue) Push(x any) {
	*cq = append(*cq, x.(Cell))
}

func (cq *CellQueue) Pop() any {
	n := len(*cq)
	x := (*cq)[n-1]
	*cq = (*cq)[:n-1]

	return x
}

func (cm CellMap) getValidNeighbours(cell Cell) [][2]int {
	x := cell.cs.position
	direction := cell.cs.direction
	steps := cell.cs.steps

	filter := func(y [2]int) bool {
		// Check invalid position
		if y[0] < 0 || y[0] >= len(cm) {
			return true
		}
		if y[1] < 0 || y[1] >= len(cm[0]) {
			return true
		}

		// Early return for starting point
		if direction == [2]int{0, 0} {
			return false
		}

		newDir := [2]int{y[0] - x[0], y[1] - x[1]}

		// Check if going backwards
		if newDir[0] == -direction[0] && newDir[1] == -direction[1] {
			return true
		}

		// Check consecutive-step limits
		if direction == newDir {
			if steps >= maxSteps {
				return true
			}
		} else {
			if steps < minSteps {
				return true
			}
		}

		return false
	}

	// Allocate worst case
	neighbours := make([][2]int, 0, 3)
	for _, neighbour := range [][2]int{
		{x[0] + 1, x[1]},
		{x[0] - 1, x[1]},
		{x[0], x[1] + 1},
		{x[0], x[1] - 1},
	} {
		if !filter(neighbour) {
			neighbours = append(neighbours, neighbour)
		}
	}

	return slices.Clip(neighbours)
}

func (cm CellMap) calcDistances(x0, y [2]int) int {
	seen := make(map[CellState]bool, maxSteps*4*len(cm)*len(cm[0]))

	cq := CellQueue{
		{cs: CellState{position: x0}},
	}
	heap.Init(&cq)

	for len(cq) > 0 {
		x := heap.Pop(&cq).(Cell)

		if seen[x.cs] {
			continue
		}

		if x.cs.position == y && x.cs.steps >= minSteps {
			return x.heatLoss
		}

		seen[x.cs] = true

		for _, neighPos := range cm.getValidNeighbours(x) {
			newDist := x.heatLoss + cm[neighPos[0]][neighPos[1]]

			newDir := [2]int{
				neighPos[0] - x.cs.position[0],
				neighPos[1] - x.cs.position[1],
			}
			isForward := x.cs.direction == newDir

			newState := Cell{
				heatLoss: newDist,
				cs: CellState{
					position:  neighPos,
					direction: newDir,
					steps:     1,
				},
			}
			if isForward {
				newState.cs.steps += x.cs.steps
			}

			heap.Push(&cq, newState)
		}
	}

	return -1
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

	cm := make(CellMap, n)

	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()

		row := make([]int, 0, len(line))
		for _, char := range line {
			val, err := strconv.Atoi(string(char))
			if err != nil {
				val = 0
			}
			row = append(row, val)
		}

		cm[i] = row
	}

	fmt.Printf("Minimum heat to end: %v",
		cm.calcDistances([2]int{0, 0}, [2]int{len(cm) - 1, len(cm[0]) - 1}))
}
