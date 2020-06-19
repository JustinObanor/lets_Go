package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	resp, err := http.Get("https://www.bsuir.by/ru/rektorat/davydov-m-v")
	if err != nil {
		log.Println(err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println(err)
	}

	fmt.Printf(strings.TrimSpace(doc.Find("a[href]").Text()))
}
