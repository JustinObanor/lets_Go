package main

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

const (
	site1 = "https://ipstack.com/"
	site2 = "https://geoip.nekudo.com"
)

var count int

func get(s string) (count int, body string) {
	resp, err := http.Get(s)
	if err != nil {
		return 0, err.Error()
	}
	defer resp.Body.Close()

	count++

	doc, _ := goquery.NewDocumentFromReader(resp.Body)

	return count, doc.Find("div.ipchecker").Text()
}

func main() {
	_, b := get(site1)
	fmt.Println(b)
}
