package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Workflow struct {
	name         string
	targets      []byte
	gt           []bool
	thresholds   []int
	destinations []string
	defaultDest  string
}

type Part map[byte]int

type WMap map[string]*Workflow

type PartSet map[byte]*[2]int

func parseStep(step string) (target byte, gt bool, thr int, dest string) {
	var cond string
	if split := strings.Split(step, ":"); len(split) == 2 {
		cond, dest = split[0], split[1]
	}

	target, gt = cond[0], cond[1] == '>'
	thr, _ = strconv.Atoi(cond[2:])

	return
}

func parseWorkflow(line string) Workflow {
	loc := strings.IndexRune(line, '{')

	steps := strings.Split(line[loc+1:len(line)-1], ",")
	n := len(steps) - 1
	w := Workflow{
		name:         line[:loc],
		targets:      make([]byte, n),
		gt:           make([]bool, n),
		thresholds:   make([]int, n),
		destinations: make([]string, n),
		defaultDest:  steps[n],
	}

	for idx, step := range steps[:n] {
		w.targets[idx],
			w.gt[idx],
			w.thresholds[idx],
			w.destinations[idx] = parseStep(step)
	}

	return w
}

func parsePart(line string) Part {
	ratings := regexp.MustCompile("([0-9]+)").FindAllString(line, 4)
	p := make(Part, 4)

	p['x'], _ = strconv.Atoi(ratings[0])
	p['m'], _ = strconv.Atoi(ratings[1])
	p['a'], _ = strconv.Atoi(ratings[2])
	p['s'], _ = strconv.Atoi(ratings[3])

	return p
}

func (w Workflow) getNext(part Part) string {
	for idx, dest := range w.destinations {
		var complies bool
		if w.gt[idx] {
			complies = part[w.targets[idx]] > w.thresholds[idx]
		} else {
			complies = part[w.targets[idx]] < w.thresholds[idx]
		}

		if complies {
			return dest
		}
	}

	return w.defaultDest
}

func (wMap WMap) sortPart(part Part) bool {
	next := "in"

	for next != "A" && next != "R" {
		w := wMap[next]
		next = w.getNext(part)
	}

	return next == "A"
}

func (p Part) getTotalRating() int {
	return p['x'] + p['m'] + p['a'] + p['s']
}

func (fr PartSet) count() uint {
	size := uint(1)
	for _, interval := range fr {
		size *= uint(interval[1] - interval[0])
	}
	return size
}

func (wMap WMap) branchFilter(subset PartSet, next string) uint {
	switch next {
	case "A":
		return subset.count()
	case "R":
		return 0
	default:
		return wMap.filterSubset(subset, next)
	}
}

func (wMap WMap) filterSubset(subset PartSet, wName string) uint {
	if subset.count() == 0 { // Early return
		return 0
	}

	w := wMap[wName]
	var accepted uint

	for rIdx, target := range w.targets {
		thr := w.thresholds[rIdx]
		targetSet := subset[target]

		if thr >= targetSet[0] && thr < targetSet[1] {
			branchSubset := make(PartSet, 4)
			for key, interval := range subset {
				if key == target {
					if w.gt[rIdx] {
						branchSubset[key] = &[2]int{thr + 1, interval[1]}
						interval[1] = thr + 1
					} else {
						branchSubset[key] = &[2]int{interval[0], thr}
						interval[0] = thr
					}
				} else {
					branchSubset[key] = &[2]int{interval[0], interval[1]}
				}
			}

			accepted += wMap.branchFilter(branchSubset, w.destinations[rIdx])

			if subset.count() == 0 { // Early return
				return accepted
			}
		} else {
			ll, ul := targetSet[0], targetSet[1]
			complies := (w.gt[rIdx] && ll > thr) || (!w.gt[rIdx] && ul <= thr)
			if complies {
				return accepted + wMap.branchFilter(subset, w.destinations[rIdx])
			}
		}
	}

	return accepted + wMap.branchFilter(subset, w.defaultDest)
}

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	wfMap := make(WMap)
	for scanner.Scan() && scanner.Text() != "" {
		w := parseWorkflow(scanner.Text())
		wfMap[w.name] = &w
	}

	var res int
	for scanner.Scan() {
		p := parsePart(scanner.Text())
		if wfMap.sortPart(p) {
			res += p.getTotalRating()
		}
	}
	fmt.Println(res)

	validRatings := wfMap.filterSubset(
		PartSet{
			'x': {1, 4001},
			'm': {1, 4001},
			'a': {1, 4001},
			's': {1, 4001},
		}, "in")

	fmt.Println(validRatings)
}
