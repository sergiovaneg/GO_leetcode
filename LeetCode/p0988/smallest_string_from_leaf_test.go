package p0988_test

import (
	"testing"

	"github.com/sergiovaneg/GoStudy/LeetCode/p0988"
	"github.com/sergiovaneg/GoStudy/bintree"
)

const null = bintree.NullInt

func TestSmallestFromLeaf(t *testing.T) {
	var res string
	MakeTree := bintree.MakeBinaryTree[int]

	res = p0988.SmallestFromLeaf(
		MakeTree([]int{0, 1, 2, 3, 4, 3, 4}, null))
	if res != "dba" {
		t.Fatalf("Expected 'dba'; got %v", res)
	}

	res = p0988.SmallestFromLeaf(
		MakeTree([]int{25, 1, 3, 1, 3, 0, 2}, null))
	if res != "adz" {
		t.Fatalf("Expected 'adz'; got %v", res)
	}

	res = p0988.SmallestFromLeaf(
		MakeTree([]int{2, 2, 1, null, 1, 0, null, 0}, null))
	if res != "abc" {
		t.Fatalf("Expected 'abc'; got %v", res)
	}
}
