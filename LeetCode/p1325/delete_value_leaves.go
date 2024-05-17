package p1325

import "github.com/sergiovaneg/GoStudy/bintree"

type TreeNode = bintree.TreeNode[int]

func RemoveLeafNodes(root *TreeNode, target int) *TreeNode {
	var shouldRemove func(*TreeNode) bool

	shouldRemove = func(tn *TreeNode) bool {
		if tn == nil {
			return true
		}

		rLeft, rRight := shouldRemove(tn.Left), shouldRemove(tn.Right)
		if rLeft && rRight && tn.Val == target {
			return true
		}
		if rLeft {
			tn.Left = nil
		}
		if rRight {
			tn.Right = nil
		}
		return false
	}

	rRoot := shouldRemove(root)
	if rRoot {
		return nil
	}
	return root
}
