package main

import (
	"bufio"
	"log"
	"os"
	"slices"
	"strconv"

	"github.com/sergiovaneg/GoStudy/utils"
)

const TARGET = 150

func getNumCombinations(
	capacities []int,
	target, currentSize int,
	minimum, counter *int) int {
	if target < 0 {
		return 0
	}

	if target == 0 {
		if currentSize < *minimum {
			*minimum = currentSize
			*counter = 1
		} else if currentSize == *minimum {
			*counter++
		}
		return 1
	}

	cnt := 0
	for idx, val := range capacities {
		cnt += getNumCombinations(
			capacities[idx+1:],
			target-val, currentSize+1,
			minimum, counter)
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

	capacities := make([]int, 0, n)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if cap, err := strconv.Atoi(scanner.Text()); err == nil {
			capacities = append(capacities, cap)
		}
	}

	slices.Sort(capacities)

	minimum, counter := n, 0
	println(getNumCombinations(capacities, TARGET, 0, &minimum, &counter))
	println(counter)
}
