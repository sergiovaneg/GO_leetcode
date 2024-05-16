package p0404_test

import (
	"testing"

	"github.com/sergiovaneg/GoStudy/LeetCode/p0404"
	"github.com/sergiovaneg/GoStudy/bintree"
)

type TreeNode = bintree.TreeNode[int]

func TestSumOfLeftLeaves(t *testing.T) {
	var res int
	var root *TreeNode

	root = bintree.MakeBinaryTree([]int{0, 2, 4, 1, bintree.NullInt, 3, -1, 5, 1, bintree.NullInt, 6, bintree.NullInt, 8}, bintree.NullInt)
	res = p0404.SumOfLeftLeaves(root)
	if res != 5 {
		t.Fatalf("Expected 5; got %v", res)
	}

	root = bintree.MakeBinaryTree([]int{0, -4, -3, bintree.NullInt, -1, 8, bintree.NullInt, bintree.NullInt, 3, bintree.NullInt, -9, -2, bintree.NullInt, 4}, bintree.NullInt)
	res = p0404.SumOfLeftLeaves(root)
	if res != 2 {
		t.Fatalf("Expected 2; got %v", res)
	}

	root = bintree.MakeBinaryTree([]int{1}, bintree.NullInt)
	res = p0404.SumOfLeftLeaves(root)
	if res != 0 {
		t.Fatalf("Expected 0; got %v", res)
	}
}
