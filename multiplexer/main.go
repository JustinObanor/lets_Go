package main

import (
	"context"
	"encoding/json"
	"errors"
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

type URL struct {
	List []string `json:"list"`
}

type Job struct {
	index int
	site  string
}

type Response struct {
	index   int
	content []byte
}

type Client struct {
	client *http.Client
}

type Worker struct {
	client  Client
	errChan chan error
	resChan chan Response
}

func newWorker() *Worker {
	return &Worker{
		client: Client{
			client: &http.Client{
				Timeout: time.Second,
			},
		},
		errChan: make(chan error, 20),
		resChan: make(chan Response, 20),
	}
}

func decodeRequest(r *http.Request) ([]string, error) {
	var url URL
	if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
		return nil, err
	}

	if len(url.List) > 20 {
		return nil, errors.New("amount of url's exceeded limit")
	}
	return url.List, nil
}

func (c Client) getContent(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
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

func (w Worker) worker(ctx context.Context, wg *sync.WaitGroup, jobs chan Job) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("request cancelled")
			return
		case url, ok := <-jobs:
			if !ok {
				return
			}
			content, err := w.client.getContent(ctx, url.site)
			if err != nil {
				w.errChan <- err
				return
			}
			w.resChan <- Response{index: url.index, content: content}
		}
	}
}

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		work := newWorker()

		jobs := make(chan Job, 20)

		var jwg sync.WaitGroup
		var swg sync.WaitGroup

		urls, err := decodeRequest(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		responses := make([]Response,0, len(urls))

		ctx, cancel := context.WithCancel(r.Context())
		var counter int

		swg.Add(1)
		go func() {
			defer swg.Done()
			for {
				if counter >= len(urls) {
					break
				}
				select {
				case err := <-work.errChan:
					w.Write([]byte(err.Error()))
					cancel()
					return
				case res, ok := <-work.resChan:
					if !ok {
						return
					}
					responses = append(responses, res)
				}
				counter++
			}
		}()

		jwg.Add(workers)
		for i := 0; i < workers; i++ {
			go work.worker(ctx, &jwg, jobs)
		}

		go func() {
			for i, url := range urls {
				jobs <- Job{index: i, site: url}
			}
			close(jobs)
		}()

		jwg.Wait()
		close(work.resChan)
		close(work.errChan)
		swg.Wait()

		sort.Slice(responses, func(i,j int) bool{
			return responses[i].index < responses[j].index
		})

		for i, cont := range responses {
			fmt.Printf("writing content of %s\n", urls[i])
			w.Write([]byte(cont.content))
		}
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
