package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/sergiovaneg/GoStudy/utils"
)

const ignoreSlopes = true

type Node struct {
	coordinates [2]int
	neighbours  []*Node
	distances   []uint
}

type Walker struct {
	last, current [2]int
	visited       []*Node
	steps         uint
}

type Graph []*Node

type HeapNode struct {
	node *Node
	dist uint
}

type PriorityQueue []HeapNode

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].dist > pq[j].dist }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }

func (pq *PriorityQueue) Push(x any) { *pq = append(*pq, x.(HeapNode)) }

func (pq *PriorityQueue) Pop() any {
	n := len(*pq)
	x := (*pq)[n-1]
	*pq = (*pq)[:n-1]
	return x
}

func (w Walker) getNext(hikingMap [][]rune) [][2]int {
	var next [][2]int

	if hikingMap[w.current[0]][w.current[1]] == '.' || ignoreSlopes {
		next = make([][2]int, 0, 4)
		for _, offset := range [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} {
			newI, newJ := w.current[0]+offset[0], w.current[1]+offset[1]
			if newI < 0 || newI >= len(hikingMap) {
				continue
			}
			if newJ < 0 || newJ >= len(hikingMap[newI]) {
				continue
			}
			if hikingMap[newI][newJ] != '#' {
				next = append(next, [2]int{newI, newJ})
			}
		}
	} else {
		switch hikingMap[w.current[0]][w.current[1]] {
		case '>':
			next = [][2]int{{
				w.current[0], w.current[1] + 1,
			}}
		case '^':
			next = [][2]int{{
				w.current[0] - 1, w.current[1],
			}}
		case '<':
			next = [][2]int{{
				w.current[0], w.current[1] - 1,
			}}
		case 'v':
			next = [][2]int{{
				w.current[0] + 1, w.current[1],
			}}
		}
	}

	return slices.DeleteFunc(
		next,
		func(x [2]int) bool {
			return x == w.last
		})
}

func (w *Walker) connect(graph Graph) ([]*Node, bool) {
	var node *Node
	last := w.visited[len(w.visited)-1]

	idx := slices.IndexFunc(graph, func(x *Node) bool {
		return x.coordinates == w.current
	})

	if idx == -1 {
		node = &Node{
			coordinates: w.current,
			neighbours:  make([]*Node, 0, 2),
			distances:   make([]uint, 0, 2),
		}
		graph = append(graph, node)
	} else {
		node = graph[idx]
	}

	last.neighbours = append(last.neighbours, node)
	last.distances = append(last.distances, w.steps)

	w.visited = append(w.visited, node)
	w.steps = 0

	return graph, idx == -1
}

func genGraph(hikingMap [][]rune, start, end *Node) Graph {
	graph := Graph{start, end}
	walkerQueue := []Walker{{
		last:    [2]int{-1, -1},
		current: start.coordinates,
		visited: []*Node{start},
	}}

	var w Walker
	for len(walkerQueue) > 0 {
		w, walkerQueue = walkerQueue[0], walkerQueue[1:]
		for {
			if w.current == end.coordinates {
				graph, _ = w.connect(graph)
				break
			}

			next := w.getNext(hikingMap)

			if len(next) == 0 { // End of the road
				break
			} else if len(next) == 1 { // Just advance
				w.last, w.current = w.current, next[0]
				w.steps++
			} else {
				idx := slices.IndexFunc(w.visited, func(x *Node) bool {
					return x.coordinates == w.current
				})
				if idx != -1 {
					break
				}
				var newNode bool
				graph, newNode = w.connect(graph)

				if !newNode {
					break
				}

				for _, toQueue := range next[1:] {
					visitedCopy := make([]*Node, len(w.visited))
					copy(visitedCopy, w.visited)

					walkerQueue = append(walkerQueue, Walker{
						last:    w.current,
						current: toQueue,
						visited: visitedCopy,
						steps:   1,
					})
				}

				w.last, w.current = w.current, next[0]
				w.steps++
			}
		}
	}

	return graph
}

func (graph Graph) simplifyGraph(start, end *Node) Graph {
	relevantMap := make(map[*Node]bool, len(graph))
	relevantMap[start] = true
	relevantMap[end] = true

	for _, node := range graph {
		uniqueMap := make(map[*Node]uint)
		for idx, neigh := range node.neighbours {
			uniqueMap[neigh] = node.distances[idx]
		}

		node.neighbours = make([]*Node, 0, len(uniqueMap))
		node.distances = make([]uint, 0, len(uniqueMap))
		for key, val := range uniqueMap {
			node.neighbours = append(node.neighbours, key)
			node.distances = append(node.distances, val)
		}
	}

	for _, node := range graph {
		for len(node.neighbours) == 1 {
			if node.neighbours[0] == end {
				break
			}

			distances := make([]uint, len(node.neighbours[0].neighbours))
			copy(distances, node.neighbours[0].distances)

			for idx := range distances {
				distances[idx] += node.distances[0]
			}
			node.neighbours = node.neighbours[0].neighbours
			node.distances = distances
		}

		// Verify relevance after simplification
		for _, neighbour := range node.neighbours {
			relevantMap[neighbour] = true
		}
	}

	return slices.DeleteFunc(graph, func(x *Node) bool {
		return !relevantMap[x]
	})
}

func (graph Graph) getLongestPath(start, end *Node) uint {
	pq := &PriorityQueue{HeapNode{
		node: start,
		dist: 0,
	}}
	heap.Init(pq)

	dist := make(map[*Node]uint, len(graph))
	for _, node := range graph {
		dist[node] = 0
	}
	dist[start] = 0

	for len(*pq) > 0 {
		heapNode := heap.Pop(pq).(HeapNode)
		for idx, neighbour := range heapNode.node.neighbours {
			alt := dist[heapNode.node] + heapNode.node.distances[idx]

			if alt > dist[neighbour] {
				dist[neighbour] = alt
				heap.Push(pq, HeapNode{
					node: neighbour,
					dist: alt,
				})
			}
		}
	}

	return uint(dist[end])
}

func (graph Graph) makeUndirected() Graph {
	for _, node := range graph {
		for idx, neigh := range node.neighbours {
			neigh.neighbours = append(neigh.neighbours, node)
			neigh.distances = append(neigh.distances, node.distances[idx])
		}
	}
	return graph
}

func (graph Graph) toDot(start, end *Node, name string) {
	file, err := os.Create("./" + name + ".dot")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if ignoreSlopes {
		file.WriteString("graph{\n")
	} else {
		file.WriteString("digraph {\n")
	}
	for _, src := range graph {
		var srcName string
		switch src {
		case start:
			srcName = "S"
			file.WriteString("\tS [color=green, shape=circle, ordering=in]\n")
		case end:
			srcName = "E"
			file.WriteString("\tE [color=red, shape=circle, ordering=out]\n")
		default:
			srcName = fmt.Sprintf("i%vj%v", src.coordinates[0], src.coordinates[1])
			file.WriteString(
				"\t" + srcName + " [shape=square]\n")
		}

		for neighIdx, dst := range src.neighbours {
			var dstName string
			switch dst {
			case start:
				dstName = "S"
			case end:
				dstName = "E"
			default:
				dstName = fmt.Sprintf("i%vj%v", dst.coordinates[0], dst.coordinates[1])
			}
			properties := fmt.Sprintf(" [label=%v]", src.distances[neighIdx])
			if ignoreSlopes {
				file.WriteString("\t" + srcName + " -- " + dstName + properties + "\n")
			} else {
				file.WriteString("\t" + srcName + " -> " + dstName + properties + "\n")
			}
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

	hikingMap := make([][]rune, n)
	var start, end *Node

	for idx := 0; scanner.Scan(); idx++ {
		line := []rune(scanner.Text())
		if idx == 0 {
			start = &Node{
				coordinates: [2]int{0, slices.Index(line, '.')},
				neighbours:  make([]*Node, 0, 1),
				distances:   make([]uint, 0, 1),
			}
		} else if idx == n-1 {
			end = &Node{
				coordinates: [2]int{n - 1, slices.Index(line, '.')},
			}
		}
		hikingMap[idx] = line
	}

	graph := genGraph(hikingMap, start, end)
	if !ignoreSlopes {
		graph = graph.simplifyGraph(start, end)
		graph.toDot(start, end, "graph")
		println(graph.getLongestPath(start, end))
	} else {
		graph = graph.makeUndirected()
		graph = graph.simplifyGraph(start, end)
		graph.toDot(start, end, "undirected")
	}
}
