package p2331

import (
	"github.com/sergiovaneg/GoStudy/bintree"
)

type TreeNode = bintree.TreeNode[int]

func EvaluateTree(root *TreeNode) bool {
	if root.Val == 2 {
		return EvaluateTree(root.Left) || EvaluateTree(root.Right)
	}
	if root.Val == 3 {
		return EvaluateTree(root.Left) && EvaluateTree(root.Right)
	}

	return root.Val == 1
}
