package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"sync"

	"github.com/sergiovaneg/GoStudy/utils"
)

func isNice(line string) bool {
	matches := regexp.MustCompile("([aeiou])").FindAllString(line, 3)
	if len(matches) < 3 {
		return false
	}

	flag := false
	for idx := 0; idx < len(line)-1; idx++ {
		if line[idx] == line[idx+1] {
			flag = true
			break
		}
	}
	if !flag {
		return false
	}

	return regexp.MustCompile("(ab|cd|pq|xy)").FindString(line) == ""
}

func isNiceAlt(line string) bool {
	n := len(line)
	record := make(map[[2]byte]int, n)
	var validOne, validTwo bool

	for currentIdx := 0; currentIdx < n-2; currentIdx++ {
		validTwo = validTwo || (line[currentIdx] == line[currentIdx+2])

		currentKey := [2]byte{line[currentIdx], line[currentIdx+1]}
		pastIdx, found := record[currentKey]
		if !found {
			record[currentKey] = currentIdx
		} else {
			validOne = validOne || (pastIdx != currentIdx-1)
		}

		if validOne && validTwo {
			break
		}
	}

	if validTwo && !validOne { // Edge case
		key := [2]byte{line[n-2], line[n-1]}
		idx, found := record[key]
		validOne = found && idx != n-3
	}

	return validOne && validTwo
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
	c, cAlt := make(chan bool, n), make(chan bool, n)
	wg.Add(n)
	for scanner.Scan() {
		go func(line string) {
			defer wg.Done()
			c <- isNice(line)
			cAlt <- isNiceAlt(line)
		}(scanner.Text())
	}
	wg.Wait()
	close(c)
	close(cAlt)

	var result int
	for isNice := range c {
		if isNice {
			result++
		}
	}
	println(result)

	result = 0
	for isNice := range cAlt {
		if isNice {
			result++
		}
	}
	println(result)
}
