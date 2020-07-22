package main

import "fmt"

type Node struct {
	value int
	next  *Node
}

type LinkedList struct {
	head   *Node
	tail   *Node
	length int
}

func newNode(value int) *Node {
	return &Node{
		value: value,
	}
}

func newLinkedList(value, length int) LinkedList {
	Node := newNode(value)
	return LinkedList{
		head:   Node,
		tail:   Node,
		length: length,
	}
}

func (l *LinkedList) search(item int) int {
	if l == nil {
		return 0
	}

	curr := l.head

	for curr != nil {
		if curr.value == item {
			return item
		}
		curr = curr.next
	}
	return 0
}

func (l *LinkedList) getList() []int {
	if l == nil {
		return nil
	}

	curr := l.head

	res := make([]int, 0)
	for curr != nil {
		res = append(res, curr.value)
		curr = curr.next
	}
	return res
}

func (l *LinkedList) prepend(value int) {
	node := newNode(value)

	node.next = l.head
	l.head = node
	l.length++
}

func (l *LinkedList) append(value int) {
	node := newNode(value)

	l.tail.next = node
	l.tail = node
	l.length++
}

func (l *LinkedList) insert(index, value int) {
	if index == 0 {
		l.prepend(value)
	} else if index == l.length {
		l.append(value)
	} else {
		node := newNode(value)

		pre := l.head
		for i := 0; i < index-1; i++ {
			pre = pre.next
		}
		aft := pre.next
		pre.next = node
		node.next = aft
		l.length++
	}
}

func (l *LinkedList) delete(index int) {
	if index > l.length {
		return
	}
	if index == 0 {
		curr := l.head
		l.head = curr.next
		curr = nil
		l.length--
	} else{
		pre := l.head
		for i := 0; i < index-1; i++{
			pre = pre.next
		}
		aft := pre.next
		pre.next = aft.next
		aft = nil
		l.length--
	}
	
}

func main() {
	l := newLinkedList(6, 1)
	l.insert(0, 3)
	l.insert(0, 3)
	l.insert(3, 6)
	l.insert(4, 6)
	l.insert(2, 5)
	l.delete(0)
	l.delete(2)
	l.delete(2)
	fmt.Println(l.getList(), l.length)
}
