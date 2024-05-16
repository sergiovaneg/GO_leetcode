package p0988

import "github.com/sergiovaneg/GoStudy/bintree"

type TreeNode = bintree.TreeNode[int]

func lexiCompare(s1, s2 string) bool {
	l1, l2 := len(s1), len(s2)
	for idx, lim := 0, min(l1, l2); idx < lim; idx++ {
		if s1[idx] != s2[idx] {
			return s1[idx] > s2[idx]
		}
	}
	return l1 > l2
}

func SmallestFromLeaf(root *TreeNode) string {
	var recursive_step func(*TreeNode, string) string
	recursive_step = func(root *TreeNode, s string) string {
		self := string(rune(root.Val+0x61)) + s
		if root.Left == nil && root.Right == nil {
			return self
		}
		var left, right string
		if root.Left != nil {
			left = recursive_step(root.Left, self)
		}
		if root.Right != nil {
			right = recursive_step(root.Right, self)
		}

		if len(left) == 0 {
			return right
		}
		if len(right) == 0 {
			return left
		}

		if lexiCompare(left, right) {
			return right
		} else {
			return left
		}
	}

	return recursive_step(root, "")
}
