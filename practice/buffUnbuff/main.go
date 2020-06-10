package main
 
import (
    "fmt"
    "time"
)
 
func main() {
    c1 := make(chan int, 3)
 
    go func(c chan int) {
		for i := 1; i <= 5; i++{
			fmt.Printf("func goroutine #%d starts sending data into the channel\n", i)
			c <- i
			fmt.Printf("func goroutine #%d after sending data into the channel\n", i) 
		}
		close(c1)
       }(c1) 
 
    fmt.Println("main goroutine sleeps for 2 seconds")
    time.Sleep(time.Second * 2)
 
	fmt.Println("main goroutine starts receiving data")
	for v := range c1{
		fmt.Println("main goroutine received data:", v)
	}
 
    // we sleep for a second to give time to the goroutine to finish
    time.Sleep(time.Second * 2)
 
    // After running the program we notice that the sender (the func goroutine) blocks on the channel
    // until the receiver (the main goroutine) receives the data from the channel.
}