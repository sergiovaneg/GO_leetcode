package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
	"sync"

	"github.com/sergiovaneg/GoStudy/utils"
)

type Brick [2][3]int
type BrickStack []Brick
type Log map[*Brick]bool

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
	n := len(stack)

	slices.SortFunc(stack, func(a, b Brick) int {
		return b[0][2] - a[0][2]
	})

	supportMap := make(SupportMap, n)

	// First offset is localized
	offset := stack[n-1][0][2] - 1
	stack[n-1][0][2] -= offset
	stack[n-1][1][2] -= offset

	// The rest depend on previous modifications
	for idx := n - 2; idx >= 0; idx-- {
		var reference int
		brick := &stack[idx]

		supporters := make([]*Brick, 0, 1)
		for supportIdx := range stack[idx+1:] {
			x := &stack[idx+supportIdx+1]

			// Skip if
			if (*x)[1][2] < reference {
				continue
			}

			// Skip if no overlap
			if (*brick)[0][0] > (*x)[1][0] || (*x)[0][0] > (*brick)[1][0] {
				continue
			}
			if (*brick)[0][1] > (*x)[1][1] || (*x)[0][1] > (*brick)[1][1] {
				continue
			}

			if (*x)[1][2] > reference { // New reference
				supporters = []*Brick{x}
				reference = (*x)[1][2]
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
		offset = (*brick)[0][2] - (reference + 1)

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

func (sm SupportMap) wouldFall(brick *Brick, log Log) bool {
	wouldFall := true

	for _, x := range sm[brick].bottom {
		if !log[x] {
			wouldFall = false
			break
		}
	}

	return wouldFall
}

func (sm SupportMap) catchFallen(base *Brick, log Log) []*Brick {
	queue := make([]*Brick, 0, len(sm[base].top))

	for _, x := range sm[base].top {
		if sm.wouldFall(x, log) {
			queue = append(queue, x)
		}
	}

	return slices.Clip(queue)
}

func (sm SupportMap) chainCount(brick *Brick) uint {
	if len(sm[brick].top) == 0 { // No reaction caused
		return 0
	}

	var result uint

	log := make(Log, len(sm))
	log[brick] = true // Logged but not counted
	queue := sm.catchFallen(brick, log)

	var x *Brick
	for len(queue) > 0 {
		// Pop from queue
		x, queue = queue[0], queue[1:]

		// Register and add
		result, log[x] = result+1, true

		// Queue newly fallen bricks
		queue = append(queue, sm.catchFallen(x, log)...)
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
	var wg sync.WaitGroup
	c := make(chan uint, len(stack))

	wg.Add(len(stack))
	for idx := range stack {
		go func(idx int) {
			defer wg.Done()
			c <- supportMap.chainCount(&stack[idx])
		}(idx)
	}
	wg.Wait()
	close(c)

	result = 0
	for val := range c {
		result += val
	}
	println(result)
}
