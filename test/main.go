package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, job <-chan int, result chan<- int) {
	for j := range job {
		fmt.Printf("worker %d started job %d", id, j)
		fmt.Println()
		time.Sleep(time.Second)
		fmt.Printf("worker %d started job %d", id, j)
		fmt.Println()
		result <- j * 2
	}
}

func main() {
	now := time.Now()
	job := make(chan int, 5)
	result := make(chan int, 5)

	for i := 0; i < 3; i++ {
		go worker(i, job, result)
	}

	for i := 0; i < 5; i++ {
		job <- i
	}
	close(job)

	var wg sync.WaitGroup
	wg.Add(11)
	for i := 0; i <= 10; i++ {

		go func(i int) {
			defer wg.Done()
			fmt.Printf("loop i is - %d\n", i)
		}(i)
	}

	fmt.Println(time.Since(now))
}
