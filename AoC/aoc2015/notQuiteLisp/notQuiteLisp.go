package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()

	nUp := strings.Count(line, "(")
	nDown := len(line) - nUp

	println(nUp - nDown)

	pos := 0
	for idx, char := range line {
		if char == '(' {
			pos++
		} else {
			pos--
		}

		if pos == -1 {
			println(idx + 1)
			break
		}
	}
}
