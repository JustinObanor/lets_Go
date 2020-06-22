package main

import "fmt"

func main() {
	sq := make([]int, 64)
	s := 1
	for i := 1; i < 64; i++{
		sq[i] = s
		s *= 2
	}
	fmt.Println(sq)
}