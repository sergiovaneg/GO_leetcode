package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
)

type Plant map[string][]string

func (p *Plant) addTransformation(line string) {
	molecules := strings.Split(line, " => ")
	(*p)[molecules[0]] = append((*p)[molecules[0]], molecules[1])
}

func (p Plant) mapSingleStep(initialMolecule string) map[string]bool {
	record := make(map[string]bool, len(initialMolecule))
	for key, candidates := range p {
		for _, loc := range regexp.MustCompile(
			"("+key+")").FindAllStringIndex(initialMolecule, -1) {
			for _, rep := range candidates {
				newMolecule := initialMolecule[:loc[0]] + rep + initialMolecule[loc[1]:]
				record[newMolecule] = true
			}
		}
	}
	return record
}

func (p Plant) getMinSteps(target string) int {

	stepsMap := map[string]int{
		"e": 0,
	}
	for steps := 0; stepsMap[target] == 0; steps++ {
		for initialMolecule, recordedSteps := range stepsMap {
			if recordedSteps < steps {
				continue
			}
			nextMolecules := p.mapSingleStep(initialMolecule)
			for next := range nextMolecules {
				if _, ok := stepsMap[next]; !ok {
					stepsMap[next] = steps + 1
				}
			}
		}
	}

	return stepsMap[target]
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
		if scanner.Text() == "" {
			break
		}
		p.addTransformation(scanner.Text())
	}

	scanner.Scan()
	println(len(p.mapSingleStep(scanner.Text())))
	println(p.getMinSteps(scanner.Text()))
}
