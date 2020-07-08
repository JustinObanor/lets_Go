package main

import "fmt"

func main() {
	s := []int{1, 2, 3, 4}
	a := [...]int{1, 2, 3, 4}

	ss := s
	aa := a
	s[0] = 2
	a[0] = 2

	fmt.Println(ss[0])
	fmt.Println(s[0])
	fmt.Println(aa[0])
	fmt.Println(a[0])

}
