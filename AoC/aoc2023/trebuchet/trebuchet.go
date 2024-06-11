package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"sync"

	"github.com/sergiovaneg/GoStudy/utils"
)

func matchNumber(num string) int {
	switch num {
	case "one":
		return 1
	case "two":
		return 2
	case "three":
		return 3
	case "four":
		return 4
	case "five":
		return 5
	case "six":
		return 6
	case "seven":
		return 7
	case "eight":
		return 8
	case "nine":
		return 9
	default:
		return 0
	}
}

func getValue(line string) int {
	var d0, d1 int

	dLoc := regexp.MustCompile("([0-9]{1})").FindAllStringIndex(line, -1)
	dN := len(dLoc)

	litRe := regexp.MustCompile(
		"(?:f(?:ive|our)|s(?:even|ix)|t(?:hree|wo)|(?:ni|o)ne|eight)")
	litLoc, j, litN := make([][]int, 0, 2), 0, 0
	for {
		if loc := litRe.FindStringIndex(line[j:]); loc == nil {
			break
		} else {
			j, loc[0], loc[1] = j+loc[0]+1, loc[0]+j, loc[1]+j
			if litN < 2 {
				litLoc = append(litLoc, loc)
				litN++
			} else {
				litLoc[1] = loc
			}
		}
	}

	if litN == 0 {
		d0, _ = strconv.Atoi(line[dLoc[0][0]:dLoc[0][1]])
		d1, _ = strconv.Atoi(line[dLoc[dN-1][0]:dLoc[dN-1][1]])
	} else {
		if dLoc[0][0] < litLoc[0][0] {
			d0, _ = strconv.Atoi(line[dLoc[0][0]:dLoc[0][1]])
		} else {
			d0 = matchNumber(line[litLoc[0][0]:litLoc[0][1]])
		}

		if dLoc[dN-1][0] > litLoc[litN-1][0] {
			d1, _ = strconv.Atoi(line[dLoc[dN-1][0]:dLoc[dN-1][1]])
		} else {
			d1 = matchNumber(line[litLoc[litN-1][0]:litLoc[litN-1][1]])
		}
	}
	return d0*10 + d1
}

func main() {
	file, err := os.Open("./trebuchet_input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var res int

	var wg sync.WaitGroup
	n, _ := utils.LineCounter(file)
	c := make(chan int, n+1)

	for scanner.Scan() {
		wg.Add(1)
		go func(line string) {
			defer wg.Done()
			c <- getValue(line)
		}(scanner.Text())
	}
	wg.Wait()
	close(c)

	for val := range c {
		res += val
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
}
