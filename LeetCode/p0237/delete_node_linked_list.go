package p0237

import "github.com/sergiovaneg/GoStudy/lists"

type ListNode = lists.SinglyLinkedNode[int]

func DeleteNode(node *ListNode) {
	aux := node.Next
	if aux == nil {
		node = nil
		return
	}

	node.Val = aux.Val
	node.Next = aux.Next
	aux = nil
}
