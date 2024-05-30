package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"github.com/sergiovaneg/GoStudy/utils"
)

const warmUpIters = 1000

type Node struct {
	name    string
	class   string
	memory  []bool
	inputs  []*Node
	outputs []*Node
}

type System []*Node

type Pulse struct {
	src    *Node
	target *Node
	value  bool
}

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

func (system *System) registerNode(name string, class byte) *Node {
	node := &Node{
		name:   name,
		inputs: make([]*Node, 0),
	}
	node.setClass(class)
	*system = append(*system, node)

	return node
}

func (system *System) getNode(name string, class byte) *Node {
	idx := slices.IndexFunc(*system, func(x *Node) bool {
		return x.name == name
	})

	if idx == -1 {
		return system.registerNode(name, class)
	} else {
		node := (*system)[idx]

		if class != 0 {
			node.setClass(class)
		}

		return node
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

func (system System) resetMemory() {
	for _, node := range system {
		node.initMemory()
	}
}

func (src *Node) emit(value bool, targets []*Node) []Pulse {
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
		return node.emit(pulse.value, node.outputs)
	case "flipflop":
		if pulse.value {
			return nil
		}
		node.memory[0] = !node.memory[0]
		return node.emit(node.memory[0], node.outputs)
	case "conjunction":
		idx := slices.Index(node.inputs, pulse.src)
		node.memory[idx] = pulse.value

		for _, mem := range node.memory {
			if !mem {
				return node.emit(true, node.outputs)
			}
		}

		return node.emit(false, node.outputs)
	default:
		return nil
	}
}

/* Too many states for memoization to be feasible
func (system System) encodeState() string {
	state := ""

	for _, node := range system {
		if node.class == "flipflop" {
			if node.memory[0] {
				state += "1"
			} else {
				state += "0"
			}
		}
	}

	return state
}
*/

func (system System) broadcast() ([2]uint, bool) {
	count := [2]uint{1, 0} // A low pulse is always sent at the start
	bcastIdx := slices.IndexFunc(
		system,
		func(x *Node) bool {
			return x.name == "broadcaster"
		})
	if bcastIdx == -1 {
		return count, false
	}

	var pulse Pulse
	var low2rx uint
	queue := system[bcastIdx].processPulse(pulse)
	for len(queue) > 0 {
		pulse, queue = queue[0], queue[1:]
		if pulse.value {
			count[1]++
		} else {
			count[0]++
		}

		if pulse.target.name == "rx" && !pulse.value {
			low2rx++
		}

		newPulses := pulse.target.processPulse(pulse)
		if newPulses != nil {
			queue = append(queue, newPulses...)
		}
	}

	return count, low2rx > 0
}

func (system System) warmUp() (uint, int) {
	result := [2]uint{0, 0}
	minPresses := -1

	for i := 0; i < warmUpIters && minPresses == -1; i++ {
		count, achieved := system.broadcast()

		if achieved {
			minPresses = i + 1
		}

		result[0] += count[0]
		result[1] += count[1]
	}

	return result[0] * result[1], minPresses
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
	system := make(System, 0, n)
	var src *Node

	for scanner.Scan() {
		srcDst := strings.Split(scanner.Text(), " -> ")

		if srcDst[0] != "broadcaster" {
			name, class := srcDst[0][1:], srcDst[0][0]
			src = system.getNode(name, class)
		} else {
			src = system.registerNode("broadcaster", 'b')
		}

		outputs := strings.Split(srcDst[1], ", ")
		src.outputs = make([]*Node, len(outputs))
		for idx, name := range outputs {
			out := system.getNode(name, 0)
			src.outputs[idx] = out
			out.inputs = append(out.inputs, src)
		}
	}

	system.resetMemory()
	fmt.Println(system.warmUp())
}