package main

import (
	"context"
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

	wg.Add(*con)
	for i := 0; i <= *con; i++ {
		go consumer(context.Background(), c, &wg, i)
	}

	ticker := time.NewTicker(*dur)
	var tickCounter int

	wg.Add(*wid)
	for range ticker.C {
		if tickCounter >= *wid {
			ticker.Stop()
			break
		}

		go func(num int) {
			c <- widget{label: "widget_" + strconv.Itoa(num), time: time.Now()}
			wg.Done()
		}(tickCounter)

		tickCounter++
	}

	go func() {
		wg.Wait()
		close(c)
	}()
}

func consumer(ctx context.Context, c chan widget, wg *sync.WaitGroup, con int) {
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			break
		case elem := <-c:
			fmt.Printf("[%s  %v] consumer_%d\n", elem.label, elem.time.Format("15: 04: 05.0000"), con)
		}
	}
}
