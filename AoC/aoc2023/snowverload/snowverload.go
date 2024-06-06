package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/sergiovaneg/GoStudy/utils"
)

type Node struct {
	name string
	adj  []*Node
}

type Vertex [2]*Node
type Graph []*Node

func (g *Graph) addNode(line string) {
	matches := regexp.MustCompile("([a-z]+)").FindAllString(line, -1)

	var root *Node

	rootIdx := slices.IndexFunc(*g, func(x *Node) bool {
		return x.name == matches[0]
	})
	if rootIdx == -1 {
		root = &Node{
			name: matches[0],
			adj:  make([]*Node, 0, len(matches)-1),
		}
		*g = append(*g, root)
	} else {
		root = (*g)[rootIdx]
	}

	var other *Node
	for _, match := range matches[1:] {
		otherIdx := slices.IndexFunc(*g, func(x *Node) bool {
			return x.name == match
		})

		if otherIdx == -1 {
			other = &Node{
				name: match,
				adj:  []*Node{root},
			}
			*g = append(*g, root)
		} else {
			other = (*g)[otherIdx]
			other.adj = append(other.adj, root)
		}

		root.adj = append(root.adj, other)
	}
}

func (g Graph) getLocalBetweenness(s, t *Node) map[Vertex]float64 {
	optimalPaths := make([][]*Node, 0)
	thr := math.MaxInt
	queue := [][]*Node{{s}}

	var currentPath []*Node
	for len(queue) > 0 {
		currentPath, queue = queue[0], queue[1:]
		for len(currentPath) < thr {
			n := len(currentPath)

			if slices.Contains(currentPath[n-1].adj, t) {
				currentPath = append(currentPath, t)
				if n == thr-1 {
					optimalPaths = append(optimalPaths, currentPath)
				} else {
					thr = n + 1
					optimalPaths = [][]*Node{currentPath}
				}
				break
			}

			nextNodes := make([]*Node, 0, len(currentPath[n-1].adj))
			for _, candidate := range currentPath[n-1].adj {
				if !slices.Contains(currentPath, candidate) {
					nextNodes = append(nextNodes, candidate)
				}
			}

			if len(nextNodes) == 0 {
				break
			}

			for _, next := range nextNodes[1:] {
				path := make([]*Node, n, n+1)
				copy(path, currentPath)
				path = append(path, next)
				queue = append(queue, path)
			}

			currentPath = append(currentPath, nextNodes[0])
		}
	}
	frequency := make(map[Vertex]float64)
	var vertex Vertex
	nPaths := float64(len(optimalPaths))
	for _, path := range optimalPaths {
		for idx, next := range path[1:] {
			if strings.Compare(path[idx].name, next.name) > 0 {
				vertex[0], vertex[1] = path[idx], next
			} else {
				vertex[0], vertex[1] = next, path[idx]
			}
			frequency[vertex] += 1. / nPaths
		}
	}

	return frequency
}

type VertexCentrality struct {
	v Vertex
	c float64
}

func (g Graph) split() (g1, g2 Graph) {
	globalCentrality := make([]VertexCentrality, 0)

	for sIdx, s := range g {
		for _, t := range g[sIdx+1:] {
			for vertex, centrality := range g.getLocalBetweenness(s, t) {
				vIdx := slices.IndexFunc(globalCentrality,
					func(x VertexCentrality) bool {
						return x.v == vertex
					})
				if vIdx == -1 {
					globalCentrality = append(globalCentrality, VertexCentrality{
						v: vertex,
						c: centrality,
					})
				} else {
					globalCentrality[vIdx].c += centrality
				}
			}
		}
	}

	slices.SortFunc(globalCentrality, func(a, b VertexCentrality) int {
		return int(math.Round(b.c - a.c))
	})

	return
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

	g := make(Graph, 0, n)
	for scanner.Scan() {
		g.addNode(scanner.Text())
	}

	g1, g2 := g.split()
	println(len(g1) * len(g2))
}
