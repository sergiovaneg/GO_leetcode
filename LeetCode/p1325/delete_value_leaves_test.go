package p1325_test

import (
	"reflect"
	"testing"

	"github.com/sergiovaneg/GoStudy/LeetCode/p1325"
	"github.com/sergiovaneg/GoStudy/bintree"
)

type TreeNode = bintree.TreeNode[int]

const null = bintree.NullInt

func TestRemoveLeafNodes(t *testing.T) {
	makeTree := bintree.MakeBinaryTree[int]
	var res, wnt *TreeNode

	res = p1325.RemoveLeafNodes(
		makeTree([]int{2, 5, 2, 7, null, null, null, null, 2}, null),
		2)
	wnt = makeTree([]int{2, 5, null, 7}, null)
	if !reflect.DeepEqual(res, wnt) {
		t.Fatalf("Expected:\n%v\nGot:\n%v\n", wnt, res)
	}
}
