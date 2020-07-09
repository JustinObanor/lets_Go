package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//URLs stores the list of sites
type URLs struct {
	List []string `json:"list"`
}

func main() {
	var urls URLs

	client := http.Client{
		Timeout: time.Second,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&urls); err != nil {
			return
		}

		for _, payLoad := range urls.List {
			resp, err := client.Get(payLoad)
			if err != nil {
				w.Write([]byte(err.Error()))
			}
			defer resp.Body.Close()

			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				w.Write([]byte(err.Error()))
			}

			w.Write(b)
		}
	})

	log.Println(http.ListenAndServe(":8080", nil))
}
