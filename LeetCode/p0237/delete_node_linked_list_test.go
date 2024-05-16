package p0237_test

import (
	"testing"

	"github.com/sergiovaneg/GoStudy/LeetCode/p0237"
	"github.com/sergiovaneg/GoStudy/lists"
)

type ListNode = lists.SinglyLinkedNode[int]

func TestDeleteNode(t *testing.T) {
	var head, expected *ListNode

	head = lists.MakeSinglyLinkedList([]int{4, 5, 1, 9})
	expected = lists.MakeSinglyLinkedList([]int{4, 1, 9})
	p0237.DeleteNode(head.Next)
	if !lists.CompareSinglyLinkedLists(head, expected) {
		t.Fatal("Wrong procedure")
	}

	head = lists.MakeSinglyLinkedList([]int{4, 5, 1, 9})
	expected = lists.MakeSinglyLinkedList([]int{4, 5, 9})
	p0237.DeleteNode(head.Next.Next)
	if !lists.CompareSinglyLinkedLists(head, expected) {
		t.Fatal("Wrong procedure")
	}
}
