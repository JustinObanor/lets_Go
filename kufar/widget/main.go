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
var con = flag.Int("c", 1, "number of consumers")
var dur = flag.Duration("d", time.Second, "a consumer taking a while to process a widget")

func main() {
	c := make(chan widget)
	var cwg sync.WaitGroup
	var pwg sync.WaitGroup
	flag.Parse()

	cwg.Add(*con)
	for i := 0; i <= *con; i++ {
		go consumer(c, i, &cwg)
	}

	ticker := time.NewTicker(*dur)
	var tickCounter int

	pwg.Add(*wid)
	for range ticker.C {
		if tickCounter >= *wid {
			ticker.Stop()
			break
		}

		go func(num int) {
			c <- widget{label: "widget_" + strconv.Itoa(num), time: time.Now()}
			pwg.Done()
		}(tickCounter)

		tickCounter++
	}

	go func() {
		cwg.Wait()
		close(c)
	}()
}

func consumer(c <-chan widget, con int, cwg *sync.WaitGroup) {
	for elem := range c {
		fmt.Printf("[%s  %v] consumer_%d\n", elem.label, elem.time.Format("15: 04: 05.0000"), con)
	}
	cwg.Done()
}
