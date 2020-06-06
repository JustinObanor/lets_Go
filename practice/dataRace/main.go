package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var counter int
	var wg sync.WaitGroup
	const gr = 100

	wg.Add(gr * 2)

	for i := 0; i < gr; i++ {
		go func() {
			time.Sleep(time.Millisecond * 100)
			counter++
			wg.Done()
		}()

		go func() {
			time.Sleep(time.Millisecond * 100)
			counter--
			wg.Done()
		}()
	}
	wg.Wait()

	//value of counter changes on which gorouotine finished first
	fmt.Println("counter: ", counter)
}
