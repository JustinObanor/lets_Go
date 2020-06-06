package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var counter int
	var wg sync.WaitGroup
	var mu sync.Mutex
	const gr = 100

	wg.Add(gr)

	for i := 0; i < gr; i++ {
		go func() {
			mu.Lock()

			num := counter
			time.Sleep(time.Millisecond * 100)
			num++
			counter = num

			mu.Unlock()
			wg.Done()
		}()

	}
	wg.Wait()

	//value of counter changes on which gorouotine finished first
	fmt.Println("counter: ", counter)
}
