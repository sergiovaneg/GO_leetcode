package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
)

func getMaterialDims(line string) (int, int) {
	dims := make([]int, 3)

	for idx, lit := range regexp.MustCompile("([0-9]+)").FindAllString(line, 3) {
		dims[idx], _ = strconv.Atoi(lit)
	}

	slices.Sort(dims)

	sides := [3]int{dims[0] * dims[1], dims[1] * dims[2], dims[0] * dims[2]}
	paper := (sides[0]+sides[1]+sides[2])<<1 + min(sides[0], sides[1], sides[2])
	ribbon := (dims[0]+dims[1])<<1 + sides[0]*dims[2]

	return paper, ribbon
}

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	var paper, ribbon int
	for scanner.Scan() {
		newPaper, newRibbon := getMaterialDims(scanner.Text())
		paper += newPaper
		ribbon += newRibbon
	}

	println(paper)
	println(ribbon)
}
