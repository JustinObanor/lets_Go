package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	x := []int{1, 3, 5, 4, 3, 6, 8, 5, 3, 5}
	c := make(chan int)
	wg.Add(len(x))
	for _, v := range x {
		go func(v int) {
			defer wg.Done()
			c <- v
		}(v)
	}
	for i := 0; i < len(x); i++ {
		fmt.Println(<-c)
	}
	wg.Wait()
}
