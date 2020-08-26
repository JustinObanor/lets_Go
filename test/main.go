package main

import "fmt"

func main() {
	x := []int{5, 1, 2, 6, 2, 0}
	selectionSort(x)
	fmt.Println(x)
}

func selectionSort(s []int) {
	length := len(s)
	for i := 1; i < length; i++ {
		for j := i; j > 0; j-- {
			if s[j] < s[j-1] {
				s[j], s[j-1] = s[j-1], s[j]
			}
		}
		fmt.Println(s)
	}
}
