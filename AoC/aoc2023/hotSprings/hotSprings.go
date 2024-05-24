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

type DP struct {
	memory  map[[3]int]int
	springs []int
	groups  []int
}

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

/* Attempt without DP
func modAndTest(springs, groups []int, cnt *int) {
	// Check if this is a valid finished configuration
	if len(groups) == 0 {
		// Check there are no broken springs later on
		if slices.Index(springs, -1) == -1 {
			*cnt++
		}
		return
	}

	// Check for potential valid placings
	groupSize := groups[0]
	if len(springs) < groupSize {
		return
	}

	for idx := 0; idx < len(springs)-groupSize+1; {
		if springs[idx] == 1 {
			idx++
			continue
		}
		isValid := true

		for _, status := range springs[idx : idx+groupSize] {
			if status == 1 { // Spring in a good condition
				isValid = false
				break
			}
		}

		if isValid {
			if idx+groupSize == len(springs) {
				modAndTest([]int{}, groups[1:], cnt)
			} else if springs[idx+groups[0]] != -1 {
				modAndTest(springs[idx+groups[0]+1:], groups[1:], cnt)
			}
		}

		// Logic based on current idx being the first in a group of broken springs
		// (Those scenarios were tested above)
		if springs[idx] == -1 {
			return
		}

		idx++
		if !isValid {
			offset := slices.IndexFunc(
				springs[idx:],
				func(x int) bool { return x <= 0 })
			if offset == -1 {
				return
			} else {
				idx += offset
			}
		}
	}
}
*/

// Key: (springIdx, groupIdx, accumulated group length)
func (dP *DP) f(key [3]int) int {
	value, isKey := dP.memory[key]

	if isKey {
		return value
	}

	if key[0] == len(dP.springs) {
		if key[1] == len(dP.groups) && key[2] == 0 {
			return 1
		} else if key[1] == len(dP.groups)-1 && dP.groups[key[1]] == key[2] {
			return 1
		} else {
			return 0
		}
	}

	var result int

	// Assume working spring
	if dP.springs[key[0]] >= 0 {
		if key[2] == 0 {
			result += dP.f([3]int{key[0] + 1, key[1], 0})
		} else if key[1] < len(dP.groups) && dP.groups[key[1]] == key[2] {
			result += dP.f([3]int{key[0] + 1, key[1] + 1, 0})
		}
	}

	// Assume broken spring
	if dP.springs[key[0]] <= 0 {
		result += dP.f([3]int{key[0] + 1, key[1], key[2] + 1})
	}

	dP.memory[key] = result
	return result
}

func countArrangements(line string) int {
	springs, groups := getArrays(line)

	//modAndTest(springs, groups, cnt)
	dP := DP{
		memory:  make(map[[3]int]int),
		springs: duplicateArrays(springs, []int{0}, numCopies),
		groups:  duplicateArrays(groups, []int{}, numCopies),
	}

	return dP.f([3]int{0, 0, 0})
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
	c := make(chan int, n)

	wg.Add(n)
	for scanner.Scan() {
		go func(line string) {
			defer wg.Done()
			c <- countArrangements(line)
		}(scanner.Text())
	}
	wg.Wait()
	close(c)

	result := 0
	for val := range c {
		result += val
	}

	fmt.Println(result)
}
