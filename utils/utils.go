package utils

import (
	"bytes"
	"io"
	"log"
	"os"
)

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func LineCounter(r *os.File) (int, error) {
	offset, err := r.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	buf := make([]byte, 512*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			r.Seek(offset, io.SeekStart)
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func GCD[I Integer](a, b I) I {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func LCM[I Integer](a, b I, integers ...I) I {
	result := a * b / GCD(a, b)

	for i := 0; i < len(integers); i++ {
		result = LCM(result, integers[i])
	}

	return result
}

func ExtGCD[I Integer](a, b I) (I, I, I) {
	rOld, r := a, b
	sOld, s := I(1), I(0)
	tOld, t := I(0), I(1)
	for r != 0 {
		quotient, remainder := rOld/r, rOld%r
		rOld, r = r, remainder
		sOld, s = s, sOld-quotient*s
		tOld, t = t, tOld-quotient*t
	}

	return rOld, sOld, tOld
}

func ExtLCM[I Integer](mu []I, lambda []I) (I, I) {
	muC, lambdaC := mu[0], lambda[0]

	for idx := range mu[1:] {
		g, s, _ := ExtGCD(lambdaC, lambda[idx+1])
		phaseDiff := muC - mu[idx+1]
		pdMult, pdRem := phaseDiff/g, phaseDiff%g

		if pdRem != 0 {
			log.Fatal("No synchronization possible")
		}

		lambdaNew := (lambdaC / g) * lambda[idx+1]
		muNew := (muC - s*pdMult*lambdaC) % lambdaNew
		lambdaC, muC = lambdaNew, muNew
	}

	return muC, lambdaC
}
