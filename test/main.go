package main

import (
	"fmt"
	"sync"
)

var counter int

func main() {
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			num := counter
			num++
			counter = num
		}()
	}

	fmt.Println(counter)
	wg.Wait()
}
