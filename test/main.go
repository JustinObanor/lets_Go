package main

import "fmt"

func main() {
	c1 := make(chan int, 3)
	c2 := make(chan int, 3)

	c1 <- 1
	c1 <- 2
	c1 <- 3
	c2 <- 6


	for i := 0; i < 3; i++ {
		select {
		case v1 := <-c1:
			fmt.Println(v1)
		case v2 := <-c2:
			fmt.Println(v2)
		}
	}

}
