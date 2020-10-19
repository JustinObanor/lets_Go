package main

import (
	"errors"
	"fmt"
)

type LinkedList struct {
	head   *Node
	tail   *Node
	length int
}

type Node struct {
	value    int
	next     *Node
	previous *Node
}

func newNode(value int) *Node {
	return &Node{
		value:    value,
		next:     nil,
		previous: nil,
	}
}

func newLinkedList(value, length int) LinkedList {
	newNode := newNode(value)

	return LinkedList{
		head:   newNode,
		tail:   newNode,
		length: length,
	}
}

func (l *LinkedList) append(value int) {
	newNode := newNode(value)

	//pointing the last node to the new one and making it our tail
	l.tail.next = newNode
	newNode.previous = l.tail
	l.tail = newNode
	l.length++
}

func (l *LinkedList) prepend(value int) {
	newNode := newNode(value)

	//pointing the newNode to the head and making it our head
	newNode.next = l.head
	l.head.previous = newNode
	l.head = newNode
	l.length++
}

func (l *LinkedList) insertAt(index, value int) error {
	if index > l.length {
		l.append(value)
		return errors.New("index out of bound, appending instead")
	}

	if index == 0 {
		l.prepend(value)
	} else if index >= l.length {
		l.append(value)
	} else {
		newNode := newNode(value)
		pre := l.head

		for i := 0; i < index-1; i++ {
			pre = pre.next
		}
		aft := pre.next

		newNode.next = aft
		newNode.previous = pre
		aft.previous = newNode
		pre.next = newNode
		l.length++
	}
	return nil
}

func (l *LinkedList) deleteAt(index int) error {
	currentNode := l.head

	if index >= l.length {
		return errors.New("index out of bound")
	}

	if index == 0 {
		l.head = currentNode.next
		currentNode = nil
		l.length--
	} else {
		pre := l.head
		for i := 0; i < index-1; i++ {
			pre = pre.next
		}
		aft := pre.next
		pre.next = aft.next
		aft.previous = pre
		aft = nil
		l.length--
	}
	return nil
}

func (l *LinkedList) lookup(item int) bool {
	currentNode := l.head
	for currentNode != nil {
		if currentNode.value == item {
			return true
		}
		currentNode = currentNode.next
	}
	return false
}

func (l *LinkedList) getList() []int {
	values := make([]int, 0, l.length)
	currentNode := l.head

	for currentNode != nil {
		values = append(values, currentNode.value)
		currentNode = currentNode.next
	}
	return values
}

func main() {
	l := newLinkedList(10, 1)
	l.append(5)
	l.append(16)
	l.prepend(1)
	if err := l.insertAt(1, 99); err != nil {
		fmt.Println(err)
	}
	if err := l.deleteAt(2); err != nil {
		fmt.Println(err)
	}
	if err := l.deleteAt(0); err != nil {
		fmt.Println(err)
	}
	if err := l.deleteAt(2); err != nil {
		fmt.Println(err)
	}
	// if err := l.deleteAt(2); err != nil {
	// 	fmt.Println(err)
	// }

	fmt.Println("length: ", l.length)
	fmt.Println("list 	", l.getList())
}