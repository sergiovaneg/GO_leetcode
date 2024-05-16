package p0206_test

import (
	"testing"

	"github.com/sergiovaneg/GoStudy/LeetCode/p0206"
	"github.com/sergiovaneg/GoStudy/lists"
)

func TestReverseList(t *testing.T) {
	makeList := lists.MakeSinglyLinkedList[int]
	compareList := lists.CompareSinglyLinkedLists[int]

	original := makeList([]int{1, 2, 3, 4, 5})
	reversed := makeList([]int{5, 4, 3, 2, 1})
	if !compareList(p0206.ReverseList(original), reversed) {
		t.Fatal()
	}
}
