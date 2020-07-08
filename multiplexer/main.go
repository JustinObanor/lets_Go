package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

//URLs stores the list of sites
type URLs struct {
	List []string `json:"list"`
}

//Client is an individual site that needs to be aggregated
type Client struct {
	url    string
	client http.Client
}

//NewClient constructor
func NewClient(url string) *Client {
	return &Client{
		url:    url,
		client: http.Client{
			// Timeout: 1,
		},
	}
}

//GetContent gets the content of the site
func (c Client) GetContent() ([]byte, error) {
	resp, err := c.client.Get(c.url)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

func main() {
	var urls URLs

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&urls); err != nil {
			return
		}

		for _, payLoad := range urls.List {
			site := NewClient(payLoad)
			body, err := site.GetContent()
			if err != nil {
				w.Write([]byte(err.Error()))
			}
			w.Write(body)
		}
	})

	log.Println(http.ListenAndServe(":8080", nil))
}
