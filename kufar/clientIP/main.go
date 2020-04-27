package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

const (
	site1 = "https://ipstack.com/"
	site2 = "https://geoip.nekudo.com"
)

func get(s string) string {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(s)
	if err != nil {
		return err.Error()
	}

	defer resp.Body.Close()

	doc, _ := goquery.NewDocumentFromReader(resp.Body)

	return doc.Find("div.ipchecker").Text()
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

	fmt.Println(get(site1))

	get2(site1)
}
