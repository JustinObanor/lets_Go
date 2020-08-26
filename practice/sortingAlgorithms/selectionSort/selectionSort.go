package main

import (
	"fmt"
)

func main() {
	s := []int{99, 44, 6, 2, 1, 5, 63, 87, 283, 4, 0}
	selectionSort(s)
	fmt.Println(s)
}

func selectionSort(s []int) {
	length := len(s)
	for i := 0; i < length; i++ {
		min := i
		for j := i + 1; j < length; j++ {
			if s[j] < s[min] {
				min = j
			}
		}
		s[i], s[min] = s[min], s[i]
	}
}
