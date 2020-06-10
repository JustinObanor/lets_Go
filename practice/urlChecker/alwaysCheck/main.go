package main

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"
)

func checkURL(url string, c chan<- string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s id Down. Error %v\n", url, err)
		c <- url
	} else {
		defer resp.Body.Close()
		fmt.Printf("%s is UP -> STATUS CODE: %d\n", url, resp.StatusCode)
		c <- url
	}
}

func main() {
	sites := []string{"https://google.com", "https://facebook.com", "https://stackoverflow.com"}

	c := make(chan string)

	for _, url := range sites {
		go checkURL(url, c)
	}

	fmt.Println("Num of grs", runtime.NumGoroutine())

	for {
		go checkURL(<-c, c)
		fmt.Println(strings.Repeat("#", 30))
		time.Sleep(time.Second)
	}

	/*
		for url := range c {
		time.Sleep(time.Second * 2)
		go checkURL(url, c)
	}
	*/

	/*
		for url := range c {
		go func(u string) {
			time.Sleep(time.Second * 2)
			checkURL(u, c)
		}(url)
	}
	*/
}
