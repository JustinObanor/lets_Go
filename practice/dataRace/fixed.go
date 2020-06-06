package main

import (
	"fmt"
	"sync"
	"time"
)

func fixed() {
	var counter int
	var wg sync.WaitGroup
	var mu sync.Mutex
	const gr = 100

	wg.Add(gr * 2)

	for i := 0; i < gr; i++ {
		go func() {
			time.Sleep(time.Millisecond * 100)

			mu.Lock()
			counter++
			mu.Unlock()

			wg.Done()
		}()

		go func() {
			time.Sleep(time.Millisecond * 100)

			mu.Lock()
			counter--
			mu.Unlock()

			wg.Done()
		}()
	}
	wg.Wait()

	//value of counter changes on which gorouotine finished first
	fmt.Println("counter: ", counter)
}
