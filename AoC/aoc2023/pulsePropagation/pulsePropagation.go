package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Node struct {
	class   string
	memory  []bool
	inputs  []*Node
	outputs []*Node
}

type System map[string]*Node

func (node *Node) updateClass(class byte) {
	switch class {
	case '%':
		node.class = "flipflop"
	case '&':
		node.class = "conjunction"
	}
}

func (system *System) registerNode(name string, class byte) *Node {
	node := &Node{
		memory:  make([]bool, 0),
		inputs:  make([]*Node, 0),
		outputs: make([]*Node, 0),
	}
	node.updateClass(class)

	(*system)[name] = node

	return node
}

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	system := make(System)
	broadcaster := Node{
		class: "broadcaster",
	}
	var src *Node

	for scanner.Scan() {
		srcDst := strings.Split(scanner.Text(), " -> ")

		if srcDst[0] != "broadcaster" {
			name, class := srcDst[0][1:], srcDst[0][0]
			src = system[name]
			if src == nil {
				src = system.registerNode(name, class)
			} else {
				src.updateClass(class)
			}
		} else {
			src = &broadcaster
		}

		outputs := strings.Split(srcDst[1], ", ")
		src.outputs = make([]*Node, len(outputs))
		for idx, name := range outputs {
			out := system[name]
			if out == nil {
				out = system.registerNode(name, 0)
			}

			src.outputs[idx] = out
			out.inputs = append(out.inputs, src)
		}
	}

	fmt.Println("ok")
}
