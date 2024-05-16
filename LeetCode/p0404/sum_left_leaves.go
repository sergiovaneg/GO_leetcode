package p0404

import "github.com/sergiovaneg/GoStudy/bintree"

type TreeNode = bintree.TreeNode[int]

func recursiveLeftSum(branch *TreeNode) int {
	if branch.Left == nil && branch.Right == nil {
		return branch.Val
	}

	var val int
	if branch.Left != nil {
		val = recursiveLeftSum(branch.Left)
	}
	if branch.Right != nil {
		if branch.Right.Left != nil || branch.Right.Right != nil {
			val += recursiveLeftSum(branch.Right)
		}
	}
	return val
}

func SumOfLeftLeaves(root *TreeNode) int {
	if root.Left == nil && root.Right == nil {
		return 0
	}

	return recursiveLeftSum(root)
}
