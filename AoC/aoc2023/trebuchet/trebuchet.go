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

func isDigit(r byte) bool {
	return r >= '0' && r <= '9'
}

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

func getValue(byteArray []byte) int {
	re_fwd := regexp.MustCompile(
		"(?:f(?:ive|our)|s(?:even|ix)|t(?:hree|wo)|(?:ni|o)ne|eight)")
	re_bwd := regexp.MustCompile(
		"(?:(?:evi|ruo)f|(?:neve|xi)s|(?:eerh|ow)t|en(?:in|o)|thgie)")

	var d0, d1 int
	var match []byte

	idx := slices.IndexFunc(byteArray, isDigit)
	if idx != -1 {
		match = re_fwd.Find(byteArray[:idx])
		if match != nil {
			d0 = matchNumber(string(match))
		} else {
			d0, _ = strconv.Atoi(string(byteArray[idx]))
		}

		slices.Reverse(byteArray)
		idx = slices.IndexFunc(byteArray, isDigit)

		match = re_bwd.Find(byteArray[:idx])
		if match != nil {
			aux := make([]byte, len(match))
			copy(aux, match)
			slices.Reverse(aux)
			d1 = matchNumber(string(aux))
		} else {
			d1, _ = strconv.Atoi(string(byteArray[idx]))
		}
	} else {
		d0 = matchNumber(string(re_fwd.Find(byteArray)))
		slices.Reverse(byteArray)

		match := re_bwd.Find(byteArray)
		aux := make([]byte, len(match))
		copy(aux, match)
		slices.Reverse(aux)
		d1 = matchNumber(string(aux))
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
			c <- getValue([]byte(line))
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
