package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"
	"sync"
)

func checkAndSaveBody(url string) (string, error) {
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

	fmt.Printf("done with %s", url)
	return "", err
}

func main() {
	sites := []string{"https://ngoogle.com", "https://facebook.com", "https://stackoverflow.com"}

	var wg sync.WaitGroup

	for _, url := range sites {
		wg.Add(1)
		go func(s string) {
			resp, err := checkAndSaveBody(s)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(resp)
			wg.Done()
		}(url)
	}

	fmt.Println("Num of grs", runtime.NumGoroutine())

	wg.Wait()
}
