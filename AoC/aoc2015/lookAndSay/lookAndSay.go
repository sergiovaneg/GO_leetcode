package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"sync"
)

const reps = 50

type Memory struct {
	sync.RWMutex
	record map[string]string
}

func (m *Memory) parallelLauncher(substr string, sink *string,
	wg *sync.WaitGroup) {
	defer wg.Done()
	*sink = m.lookSay(substr)
}

func (m *Memory) lookSay(input string) string {
	m.RLock()
	result, ok := m.record[input]
	m.RUnlock()
	if ok {
		return result
	}

	n := len(input)
	if n_halves := n >> 1; n&0x01 == 0 && input[n_halves] != input[n_halves-1] {
		var left, right string
		var wg sync.WaitGroup

		wg.Add(2)
		go m.parallelLauncher(input[:n_halves], &left, &wg)
		go m.parallelLauncher(input[n_halves:], &right, &wg)
		wg.Wait()

		result = left + right

		m.Lock()
		m.record[input] = result
		m.Unlock()

		return result
	}

	result = ""
	for i := 0; i < n; {
		j := i + 1
		for j < n {
			if input[i] != input[j] {
				break
			}
			j++
		}
		result += strconv.Itoa(j-i) + string(input[i])
		i = j
	}

	m.Lock()
	m.record[input] = result
	m.Unlock()

	return result
}

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	input := scanner.Text()

	m := Memory{
		record: make(map[string]string),
	}

	for iter := 0; iter < reps; iter++ {
		println(iter)
		input = m.lookSay(input)
	}
	println(len(input))
}
