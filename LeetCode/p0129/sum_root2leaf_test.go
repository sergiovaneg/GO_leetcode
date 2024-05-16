package p0129_test

import (
	"testing"

	"github.com/sergiovaneg/GoStudy/LeetCode/p0129"
	"github.com/sergiovaneg/GoStudy/bintree"
)

func TestSumNumbers(t *testing.T) {
	var res int

	MakeBinaryTree := bintree.MakeBinaryTree[int]

	res = p0129.SumNumbers(MakeBinaryTree([]int{1, 2, 3}, bintree.NullInt))
	if res != 25 {
		t.Fatalf("Expected 25; got %v", res)
	}

	res = p0129.SumNumbers(MakeBinaryTree([]int{4, 9, 0, 5, 1},
		bintree.NullInt))
	if res != 1026 {
		t.Fatalf("Expected 1026; got %v", res)
	}
}
