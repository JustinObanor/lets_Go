package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"
)

type widget struct {
	label string
	time  time.Time
}

var n = flag.Int("n", 1, "how many widgets produced by producer")

func main() {
	c := make(chan widget)

	flag.Parse()

	for i := 0; i < *n; i++ {
		time.Sleep(time.Second)
		go func(num int) {
			c <- widget{"widget_" + strconv.Itoa(num), time.Now()}
		}(i)
	}

	consumer(c)
}

func consumer(c chan widget) {
	for i := 0; i < *n; i++ {
		res := <-c
		fmt.Printf("[%s  %v]\n", res.label, res.time.Format("15: 04: 05.0000"))
	}

}
