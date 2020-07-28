package main

import "fmt"

type Node struct {
	value string
	next  *Node
}

type Stack struct {
	top    *Node
	bottom *Node
	length int
}

func newNode(value string) Node {
	return Node{
		value: value,
		next:  nil,
	}
}

func newStack(value string, length int) Stack {
	node := newNode(value)
	return Stack{
		top:    &node,
		bottom: &node,
		length: length,
	}
}

func (s Stack) peek() string {
	return s.top.value
}

//making new item our top, and pointing it to prev top
func (s *Stack) push(value string) {
	node := newNode(value)
	currNode := s.top
	s.top = &node
	s.top.next = currNode
	s.length++
}

func (s *Stack) pop() string {
	if s.isEmpty() {
		return ""
	}

	if s.bottom == s.top{
		s.bottom = nil
	}

	currNode := s.top
	s.top = currNode.next
	s.length--
	return currNode.value
}

func (s *Stack) isEmpty() bool {
	return s.length == 0
}

func (s Stack) getValue() []string {
	res := make([]string, 0, s.length)
	currNode := s.top
	for currNode != nil {
		res = append(res, currNode.value)
		currNode = currNode.next
	}
	return res
}

func main() {
	s := newStack("google", 1)
	s.push("udemy")
	s.push("discord")
	fmt.Println(s.pop())
	fmt.Println(s.pop())
	fmt.Println(s.getValue())
}
