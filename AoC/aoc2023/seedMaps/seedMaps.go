package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
)

func getSeeds(line string) [][2]int {
	matches := regexp.MustCompile("([0-9]+)").FindAllString(line, -1)
	seed_ranges := make([]int, len(matches))

	for idx, match := range matches {
		val, err := strconv.Atoi(match)
		if err == nil {
			seed_ranges[idx] = val
		}
	}

	seeds := make([][2]int, 0, len(matches)/2)

	for idx := 0; idx < len(seed_ranges); idx += 2 {
		seeds = append(seeds, [2]int{seed_ranges[idx], seed_ranges[idx+1]})
	}

	slices.SortFunc(seeds, func(a, b [2]int) int { return a[0] - b[0] })
	return seeds
}

func getInterval(line string) ([3]int, error) {
	matches := regexp.MustCompile("([0-9]+)").FindAllString(line, 3)
	result := [3]int{}
	if matches == nil {
		return result, errors.New("no numbers")
	}

	for idx, match := range matches {
		val, err := strconv.Atoi(match)
		if err == nil {
			result[idx] = val
		}
	}

	return result, nil
}

func seedDefrag(seeds [][2]int) [][2]int {
	slices.SortFunc(seeds, func(a, b [2]int) int { return a[0] - b[0] })

	for idx := 0; idx < len(seeds)-1; {
		if seeds[idx][0]+seeds[idx][1] == seeds[idx+1][0] {
			aux := [2]int{seeds[idx][0], seeds[idx][1] + seeds[idx+1][1]}
			seeds = slices.Replace(seeds, idx, idx+2, aux)
		} else {
			idx++
		}
	}

	return seeds
}

func mapSeeds(seeds [][2]int, seed_map [][3]int) [][2]int {
	nMap := len(seed_map)

	slices.SortFunc(seed_map, func(a, b [3]int) int { return a[1] - b[1] })

	for idxSeed, idxMap := 0, 0; idxSeed < len(seeds) && idxMap < nMap; {
		sStt, sEnd := seeds[idxSeed][0], seeds[idxSeed][0]+seeds[idxSeed][1]
		mStt, mEnd := seed_map[idxMap][1], seed_map[idxMap][1]+seed_map[idxMap][2]

		if sEnd <= mStt {
			idxSeed++
		} else if sStt >= mEnd {
			idxMap++
		} else {
			if sStt < mStt && sEnd > mEnd {
				aux := [][2]int{
					{sStt, mStt - sStt},
					{seed_map[idxMap][0], seed_map[idxMap][2]},
					{mEnd, sEnd - mEnd},
				}
				seeds = slices.Replace(seeds, idxSeed, idxSeed+1, aux...)
				idxSeed += 2
				idxMap++
			} else if sStt < mStt {
				aux := [][2]int{
					{sStt, mStt - sStt},
					{seed_map[idxMap][0], sEnd - mStt},
				}
				seeds = slices.Replace(seeds, idxSeed, idxSeed+1, aux...)
				idxSeed += 2
			} else if sEnd > mEnd {
				aux := [][2]int{
					{seed_map[idxMap][0] + (sStt - mStt), mEnd - sStt},
					{mEnd, sEnd - mEnd},
				}
				seeds = slices.Replace(seeds, idxSeed, idxSeed+1, aux...)
				idxSeed++
				idxMap++
			} else {
				aux := [2]int{seed_map[idxMap][0] + (sStt - mStt), seeds[idxSeed][1]}
				seeds[idxSeed] = aux
				idxSeed++
			}
		}
	}

	return seedDefrag(seeds)
}

func main() {
	file, err := os.Open("./input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	seeds, seed_map := getSeeds(scanner.Text()), make([][3]int, 0)
	scanner.Scan()

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			seeds = mapSeeds(seeds, seed_map)
			seed_map = make([][3]int, 0)
		} else {
			interval, err := getInterval(line)
			if err == nil {
				seed_map = append(seed_map, interval)
			}
		}
	}

	if len(seed_map) > 0 {
		seeds = mapSeeds(seeds, seed_map)
		seed_map = nil
	}

	fmt.Println(slices.MinFunc(seeds, func(a, b [2]int) int {
		return a[0] - b[0]
	})[0])
}
