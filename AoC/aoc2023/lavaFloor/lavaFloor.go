package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/sergiovaneg/GoStudy/utils"
)

/* Directions:
- 0: East
- 1: North
- 2: West
- 3: South
*/

type Tile struct {
	tileType    rune
	isEnergized bool
	isCrossed   [4]bool
	mu          [4]sync.Mutex
}

type Layout [][]Tile

func (layout Layout) forward(state [3]int) ([3]int, bool) {
	switch state[2] {
	case 0:
		if state[1] < len(layout[0])-1 {
			return [3]int{state[0], state[1] + 1, state[2]}, true
		} else {
			return state, false
		}
	case 1:
		if state[0] > 0 {
			return [3]int{state[0] - 1, state[1], state[2]}, true
		} else {
			return state, false
		}
	case 2:
		if state[1] > 0 {
			return [3]int{state[0], state[1] - 1, state[2]}, true
		} else {
			return state, false
		}
	case 3:
		if state[0] < len(layout)-1 {
			return [3]int{state[0] + 1, state[1], state[2]}, true
		} else {
			return state, false
		}
	}

	return state, false
}

func (layout Layout) String() string {
	var out string

	for i := range layout {
		for j := range layout[i] {
			if layout[i][j].isEnergized {
				out += "#"
			} else {
				out += "."
			}
		}
		out += "\n"
	}

	return out
}

func (layout Layout) advance(state [3]int) [][3]int {
	i, j, currentDir := state[0], state[1], state[2]
	newPaths := make([][3]int, 0, 2)

	// Path-uniqueness logic
	layout[i][j].mu[currentDir].Lock()
	if layout[i][j].isCrossed[currentDir] {
		layout[i][j].mu[currentDir].Unlock()
		return newPaths
	}
	layout[i][j].isCrossed[currentDir] = true
	layout[i][j].mu[currentDir].Unlock()

	// Mark as energized
	layout[i][j].isEnergized = true

	var path [3]int
	var isValid bool

	switch layout[i][j].tileType {
	case '|':
		if currentDir == 0 || currentDir == 2 {
			path, isValid = layout.forward([3]int{i, j, 1})
			if isValid {
				newPaths = append(newPaths, path)
			}

			path, isValid = layout.forward([3]int{i, j, 3})
			if isValid {
				newPaths = append(newPaths, path)
			}
		} else {
			path, isValid = layout.forward(state)
			if isValid {
				newPaths = append(newPaths, path)
			}
		}
	case '-':
		if currentDir == 1 || currentDir == 3 {
			path, isValid = layout.forward([3]int{i, j, 0})
			if isValid {
				newPaths = append(newPaths, path)
			}

			path, isValid = layout.forward([3]int{i, j, 2})
			if isValid {
				newPaths = append(newPaths, path)
			}
		} else {
			path, isValid = layout.forward(state)
			if isValid {
				newPaths = append(newPaths, path)
			}
		}
	case '/':
		switch currentDir {
		case 0:
			path, isValid = layout.forward([3]int{i, j, 1})
			if isValid {
				newPaths = append(newPaths, path)
			}
		case 1:
			path, isValid = layout.forward([3]int{i, j, 0})
			if isValid {
				newPaths = append(newPaths, path)
			}
		case 2:
			path, isValid = layout.forward([3]int{i, j, 3})
			if isValid {
				newPaths = append(newPaths, path)
			}
		case 3:
			path, isValid = layout.forward([3]int{i, j, 2})
			if isValid {
				newPaths = append(newPaths, path)
			}
		}
	case '\\':
		switch currentDir {
		case 0:
			path, isValid = layout.forward([3]int{i, j, 3})
			if isValid {
				newPaths = append(newPaths, path)
			}
		case 1:
			path, isValid = layout.forward([3]int{i, j, 2})
			if isValid {
				newPaths = append(newPaths, path)
			}
		case 2:
			path, isValid = layout.forward([3]int{i, j, 1})
			if isValid {
				newPaths = append(newPaths, path)
			}
		case 3:
			path, isValid = layout.forward([3]int{i, j, 0})
			if isValid {
				newPaths = append(newPaths, path)
			}
		}
	case '.':
		path, isValid = layout.forward(state)
		if isValid {
			newPaths = append(newPaths, path)
		}
	}

	return newPaths
}

func (layout Layout) cast(state [3]int) {
	var wg sync.WaitGroup
	for {
		newPaths := layout.advance(state)

		if len(newPaths) == 0 {
			break
		}

		if len(newPaths) == 2 {
			wg.Add(1)
			func(state [3]int) {
				defer wg.Done()
				layout.cast(state)
			}(newPaths[1])
		}

		state = newPaths[0]
	}
	wg.Wait()
}

func (layout Layout) countEnergized() uint {
	var res uint

	for i := range layout {
		for j := range layout[i] {
			if layout[i][j].isEnergized {
				res++
			}
		}
	}

	return res
}

func createLayout(text []string) Layout {
	layout := make(Layout, 0, len(text))

	for _, line := range text {
		row := make([]Tile, len(line))

		for idx, char := range line {
			row[idx].tileType = char
		}

		layout = append(layout, row)
	}

	return layout
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

	text := make([]string, 0, n)
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	var maximum uint

	// Left->Right
	for i := 0; i < len(text); i++ {
		layout := createLayout(text)
		layout.cast([3]int{i, 0, 0})

		if energized := layout.countEnergized(); energized > maximum {
			maximum = energized
		}
	}

	// Top->Bottom
	for j := 0; j < len(text[0]); j++ {
		layout := createLayout(text)
		layout.cast([3]int{0, j, 3})

		if energized := layout.countEnergized(); energized > maximum {
			maximum = energized
		}
	}

	// Right->Left
	for i := 0; i < len(text); i++ {
		layout := createLayout(text)
		layout.cast([3]int{i, len(layout[0]) - 1, 2})

		if energized := layout.countEnergized(); energized > maximum {
			maximum = energized
		}
	}

	// Bottom->Top
	for j := 0; j < len(text[0]); j++ {
		layout := createLayout(text)
		layout.cast([3]int{len(layout) - 1, j, 3})

		if energized := layout.countEnergized(); energized > maximum {
			maximum = energized
		}
	}

	fmt.Println(maximum)
}
