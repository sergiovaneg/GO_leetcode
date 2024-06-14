package main

import (
	"bufio"
	"log"
	"os"
)

func countMemory(line string) int {
	// Remove outer quotes
	line = line[1 : len(line)-1]
	count := 0

	for idx := 0; idx < len(line); {
		count++
		if line[idx] == '\\' {
			if line[idx+1] == 'x' {
				idx += 4
			} else {
				idx += 2
			}
		} else {
			idx++
		}
	}

	return count
}

func countEncoded(line string) int {
	count := 2

	for _, c := range line {
		if c == '\\' || c == '"' {
			count += 2
		} else {
			count++
		}
	}

	return count
}

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var resultA, resultB int
	for scanner.Scan() {
		line := scanner.Text()
		resultA += len(line) - countMemory(line)
		resultB += countEncoded(line) - len(line)
	}

	println(resultA)
	println(resultB)
}
