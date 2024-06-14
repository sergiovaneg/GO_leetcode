package main

import (
	"bufio"
	"log"
	"os"
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
}
