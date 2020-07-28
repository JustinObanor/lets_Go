package main

import "fmt"

type Stack struct {
	data []string
}

func newStack(length int) Stack {
	return Stack{
		data: make([]string, 0, length),
	}
}

func (s Stack) peek() string {
	return s.data[len(s.data)-1]
}

func (s *Stack) push(value string) {
	s.data = append(s.data, value)
}

func (s *Stack) pop() string {
	if len(s.data) == 0 {
		return ""
	}

	res := s.data[len(s.data)-1]

	s.data = s.data[:len(s.data)-1]

	return res
}

func (s *Stack) isEmpty() bool {
	return len(s.data) == 0
}

func (s Stack) getValue() []string {
	return s.data
}

func main() {
	s := newStack(10)
	s.push("google")
	s.push("udemy")
	s.push("discord")
	fmt.Println(s.peek())
	fmt.Println(s.getValue())
}
