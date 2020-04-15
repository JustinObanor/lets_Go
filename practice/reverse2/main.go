package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	slice := generator(10)
	fmt.Println(slice)
	//rev := reverser(slice)
	//fmt.Println(rev)
	rev2 := reverser2(slice)
	fmt.Println(rev2)
}

func reverser(x []int) []int {
	for i := 0; i < len(x)/2; i++ {
		j := len(x) - i - 1
		x[i], x[j] = x[j], x[i]
	}
	return x
}

func reverser2(x []int) []int {
	for i, j := 0, len(x)-1; i < j; i, j = i+1, j-1 {
		x[i], x[j] = x[j], x[i]
	}
	return x
}

func generator(size int) []int {
	slice := make([]int, size)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		slice[i] = rand.Intn(10) - rand.Intn(10)
	}
	return slice
}
