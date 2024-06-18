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

type Node struct {
	adj  []*Node
	dist []uint
}

type Graph map[string]*Node

func (g *Graph) addEdge(line string) {
	words := strings.Split(line, " ")
	a, b := words[0], words[2]
	dist, _ := strconv.Atoi(words[4])

	nodeA, okA := (*g)[a]
	if !okA {
		nodeA = &Node{
			adj:  make([]*Node, 0, 1),
			dist: make([]uint, 0, 1),
		}
		(*g)[a] = nodeA
	}

	nodeB, okB := (*g)[b]
	if !okB {
		nodeB = &Node{
			adj:  make([]*Node, 0, 1),
			dist: make([]uint, 0, 1),
		}
		(*g)[b] = nodeB
	}

	nodeA.adj = append(nodeA.adj, nodeB)
	nodeA.dist = append(nodeA.dist, uint(dist))

	nodeB.adj = append(nodeB.adj, nodeA)
	nodeB.dist = append(nodeB.dist, uint(dist))
}

func (g Graph) recursiveMinimize(path []*Node,
	pathLen uint, currentBest *uint) {
	if pathLen >= *currentBest {
		return
	}

	n := len(path)
	if n == len(g) {
		*currentBest = pathLen
		return
	}

	current := path[n-1]
	for idx, next := range current.adj {
		if slices.Contains(path, next) {
			continue
		}

		g.recursiveMinimize(
			append(path, next),
			pathLen+current.dist[idx],
			currentBest)
	}
}

func (g Graph) minimize() uint {
	best := uint(math.MaxUint)
	path := make([]*Node, 1, len(g))
	for _, node := range g {
		path[0] = node
		g.recursiveMinimize(path, 0, &best)
	}

	return best
}

func (g Graph) recursiveMaximize(path []*Node,
	pathLen uint, currentBest *uint) {
	n := len(path)
	if n == len(g) {
		if pathLen > *currentBest {
			*currentBest = pathLen
		}
		return
	}

	current := path[n-1]
	for idx, next := range current.adj {
		if slices.Contains(path, next) {
			continue
		}

		g.recursiveMaximize(
			append(path, next),
			pathLen+current.dist[idx],
			currentBest)
	}
}

func (g Graph) maximize() uint {
	best := uint(0)
	path := make([]*Node, 1, len(g))
	for _, node := range g {
		path[0] = node
		g.recursiveMaximize(path, 0, &best)
	}

	return best
}

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	g := make(Graph)
	for scanner.Scan() {
		g.addEdge(scanner.Text())
	}

	println(g.minimize())
	println(g.maximize())
}
