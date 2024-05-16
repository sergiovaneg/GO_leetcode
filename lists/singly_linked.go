package lists

import "fmt"

type SinglyLinkedNode[T any] struct {
	Val  T
	Next *SinglyLinkedNode[T]
}

func MakeSinglyLinkedList[T any](elements []T) *SinglyLinkedNode[T] {
	list := new(SinglyLinkedNode[T])
	current := list
	for _, e := range elements {
		current.Next = &SinglyLinkedNode[T]{Val: e, Next: nil}
		current = current.Next
	}
	return list.Next
}

func CompareSinglyLinkedLists[T comparable](
	list1 *SinglyLinkedNode[T],
	list2 *SinglyLinkedNode[T]) bool {
	for list1 != nil && list2 != nil {
		if list1.Val != list2.Val {
			return false
		}
		list1 = list1.Next
		list2 = list2.Next
	}

	if list1 != nil || list2 != nil {
		return false
	}
	return true
}

func (node *SinglyLinkedNode[T]) String() string {
	s, first := "[", true
	for node != nil {
		if first {
			first = false
		} else {
			s += ","
		}
		s += fmt.Sprint(node.Val)
		node = node.Next
	}
	return s + "]"
}
