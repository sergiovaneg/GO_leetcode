package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/sergiovaneg/GoStudy/utils"
)

type Aunt struct {
	id    int
	stats map[string]int
}

func parseReference(path string) Aunt {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	n, err := utils.LineCounter(file)
	if err != nil {
		log.Fatal(err)
	}

	trueAunt := Aunt{
		id:    0,
		stats: make(map[string]int, n),
	}

	for scanner.Scan() {
		pair := strings.Split(scanner.Text(), ": ")
		count, _ := strconv.Atoi(pair[1])
		trueAunt.stats[pair[0]] = count
	}

	return trueAunt
}

func parseAunt(line string) Aunt {
	idLoc := regexp.MustCompile("([0-9]+)").FindStringIndex(line)
	id, _ := strconv.Atoi(line[idLoc[0]:idLoc[1]])

	line = line[idLoc[1]+2:]
	stats := strings.Split(line, ", ")

	candidate := Aunt{
		id:    id,
		stats: make(map[string]int, len(stats)),
	}

	for _, stat := range stats {
		nameVal := strings.Split(stat, ": ")
		val, _ := strconv.Atoi(nameVal[1])
		candidate.stats[nameVal[0]] = val
	}

	return candidate
}

func (candidate Aunt) isValid(reference Aunt) bool {
	for key, val := range candidate.stats {
		if val != reference.stats[key] {
			return false
		}
	}
	return true
}

func (candidate Aunt) isValidAlt(reference Aunt) bool {
	for key, val := range candidate.stats {
		switch key {
		case "cats":
			fallthrough
		case "trees":
			if !(reference.stats[key] < val) {
				return false
			}
		case "pomeranians":
			fallthrough
		case "goldfish":
			if !(reference.stats[key] > val) {
				return false
			}
		default:
			if val != reference.stats[key] {
				return false
			}
		}
	}
	return true
}

func main() {
	ref := parseReference("./reference.txt")

	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if aunt := parseAunt(scanner.Text()); aunt.isValid(ref) {
			fmt.Printf("Part 1: %v\n", aunt.id)
		} else if aunt.isValidAlt(ref) {
			fmt.Printf("Part 2: %v\n", aunt.id)
		}
	}
}
