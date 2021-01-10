package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var apiURLPre = "https://hacker-news.firebaseio.com/v0/"
var apiURLSuf = ".json?print=pretty"

type Client struct {
	apiURLPre string
	apiURLSuf string
}

type Item struct {
	By          string `json:"by"`
	Descendants int    `json:"descendants"`
	ID          int    `json:"id"`
	Kids        []int  `json:"kids"`
	Score       int    `json:"score"`
	Time        int    `json:"time"`
	Title       string `json:"title"`
	Type        string `json:"type"`

	// Only one of these should exist
	Text string `json:"text"`
	URL  string `json:"url"`
}

func (c *Client) instantiate() {
	if c.apiURLPre == "" {
		c.apiURLPre = apiURLPre
	}

	if c.apiURLSuf == "" {
		c.apiURLSuf = apiURLSuf
	}
}

func (c *Client) TopStories() ([]int, error) {
	c.instantiate()

	subpath := "topstories"
	resp, err := http.Get(c.apiURLPre + subpath + c.apiURLSuf)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ids []int
	if err := json.NewDecoder(resp.Body).Decode(&ids); err != nil {
		return nil, err
	}

	return ids, nil
}

func (c *Client) GetItem(id int) (Item, error) {
	c.instantiate()

	subpath := fmt.Sprintf("item/%d", id)
	resp, err := http.Get(c.apiURLPre + subpath + c.apiURLSuf)
	if err != nil {
		return Item{}, err
	}
	defer resp.Body.Close()

	var item Item
	if err := json.NewDecoder(resp.Body).Decode(&item); err != nil {
		return Item{}, err
	}

	return item, nil
}
