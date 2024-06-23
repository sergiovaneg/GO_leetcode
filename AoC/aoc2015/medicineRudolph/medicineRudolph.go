package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type Plant map[string][]string

func (p *Plant) addTransformation(line string) {
	molecules := strings.Split(line, " => ")
	if _, exists := (*p)[molecules[0]]; exists {
		(*p)[molecules[0]] = append((*p)[molecules[0]], molecules[1])
	} else {
		(*p)[molecules[0]] = []string{molecules[1]}
	}
}

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	p := make(Plant)
	for scanner.Scan() {
		if scanner.Text() == "\n" {
			break
		}
		p.addTransformation(scanner.Text())
	}
}
