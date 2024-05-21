package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
	"sync"

	"github.com/sergiovaneg/GoStudy/utils"
)

func processLine(line string) int {
	matches := regexp.MustCompile("([-|0-9]+)").FindAllString(line, -1)
	history, cache := make([]int, len(matches)), make([]int, 0, len(matches))

	for idx, match := range matches {
		if val, err := strconv.Atoi(match); err == nil {
			history[idx] = val
		}
	}

	for len(history) > 0 {
		done := true
		cache = append(cache, history[0])
		for idx, num := range history[1:] {
			history[idx] = num - history[idx]
			if done && history[idx] != 0 {
				done = false
			}
		}

		if done {
			break
		}
		history = history[:len(history)-1]
	}

	result := 0
	slices.Reverse(cache)
	for _, val := range cache {
		result = val - result
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

	var wg sync.WaitGroup
	n, _ := utils.LineCounter(file)
	c := make(chan int, n)

	wg.Add(n)
	for scanner.Scan() {
		go func(line string) {
			defer wg.Done()
			c <- processLine(line)
		}(scanner.Text())
	}
	wg.Wait()
	close(c)

	var result int
	for val := range c {
		result += val
	}

	fmt.Println(result)
}
