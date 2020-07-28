package main

import "fmt"

type Node struct {
	value string
	next  *Node
}

type Queue struct {
	first  *Node
	last   *Node
	length int
}

func newNode(value string) *Node {
	return &Node{
		value: value,
		next:  nil,
	}
}

func newQueue(value string, length int) Queue {
	node := newNode(value)
	return Queue{
		first:  node,
		last:   node,
		length: length,
	}
}

func (q *Queue) peek() string {
	return q.first.value
}

func (q *Queue) dequeue() string {
	if q.isEmpty() {
		return ""
	}

	if q.first == q.last {
		q.last = nil
	}

	currNode := q.first
	q.first = currNode.next
	q.length--
	return currNode.value
}

func (q *Queue) enqueue(value string) {
	node := newNode(value)
	q.last.next = node
	q.last = node
	q.length++
}

func (q *Queue) getItems() []string {
	items := make([]string, 0, q.length)
	currNode := q.first
	for currNode != nil {
		items = append(items, currNode.value)
		currNode = currNode.next
	}
	return items
}

func (q *Queue) isEmpty() bool {
	return q.length == 0
}

func main() {
	q := newQueue("joy", 1)
	q.enqueue("matt")
	q.enqueue("pavel")
	q.enqueue("samir")
	fmt.Println(q.dequeue())
	fmt.Println(q.dequeue())
	fmt.Println(q.dequeue())
	fmt.Println(q.getItems())
}
