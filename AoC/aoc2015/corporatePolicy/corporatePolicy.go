package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
)

func incrementPassword(current string) string {
	next := []byte(current)
	for idx := len(current) - 1; idx >= 0; idx-- {
		if next[idx] == 'z' {
			next[idx] = 'a'
		} else {
			next[idx]++
			break
		}
	}
	return string(next)
}

func isValidPassword(current string) bool {
	if regexp.MustCompile("[iol]").MatchString(current) {
		return false
	}

	l := len(current)
	couples := make(map[string]bool)
	for idx := 0; idx < l-1; idx++ {
		if current[idx] == current[idx+1] {
			couples[current[idx:idx+2]] = true
			idx++
		}
	}
	if len(couples) < 2 {
		return false
	}

	for idx := 0; idx < l-2; idx++ {
		if current[idx]+1 == current[idx+1] && current[idx]+2 == current[idx+2] {
			return true
		}
	}

	return false
}

func getNextPassword(current string) string {
	current = incrementPassword(current)

	for !isValidPassword(current) {
		current = incrementPassword(current)
	}

	return current
}

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	pass := getNextPassword(scanner.Text())
	println(pass)
	pass = getNextPassword(pass)
	println(pass)
}
