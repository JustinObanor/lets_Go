package main

import (
	"bufio"
	"encoding/json"
	"fmt"
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
		colly.DisallowedDomains("www.facebook.com", "www.vk.com", "www.instagram.com", "www.twitter.com"),
	)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		if e.Attr("class") == "Button_1qxkboh-o_O-primary_cv02ee-o_O-md_28awn8-o_O-primaryLink_109aggg" {
			return
		}

		link := e.Attr("href")
		link = strings.ToLower(link)

		if len(link) > 71 {
			return
		}

		if strings.Index(link, "=signup") > -1 || strings.Index(link, "=login") > -1 ||
			strings.HasSuffix(link, "jpg") || strings.Contains(link, "video") ||
			strings.HasSuffix(link, "djvu") || strings.HasSuffix(link, "pps") ||
			strings.HasSuffix(link, "pptx") || strings.Contains(link, "blog") ||
			strings.Contains(link, "fotogalereya") || strings.Contains(link, "webinar") ||
			strings.Contains(link, "moodle") || strings.Contains(link, "impuls") ||
			strings.HasSuffix(link, "rar") || strings.HasSuffix(link, "doc") ||
			strings.HasSuffix(link, "pdf") || strings.Contains(link, "phones") ||
			strings.HasSuffix(link, "docx") || strings.HasSuffix(link, "zip") ||
			strings.HasSuffix(link, "ppt") || strings.Contains(link, "news") ||
			strings.Contains(link, "libeldoc") || strings.Contains(link, "doklady") ||
			strings.Contains(link, "e-lib") || strings.Contains(link, "elearning") ||
			strings.Contains(link, "elib") {
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

		if strings.Contains(e.Request.AbsoluteURL(link), url) && !strings.HasPrefix(e.Request.AbsoluteURL(link), "mailto") {
			c.Visit(e.Request.AbsoluteURL(link))
		}

	})

	c.OnRequest(func(r *colly.Request) {
		if strings.Contains(r.URL.String(), url) && !strings.HasPrefix(r.URL.String(), "mailto") {
			fmt.Println("visiting", r.URL.String())
		}
	})

	var b strings.Builder
	b.WriteString("https://www.")
	b.WriteString(url)
	b.WriteString("/")

	if err := c.Visit(b.String()); err != nil {
		return err
	}

	b.Reset()

	return nil
}

func main() {
	url := ""
	fmt.Print("Site to visit: ")
	fmt.Scan(&url)

	os.Remove("emails.txt")

	getEmail(url)
}
