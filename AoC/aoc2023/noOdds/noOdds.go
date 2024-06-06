package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/sergiovaneg/GoStudy/utils"
	"gonum.org/v1/gonum/mat"
)

const lowerBound = 200000000000000
const upperBound = 400000000000000

func parseHailstone(line string) [2][3]int {
	matches := regexp.MustCompile("(-*[0-9]+)").FindAllString(line, 6)
	var numbers [6]int

	for idx, match := range matches {
		numbers[idx], _ = strconv.Atoi(match)
	}

	return [2][3]int{
		{numbers[0], numbers[1], numbers[2]},
		{numbers[3], numbers[4], numbers[5]},
	}
}

func countCollisionsXY(hailstones [][2][3]int) uint {
	var res uint

	for i, ha := range hailstones {
		pa, va := ha[0], ha[1]
		for _, hb := range hailstones[i+1:] {
			pb, vb := hb[0], hb[1]

			det := va[1]*vb[0] - va[0]*vb[1]
			if det == 0 {
				continue
			}

			ta := float64(vb[0]*(pb[1]-pa[1])-vb[1]*(pb[0]-pa[0])) / float64(det)
			tb := float64(va[0]*(pb[1]-pa[1])-va[1]*(pb[0]-pa[0])) / float64(det)

			if ta < 0 || tb < 0 {
				continue
			}

			x, y := float64(pa[0])+float64(va[0])*ta, float64(pa[1])+float64(va[1])*ta
			if x < lowerBound || x > upperBound || y < lowerBound || y > upperBound {
				continue
			}

			res++
		}
	}

	return res
}

func calcInitialPosition(h [][2][3]int) [3]int {
	aData, bData := make([]float64, 0, 6*6), make([]float64, 0, 6)
	for idx := 1; idx <= 3; idx++ {
		aData = append(aData,
			float64(h[0][1][1]-h[idx][1][1]),
			float64(h[idx][1][0]-h[0][1][0]),
			0.,
			float64(h[idx][0][1]-h[0][0][1]),
			float64(h[0][0][0]-h[idx][0][0]),
			0.,
			float64(h[0][1][2]-h[idx][1][2]),
			0.,
			float64(h[idx][1][0]-h[0][1][0]),
			float64(h[idx][0][2]-h[0][0][2]),
			0.,
			float64(h[0][0][0]-h[idx][0][0]))

		bData = append(bData,
			float64(h[0][0][0]*h[0][1][1]-
				h[0][0][1]*h[0][1][0]-
				h[idx][0][0]*h[idx][1][1]+
				h[idx][0][1]*h[idx][1][0]),
			float64(h[0][0][0]*h[0][1][2]-
				h[0][0][2]*h[0][1][0]-
				h[idx][0][0]*h[idx][1][2]+
				h[idx][0][2]*h[idx][1][0]))
	}

	a := mat.NewDense(6, 6, aData)
	b := mat.NewVecDense(6, bData)

	var x mat.VecDense
	err := x.SolveVec(a, b)
	if err != nil {
		log.Fatal(err)
	}

	return [3]int{
		int(x.AtVec(0)),
		int(x.AtVec(1)),
		int(x.AtVec(2)),
	}
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

	hailstones := make([][2][3]int, n)
	for idx := 0; scanner.Scan(); idx++ {
		hailstones[idx] = parseHailstone(scanner.Text())
	}

	println(countCollisionsXY(hailstones))

	initialPosition := calcInitialPosition(hailstones)
	println(initialPosition[0] + initialPosition[1] + initialPosition[2])
}
