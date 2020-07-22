package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"sync"
	"syscall"
	"time"

	"golang.org/x/net/netutil"
)

const workers = 5
const connectionLimit = 100

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

func (c Client) getContent(_ context.Context, url string) ([]byte, error) {
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
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	client := newClient()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var urls URLs
		var jwg sync.WaitGroup
		var rwg sync.WaitGroup
		//Limitation#5
		// var ctxReq = r.Context
		jobs := make(chan Job, 20)
		result := make(chan Response, 20)

		if err := json.NewDecoder(r.Body).Decode(&urls); err != nil {
			return
		}

		if len(urls.List) > 20 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("amount of url's exceeded limit"))
			return
		}

		//reciving result, sorting, and then sending to client
		rwg.Add(1)
		go func() {
			responses := make([]string, 0, len(urls.List))
			for r := range result {
				responses = append(responses, string(r.content))
			}

			sort.Sort(sort.StringSlice(responses))

			for i, content := range responses {
				w.Write([]byte(content))
				fmt.Printf("wrote content of %s\n", urls.List[i])
			}
			rwg.Done()
		}()

		//pulling from channel and working, and then sending to be sorted
		jwg.Add(workers)
		for i := 0; i < workers; i++ {
			ctx, cancel := context.WithCancel(context.Background())

			go func() {
				for url := range jobs {
					content, err := client.getContent(ctx, url.site)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte(err.Error()))
						cancel()

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

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	l = netutil.LimitListener(l, connectionLimit)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	go func() {
		if err := srv.Serve(l); err != nil {
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
