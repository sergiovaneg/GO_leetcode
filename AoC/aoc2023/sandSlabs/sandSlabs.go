package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"

	"github.com/sergiovaneg/GoStudy/utils"
)

type Brick [2][3]int
type BrickStack []Brick

type Support struct {
	bottom []*Brick
	top    []*Brick
}

type SupportMap map[*Brick]Support

func parseBrick(line string) Brick {
	coordinates := regexp.MustCompile("([0-9]+)").FindAllString(line, 6)
	var vals [6]int

	for idx, val := range coordinates {
		vals[idx], _ = strconv.Atoi(val)
	}

	return Brick{
		[3]int(vals[:3]),
		[3]int(vals[3:]),
	}
}

func (stack BrickStack) settleStack() SupportMap {
	slices.SortFunc(stack, func(a, b Brick) int {
		return b[0][2] - a[0][2]
	})

	supportMap := make(SupportMap, len(stack))

	for idx := len(stack) - 1; idx >= 0; idx-- {
		var offset int
		brick := &stack[idx]

		supporters := make([]*Brick, 0, 1)
		for supportIdx := range stack[idx:] {
			x := &stack[idx+supportIdx]

			if supportIdx == 0 || (*x)[1][2] < offset { // Early skips
				continue
			}

			// Skip if no overlap
			if (*brick)[0][0] > (*x)[1][0] || (*x)[0][0] > (*brick)[1][0] {
				continue
			}
			if (*brick)[0][1] > (*x)[1][1] || (*x)[0][1] > (*brick)[1][1] {
				continue
			}

			if (*x)[1][2] > offset { // New reference
				supporters = []*Brick{x}
				offset = (*x)[1][2]
			} else { // Append to supporters
				supporters = append(supporters, x)
			}
		}

		if len(supporters) > 0 {
			supportMap[brick] = Support{
				bottom: supporters,
				top:    make([]*Brick, 0),
			}
			for _, x := range supporters {
				supportMap[x] = Support{
					bottom: supportMap[x].bottom,
					top:    append(supportMap[x].top, brick),
				}
			}

		}

		// Reference -> Offset
		offset = (*brick)[0][2] - (offset + 1)

		(*brick)[0][2] -= offset
		(*brick)[1][2] -= offset
	}

	return supportMap
}

func (sm SupportMap) isSafe(brick *Brick) bool {
	for _, topBrick := range sm[brick].top {
		if len(sm[topBrick].bottom) == 1 {
			return false
		}
	}

	return true
}

func (sm SupportMap) wouldFall(brick *Brick, log map[*Brick]bool) bool {
	wouldFall := true

	for _, x := range sm[brick].bottom {
		if !log[x] {
			wouldFall = false
			break
		}
	}

	return wouldFall
}

func (sm SupportMap) chainCount(brick *Brick) uint {
	if len(sm[brick].top) == 0 { // No reaction caused
		return 0
	}

	var result uint

	// 'brick' logged but not counted
	log := make(map[*Brick]bool)
	log[brick] = true

	queue := make([]*Brick, 0, len(sm[brick].top))
	for _, x := range sm[brick].top {
		if sm.wouldFall(x, log) {
			queue = append(queue, x)
		}
	}

	var x *Brick
	for len(queue) > 0 {
		// Pop from queue
		x, queue = queue[0], queue[1:]

		if log[x] { // Already accounted for
			continue
		}

		// Register and add
		result, log[x] = result+1, true

		for _, y := range sm[x].top {
			if sm.wouldFall(y, log) {
				queue = append(queue, y)
			}
		}
	}

	return result
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

	stack := make(BrickStack, n)
	for idx := 0; scanner.Scan(); idx++ {
		stack[idx] = parseBrick(scanner.Text())
	}

	// I can see the future
	supportMap := stack.settleStack()

	// Part 1
	var result uint
	for idx := range stack {
		if supportMap.isSafe(&stack[idx]) {
			result++
		}
	}
	println(result)

	// Part 2
	result = 0
	for idx := range stack {
		result += supportMap.chainCount(&stack[idx])
	}
	println(result)
}
