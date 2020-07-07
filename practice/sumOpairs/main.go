package main

import "fmt"

func main() {
	// s := []int{1, 2, 4, 4}
	//s := []int{6,4,2,2,2,2,2,1,2} 5
	s := []int{1, 2, 4, 5, 6, 7, 8, 9, 11}
	fmt.Println(checkPair(s, 8))
}

//linear time for sorted list
func checkPair(data []int, sum int) bool {
	low, high := 0, len(data)-1
	for low < high {
		s := data[low] + data[high]
		if s == sum {
			fmt.Println(data[low], data[high])
			return true
		}
		if s > sum {
			high--
		} else if s < sum {
			low++
		}
	}
	return false
}

//linear time for unsorted list 
//space complexity O(a)
func isPair(data []int, sum int) bool {
	comp := make(map[int]bool)
	for i := 0; i < len(data); i++ {
		if _, ok := comp[data[i]]; ok {
			return true
		}
		comp[sum-data[i]] = true
	}
	return false
}
