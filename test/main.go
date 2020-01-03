package main

import "time"

import "fmt"

func main() {
	c := make(chan string)

	go func() {
		time.Sleep(time.Second)
		c <- "hello there"
	}()

	go func() {
		time.Sleep(time.Second * 2)
		c <- "general kenobi"
	}()

	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-c:
			fmt.Println(msg1)
		case msg1 := <-c:
			fmt.Println(msg1)
		case <-time.After(time.Second * 3):
			fmt.Println("timed out")
		default:
			fmt.Println("nothing received")
		}
	}
}
