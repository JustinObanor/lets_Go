package main

func main() {
	deleteDuplicates(&ListNode{1, &ListNode{1, &ListNode{2, nil}}})
}

type ListNode struct {
	Val  int
	Next *ListNode
}

/*
	h
	c		1
	1	->		->	2
*/

func deleteDuplicates(head *ListNode) *ListNode {
	currNode := head

	for currNode != nil {
		if currNode.Val == currNode.Next.Val {
			currNode.Next = currNode.Next.Next
			currNode = currNode.Next
			break
		} else {
			currNode = currNode.Next
		}
	}
	return head
}
