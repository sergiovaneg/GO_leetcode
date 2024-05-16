package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
)

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func matchNumber(num string) int {
	switch num {
	case "one":
		return 1
	case "two":
		return 2
	case "three":
		return 3
	case "four":
		return 4
	case "five":
		return 5
	case "six":
		return 6
	case "seven":
		return 7
	case "eight":
		return 8
	case "nine":
		return 9
	default:
		return 0
	}
}

func getValue(line string) int {
	re := regexp.MustCompile(
		"(?:f(?:ive|our)|s(?:even|ix)|t(?:hree|wo)|(?:ni|o)ne|eight)")

	idx0, idx1 := -1, -1
	var num0, num1 string

	runeArray := []rune(line)
	idx0 = slices.IndexFunc(runeArray, isDigit)
	if idx0 != -1 {
		slices.Reverse(runeArray)
		idx1 = len(line) - slices.IndexFunc(runeArray, isDigit) - 1

		num0 = re.FindString(line[:idx0])
		if matches := re.FindAllString(line[idx1+1:], -1); matches != nil {
			num1 = matches[len(matches)-1]
		}
	} else {
		if matches := re.FindAllString(line, -1); matches != nil {
			num0 = matches[0]
			num1 = matches[len(matches)-1]
		}
	}

	var d0, d1 int
	if num0 != "" {
		d0 = matchNumber(num0)
	} else {
		d0, _ = strconv.Atoi(string(line[idx0]))
	}
	if num1 != "" {
		d1 = matchNumber(num1)
	} else {
		d1, _ = strconv.Atoi(string(line[idx1]))
	}

	return d0*10 + d1
}

func main() {
	file, err := os.Open("./trebuchet_input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var res int

	for scanner.Scan() {
		val := getValue(scanner.Text())
		fmt.Println(val)
		res += val
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
}
