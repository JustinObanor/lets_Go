package main

import (
	"fmt"
)

func main() {
	s := []int{99, 44, 6, 2, 1, 5, 63, 87, 283, 4, 0}
	bubbleSort(s)
	fmt.Println(s)
}

func bubbleSort(s []int) {
	length := len(s)
	for i := 0; i < length; i++ {
		for j := 0; j < length-1; j++ {
			if s[j] > s[j+1] {
				s[j], s[j+1] = s[j+1], s[j]
			}
		}
	}
}
