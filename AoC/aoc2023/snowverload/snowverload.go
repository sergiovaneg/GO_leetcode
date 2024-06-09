package main

import (
	"bufio"
	"errors"
	"log"
	"math"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/sergiovaneg/GoStudy/utils"
)

const nRemove = 3

type Node struct {
	name string
	adj  []*Node
}
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
			*g = append(*g, other)
		} else {
			other = (*g)[otherIdx]
			other.adj = append(other.adj, root)
		}

		root.adj = append(root.adj, other)
	}
}

type Edge [2]*Node

func (g Graph) getLocalBetweenness(s, t *Node) map[Edge]float64 {
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

	frequency := make(map[Edge]float64)
	var vertex Edge
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

func (g Graph) removeEdge(e Edge) error {
	nodeA, nodeB := e[0], e[1]
	idxA, idxB := slices.Index(nodeA.adj, e[1]), slices.Index(nodeB.adj, e[0])
	if idxA == -1 || idxB == -1 {
		return errors.New("edge not in graph")
	}

	nodeA.adj = slices.Delete(nodeA.adj, idxA, idxA+1)
	nodeB.adj = slices.Delete(nodeB.adj, idxB, idxB+1)

	return nil
}

func (g Graph) getFullyConnected(node *Node) Graph {
	gFC := make(Graph, 1, len(g))
	gFC[0] = node

	for idx := 0; idx < len(gFC); idx++ {
		for _, candidate := range gFC[idx].adj {
			if !slices.Contains(gFC, candidate) {
				gFC = append(gFC, candidate)
			}
		}
	}

	return slices.Clip(gFC)
}

type EdgeBetweenness struct {
	edge        Edge
	betweenness float64
}

func (g Graph) getGlobalBetweenness() []EdgeBetweenness {
	globalCentrality := make([]EdgeBetweenness, 0)

	for sIdx, s := range g {
		for _, t := range g[sIdx+1:] {
			for vertex, centrality := range g.getLocalBetweenness(s, t) {
				vIdx := slices.IndexFunc(globalCentrality,
					func(x EdgeBetweenness) bool {
						return x.edge == vertex
					})

				if vIdx == -1 {
					globalCentrality = append(globalCentrality, EdgeBetweenness{
						edge:        vertex,
						betweenness: centrality,
					})
				} else {
					globalCentrality[vIdx].betweenness += centrality
				}
			}
		}
	}

	return globalCentrality
}

func (g Graph) split() (g1, g2 Graph) {

	for n := 0; n < nRemove; n++ {
		globalCentrality := g.getGlobalBetweenness()

		centralEdge := slices.MaxFunc(globalCentrality, func(a, b EdgeBetweenness) int {
			return int(math.Round(a.betweenness - b.betweenness))
		}).edge

		err := g.removeEdge(centralEdge)
		if err != nil {
			log.Fatal(err)
		}
	}

	g1 = g.getFullyConnected(g[0])
	g2 = make(Graph, 0, len(g)-len(g1))
	for _, node := range g {
		if !slices.Contains(g1, node) {
			g2 = append(g2, node)
		}
	}

	return
}

func (g Graph) toDot(g1, g2 Graph) {
	file, err := os.Create("./graph.dot")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file.WriteString("strict graph {\n")
	for _, src := range g {
		if slices.Contains(g1, src) {
			file.WriteString("\t" + src.name + " [color=green, shape=circle]\n")
		} else if slices.Contains(g2, src) {
			file.WriteString("\t" + src.name + " [color=red, shape=circle]\n")
		} else {
			file.WriteString("\t" + src.name + " [color=blue, shape=circle]\n")
		}
		for _, dst := range src.adj {
			file.WriteString("\t" + src.name + "--" + dst.name)
		}
	}
	file.WriteString("}\n")
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

	g.toDot(g1, g2)
}
