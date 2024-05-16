package p0206

import "github.com/sergiovaneg/GoStudy/lists"

type ListNode = lists.SinglyLinkedNode[int]

/* Iterative version
func ReverseList(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}

	last := head
	current := head.Next
	next := current.Next
	head.Next = nil

	for next != nil {
		current.Next = last
		last = current
		current = next
		next = next.Next
	}

	current.Next = last
	head.Next = nil
	return current
}
*/

func reversePair(node1, node2 *ListNode) *ListNode {
	aux := node2.Next
	node2.Next = node1
	if aux == nil {
		return node2
	} else {
		return reversePair(node2, aux)
	}
}

func ReverseList(head *ListNode) *ListNode {
	if head == nil {
		return head
	}
	return reversePair(nil, head)
}
