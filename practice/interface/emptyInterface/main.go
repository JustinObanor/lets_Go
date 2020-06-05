package main

import "fmt"

//emptyInterface
type emptyInterface interface{}

type person struct {
	info interface{}
}

func main() {
	var e emptyInterface
	e = 12
	fmt.Println(e)
	e = "Yo"
	fmt.Println(e)
	e = []int{1, 2, 3, 4}
	fmt.Println(len(e.([]int)))

	p1 := person{info: 1}
	fmt.Println(p1)
	p2 := person{info: "james"}
	fmt.Println(p2)
	p3 := person{info: []int{1, 2, 3}}
	fmt.Println(p3)
}
