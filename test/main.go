package main

import "fmt"

func main() {
fmt.Println(reverse(123))
}

func reverse(x int) int {
	var rev int

	for x != 0 {
		pop := x % 10
		x /= 10

		rev = rev * 10 + pop
	}
	return rev
}
