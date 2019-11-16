package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	x := generate(10)
	info(x)
	xx := make([]int, len(x)*2)
	for i := range x {
		xx[i] = x[i]
		//or copy(dst, src)
	}
	x = xx
	info(x)
}

func generate(i int) []int {
	x := make([]int, i)
	rand.Seed(time.Now().UnixNano())
	for i := range x {
		x[i] = rand.Intn(10) - rand.Intn(10)
	}
	return x
}

func info(x []int) {
	fmt.Printf("%d %d %v\n", len(x), cap(x), x)
}
