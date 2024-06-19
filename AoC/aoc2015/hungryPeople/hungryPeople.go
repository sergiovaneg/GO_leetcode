package main

import (
	"bufio"
	"log"
	"math"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"sync"

	"github.com/sergiovaneg/GoStudy/utils"
)

const CAPACITY = 100
const CALORYTARGET = 500
const LAMBDA = 1e18
const WORKERS = 12

type Ingredient []int

func parseIngredient(line string) Ingredient {
	ingredient := make(Ingredient, 0, 5)
	for _, match := range regexp.MustCompile(
		"(-*[0-9]+)").FindAllString(line, 5) {
		val, _ := strconv.Atoi(match)
		ingredient = append(ingredient, val)
	}
	return ingredient
}

type Recipe map[*Ingredient]int

func scoreNormal(r Recipe) int {
	var individualScore []int
	for ingredient, amount := range r {
		if individualScore == nil {
			// Ignore calories
			individualScore = make([]int, len(*ingredient)-1)
		}

		for idx := range individualScore {
			individualScore[idx] += amount * (*ingredient)[idx]
		}
	}

	totalScore := 1
	for _, score := range individualScore {
		totalScore *= max(score, 0)
	}

	return totalScore
}

func scoreCaloric(r Recipe) int {
	var individualScore []int
	for ingredient, amount := range r {
		if individualScore == nil {
			// Ignore calories
			individualScore = make([]int, len(*ingredient)-1)
		}

		for idx := range individualScore {
			individualScore[idx] += amount * (*ingredient)[idx]
		}
	}

	totalScore := 1
	for _, score := range individualScore {
		totalScore *= max(score, 0)
	}

	calories, nProps := 0, len(individualScore)
	for ingredient, amount := range r {
		calories += amount * (*ingredient)[nProps]
	}
	penalization := int(LAMBDA * math.Abs(float64(calories-CALORYTARGET)))

	return totalScore - penalization
}

func getRandomRecipe(ingredients []Ingredient) Recipe {
	r, acc := make(Recipe, len(ingredients)), 0
	for idx := range ingredients {
		r[&ingredients[idx]] = 0
	}

	for idx := range ingredients {
		amount := 1 + rand.Intn(CAPACITY-acc)
		r[&ingredients[idx]] = amount
		acc += amount
		if acc == CAPACITY {
			break
		}
	}
	r[&ingredients[0]] += CAPACITY - acc // Enforce capacity

	return r
}

func (r Recipe) copy() Recipe {
	newR := make(Recipe, len(r))
	for ingredient, amount := range r {
		newR[ingredient] = amount
	}
	return newR
}

func (r Recipe) swap(a, b *Ingredient) Recipe {
	newR := r.copy()
	if newR[b] == 0 {
		return newR
	}

	newR[a]++
	newR[b]--

	return newR
}

func (r Recipe) isEqual(other Recipe) bool {
	for ingredient := range r {
		if r[ingredient] != other[ingredient] {
			return false
		}
	}
	return true
}

func (r Recipe) getNext(scoringFunction func(r Recipe) int) Recipe {
	bestR := r.copy()
	for a := range r {
		for b := range r {
			if a == b {
				continue
			}
			if newR := r.swap(a, b); scoringFunction(newR) >= scoringFunction(bestR) {
				bestR = newR
			}
		}
	}

	if r.isEqual(bestR) {
		return nil
	}

	return bestR
}

func getLocalOptimum(ingredients []Ingredient, localPatience int,
	scoringFunction func(Recipe) int) Recipe {
	localR := getRandomRecipe(ingredients)
	for localCounter := 0; localCounter <= localPatience; {
		nextR := localR.getNext(scoringFunction)
		if nextR == nil {
			break
		}
		if scoringFunction(nextR) == scoringFunction(localR) {
			localCounter++
		} else {
			localCounter = 0
		}
		localR = nextR
	}
	return localR
}

func getOptimalRecipe(
	ingredients []Ingredient,
	globalPatience, localPatience int,
	scoringFunction func(Recipe) int) Recipe {
	bestR := getRandomRecipe(ingredients)

	var wg sync.WaitGroup

	for globalCounter := 0; globalCounter <= globalPatience; {
		c := make(chan Recipe, WORKERS)
		wg.Add(WORKERS)
		for range WORKERS {
			go func() {
				defer wg.Done()
				c <- getLocalOptimum(ingredients, localPatience, scoringFunction)
			}()
		}
		wg.Wait()
		close(c)

		for localR := range c {
			if scoringFunction(localR) > scoringFunction(bestR) {
				bestR = localR
				globalCounter = 0
			} else {
				globalCounter++
			}
		}
	}

	return bestR
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

	ingredients := make([]Ingredient, 0, n)
	for scanner.Scan() {
		ingredients = append(ingredients, parseIngredient(scanner.Text()))
	}
	r := getOptimalRecipe(ingredients, 100, 10, scoreNormal)
	println(scoreNormal(r))

	rCaloric := getOptimalRecipe(ingredients, 50000, 10, scoreCaloric)
	println(scoreNormal(rCaloric))
}
