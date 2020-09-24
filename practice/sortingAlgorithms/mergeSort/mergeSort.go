package main

import "fmt"

func main() {
	x := []int{6, 3, 4, 8, 2, 1, 5, 7}
	fmt.Println(mergeSort(x))
}

func mergeSort(s []int) []int {
	if len(s) < 2 {
		return s
	}

	mid := len(s) / 2

	return merge(mergeSort(s[:mid]), mergeSort(s[mid:]))
}

func merge(left, right []int) []int {
	size, leftIdx, rightIdx := len(left)+len(right), 0, 0
	slice := make([]int, size, size)

	for k := 0; k < size; k++ {
		if leftIdx > len(left)-1 && rightIdx <= len(right)-1 {
			slice[k] = right[rightIdx]
			rightIdx++
		} else if rightIdx > len(right)-1 && leftIdx <= len(left)-1 {
			slice[k] = left[leftIdx]
			leftIdx++
		} else if left[leftIdx] < right[rightIdx] {
			slice[k] = left[leftIdx]
			leftIdx++
		} else {
			slice[k] = right[rightIdx]
			rightIdx++
		}
	}
	return slice
}
