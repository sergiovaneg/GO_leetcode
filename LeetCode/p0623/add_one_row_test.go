package p0623_test

import (
	"testing"

	"github.com/sergiovaneg/GoStudy/LeetCode/p0623"
	"github.com/sergiovaneg/GoStudy/bintree"
)

type TreeNode = bintree.TreeNode[int]

func TestAddOneRow(t *testing.T) {
	null := bintree.NullInt
	var expected, res *TreeNode

	res = p0623.AddOneRow(bintree.MakeBinaryTree(
		[]int{4, 2, 6, 3, 1, 5}, null), 1, 2)
	expected = bintree.MakeBinaryTree(
		[]int{4, 1, 1, 2, null, null, 6, 3, 1, 5}, null)
	if !bintree.CompareBinaryTree(res, expected) {
		t.Fatalf("Expected:\n%v\nGot:\n%v", expected, res)
	}

	t.Logf("\n%v", res)

	res = p0623.AddOneRow(bintree.MakeBinaryTree(
		[]int{4, 2, null, 3, 1}, null), 1, 3)
	expected = bintree.MakeBinaryTree(
		[]int{4, 2, null, 1, 1, 3, null, null, 1}, null)
	if !bintree.CompareBinaryTree(res, expected) {
		t.Fatalf("Expected:\n%v\nGot:\n%v", expected, res)
	}
}
