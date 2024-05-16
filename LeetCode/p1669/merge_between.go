package p1669

import "github.com/sergiovaneg/GoStudy/lists"

type ListNode = lists.SinglyLinkedNode[int]

func MergeInBetween(list1 *ListNode, a int, b int, list2 *ListNode) *ListNode {
	idx := 0
	ptr := list1
	for idx < a-1 {
		idx++
		ptr = ptr.Next
	}

	idx++
	qtr := ptr.Next
	for idx < b {
		idx++
		qtr = qtr.Next
	}

	ptr.Next = list2
	for list2.Next != nil {
		list2 = list2.Next
	}
	list2.Next = qtr.Next

	return list1
}
