package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
)

var totalWords = make(map[string]int)

type info struct {
	Email string
}

func accumulateCount(counter map[string]int) {
	for key, value := range counter {
		totalWords[key] += value
	}
}

func writeToFile(name string, data []byte) error {
	file, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	b := bufio.NewWriter(file)

	if _, err := b.Write(data); err != nil {
		return err
	}

	if err := b.Flush(); err != nil {
		return err
	}

	return nil
}

func getEmail(url string) error {
	c := colly.NewCollector(
		colly.AllowedDomains("www.bsuir.by", "www.bsuir.by/ru"),
	)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		link = strings.ToLower(link)

		if strings.HasSuffix(link, "jpg") || strings.HasSuffix(link, "html") || strings.HasSuffix(link, "djvu") || strings.HasSuffix(link, "pps") || strings.HasSuffix(link, "pptx") || strings.Contains(link, "fotogalereya") || strings.Contains(link, "impuls") || strings.HasSuffix(link, "rar") || strings.HasSuffix(link, "doc") || strings.HasSuffix(link, "pdf") || strings.Contains(link, "phones") || strings.HasSuffix(link, "docx") || strings.HasSuffix(link, "zip") || strings.HasSuffix(link, "ppt") || strings.Contains(link, "news") {
			return
		}

		var email info

		counter := make(map[string]int)
		if strings.HasPrefix(link, "mailto") {
			email.Email = strings.TrimPrefix(link, "mailto:")

			counter[email.Email]++
			accumulateCount(counter)

			if totalWords[email.Email] < 2 {
				bs, err := json.Marshal(email)
				if err != nil {
					return
				}
				writeToFile("emails.txt", bs)
			}
		}

		e.Request.Visit(link)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("visiting", r.URL.String())
	})

	if err := c.Visit(url); err != nil {
		return err
	}
	return nil
}

func main() {
	url := "https://www.bsuir.by/"

	if err := os.Remove("emails.txt"); err != nil {
		log.Println(err)
	}

	getEmail(url)
}
