package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/sergiovaneg/GoStudy/utils"
)

/*
func isSymbol(r byte) bool {
	return r > '9' || (r < '0' && r != '.')
}

func processBlock(block [3]string) int {
	re, result := regexp.MustCompile("([0-9]+)"), 0

	for _, matchIdx := range re.FindAllStringIndex(block[1], -1) {
		var valid bool
		atStart, atEnd := matchIdx[0] == 0, matchIdx[1] == len(block[1])

		if !atStart && isSymbol(block[1][matchIdx[0]-1]) {
			valid = true
		} else if !atEnd && isSymbol(block[1][matchIdx[1]]) {
			valid = true
		} else {
			var bArr0, bArr2 []byte

			if atStart && atEnd {
				bArr0 = []byte(block[0][matchIdx[0]:matchIdx[1]])
				bArr2 = []byte(block[2][matchIdx[0]:matchIdx[1]])
			} else if atStart {
				bArr0 = []byte(block[0][matchIdx[0] : matchIdx[1]+1])
				bArr2 = []byte(block[2][matchIdx[0] : matchIdx[1]+1])
			} else if atEnd {
				bArr0 = []byte(block[0][matchIdx[0]-1 : matchIdx[1]])
				bArr2 = []byte(block[2][matchIdx[0]-1 : matchIdx[1]])
			} else {
				bArr0 = []byte(block[0][matchIdx[0]-1 : matchIdx[1]+1])
				bArr2 = []byte(block[2][matchIdx[0]-1 : matchIdx[1]+1])
			}

			valid = slices.IndexFunc(
				bArr0, isSymbol) != -1 || slices.IndexFunc(
				bArr2, isSymbol) != -1
		}

		if valid {
			val, err := strconv.Atoi(block[1][matchIdx[0]:matchIdx[1]])
			if err == nil {
				result += val
			}
		}
	}

	return result
}
*/

func isAdjacent(gearIdx int, matchIdx []int) bool {
	return (gearIdx+1 >= matchIdx[0]) && (gearIdx <= matchIdx[1])
}

func processBlock(block [3]string) int {
	re, result := regexp.MustCompile("([0-9]+)"), 0
	block_matches := [3][][]int{
		re.FindAllStringIndex(block[0], -1),
		re.FindAllStringIndex(block[1], -1),
		re.FindAllStringIndex(block[2], -1),
	}

	for gearIdx, char := range block[1] {
		if char != '*' {
			continue
		}

		buff := make([]int, 0, 4)
		for blockIdx, matches := range block_matches {
			for _, matchIdx := range matches {
				if isAdjacent(gearIdx, matchIdx) {
					val, err := strconv.Atoi(block[blockIdx][matchIdx[0]:matchIdx[1]])
					if err == nil {
						buff = append(buff, val)
					}
				}
			}
			if len(buff) > 2 {
				break
			}
		}
		if len(buff) == 2 {
			result += buff[0] * buff[1]
		}
	}

	return result
}

func initBlock(line string) [3]string {
	aux := strings.Repeat(".", len(line))
	return [3]string{aux, aux, line}
}

func main() {
	file, err := os.Open("./gear_ratios_input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var res uint

	var wg sync.WaitGroup
	n, _ := utils.LineCounter(file)
	c := make(chan int, n+1)

	scanner.Scan()
	line := scanner.Text()
	block := initBlock(line)

	wg.Add(n)
	f := func(block [3]string) {
		defer wg.Done()
		c <- processBlock(block)
	}

	for scanner.Scan() {
		line = scanner.Text()
		block[0], block[1], block[2] = block[1], block[2], line

		go f(block)
	}

	block[0], block[1], block[2] = block[1], block[2], strings.Repeat(
		".", len(block[2]))
	go f(block)

	wg.Wait()
	close(c)

	for val := range c {
		res += uint(val)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
}
