package main

import (
	"bufio"
	"log"
	"math"
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

func countAtoms(molecule string) int {
	return len(regexp.MustCompile("([A-Z][a-z]*)").FindAllString(molecule, -1))
}

func pruneRecord(queue []string, record map[string]bool) map[string]bool {
	thr := math.MaxInt
	for _, molecule := range queue {
		if l := countAtoms(molecule); l < thr {
			thr = l
		}
	}

	for molecule := range record {
		if countAtoms(molecule) < thr {
			delete(record, molecule)
		}
	}

	return record
}

func (p Plant) getMinSteps(target string) int {
	moleculeLengthLimit := countAtoms(target)

	var mu sync.RWMutex
	var wg sync.WaitGroup
	record := map[string]bool{
		"e": true,
	}
	moleculeQueue := []string{"e"}

	minSteps := new(int)

	for steps := 0; *minSteps == 0; steps++ {
		// Remove unreachable from record
		record = pruneRecord(moleculeQueue, record)

		wg.Add(len(moleculeQueue))
		nextQueue := new([]string)

		for _, initialMolecule := range moleculeQueue {
			go func(initialMolecule string) {
				defer wg.Done()
				for next := range p.mapSingleStep(initialMolecule) {
					mu.RLock()
					ok := record[next]
					mu.RUnlock()
					if ok {
						continue
					}

					mu.Lock()
					record[next] = true
					if countAtoms(next) <= moleculeLengthLimit {
						*nextQueue = append(*nextQueue, next)
					}
					if next == target {
						*minSteps = steps + 1
					}
					mu.Unlock()
				}
			}(initialMolecule)
		}
		wg.Wait()
		moleculeQueue = *nextQueue
	}

	return *minSteps
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
