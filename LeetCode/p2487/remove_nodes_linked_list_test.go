package p2487_test

import (
	"testing"

	"github.com/sergiovaneg/GoStudy/LeetCode/p2487"
	"github.com/sergiovaneg/GoStudy/lists"
)

type ListNode = lists.SinglyLinkedNode[int]

func TestRemoveNodes(t *testing.T) {
	makeList := lists.MakeSinglyLinkedList[int]
	compareLists := lists.CompareSinglyLinkedLists[int]

	var result, expected *ListNode

	result = p2487.RemoveNodes(makeList([]int{5, 2, 13, 3, 8}))
	expected = makeList([]int{13, 8})
	if !compareLists(result, expected) {
		t.Fatalf("Expected %v; got %v", expected, result)
	}

	result = p2487.RemoveNodes(makeList([]int{1, 1, 1, 1}))
	expected = makeList([]int{1, 1, 1, 1})
	if !compareLists(result, expected) {
		t.Fatalf("Expected %v; got %v", expected, result)
	}
}
