package p0234

import lists "github.com/sergiovaneg/GO_leetcode/Lists"

type ListNode = lists.SinglyLinkedNode

func IsPalindrome(head *ListNode) bool {
	// Get number of nodes
	N := 1
	for node := head; node.Next != nil; node = node.Next {
		N++
	}

	// Early returns
	if N == 1 {
		return true
	}
	if N == 2 {
		return head.Val == head.Next.Val
	}
	if N == 3 {
		return head.Val == head.Next.Next.Val
	}

	// Invert half the list
	var prev *ListNode = nil
	current := head
	for idx := 0; idx < N/2; idx++ {
		tmp := current.Next
		current.Next = prev
		prev = current
		current = tmp
	}

	// Skip middle node if odd node count
	if N%2 == 1 {
		current = current.Next
	}

	// Check symmetry
	for current.Next != nil {
		if current.Val != prev.Val {
			return false
		}
		current = current.Next
		prev = prev.Next
	}

	return current.Val == prev.Val
}
