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

func processNumbers(arr []byte) []int {
	re, result := regexp.MustCompile("([0-9]+)"), []int{}

	matches := re.FindAllIndex(arr, -1)
	for _, match := range matches {
		val, err := strconv.Atoi(string(arr[match[0]:match[1]]))
		if err == nil {
			result = append(result, val)
		}
	}
	slices.Sort(result)

	return result
}

func getPoints(line string) int {
	arr := []byte(line)

	arr = arr[slices.Index(arr, ':'):]
	barIdx, result := slices.Index(arr, '|'), 0

	winningNums := processNumbers(arr[:barIdx])
	nums := processNumbers(arr[barIdx+1:])

	for idx0, idx1 := 0, 0; idx0 < len(winningNums) && idx1 < len(nums); {
		if winningNums[idx0] == nums[idx1] {
			result++
			idx0++
			idx1++
		} else if winningNums[idx0] > nums[idx1] {
			idx1++
		} else {
			idx0++
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
	var res int

	var wg sync.WaitGroup
	n, _ := utils.LineCounter(file)
	c, count, i := make([]int, n), make([]int, n), 0

	wg.Add(n)
	for scanner.Scan() {
		go func(line string, i int) {
			defer wg.Done()
			c[i] = getPoints(line)
		}(scanner.Text(), i)

		count[i] = 1
		i++
	}
	wg.Wait()

	for i, matches := range c {
		for j := i + 1; matches > 0 && j < n; j, matches = j+1, matches-1 {
			count[j] += count[i]
		}

		res += count[i]
		i++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
}
