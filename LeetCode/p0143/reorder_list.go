package p0143

import "github.com/sergiovaneg/GoStudy/lists"

type ListNode = lists.SinglyLinkedNode[int]

func ReorderList(head *ListNode) {
	// Get list length
	N := 1
	for aux := head; aux.Next != nil; aux = aux.Next {
		N++
	}

	// Early return
	if N < 3 {
		return
	}

	// Reverse second half of the list
	aux := head
	for idx := 0; idx < N>>1; idx++ {
		aux = aux.Next
	}
	tail := aux.Next
	aux.Next = nil
	for tail.Next != nil {
		tmp := tail.Next
		tail.Next = aux
		aux = tail
		tail = tmp
	}
	tail.Next = aux

	// Reorder list
	aux = head
	for tail.Next != nil {
		tmp_1 := aux.Next
		tmp_2 := tail.Next

		aux.Next = tail
		tail.Next = tmp_1

		aux = tmp_1
		tail = tmp_2
	}
}
