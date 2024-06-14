package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/sergiovaneg/GoStudy/utils"
)

type Instruction struct {
	operation string
	operands  []string
	dst       string
}

type Circuit map[string]uint16

func parseInstruction(line string) Instruction {
	arrowIdx := strings.Index(line, "->")
	dst := line[arrowIdx+3:]
	operands := strings.Split(line[:arrowIdx-1], " ")

	re := regexp.MustCompile("([A-Z]+)")
	operation := ""
	for idx, x := range operands {
		if re.FindString(x) == x {
			operation = x
			operands = slices.Delete(operands, idx, idx+1)
			break
		}
	}

	return Instruction{
		operation: operation,
		operands:  operands,
		dst:       dst,
	}
}

func install(instSet []Instruction) (circuit Circuit) {
	circuit = make(Circuit, len(instSet))

	queue := make([]Instruction, len(instSet))
	copy(queue, instSet)

	reDigit := regexp.MustCompile("([0-9]+)")
	var current Instruction

	for len(queue) > 0 {
		current, queue = queue[0], queue[1:]

		operands, ready := make([]uint16, len(current.operands)), true
		for idx, tag := range current.operands {
			if reDigit.MatchString(tag) {
				aux, _ := strconv.Atoi(tag)
				operands[idx] = uint16(aux)
			} else if aux, ok := circuit[tag]; ok {
				operands[idx] = aux
			} else {
				ready = false
				break
			}
		}

		if ready {
			var result uint16
			switch current.operation {
			case "AND":
				result = operands[0] & operands[1]
			case "OR":
				result = operands[0] | operands[1]
			case "NOT":
				result = operands[0] ^ 0xFFFF
			case "LSHIFT":
				result = operands[0] << operands[1]
			case "RSHIFT":
				result = operands[0] >> operands[1]
			default:
				result = operands[0]
			}
			circuit[current.dst] = result
		} else {
			queue = append(queue, current)
		}
	}

	return
}

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	n, err := utils.LineCounter(file)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	instSet := make([]Instruction, 0, n)
	for scanner.Scan() {
		instSet = append(
			instSet,
			parseInstruction(scanner.Text()))
	}

	// Part 1
	aValue := install(instSet)["a"]
	println(aValue)

	// Part 2
	overrideIdx := slices.IndexFunc(instSet, func(x Instruction) bool {
		return x.dst == "b"
	})
	instSet[overrideIdx] = Instruction{
		dst:       "b",
		operands:  []string{strconv.Itoa(int(aValue))},
		operation: "",
	}

	aValue = install(instSet)["a"]
	println(aValue)
}
