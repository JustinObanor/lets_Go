package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	x := []int{1, 5, 3, 2, 6, 8, 9, 7, 5}

	wg.Add(len(x))
	for _, v := range x {
		go func(i int) {
			defer wg.Done()
			fmt.Println(i)
		}(v)
	}
	wg.Wait()
}
