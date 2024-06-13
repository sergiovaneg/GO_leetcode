package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

type Instruction struct {
	operation string
	operands  []string
	dst       string
}

type Circuit map[string]func() uint16

func parseInstruction(line string) Instruction {
	arrowIdx := strings.Index(line, "->")
	dst := line[arrowIdx+3:]
	operands := strings.Split(line[:arrowIdx], " ")

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

func (circuit *Circuit) install(inst Instruction) {
	reDigits := regexp.MustCompile("([0-9]+)")
	switch inst.operation {
	case "AND":
		isDigit := [2]bool{
			reDigits.MatchString(inst.operands[0]),
			reDigits.MatchString(inst.operands[1]),
		}
		if isDigit[0] && isDigit[1] {
			a, _ := strconv.Atoi(inst.operands[0])
			b, _ := strconv.Atoi(inst.operands[1])
			val := uint16(a) & uint16(b)
			(*circuit)[inst.dst] = func() uint16 {
				return val
			}
		} else if isDigit[0] {
			a, _ := strconv.Atoi(inst.operands[0])
			val := uint16(a)
			(*circuit)[inst.dst] = func() uint16 {
				return val & (*circuit)[inst.operands[1]]()
			}
		} else if isDigit[1] {
			a, _ := strconv.Atoi(inst.operands[1])
			val := uint16(a)
			(*circuit)[inst.dst] = func() uint16 {
				return val & (*circuit)[inst.operands[0]]()
			}
		} else {
			(*circuit)[inst.dst] = func() uint16 {
				return (*circuit)[inst.operands[0]]() & (*circuit)[inst.operands[1]]()
			}
		}
	case "OR":
		isDigit := [2]bool{
			reDigits.MatchString(inst.operands[0]),
			reDigits.MatchString(inst.operands[1]),
		}
		if isDigit[0] && isDigit[1] {
			a, _ := strconv.Atoi(inst.operands[0])
			b, _ := strconv.Atoi(inst.operands[1])
			val := uint16(a) | uint16(b)
			(*circuit)[inst.dst] = func() uint16 {
				return val
			}
		} else if isDigit[0] {
			a, _ := strconv.Atoi(inst.operands[0])
			val := uint16(a)
			(*circuit)[inst.dst] = func() uint16 {
				return val | (*circuit)[inst.operands[1]]()
			}
		} else if isDigit[1] {
			a, _ := strconv.Atoi(inst.operands[1])
			val := uint16(a)
			(*circuit)[inst.dst] = func() uint16 {
				return val | (*circuit)[inst.operands[0]]()
			}
		} else {
			(*circuit)[inst.dst] = func() uint16 {
				return (*circuit)[inst.operands[0]]() | (*circuit)[inst.operands[1]]()
			}
		}
	case "LSHIFT":
		isDigit := [2]bool{
			reDigits.MatchString(inst.operands[0]),
			reDigits.MatchString(inst.operands[1]),
		}
		if isDigit[0] && isDigit[1] {
			a, _ := strconv.Atoi(inst.operands[0])
			b, _ := strconv.Atoi(inst.operands[1])
			val := uint16(a) << uint16(b)
			(*circuit)[inst.dst] = func() uint16 {
				return val
			}
		} else if isDigit[0] {
			a, _ := strconv.Atoi(inst.operands[0])
			val := uint16(a)
			(*circuit)[inst.dst] = func() uint16 {
				return val << (*circuit)[inst.operands[1]]()
			}
		} else if isDigit[1] {
			a, _ := strconv.Atoi(inst.operands[1])
			val := uint16(a)
			(*circuit)[inst.dst] = func() uint16 {
				return val << (*circuit)[inst.operands[0]]()
			}
		} else {
			(*circuit)[inst.dst] = func() uint16 {
				return (*circuit)[inst.operands[0]]() << (*circuit)[inst.operands[1]]()
			}
		}
	case "RSHIFT":
		isDigit := [2]bool{
			reDigits.MatchString(inst.operands[0]),
			reDigits.MatchString(inst.operands[1]),
		}
		if isDigit[0] && isDigit[1] {
			a, _ := strconv.Atoi(inst.operands[0])
			b, _ := strconv.Atoi(inst.operands[1])
			val := uint16(a) >> uint16(b)
			(*circuit)[inst.dst] = func() uint16 {
				return val
			}
		} else if isDigit[0] {
			a, _ := strconv.Atoi(inst.operands[0])
			val := uint16(a)
			(*circuit)[inst.dst] = func() uint16 {
				return val >> (*circuit)[inst.operands[1]]()
			}
		} else if isDigit[1] {
			a, _ := strconv.Atoi(inst.operands[1])
			val := uint16(a)
			(*circuit)[inst.dst] = func() uint16 {
				return val >> (*circuit)[inst.operands[0]]()
			}
		} else {
			(*circuit)[inst.dst] = func() uint16 {
				return (*circuit)[inst.operands[0]]() >> (*circuit)[inst.operands[1]]()
			}
		}
	case "NOT":
		isDigit := reDigits.MatchString(inst.operands[0])
		if isDigit {
			a, _ := strconv.Atoi(inst.operands[0])
			val := uint16(a) ^ 0xFFFF
			(*circuit)[inst.dst] = func() uint16 {
				return val
			}
		} else {
			(*circuit)[inst.dst] = func() uint16 {
				return (*circuit)[inst.operands[0]]() ^ 0xFFFF
			}
		}
	default:
		isDigit := reDigits.MatchString(inst.operands[0])
		if isDigit {
			a, _ := strconv.Atoi(inst.operands[0])
			val := uint16(a)
			(*circuit)[inst.dst] = func() uint16 {
				return val
			}
		} else {
			(*circuit)[inst.dst] = func() uint16 {
				return (*circuit)[inst.operands[0]]()
			}
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
	circuit := make(Circuit)
	for scanner.Scan() {
		circuit.install(parseInstruction(scanner.Text()))
	}
	println(circuit["i"]())
}
