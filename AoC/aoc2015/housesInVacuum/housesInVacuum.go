package main

import (
	"bufio"
	"log"
	"os"
)

func move(x [2]int, dir byte) [2]int {
	switch dir {
	case '^':
		return [2]int{x[0] - 1, x[1]}
	case 'v':
		return [2]int{x[0] + 1, x[1]}
	case '<':
		return [2]int{x[0], x[1] - 1}
	case '>':
		return [2]int{x[0], x[1] + 1}
	default:
		return x
	}
}

func processDirections(line string, nSanta int) int {
	x0 := [2]int{0, 0}
	record := map[[2]int]bool{
		x0: true,
	}

	x := make([][2]int, nSanta)
	for idx := range x {
		x[idx] = x0
	}

	l := len(line)
	for instIdx := 0; instIdx < l; {
		for santaIdx := range x {
			x[santaIdx] = move(x[santaIdx], line[instIdx])
			record[x[santaIdx]] = true
			instIdx++
		}
	}

	return len(record)
}

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	scanner.Scan()
	println(processDirections(scanner.Text(), 1))
	println(processDirections(scanner.Text(), 2))
}
