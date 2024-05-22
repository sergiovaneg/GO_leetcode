package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"sync"

	"github.com/sergiovaneg/GoStudy/utils"
)

type Maze [][]rune

func canGo(char, direction rune) bool {
	aux := string(char)
	switch direction {
	case 'N':
		return strings.Contains("|LJ", aux)
	case 'S':
		return strings.Contains("|7F", aux)
	case 'E':
		return strings.Contains("-LF", aux)
	case 'W':
		return strings.Contains("-J7", aux)
	default:
		return false
	}
}

func (maze Maze) getNext(x, x0 [2]int) [2]int {
	char := maze[x[0]][x[1]]
	if canGo(char, 'N') && x0[0] != x[0]-1 {
		return [2]int{x[0] - 1, x[1]}
	}
	if canGo(char, 'S') && x0[0] != x[0]+1 {
		return [2]int{x[0] + 1, x[1]}
	}
	if canGo(char, 'W') && x0[1] != x[1]-1 {
		return [2]int{x[0], x[1] - 1}
	}
	if canGo(char, 'E') && x0[1] != x[1]+1 {
		return [2]int{x[0], x[1] + 1}
	}
	return [2]int{}
}

func (maze Maze) initExplorers(x0 [2]int) [2][2]int {
	m, n, idx := len(maze), len(maze[0]), 0
	var explorers [2][2]int
	var newS [2]rune

	if x0[0] > 0 {
		char := maze[x0[0]-1][x0[1]]
		if canGo(char, 'S') {
			explorers[idx] = [2]int{x0[0] - 1, x0[1]}
			newS[idx] = 'N'
			idx++
		}
	}

	if x0[1] > 0 {
		char := maze[x0[0]][x0[1]-1]
		if canGo(char, 'E') {
			explorers[idx] = [2]int{x0[0], x0[1] - 1}
			newS[idx] = 'W'
			idx++
		}
	}

	if x0[0] < m-1 {
		char := maze[x0[0]+1][x0[1]]
		if canGo(char, 'N') {
			explorers[idx] = [2]int{x0[0] + 1, x0[1]}
			newS[idx] = 'S'
			idx++
		}
	}

	if x0[1] < n-1 {
		char := maze[x0[0]][x0[1]+1]
		if canGo(char, 'W') {
			explorers[idx] = [2]int{x0[0], x0[1] + 1}
			newS[idx] = 'E'
			idx++
		}
	}

	switch newS {
	case [2]rune{'N', 'W'}:
		maze[x0[0]][x0[1]] = 'J'
	case [2]rune{'N', 'S'}:
		maze[x0[0]][x0[1]] = '|'
	case [2]rune{'N', 'E'}:
		maze[x0[0]][x0[1]] = 'L'
	case [2]rune{'W', 'S'}:
		maze[x0[0]][x0[1]] = '7'
	case [2]rune{'W', 'E'}:
		maze[x0[0]][x0[1]] = '-'
	case [2]rune{'S', 'E'}:
		maze[x0[0]][x0[1]] = 'F'
	}

	return explorers
}

func (maze Maze) exploreMaze(x0 [2]int) ([][]bool, int) {
	m, n := len(maze), len(maze[0])

	path := make([][]bool, m)
	for i := range path {
		path[i] = make([]bool, n)
	}
	path[x0[0]][x0[1]] = true

	var explorers [2][2]int
	explorers = maze.initExplorers(x0)

	path[explorers[0][0]][explorers[0][1]] = true
	path[explorers[1][0]][explorers[1][1]] = true

	steps := 1
	prev := [2][2]int{x0, x0}
	for explorers[0] != explorers[1] {
		prev, explorers = explorers, [2][2]int{
			maze.getNext(explorers[0], prev[0]),
			maze.getNext(explorers[1], prev[1]),
		}

		path[explorers[0][0]][explorers[0][1]] = true
		path[explorers[1][0]][explorers[1][1]] = true

		steps++
	}

	return path, steps
}

func processMazeRow(mazeRow []rune, pathRow []bool) {
	inside := false
	flagBuffer := [2]bool{false, false}

	for j, isPath := range pathRow {
		if isPath {
			flagBuffer[0] = flagBuffer[0] != canGo(mazeRow[j], 'N')
			flagBuffer[1] = flagBuffer[1] != canGo(mazeRow[j], 'S')
			if flagBuffer[0] && flagBuffer[1] {
				/*
					I have only gone inside the polygon once I have encountered a pipe
					going up and another one going down (projected as a vertical pipe).
					Consecutive pipes going in the same direction take me out of the
					enclosure.
				*/
				inside = !inside
				flagBuffer[0] = false
				flagBuffer[1] = false
			}
		} else if inside {
			mazeRow[j] = 'I'
		}
	}
}

func (maze Maze) markEnclosed(path [][]bool) {
	// Using Winding-Number algorithm
	var wg sync.WaitGroup

	wg.Add(len(maze))
	for i := range path {
		go func(i int) {
			defer wg.Done()
			processMazeRow(maze[i], path[i])
		}(i)
	}
	wg.Wait()
}

func (maze Maze) countEnclosed() int {
	var res int

	for _, row := range maze {
		for _, char := range row {
			if char == 'I' {
				res++
			}
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
	maze := make(Maze, 0, n)
	var x0 [2]int

	for scanner.Scan() {
		line := []rune(scanner.Text())
		if idx := slices.Index(line, 'S'); idx != -1 {
			x0 = [2]int{len(maze), idx}
		}
		maze = append(maze, line)
	}

	path, steps := maze.exploreMaze(x0)

	fmt.Printf("Steps: %v\n", steps)

	maze.markEnclosed(path)

	fmt.Printf("Enclosed: %v\n", maze.countEnclosed())
}
