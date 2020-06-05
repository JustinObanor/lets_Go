package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"
	"sync"
)

var sites = []string{"https://google.com", "https://facebook.com", "https://stackoverflow.com", "https://golang.org", "https://amazon.com"}

func checkAndSaveBody(url string, wg *sync.WaitGroup) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("%s is down: %v", url, err)
	}
	defer resp.Body.Close()

	fmt.Printf("%s -> STATUS CODE: %d\n", url, resp.StatusCode)

	if resp.StatusCode == 200 {
		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		fileName := strings.Split(url, "//")
		fmt.Printf("Writing response body to %s.txt\n", strings.Trim(fileName[1], "[]"))

		if err := ioutil.WriteFile(fmt.Sprintf("%s.txt", fileName[1]), bs, 0664); err != nil {
			return "", err
		}
	}

	wg.Done()
	return "", err
}

func main() {
	var wg sync.WaitGroup

	wg.Add(len(sites))

	for _, url := range sites {
		go func(s string) {
			resp, err := checkAndSaveBody(s, &wg)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(resp)
		}(url)
	}

	fmt.Println("Num of grs", runtime.NumGoroutine())

	wg.Wait()
}
