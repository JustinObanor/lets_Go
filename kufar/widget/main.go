package main

import (
	"flag"
	"fmt"
	"strconv"
	"sync"
	"time"
)

type widget struct {
	label string
	time  time.Time
}

var n = flag.Int("n", 1, "how many widgets produced by producer")

func main() {
	c := make(chan widget)
	var wg sync.WaitGroup
	flag.Parse()

	for i := 1; i <= *n; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			c <- widget{label: "widget_" + strconv.Itoa(num), time: time.Now()}
		}(i)
	}

	go func(){
		wg.Wait()
		close(c)
	}()

	consumer(c)
}

func consumer(c chan widget) {
	for elem := range c {
		fmt.Printf("[%s  %v]\n", elem.label, elem.time.Format("15: 04: 05.0000"))
	}
	// close(c)
}
