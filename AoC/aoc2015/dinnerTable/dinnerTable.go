package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Person struct {
	name      string
	happiness map[string]int
}

type Guests map[string]*Person

func (g *Guests) updateGuests(line string) {
	words := strings.Split(line[:len(line)-1], " ")
	nWords := len(words)

	newHappiness, _ := strconv.Atoi(words[3])
	if words[2] == "lose" {
		newHappiness = -newHappiness
	}

	n1, n2 := words[0], words[nWords-1]
	var p1, p2 *Person

	if p1 = (*g)[n1]; p1 == nil {
		p1 = &Person{
			name:      n1,
			happiness: make(map[string]int),
		}
		(*g)[n1] = p1
	}

	if p2 = (*g)[n2]; p2 == nil {
		p2 = &Person{
			name:      n2,
			happiness: make(map[string]int),
		}
		(*g)[n2] = p2
	}

	if _, ok := p1.happiness[n2]; ok {
		p1.happiness[n2] += newHappiness
		p2.happiness[n1] += newHappiness
	} else {
		p1.happiness[n2] = newHappiness
		p2.happiness[n1] = newHappiness
	}
}

func getHappiness(order []*Person) int {
	var result int
	for idx := 0; idx < len(order)-1; idx++ {
		result += order[idx].happiness[order[idx+1].name]
	}
	return result + order[len(order)-1].happiness[order[0].name]
}

func (g Guests) recursiveAdd(order []*Person) int {
	if len(order) == len(g) {
		return getHappiness(order)
	}

	best := math.MinInt

	for candidate := range g {
		if slices.ContainsFunc(order, func(x *Person) bool {
			return x.name == candidate
		}) {
			continue
		}

		if h := g.recursiveAdd(append(order, g[candidate])); h > best {
			best = h
		}
	}

	return best
}

func (g Guests) getOptimalHappiness() int {
	order := make([]*Person, 1, len(g))
	for person := range g {
		order[0] = g[person]
		return g.recursiveAdd(order)
	}
	return 0
}

func (g *Guests) addMyself() {
	myself := &Person{
		name:      "myself",
		happiness: make(map[string]int),
	}

	for name, guest := range *g {
		myself.happiness[name] = 0
		guest.happiness["myself"] = 0
	}

	(*g)["myself"] = myself
}

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	guests := make(Guests)
	for scanner.Scan() {
		guests.updateGuests(scanner.Text())
	}

	println(guests.getOptimalHappiness())
	guests.addMyself()
	println(guests.getOptimalHappiness())
}
