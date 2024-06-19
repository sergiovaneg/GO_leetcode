package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/sergiovaneg/GoStudy/utils"
)

const DURATION = 2503

type Reindeer struct {
	name              string
	speed             int
	flyTime, restTime int
}

func parseReindeer(line string) Reindeer {
	matches := regexp.MustCompile("([0-9]+)").FindAllString(line, 3)
	speed, _ := strconv.Atoi(matches[0])
	flyTime, _ := strconv.Atoi(matches[1])
	restTime, _ := strconv.Atoi(matches[2])

	return Reindeer{
		name:     strings.SplitN(line, " ", 1)[0],
		speed:    speed,
		flyTime:  flyTime,
		restTime: restTime,
	}
}

func (r Reindeer) getDistanceTraveled(duration int) int {
	qtt := duration / (r.flyTime + r.restTime)
	mod := duration % (r.flyTime + r.restTime)
	distance := r.speed * r.flyTime * qtt
	distance += r.speed * min(r.flyTime, mod)

	return distance
}

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	n, err := utils.LineCounter(file)
	if err != nil {
		log.Fatal(err)
	}

	longestDistance := math.MinInt
	reindeers := make([]Reindeer, 0, n)
	for scanner.Scan() {
		r := parseReindeer(scanner.Text())
		if candidate := r.getDistanceTraveled(DURATION); candidate > longestDistance {
			longestDistance = candidate
		}
		reindeers = append(reindeers, r)
	}

	println(longestDistance)

	score := make(map[string]int, len(reindeers))
	for t := 1; t <= DURATION; t++ {
		distances := make([]int, len(reindeers))

		for idx, r := range reindeers {
			distances[idx] = r.getDistanceTraveled(t)
		}

		thr := slices.Max(distances)
		for idx, r := range reindeers {
			if distances[idx] == thr {
				score[r.name]++
			}
		}
	}

	mostPoints := math.MinInt
	for _, points := range score {
		if points > mostPoints {
			mostPoints = points
		}
	}
	println(mostPoints)
}
