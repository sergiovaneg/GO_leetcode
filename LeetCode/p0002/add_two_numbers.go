package p0002

import "github.com/sergiovaneg/GoStudy/lists"

type ListNode = lists.SinglyLinkedNode[int]

func addAndUpdate(new_tmp int, current *ListNode) int {
	if new_tmp >= 10 {
		current.Next = &ListNode{Val: new_tmp - 10, Next: nil}
		return 1
	} else {
		current.Next = &ListNode{Val: new_tmp, Next: nil}
		return 0
	}
}

func AddTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	current := new(ListNode)
	res := current

	tmp := 0
	for (l1 != nil) && (l2 != nil) {
		tmp = addAndUpdate(tmp+l1.Val+l2.Val, current)

		l1 = l1.Next
		l2 = l2.Next

		current = current.Next
	}
	for l1 != nil {
		tmp = addAndUpdate(tmp+l1.Val, current)
		l1 = l1.Next
		current = current.Next
	}
	for l2 != nil {
		tmp = addAndUpdate(tmp+l2.Val, current)
		l2 = l2.Next
		current = current.Next
	}
	if tmp == 1 {
		current.Next = &ListNode{Val: 1, Next: nil}
	}

	return res.Next
}
