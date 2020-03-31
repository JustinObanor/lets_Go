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
var con = flag.Int("c", 1, "how many widgets consumed by consumer")
var d = flag.Int64("d", 1000, "a consumer taking a while to process a widget")

func main() {
	c := make(chan widget)
	var wg sync.WaitGroup
	flag.Parse()

	for i := 0 ; i <= *n; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			c <- widget{label: "widget_" + strconv.Itoa(num), time: time.Now()}

			tick := time.NewTicker(time.Millisecond * time.Duration(*d))
			for range tick.C {
				c <- widget{label: "widget_" + strconv.Itoa(num), time: time.Now()}
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	for i := 0; i <= *con; i++ {
		wg.Add(1)
		go func(num int){
			defer wg.Done()
			consumer(c, num)
		}(i)
	}
	wg.Wait()
}

func consumer(c chan widget, con int) {
	for elem := range c {
		fmt.Printf("[%s  %v] consumer_%d\n", elem.label, elem.time.Format("15: 04: 05.0000"), con)
	}
}
