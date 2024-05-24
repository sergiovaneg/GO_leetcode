package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/sergiovaneg/GoStudy/utils"
)

const numCopies = 4

func getArrays(line string) ([]int, []int) {
	splits := strings.Split(line, " ")

	springs := make([]int, len(splits[0]))
	for idx, char := range splits[0] {
		switch char {
		case '.':
			springs[idx] = 1
		case '?':
			springs[idx] = 0
		case '#':
			springs[idx] = -1
		}
	}

	consecutiveStr := strings.Split(splits[1], ",")
	groups := make([]int, len(consecutiveStr))
	for idx, num := range consecutiveStr {
		if val, err := strconv.Atoi(num); err == nil {
			groups[idx] = val
		}
	}

	return springs, groups
}

func duplicateArrays(arr, sep []int, n int) []int {
	nArr, nSep := len(arr), len(sep)
	newArr := make([]int, nArr+n*(nArr+nSep))
	copy(newArr, arr)

	for idxArr, idxCopy := nArr, 0; idxCopy < n; idxCopy++ {
		copy(newArr[idxArr:], sep)
		idxArr += nSep

		copy(newArr[idxArr:], arr)
		idxArr += nArr
	}

	return newArr
}

type State struct {
	dot  *State // operational
	hash *State // broken
}

type DFA struct {
	states []State
}

func NewDFA(groups []int) DFA {
	var nStates int
	for _, val := range groups {
		nStates++
		nStates += val
	}

	var dfa DFA
	dfa.states = make([]State, nStates)

	dfa.states[0].dot = &dfa.states[0]
	dfa.states[0].hash = &dfa.states[1]

	i := 1
	for _, size := range groups {
		for j := 0; j < size-1; i, j = i+1, j+1 {
			dfa.states[i].hash = &dfa.states[i+1]
		}

		if i+2 < nStates {
			dfa.states[i].dot = &dfa.states[i+1]
			i++
			dfa.states[i].dot = &dfa.states[i]
			dfa.states[i].hash = &dfa.states[i+1]
		}
		i++
	}

	dfa.states[nStates-1].dot = &dfa.states[nStates-1]

	return dfa
}

func (dfa DFA) countArrangements(springs []int) uint {
	nStates := len(dfa.states)
	current := make(map[*State]uint, nStates)

	current[&dfa.states[0]] = 1

	for _, springCondition := range springs {
		next := make(map[*State]uint, nStates)

		for key, value := range current {
			if springCondition >= 0 && key.dot != nil {
				next[key.dot] += value
			}
			if springCondition <= 0 && key.hash != nil {
				next[key.hash] += value
			}
		}

		current = next
	}

	return current[&dfa.states[nStates-1]]
}

func processLine(line string) uint {
	springs, groups := getArrays(line)
	springs = duplicateArrays(springs, []int{0}, numCopies)
	groups = duplicateArrays(groups, []int{}, numCopies)

	dfa := NewDFA(groups)

	return dfa.countArrangements(springs)
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

	var wg sync.WaitGroup
	c := make(chan uint, n)

	wg.Add(n)
	for scanner.Scan() {
		go func(line string) {
			defer wg.Done()
			c <- processLine(line)
		}(scanner.Text())
	}
	wg.Wait()
	close(c)

	var result uint
	for val := range c {
		result += val
	}

	fmt.Println(result)
}
