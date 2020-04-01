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

var wid = flag.Int("w", 1, "how many widgets produced by producer")
var con = flag.Int("c", 1, "widgets consumed")
var dur = flag.Duration("d", time.Second, "a consumer taking a while to process a widget")

func main() {
	c := make(chan widget)
	var wg sync.WaitGroup
	flag.Parse()


	for i := 0; i <= *con; i++ {
		 consumer(c, i)
	}

	ticker := time.NewTicker(*dur)
	var tickCounter int

	wg.Add(*wid)

	for i := 0 ; i <= *wid; i++ {
		if tickCounter >= *wid{
			ticker.Stop()
			return
		}

		go func(num int) {
			for range ticker.C {
				c <- widget{label: "widget_" + strconv.Itoa(num), time: time.Now()}
			}
		}(tickCounter)
		tickCounter++
	}


	wg.Wait()
	close(c)
}

func consumer(c chan widget, con int) {
	for elem := range c {
		fmt.Printf("[%s  %v] consumer_%d\n", elem.label, elem.time.Format("15: 04: 05.0000"), con)
	}
}
