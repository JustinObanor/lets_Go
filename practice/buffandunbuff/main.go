package main
 
import (
    "fmt"
    "time"
)
 
func main() {
    c1 := make(chan int, 1) 
 
    // Launching a goroutine
    go func(c chan int) {
        fmt.Println("func goroutine starts sending data into the channel")
        c <- 10
        fmt.Println("func goroutine after sending data into the channel")
    }(c1) // calling the anonymous func and passing c1 as argument
 
    fmt.Println("main goroutine sleeps for 2 seconds")
    time.Sleep(time.Second * 2)
 
    fmt.Println("main goroutine starts receiving data")
    d := <-c1
    fmt.Println("main goroutine received data:", d)
 
    // we sleep for a second to give time to the goroutine to finish
    time.Sleep(time.Second * 2)
 
    // After running the program we notice that the sender (the func goroutine) blocks on the channel
    // until the receiver (the main goroutine) receives the data from the channel.
    /*
    An unbuffered channel(snychronous channel) is used to 
    perform synchronous communication 
    between goroutines while a buffered
     channel is used for perform 
     asynchronous communication. 
     An unbuffered channel provides 
     a guarantee that an exchange 
     between two goroutines is performed 
     at the instant the send and receive take place.
    */
    //every send operation is syncrhronized with its receive operation for unbuffered
}