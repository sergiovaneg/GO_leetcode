package p2591_test

import (
	"testing"

	"github.com/sergiovaneg/GoStudy/LeetCode/p2591"
)

func TestDistMoney(t *testing.T) {
	var res int

	res = p2591.DistMoney(20, 3)
	if res != 1 {
		t.Fatalf("Expected 1; got %v", res)
	}

	res = p2591.DistMoney(16, 2)
	if res != 2 {
		t.Fatalf("Expected 2; got %v", res)
	}
}
