package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func getValue(line string) int {
	matches := regexp.MustCompile("([0-9]+)").FindAllString(line, -1)

	result, err := strconv.Atoi(strings.Join(matches, ""))
	if err != nil {
		return 0
	}

	return result
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func sqrtInt(x int) int { // Undershoot root
	y, yk := -1, 1
	for absInt(y-yk) > 1 {
		y = yk
		yk = (y + x/y) >> 1
	}

	if yk*yk > x {
		return yk - 1
	}
	return yk
}

func getVictories(duration, record int) int {
	inRoot := duration*duration - record<<2
	if inRoot < 0 {
		return 0
	}
	root := (duration - sqrtInt(inRoot)) >> 1

	// fDuration, fInRoot := float64(duration), float64(inRoot)
	// fRoot := (fDuration - math.Sqrt(fInRoot)) / 2.

	return duration - root<<1 + 1
}

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	time := getValue(scanner.Text())
	scanner.Scan()
	distance := getValue(scanner.Text())

	fmt.Println(getVictories(time, distance))
}
