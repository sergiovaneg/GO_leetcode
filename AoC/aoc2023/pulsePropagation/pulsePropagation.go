package main

import (
	"bufio"
	"log"
	"os"
	"slices"
	"strings"
)

type Node struct {
	class   string
	memory  []bool
	inputs  []*Node
	outputs []*Node
}

type Pulse struct {
	src    *Node
	target *Node
	value  bool
}

type System map[string]*Node

func (node *Node) setClass(class byte) {
	switch class {
	case '%':
		node.class = "flipflop"
	case '&':
		node.class = "conjunction"
	case 'b':
		node.class = "broadcaster"
	}
}

func (node *Node) initMemory() {
	switch node.class {
	case "flipflop":
		node.memory = make([]bool, 1)
	case "conjunction":
		node.memory = make([]bool, len(node.inputs))
	}
}

func (system *System) registerNode(name string, class byte) *Node {
	node := &Node{
		inputs:  make([]*Node, 0),
		outputs: make([]*Node, 0),
	}
	node.setClass(class)

	(*system)[name] = node

	return node
}

func emit(value bool, src *Node, targets []*Node) []Pulse {
	pulses := make([]Pulse, len(targets))

	for idx, target := range targets {
		pulses[idx] = Pulse{
			src:    src,
			target: target,
			value:  value,
		}
	}

	return pulses
}

func (node *Node) processPulse(pulse Pulse) []Pulse {
	switch node.class {
	case "broadcaster":
		return emit(pulse.value, node, node.outputs)
	case "flipflop":
		if pulse.value {
			break
		}
		node.memory[0] = !node.memory[0]
		return emit(node.memory[0], node, node.outputs)
	case "conjunction":
		idx := slices.Index(node.inputs, pulse.src)
		node.memory[idx] = pulse.value

		var value bool
		for _, mem := range node.memory {
			if !mem {
				value = true
				break
			}
		}

		return emit(value, node, node.outputs)
	}

	return nil
}

func (system System) broadcast() {
	var pulse Pulse
	queue := system["broadcaster"].processPulse(pulse)
	for len(queue) > 0 {
		pulse, queue = queue[0], queue[1:]
		newPulses := pulse.target.processPulse(pulse)
		if newPulses != nil {
			queue = append(queue, newPulses...)
		}
	}
}

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	system := make(System)
	var src *Node

	for scanner.Scan() {
		srcDst := strings.Split(scanner.Text(), " -> ")

		if srcDst[0] != "broadcaster" {
			name, class := srcDst[0][1:], srcDst[0][0]
			src = system[name]
			if src == nil {
				src = system.registerNode(name, class)
			} else {
				src.setClass(class)
			}
		} else {
			src = system.registerNode("broadcaster", 'b')
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

	for _, node := range system {
		node.initMemory()
	}

	system.broadcast()
}
