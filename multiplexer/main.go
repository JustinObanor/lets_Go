package main

import (
	"context"
	"encoding/json"
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

func getContent(client *http.Client, url string) ([]byte, error) {
	resp, err := client.Get(url)
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
	var wg sync.WaitGroup
	var mu sync.RWMutex
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	client := http.Client{
		Timeout: time.Second,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&urls); err != nil {
			return
		}

		if len(urls.List) > 20 {
			w.Write([]byte("amount of url's exceeded limit"))
			return
		}

		for _, url := range urls.List {
			wg.Add(1)
			go func(u string) {
				mu.RLock()
				content, err := getContent(&client, u)
				mu.RUnlock()
				
				if err != nil {
					w.Write([]byte(err.Error()))
				}

				mu.Lock()
				w.Write(content)
				mu.Unlock()

				wg.Done()
			}(url)
		}
		wg.Wait()
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen:%+s\n", err)
		}
	}()

	log.Printf("server started")

	<-done

	log.Printf("server stopped")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server Shutdown Failed:%+s", err)
	}

	log.Printf("server exited properly")

}
