package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/sergiovaneg/GoStudy/utils"
)

func getArrays(line string) ([]int, []int) {
	splits := strings.Split(line, " ")

	springs := make([]int, len(splits[0]))
	for idx, char := range splits[0] {
		switch char {
		case '.':
			springs[idx] = 1
		case '#':
			springs[idx] = -1
		case '?':
			springs[idx] = 0
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

func modAndTest(springs, groups []int, cnt *int) {
	// Check if this is a valid finished configuration
	if len(groups) == 0 {
		if slices.Index(springs, -1) == -1 {
			*cnt++
		}
		return
	}

	// Check for potential valid

}

func countArrangements(line string) int {
	springs, groups := getArrays(line)
	var cnt int

	modAndTest(springs, groups, &cnt)

	return cnt
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
