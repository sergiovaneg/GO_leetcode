package p0143_test

import (
	"testing"

	"github.com/sergiovaneg/GoStudy/LeetCode/p0143"
	"github.com/sergiovaneg/GoStudy/lists"
)

type ListNode = lists.SinglyLinkedNode[int]

func TestReorderList(t *testing.T) {
	list_1 := lists.MakeSinglyLinkedList([]int{1, 2, 3, 4, 5})
	p0143.ReorderList(list_1)
	if !lists.CompareSinglyLinkedLists(list_1,
		lists.MakeSinglyLinkedList([]int{1, 5, 2, 4, 3})) {
		t.Fatal()
	}

	list_2 := lists.MakeSinglyLinkedList([]int{1, 2, 3, 4})
	p0143.ReorderList(list_2)
	if !lists.CompareSinglyLinkedLists(list_2,
		lists.MakeSinglyLinkedList([]int{1, 4, 2, 3})) {
		t.Fatal()
	}
}
