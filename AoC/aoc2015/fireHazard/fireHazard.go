package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const size = 1000

type Instruction struct {
	code       int // -1: off, 0: toggle, 1: on
	start, end [2]int
}

type LightArray [][]int

func parseInstruction(line string) Instruction {
	var code int
	if strings.HasPrefix(line, "turn off") {
		code = -1
	} else if strings.HasPrefix(line, "turn on") {
		code = 1
	} // Default is 0, so no 'else' statement

	nums := make([]int, 4)
	for idx, match := range regexp.MustCompile("([0-9]+)").FindAllString(line, 4) {
		nums[idx], _ = strconv.Atoi(match)
	}
	start := [2]int{nums[0], nums[1]}
	end := [2]int{nums[2], nums[3]}

	return Instruction{
		code:  code,
		start: start,
		end:   end,
	}
}

func (lightArr LightArray) execute(inst Instruction) {
	var f func(*int)
	switch inst.code {
	case -1:
		f = func(b *int) { *b = max(0, *b-1) }
	case 0:
		f = func(b *int) { *b += 2 }
	case 1:
		f = func(b *int) { *b += 1 }
	}

	for i := inst.start[0]; i <= inst.end[0]; i++ {
		for j := inst.start[1]; j <= inst.end[1]; j++ {
			f(&lightArr[i][j])
		}
	}
}

func (lightArr LightArray) countOn() (count uint) {
	for i := range lightArr {
		for j := range lightArr[i] {
			count += uint(lightArr[i][j])
		}
	}

	return
}

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	/* n, err := utils.LineCounter(file)
	if err != nil {
		log.Fatal(err)
	} */

	lightArr := make(LightArray, size)
	for i := range lightArr {
		lightArr[i] = make([]int, size)
	}

	for scanner.Scan() {
		lightArr.execute(parseInstruction(scanner.Text()))
	}
	println(lightArr.countOn())
}
