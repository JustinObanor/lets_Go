package main

import (
	"fmt"
	"time"
)

func main() {
	x := []int{1, 3, 2, 8, 7, 5, 6}

	// limiter := time.Tick(1 * time.Second)

	burstyLimiter := make(chan time.Time, 3)

	for i := 0; i < 3; i++ {
		burstyLimiter <- time.Now()
	}

	var count int
	go func() {
		for t := range time.Tick(1 * time.Second) {
			if count >= len(x) {
				break
			}
			burstyLimiter <- t
			count++
		}
	}()

	jobs := make(chan int)

	for i := 1; i <= 3; i++ {
		go work(i, jobs)
	}

	for _, v := range x {
		// <-limiter
		<-burstyLimiter
		fmt.Printf("request %v\n", time.Now())
		jobs <- v
	}
	close(jobs)
}

func work(id int, jobs <-chan int) {
	for j := range jobs {
		fmt.Printf("worker %d received job %d\n", id, j)
		fmt.Printf("worker %d finished job %d\n", id, j)
	}
}
