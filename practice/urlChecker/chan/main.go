package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"
)

func checkAndSaveBody(url string, c chan<- string) {
	resp, err := http.Get(url)
	if err != nil {
		c <- fmt.Sprintf("%s is down: %v", url, err)
	} else {
		defer resp.Body.Close()

		s := fmt.Sprintf("%s -> STATUS CODE: %d\n", url, resp.StatusCode)

		if resp.StatusCode == 200 {
			bs, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				s += fmt.Sprintf("%v", err)
			}

			fileName := strings.Split(url, "//")
			s += fmt.Sprintf("Writing response body to %s.txt\n", strings.Trim(fileName[1], "[]"))

			if err := ioutil.WriteFile(fmt.Sprintf("%s.txt", fileName[1]), bs, 0664); err != nil {
				s += "Error writing file"
				c <- s
			}
		}

		s += fmt.Sprintf("done with %s", url)
		c <- s
	}
}

func main() {
	sites := []string{"https://google.com", "https://facebook.com", "https://stackoverflow.com"}

	// var wg sync.WaitGroup

	c := make(chan string)

	for _, url := range sites {
		// wg.Add(1)
		// go func(s string) {
		go checkAndSaveBody(url, c)

		// wg.Done()
		// }(url)
	}

	fmt.Println("Num of grs", runtime.NumGoroutine())

	for i := 0; i < len(sites); i++ {
		fmt.Println(<-c)
	}

	// wg.Wait()
}
