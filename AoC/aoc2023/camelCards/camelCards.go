package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/sergiovaneg/GoStudy/utils"
)

type Hand struct {
	cards    [5]int
	bid      int
	strength int
}

type HandHeap []Hand

func (hand *Hand) setStrength() {
	cards, count := make([]int, 5), make([]int, 0, 5)
	copy(cards, hand.cards[:])
	slices.Sort(cards)

	count = append(count, 1)
	for idx, card := range cards[1:] {
		if card == cards[idx] {
			count[len(count)-1]++
		} else {
			count = append(count, 1)
		}
	}

	// Joker rule
	n := len(count)
	if cards[0] == 1 && n > 1 {
		aux := count[0]

		count = count[1:]
		n--

		slices.Sort(count)
		count[n-1] += aux
	}

	slices.Sort(count)

	if n == 5 { // High card
		hand.strength = 1
	} else if n == 4 { // One pair
		hand.strength = 2
	} else if n == 3 {
		if count[1] == 2 { // Two pairs
			hand.strength = 3
		} else { // Three of a kind
			hand.strength = 4
		}
	} else if n == 2 {
		if count[1] == 3 { // Full house
			hand.strength = 5
		} else { // Four of a kind
			hand.strength = 6
		}
	} else { // Five of a kind
		hand.strength = 7
	}
}

func processHand(line string) Hand {
	substr := strings.Split(line, " ")
	bid, _ := strconv.Atoi(substr[1])
	cards := [5]int{}

	for idx, char := range substr[0] {
		if char >= '2' && char <= '9' {
			cards[idx] = int(char - '0')
		} else {
			switch char {
			case 'T':
				cards[idx] = 10
			case 'J':
				cards[idx] = 1
			case 'Q':
				cards[idx] = 12
			case 'K':
				cards[idx] = 13
			case 'A':
				cards[idx] = 14
			}
		}
	}

	hand := Hand{bid: bid, cards: cards}
	hand.setStrength()

	return hand
}

func compareHands(a, b Hand) int {
	if a.strength != b.strength {
		return a.strength - b.strength
	}

	for idx := 0; idx < 5; idx++ {
		if a.cards[idx] != b.cards[idx] {
			return a.cards[idx] - b.cards[idx]
		}
	}

	return a.bid - b.bid
}

func (h HandHeap) Len() int           { return len(h) }
func (h HandHeap) Less(i, j int) bool { return compareHands(h[i], h[j]) < 0 }
func (h HandHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *HandHeap) Push(x any)        { *h = append(*h, x.(Hand)) }

func (h *HandHeap) Pop() any {
	n := len(*h)
	x := (*h)[n-1]

	*h = (*h)[:n-1]

	return x
}

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var wg sync.WaitGroup
	var mu sync.Mutex
	n, _ := utils.LineCounter(file)

	hands := make(HandHeap, 0, n)
	heap.Init(&hands)

	wg.Add(n)
	for scanner.Scan() {
		go func(line string) {
			defer wg.Done()
			defer mu.Unlock()

			hand := processHand(line)
			mu.Lock()
			heap.Push(&hands, hand)
		}(scanner.Text())
	}
	wg.Wait()

	var res int
	for rank := 1; hands.Len() > 0; rank++ {
		res += rank * heap.Pop(&hands).(Hand).bid
	}

	fmt.Println(res)
}
