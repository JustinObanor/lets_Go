package main

import (
	"fmt"
	"sync"
)

var x = []int{1, 3, 6, 2, 9, 7, 5}
var wg sync.WaitGroup

func main() {
	wg.Add(len(x))
	for _, v := range x {
		go func(i int) {
			print(i)
		}(v)
	}
	wg.Wait()
}

func print(i int) {
	fmt.Println(i)
	wg.Done()
}
