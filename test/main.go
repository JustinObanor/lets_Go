package main

import "fmt"

func main() {
	fmt.Println(moveZerous(123))
}

func moveZerous(nums int) int {
	var rev int

	for nums != 0 {
		rev = rev * 10 + nums % 10
		nums /= 10
	}
	return rev
}
