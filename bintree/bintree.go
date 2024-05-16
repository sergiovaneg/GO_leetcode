package bintree

import (
	"fmt"
	"math"
	"slices"
	"strings"
)

const NullInt = math.MinInt

type TreeNode[T any] struct {
	Val   T
	Left  *TreeNode[T]
	Right *TreeNode[T]
}

func countNulls[T comparable](elements []T, null T) int {
	var acc int
	for _, e := range elements {
		if e == null {
			acc++
		}
	}
	return acc
}

func MakeBinaryTree[T comparable](elements []T, null T) *TreeNode[T] {
	l := len(elements)
	if l == 0 || elements[0] == null {
		return nil
	}

	root := new(TreeNode[T])
	root.Val = elements[0]

	block_size := 1
	var left_skips, right_skips int
	left := make([]T, 0, l)
	right := make([]T, 0, l)

	for idx := 1; idx < l; {
		limit := min(idx+block_size-left_skips, l)
		left = append(left, elements[idx:limit]...)
		left_skips += countNulls(elements[idx:limit], null)
		idx = limit

		if idx == l {
			break
		}

		limit = min(idx+block_size-right_skips, l)
		right = append(right, elements[idx:limit]...)
		right_skips += countNulls(elements[idx:limit], null)
		idx = limit

		block_size *= 2
		left_skips *= 2
		right_skips *= 2
	}

	// Memory optimization
	elements = nil
	left = slices.Clip(left)
	right = slices.Clip(right)

	root.Left = MakeBinaryTree(left, null)
	root.Right = MakeBinaryTree(right, null)

	return root
}

func CompareBinaryTree[T comparable](root1, root2 *TreeNode[T]) bool {
	if root1 == nil || root2 == nil {
		return root1 == root2
	}
	return root1.Val == root2.Val && CompareBinaryTree(root1.Left, root2.Left) && CompareBinaryTree(root1.Right, root2.Right)
}

func GetMaxDepth[T any](root *TreeNode[T]) int {
	var rec_level func(*TreeNode[T], int) int
	rec_level = func(root *TreeNode[T], depth int) int {
		if root == nil {
			return depth
		}
		return max(rec_level(root.Left, depth+1), rec_level(root.Right, depth+1))
	}

	return rec_level(root, 0)
}

func (root *TreeNode[T]) String() string {
	max_depth := GetMaxDepth(root)
	result := make([]string, max_depth)

	var rec_level func(*TreeNode[T], int, []string)
	rec_level = func(root *TreeNode[T], level int, result []string) {
		if root == nil {
			for i, reps := level, 1; i < max_depth; i, reps = i+1, reps*2 {
				level_separator := strings.Repeat("\t",
					int(math.Pow(2, float64(max_depth-i))))
				for j := 0; j < reps; j++ {
					result[i] += "null" + level_separator
				}
			}
		} else {
			level_separator := strings.Repeat("\t",
				int(math.Pow(2, float64(max_depth-level))))
			result[level] += fmt.Sprint(root.Val) + level_separator
			rec_level(root.Left, level+1, result)
			rec_level(root.Right, level+1, result)
		}
	}

	rec_level(root, 0, result)

	for level := 0; level < max_depth; level++ {
		pad := strings.Repeat("\t", int(math.Pow(2, float64(max_depth-level-1)))-1)
		result[level] = pad + result[level]
	}

	return strings.Join(result, "\n")
}
