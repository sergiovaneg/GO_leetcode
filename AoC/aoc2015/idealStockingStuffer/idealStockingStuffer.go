package main

import (
	"bufio"
	"crypto/md5"
	"log"
	"math"
	"os"
	"strconv"
	"sync"
)

const Workers = 32

func testInput(input []byte) bool {
	sum := md5.Sum(input)
	//return (sum[0]&0xFF == 0x00) && (sum[1]&0xFF == 0x00) && (sum[2]&0xF0 == 0x00)
	return (sum[0]&0xFF == 0x00) && (sum[1]&0xFF == 0x00) && (sum[2]&0xFF == 0x00)
}

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	scanner.Scan()
	prefix := scanner.Text()

	var wg sync.WaitGroup
	c := make(chan int, Workers)

	for blockStart := 0; len(c) == 0; blockStart += Workers {
		wg.Add(Workers)
		for n := blockStart; n < blockStart+Workers; n++ {
			go func(n int) {
				defer wg.Done()
				input := []byte(prefix + strconv.Itoa(n))
				if testInput(input) {
					c <- n
				}
			}(n)
		}
		wg.Wait()
	}

	close(c)
	minVal := math.MaxInt
	for val := range c {
		if val < minVal {
			minVal = val
		}
	}

	println(minVal)
}
