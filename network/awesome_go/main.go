package main

import (
	"awesome_go/lc"
	"fmt"
)

func main() {
	head := lc.ArrayToListNode([]int{-10, -3, 0, 5, 9})
	fmt.Println(head)
	lc.SortedListToBST(head)
}

func detectCycle(head *lc.ListNode) *lc.ListNode {
	set := make(map[*lc.ListNode]struct{})
	for head != nil {
		if _, ok := set[head]; ok {
			return head
		}
		set[head] = struct{}{}
		head = head.Next
	}
	return head
}

func reverseList(head *lc.ListNode) *lc.ListNode {
	if head != nil || head.Next == nil {
		return head
	}
	var tmp = reverseList(head.Next)
	head.Next.Next = head
	head.Next = nil
	return tmp
}
