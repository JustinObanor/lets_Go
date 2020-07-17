package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

//URLs stores the list of sites
type URLs struct {
	List []string `json:"list"`
}

//Client is a http clieng
type Client struct {
	client *http.Client
}

//Response is the response sent to the client
type Response struct {
	index   int
	content []byte
}

//Job is a single url with its index
type Job struct {
	index int
	site  string
}

func newClient() *Client {
	return &Client{
		client: &http.Client{
			Timeout: time.Second,
		},
	}
}

func (c Client) getContent(url string) ([]byte, error) {
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func main() {
	var urls URLs
	var jwg sync.WaitGroup
	var rwg sync.WaitGroup
	jobs := make(chan Job)
	result := make(chan Response)
	workers := 5

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	client := newClient()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&urls); err != nil {
			return
		}

		if len(urls.List) > 20 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("amount of url's exceeded limit"))
			return
		}

		//reciving result, sorting, and then sending to client
		rwg.Add(1)
		go func() {
			responses := make([]string, len(urls.List))
			for r := range result {
				responses[r.index] = string(r.content)
			}

			for i, content := range responses {
				w.Write([]byte(content))
				fmt.Printf("wrote content of %s\n", urls.List[i])
			}
			rwg.Done()
		}()

		//pulling from channel and working, and then sending to be sorted
		jwg.Add(workers)
		for i := 0; i < workers; i++ {
			go func() {
				for url := range jobs {
					content, err := client.getContent(url.site)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte(err.Error()))
						return
					}
					result <- Response{index: url.index, content: content}
				}
				jwg.Done()
			}()
		}

		//sending the url's to the channel to start aggregating
		go func() {
			for i, url := range urls.List {
				jobs <- Job{index: i, site: url}
			}
			close(jobs)
		}()

		jwg.Wait()
		close(result)
		rwg.Wait()
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("listen:%+s\n", err)
		}
	}()

	log.Printf("server started")

	<-done

	log.Printf("server stopped")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("server Shutdown Failed:%+s", err)
	}

	log.Printf("server exited properly")
}
