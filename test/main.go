package main

import "fmt"

func main() {
	s := []int{1, 4, 3, 6, 7, 8, 5}
	bubbleSort(s)
	fmt.Println(s)
}

func bubbleSort(s []int) {
var n = len(s)
for i := 0; i < n; i++{
var minIdx = i
for j := i; j < n; j++{
	if s[j] < s[minIdx]{
		minIdx = j
	}
}
}
}
