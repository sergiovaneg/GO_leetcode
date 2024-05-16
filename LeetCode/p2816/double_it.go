package p2816

import lists "github.com/sergiovaneg/GoStudy/Lists"

type ListNode = lists.SinglyLinkedNode[int]

func doubleCarry(node *ListNode) bool {
	node.Val <<= 1
	if node.Next != nil && doubleCarry(node.Next) {
		node.Val++
	}

	carry := node.Val >= 10
	if carry {
		node.Val -= 10
	}
	return carry
}

func DoubleIt(head *ListNode) *ListNode {
	carry := doubleCarry(head)
	if carry {
		return &ListNode{Val: 1, Next: head}
	}
	return head
}
