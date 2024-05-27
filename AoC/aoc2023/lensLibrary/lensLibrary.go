package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Lens struct {
	label       string
	focalLength int
}

type BoxMap map[byte]*[]Lens

func getHash(label string) byte {
	var result byte
	for _, char := range label {
		result += byte(char)
		result *= 17
	}
	return result
}

func (bm *BoxMap) update(step string) {
	sepIdx := strings.IndexAny(step, "-=")

	if sepIdx == -1 {
		return
	}

	separator := step[sepIdx]
	values := strings.Split(step, string(separator))

	hash := getHash(values[0])
	lenses, ok := (*bm)[hash]

	lensIdx := -1
	if ok {
		lensIdx = slices.IndexFunc(*lenses, func(x Lens) bool {
			return x.label == values[0]
		})
	}

	if separator == '-' {
		if lensIdx != -1 {
			*lenses = slices.Delete(*lenses, lensIdx, lensIdx+1)
		}
		return
	} else {
		focalLength, err := strconv.Atoi(values[1])
		if err != nil {
			return
		}

		aux := Lens{
			label:       step[:sepIdx],
			focalLength: focalLength,
		}

		if ok {
			if lensIdx != -1 {
				(*lenses)[lensIdx] = aux
			} else {
				(*lenses) = append((*lenses), aux)
			}
		} else {
			(*bm)[hash] = &[]Lens{aux}
		}
	}
}

func (bm BoxMap) getPower() uint {
	var power uint

	for box, lenses := range bm {
		for slot, lens := range *lenses {
			power += (1 + uint(box)) * (1 + uint(slot)) * uint(lens.focalLength)
		}
	}

	return power
}

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	sequence := scanner.Text()

	boxes := make(BoxMap, 256)
	for _, step := range strings.Split(sequence, ",") {
		boxes.update(step)
	}

	fmt.Println(boxes.getPower())
}
