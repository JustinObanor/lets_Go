//usefull when list is almost or already sorted
//good for small data sets

package main

import "fmt"

func main() {
	s := []int{99, 44, 6, 2, 1, 5, 63, 87, 283, 4, 0}
	insertionSort(s)
	fmt.Println(s)
}

func insertionSort(s []int) {
	length := len(s)
	for i := 1; i < length; i++ {
		for j := i; j > 0; j-- {
			if s[j] < s[j-1] {
				s[j], s[j-1] = s[j-1], s[j]
			}
		}
	}
}
