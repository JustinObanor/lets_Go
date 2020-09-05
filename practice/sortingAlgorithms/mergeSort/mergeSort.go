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

	fmt.Println("splitting", s)
	mid := len(s) / 2

	left := s[:mid]
	right := s[mid:]
	fmt.Println("left", left, "right", right)

	return merge(mergeSort(left), mergeSort(right))
}

func merge(left, right []int) []int {
	size, leftIdx, rightIdx := len(left)+len(right), 0, 0
	slice := make([]int, size, size)

	fmt.Println("merging", left, "and", right)

	fmt.Println("size", size, "lftIdx", leftIdx, "rgtIdx", rightIdx)
	fmt.Println("left", left, "right", right)

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
