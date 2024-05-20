package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
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

func getVictories(duration, record int) int {
	inRoot := duration*duration - 4*record
	if inRoot < 0 {
		return 0
	}
	fDuration, fInRoot := float64(duration), float64(inRoot)
	root := (fDuration - math.Sqrt(fInRoot)) / 2.

	return duration - int(math.Floor(root)+1)<<1 + 1
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
