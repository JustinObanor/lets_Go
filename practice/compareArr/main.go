package main

import "fmt"

//2 arrays to check if the contain similar items
func main() {
	arr1 := []rune{'a', 'b', 'c', 'x'}
	arr2 := []rune{'z', 'y', 'x'}
	fmt.Println(commonItems(arr1, arr2))
}

func commonItems(arr1, arr2 []rune) bool {
	m := make(map[rune]bool)
	for _, v := range arr1 {
		if _, ok := m[v]; !ok {
			m[v] = true
		}
	}

	for _, v := range arr2 {
		if _, ok := m[v]; ok {
			fmt.Println("FOUND", string(v))
			return true
		}
	}

	return false
}
