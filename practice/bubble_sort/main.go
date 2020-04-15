package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	slice := generateSlice(20)
	fmt.Println(slice)
	bubbleSort(slice)
	fmt.Println(slice)

}

func generateSlice(size int) []int {
	slice := make([]int, size, size)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		slice[i] = rand.Intn(10) - rand.Intn(10)
	}
	return slice
}

func bubbleSort(slice []int) {
	size := len(slice)
	if size < 2 {
		return
	}
	for i := 0; i < size; i++ {
		for j := size - 1; j >= 1; j-- {
			if slice[j] < slice[j-1] {
				slice[j], slice[j-1] = slice[j-1], slice[j]
			}
		}
	}
}
