package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
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
	var mu sync.RWMutex
	var wg sync.WaitGroup
	stepsMap := map[string]int{
		"e": 0,
	}
	moleculeQueue := []string{"e"}
	for steps := 0; stepsMap[target] == 0; steps++ {
		wg.Add(len(moleculeQueue))
		nextQueue := new([]string)
		for _, initialMolecule := range moleculeQueue {
			go func(initialMolecule string) {
				defer wg.Done()
				for next := range p.mapSingleStep(initialMolecule) {
					mu.RLock()
					_, ok := stepsMap[next]
					mu.RUnlock()
					if ok {
						continue
					}
					mu.Lock()
					stepsMap[next] = steps + 1
					*nextQueue = append(*nextQueue, next)
					mu.Unlock()
				}
			}(initialMolecule)
		}
		wg.Wait()
		moleculeQueue = *nextQueue
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
