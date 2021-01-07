package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

var (
	numberOfStories = flag.Int("num", 30, "number of stories to retrieve")

	apiURLPre = "https://hacker-news.firebaseio.com/v0/"
	apiURLSuf = ".json?print=pretty"
)

type Item struct {
	client *http.Client

	Title    string `json:"title"`
	URL      string `json:"url"`
	ItemType string `json:"type"`
}

func main() {
	flag.Parse()

	item := newItem()

	ids, err := item.getIDs()
	if err != nil {
		panic(err)
	}

	now := time.Now()

	items, err := item.getItems(ids)
	if err != nil {
		panic(err)
	}

	elpased := time.Since(now)
	for _, itm := range items {
		fmt.Printf("%v - %v\n", itm.Title, item.URL)
	}

	log.Printf("Took %v", elpased)
}

func (item *Item) getItems(ids []int) ([]Item, error) {
	var count int
	items := make([]Item, 0, *numberOfStories)

	for _, id := range ids {
		if count >= *numberOfStories {
			break
		}

		subPath := fmt.Sprintf("item/%d", id)
		path := fmt.Sprint(apiURLPre + subPath + apiURLSuf)

		resp, err := item.client.Get(path)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if err = json.NewDecoder(resp.Body).Decode(&item); err != nil {
			return nil, err
		}

		if item.ItemType != "story" || item.URL == "" || item.Title == "" {
			continue
		}

		items = append(items, *item)
		count++
	}
	return items, nil
}

func (item *Item) getIDs() ([]int, error) {
	path := fmt.Sprint(apiURLPre + "topstories" + apiURLSuf)

	resp, err := item.client.Get(path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ids []int
	if err = json.NewDecoder(resp.Body).Decode(&ids); err != nil {
		return nil, err
	}

	return ids, nil
}

func newItem() *Item {

	return &Item{
		client: &http.Client{
			Timeout: (time.Second * 15),
		},
	}
}
