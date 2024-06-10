package main

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"regexp"
	"slices"
	"sync"

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

func (g Graph) copy() Graph {
	gCopy := make(Graph, len(g))

	var wg sync.WaitGroup
	wg.Add(len(g))
	for idx, node := range g {
		go func(idx int, node *Node) {
			defer wg.Done()
			gCopy[idx] = &Node{
				name: node.name,
				adj:  make([]*Node, 0, len(node.adj)),
			}
		}(idx, node)
	}
	wg.Wait()

	wg.Add(len(g))
	for srcIdx, src := range g {
		go func(srcIdx int, src *Node) {
			defer wg.Done()
			for _, dst := range src.adj {
				dstIdx := slices.IndexFunc(gCopy, func(x *Node) bool {
					return x.name == dst.name
				})
				gCopy[srcIdx].adj = append(
					gCopy[srcIdx].adj,
					gCopy[dstIdx])
			}
		}(srcIdx, src)
	}
	wg.Wait()

	return gCopy
}

func (g *Graph) contract(idx int) {
	src := (*g)[idx]                        // src == u
	dst := src.adj[rand.Intn(len(src.adj))] // dst == v

	srcdst := &Node{ // srcdst == uv
		name: src.name + dst.name,
		adj:  make([]*Node, 0, len(src.adj)+len(dst.adj)-2),
	}

	// {w,u} -> {w,uv}
	for _, aux := range src.adj {
		if aux == dst {
			continue
		}

		srcdst.adj = append(srcdst.adj, aux)
		replaceIdx := slices.Index(aux.adj, src)
		aux.adj[replaceIdx] = srcdst
	}

	// {w,v} -> {w,uv}
	for _, aux := range dst.adj {
		if aux == src {
			continue
		}

		srcdst.adj = append(srcdst.adj, aux)
		replaceIdx := slices.Index(aux.adj, dst)
		aux.adj[replaceIdx] = srcdst
	}

	// Clean-up
	(*g)[idx] = srcdst                              // u -> uv
	*g = slices.DeleteFunc(*g, func(x *Node) bool { // v -> nil
		return x == dst
	})
	src = nil
	dst = nil
}

func (g Graph) split() (g1, g2 Graph) {
	var gAux Graph
	for {
		gAux = g.copy()
		for len(gAux) > 2 {
			gAux.contract(rand.Intn(len(gAux)))
		}

		if len(gAux[0].adj) == nRemove {
			break
		}
	}

	names1 := regexp.MustCompile("([a-z]{3})").FindAllString(gAux[0].name, -1)
	g1, g2 = make(Graph, 0, len(names1)), make(Graph, 0, len(g)-len(names1))
	for _, node := range g {
		if slices.Contains(names1, node.name) {
			g1 = append(g1, node)
		} else {
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
		var srcGrp, dstGrp int
		if slices.Contains(g1, src) {
			srcGrp = 1
			file.WriteString("\t" + src.name + " [color=green, shape=circle]\n")
		} else if slices.Contains(g2, src) {
			srcGrp = 2
			file.WriteString("\t" + src.name + " [color=red, shape=circle]\n")
		} else {
			srcGrp = 0
			file.WriteString("\t" + src.name + " [color=blue, shape=circle]\n")
		}
		for _, dst := range src.adj {
			if slices.Contains(g1, dst) {
				dstGrp = 1
			} else if slices.Contains(g2, src) {
				dstGrp = 2
			} else {
				dstGrp = 0
			}

			file.WriteString("\t" + src.name + "--" + dst.name)
			if srcGrp != dstGrp {
				file.WriteString(" [style=\"dotted\"]")
			}
			file.WriteString("\n")
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
