package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"sort"
	"sync"
	"time"

	"github.com/lets_Go/HackerNewsApi/client"
)

var (
	numOfStories = flag.Int("n", 30, "number of hn stories to print")
	port         = flag.Int("p", 8000, "port to listen on")

	tpl = template.Must(template.ParseFiles("./index.gohtml"))
)

type item struct {
	client.Item
	Host string
}

type tplData struct {
	Items []item
	Time  time.Duration
}

func main() {
	flag.Parse()
	mux := http.NewServeMux()
	cs := newCacheStories()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()

		if err := cs.stories(); err != nil {
			http.Error(w, "error getting top stories", http.StatusInternalServerError)
			return
		}

		var data tplData

		cs.RLock()
		data.Items = cs.cache
		cs.RUnlock()

		data.Time = time.Since(now)

		if err := tpl.Execute(w, data); err != nil {
			http.Error(w, "error executing template", http.StatusInternalServerError)
			return
		}
	})

	go func() {
		ticker := time.NewTicker(time.Minute * 1)
		for {
			stories, err := topStories()
			if err != nil {
				return
			}

			cs.Lock()
			cs.cache = stories
			cs.expTime = time.Now().Add(time.Minute * cs.duration)
			cs.Unlock()

			<-ticker.C
		}
	}()

	srv := &http.Server{
		Handler: mux,
		Addr:    fmt.Sprintf(":%d", *port),
	}

	log.Printf("listening on port :%d\n", *port)
	log.Fatal(srv.ListenAndServe())
}

type cacheStories struct {
	sync.RWMutex
	cache    []item
	expTime  time.Time
	duration time.Duration
}

func newCacheStories() *cacheStories {
	return &cacheStories{
		cache:    make([]item, 0, *numOfStories),
		duration: time.Minute * 3,
	}
}

func (cs *cacheStories) stories() error {
	cs.Lock()
	defer cs.Unlock()

	if time.Since(cs.expTime) < 0 {
		return nil
	}

	stories, err := topStories()
	if err != nil {
		return err
	}

	cs.cache = stories
	cs.expTime = time.Now().Add(time.Minute * cs.duration)

	return nil
}

func topStories() ([]item, error) {
	var c client.Client
	ids, err := c.TopStories()
	if err != nil {
		return nil, err
	}

	items := make([]item, 0, *numOfStories)

	at := 0
	for len(items) < *numOfStories {
		need := (*numOfStories - len(items))
		items = append(items, getItems(ids[at:at+need])...)
		at += need
	}

	return items, nil
}

type result struct {
	id   int
	item item
	err  error
}

func getItems(ids []int) []item {
	resChan := make(chan result)
	results := make([]result, 0, *numOfStories)

	rwg := &sync.WaitGroup{}
	cwg := &sync.WaitGroup{}

	cwg.Add(1)
	go func() {
		defer cwg.Done()

		for res := range resChan {
			if res.err != nil {
				continue
			}

			if isStory(res.item) {
				results = append(results, res)
			}
		}
	}()

	rwg.Add(len(ids))
	for idx, id := range ids {
		go func(orderID, itemId int) {
			var c client.Client

			defer rwg.Done()

			hnItem, err := c.GetItem(itemId)
			if err != nil {
				resChan <- result{id: orderID, err: err}
			}

			item := parseHNItem(hnItem)

			resChan <- result{id: orderID, item: item}
		}(idx, id)
	}

	rwg.Wait()
	close(resChan)
	cwg.Wait()

	sort.Slice(results, func(i, j int) bool {
		return results[i].id < results[j].id
	})

	items := make([]item, 0, *numOfStories)
	for _, res := range results {
		items = append(items, res.item)
	}

	return items
}

func isStory(it item) bool {
	return it.Type == "story" && it.URL != ""
}

func parseHNItem(hnItem client.Item) item {
	ret := item{Item: hnItem}

	u, err := url.Parse(hnItem.URL)
	if err == nil {
		ret.Host = u.Hostname()
		return ret
	}

	return ret
}
