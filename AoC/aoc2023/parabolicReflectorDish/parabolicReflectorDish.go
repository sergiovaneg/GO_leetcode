package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/sergiovaneg/GoStudy/utils"
)

const numIters = 1000000000
const bitsPerRune = 8

type Symbols [][]rune

/* Part 1
func processSymbols(symbols [][]rune) uint {
	var result uint

	m, n := len(symbols), len(symbols[0])
	count := make([]uint, n)

	for i, row := range symbols {
		for j, symbol := range row {
			switch symbol {
			case '#':
				count[j] = uint(i + 1)
			case 'O':
				result += uint(m) - count[j]
				count[j]++
			}
		}
	}

	return result
}
*/

func (symbols Symbols) encodeState() string {
	var state string

	for _, row := range symbols {
		var rowBits uint
		mask := uint(0x01)

		for idx, char := range row {
			if char == 'O' {
				rowBits += mask
			}

			if (idx+1)%bitsPerRune == 0 || idx == len(row)-1 {
				state += string(rune(rowBits))
				rowBits = 0
				mask = 0x01
			} else {
				mask <<= 1
			}
		}
	}

	return state
}

func (symbols *Symbols) loadState(state string) {
	m, n := len(*symbols), len((*symbols)[0])
	stateRunes, stateIdx := []rune(state), 0

	for i, mask := 0, uint(0x01); i < m; i++ {
		for j := 0; j < n; j++ {
			isStone := mask&uint(stateRunes[stateIdx]) != 0

			if isStone {
				(*symbols)[i][j] = 'O'
			} else if (*symbols)[i][j] != '#' {
				(*symbols)[i][j] = '.'
			}

			if (j+1)%bitsPerRune == 0 || j == n-1 {
				stateIdx++
				mask = 0x01
			} else {
				mask <<= 1
			}
		}
	}
}

func (symbols *Symbols) tilt(direction rune) {
	m, n := len(*symbols), len((*symbols)[0])

	var groups [][]uint
	var i, j int

	var outStart, outLim, outMod int
	var inStart, inLim, inMod int
	var outIdx, inIdx *int

	switch direction {
	case 'N':
		outStart, outLim, outMod = 0, n, 1
		inStart, inLim, inMod = 0, m, 1
		outIdx, inIdx = &j, &i
		groups = make([][]uint, n)
	case 'S':
		outStart, outLim, outMod = n-1, -1, -1
		inStart, inLim, inMod = m-1, -1, -1
		outIdx, inIdx = &j, &i
		groups = make([][]uint, n)
	case 'W':
		outStart, outLim, outMod = 0, m, 1
		inStart, inLim, inMod = 0, n, 1
		outIdx, inIdx = &i, &j
		groups = make([][]uint, m)
	case 'E':
		outStart, outLim, outMod = m-1, -1, -1
		inStart, inLim, inMod = n-1, -1, -1
		outIdx, inIdx = &i, &j
		groups = make([][]uint, m)
	}

	for *outIdx = outStart; *outIdx != outLim; *outIdx += outMod {
		subgroups, count := make([]uint, 0), uint(0)

		for *inIdx = inStart; *inIdx != inLim; *inIdx += inMod {
			switch (*symbols)[i][j] {
			case 'O':
				count++
			case '#':
				subgroups = append(subgroups, count)
				count = 0
			}
		}

		subgroups = append(subgroups, count)
		groups[*outIdx] = subgroups
	}

	for *outIdx = outStart; *outIdx != outLim; *outIdx += outMod {
		for *inIdx = inStart; *inIdx != inLim; *inIdx += inMod {
			if (*symbols)[i][j] == '#' {
				groups[*outIdx] = groups[*outIdx][1:]
			} else if groups[*outIdx][0] > 0 {
				(*symbols)[i][j] = 'O'
				groups[*outIdx][0]--
			} else {
				(*symbols)[i][j] = '.'
			}
		}
	}
}

func (symbols Symbols) weigh() uint {
	var result uint
	m := len(symbols)

	for i, row := range symbols {
		for _, symbol := range row {
			if symbol == 'O' {
				result += uint(m - i)
			}
		}
	}

	return result
}

func getPeriod(iterMap map[string]string, initialState string) int {
	counter, nextState := 1, iterMap[initialState]

	for nextState != initialState {
		counter++
		nextState = iterMap[nextState]
	}

	return counter
}

func (symbols Symbols) process() uint {
	iterMap := make(map[string]string, 1)
	currentState := symbols.encodeState()
	period := -1

	for iter := 0; iter < numIters; iter++ {
		nextState, ok := iterMap[currentState]

		if ok {
			if period == -1 {
				period = getPeriod(iterMap, currentState)
				iter += period*(max(numIters-1-iter, 0)/period) - 1
			} else {
				currentState = nextState
			}
		} else {
			symbols.tilt('N')
			symbols.tilt('W')
			symbols.tilt('S')
			symbols.tilt('E')
			iterMap[currentState] = symbols.encodeState()
			currentState = iterMap[currentState]
		}
	}

	symbols.loadState(currentState)
	return symbols.weigh()
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

	symbols := make(Symbols, 0, n)
	for scanner.Scan() {
		symbols = append(symbols, []rune(scanner.Text()))
	}

	fmt.Printf("Word size %v: %v\n", bitsPerRune, symbols.process())
}
