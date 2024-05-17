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

// const limit_red, limit_green, limit_blue = 12, 13, 14

func extractNumber(str string) int {
	re := regexp.MustCompile("([0-9]+)")
	res, err := strconv.Atoi(re.FindString(str))
	if err != nil {
		return 0
	}
	return res
}

/*
func validateDraw(draw string) bool {
	for _, colorCnt := range strings.Split(draw, ",") {
		cnt := extractNumber(colorCnt)
		if strings.Contains(colorCnt, "red") {
			if cnt > limit_red {
				return false
			}
		} else if strings.Contains(colorCnt, "green") {
			if cnt > limit_green {
				return false
			}
		} else {
			if cnt > limit_blue {
				return false
			}
		}
	}
	return true
}
*/

func processDraw(draw string) [3]int {
	res := [3]int{0, 0, 0}

	for _, colorCnt := range strings.Split(draw, ",") {
		cnt := extractNumber(colorCnt)
		if strings.Contains(colorCnt, "red") {
			res[0] = cnt
		} else if strings.Contains(colorCnt, "green") {
			res[1] = cnt
		} else {
			res[2] = cnt
		}
	}

	return res
}

func processGame(game string) int {
	var minCnt [3]int
	for _, draw := range strings.Split(strings.Split(game, ":")[1], ";") {
		candidates := processDraw(draw)

		for idx, candidate := range candidates {
			if candidate > minCnt[idx] {
				minCnt[idx] = candidate
			}
		}
	}

	return minCnt[0] * minCnt[1] * minCnt[2]
}

func main() {
	file, err := os.Open("./cube_conundrum_input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var res int

	var wg sync.WaitGroup
	n, _ := utils.LineCounter(file)
	file.Seek(0, 0)
	c := make(chan int, n+1)

	for scanner.Scan() {
		wg.Add(1)
		go func(line string) {
			defer wg.Done()
			c <- processGame(line)
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
