package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	s := generate(10)
	fmt.Println(s)
}

func generate(i int) []int {
	s := make([]int, i)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < len(s); i++ {
		s[i] = rand.Intn(10) - rand.Intn(10)
	}
	return s
}

func info(x []int) {
	fmt.Printf("%d %d %v\n", len(x), cap(x), x)
}
