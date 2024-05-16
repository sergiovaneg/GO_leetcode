package p2331_test

import (
	"testing"

	"github.com/sergiovaneg/GoStudy/LeetCode/p2331"
	"github.com/sergiovaneg/GoStudy/bintree"
)

func TestEvaluateTree(t *testing.T) {
	const null = bintree.NullInt

	if !p2331.EvaluateTree(
		bintree.MakeBinaryTree([]int{2, 1, 3, null, null, 0, 1}, null)) {
		t.Fatal("Expected True; got False")
	}

	if p2331.EvaluateTree(
		bintree.MakeBinaryTree([]int{0}, null)) {
		t.Fatal("Expected False; got True")
	}
}
