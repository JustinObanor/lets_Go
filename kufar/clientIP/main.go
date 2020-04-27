package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
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

func get2(s string) {
	c := colly.NewCollector()

	log.Println(c.Visit(site1))
	c.OnHTML(`div.ipchecker`, func(e *colly.HTMLElement) {
		fmt.Println(e.ChildText(`.ip`))
		fmt.Println(e.ChildText(`div.row.string`))
	})
}

func main() {
	_, b := get(site1)
	fmt.Println(b)

	get2(site1)
}
