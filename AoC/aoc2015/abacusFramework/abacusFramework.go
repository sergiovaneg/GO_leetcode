package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
)

func cleanRed(line string) string {
	re := regexp.MustCompile("(\".\":\"red\")")
	arr := []byte(line)
	for {
		matchIdx := re.FindIndex(arr)
		if matchIdx == nil {
			break
		}

		i, j, level := matchIdx[0], matchIdx[1], 0

		for level < 1 {
			i--
			if arr[i] == '{' {
				level++
			} else if arr[i] == '}' {
				level--
			}
		}

		for level > 0 {
			if arr[j] == '{' {
				level++
			} else if arr[j] == '}' {
				level--
			}
			j++
		}

		arr = slices.Delete(arr, i, j+1)
	}

	return string(arr)
}

func sumAll(line string) int {
	result := 0

	for _, match := range regexp.MustCompile(
		"(-*[0-9]+)").FindAllString(line, -1) {
		num, _ := strconv.Atoi(match)
		result += num
	}

	return result
}

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()

	println(sumAll(line))
	cleanLine := cleanRed(line)
	println(sumAll(cleanLine))
}
