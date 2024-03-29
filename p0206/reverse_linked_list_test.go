package p0206_test

import (
	"testing"

	lists "github.com/sergiovaneg/GO_leetcode/Lists"
	"github.com/sergiovaneg/GO_leetcode/p0206"
)

func TestReverseList(t *testing.T) {
	makeList := lists.MakeSinglyLinkedList
	compareList := lists.CompareSinglyLinkedLists

	original := makeList([]int{1, 2, 3, 4, 5})
	reversed := makeList([]int{5, 4, 3, 2, 1})
	if !compareList(p0206.ReverseList(original), reversed) {
		t.Fatal()
	}
}
