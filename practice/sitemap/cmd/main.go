package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/lets_Go/practice/sitemap/parser"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type location struct {
	Value string `xml:"loc"`
}

type urlset struct {
	XMLName xml.Name   `xml:"urlset"`
	XMLns   string     `xml:"xmlns,attr"`
	URLs    []location `xml:"url"`
}

func main() {
	urlFlag := flag.String("url", "https://gophercises.com", "the site we want to build a site map for")
	maxBFSDepth := flag.Int("depth", 10, "the maximum depth we want to go while visiting links")
	xmlFile := flag.String("file", "urlset.xml", "file where we want to store the parsed url")
	flag.Parse()

	pages := bfs(*urlFlag, *maxBFSDepth)
	toXML := urlset{
		XMLns: xmlns,
		URLs:  make([]location, len(pages)),
	}

	for i, link := range pages {
		toXML.URLs[i] = location{
			Value: link,
		}
	}

	if err := marshallXML(*xmlFile, toXML); err != nil {
		log.Println(err)
	}
}

func getLinks(site string) []string {
	resp, err := http.Get(site)
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()

	links, err := parser.ParsePage(resp.Body)
	if err != nil {
		return []string{}
	}

	reqURL := resp.Request.URL
	baseURL := &url.URL{
		Scheme: reqURL.Scheme,
		Host:   reqURL.Host,
	}

	base := baseURL.String()

	return filter(buildLink(base, links), withPrefix(base))
}

func bfs(base string, maxDepth int) []string {
	seen := make(map[string]struct{})
	q := map[string]struct{}{}
	nq := map[string]struct{}{
		base: struct{}{},
	}

	for i := 0; i <= maxDepth; i++ {
		fmt.Println("depth = ", i)
		q, nq = nq, make(map[string]struct{})
		if len(q) == 0 {
			break
		}

		for link := range q {
			if _, ok := seen[link]; ok {
				continue
			}

			seen[link] = struct{}{}

			for _, link := range getLinks(link) {
				if _, ok := seen[link]; !ok {
					nq[link] = struct{}{}
				}
			}
		}
	}

	ret := make([]string, 0, len(seen))
	for link := range seen {
		ret = append(ret, link)
	}

	return ret
}

func marshallXML(file string, data urlset) error {
	output, err := xml.MarshalIndent(data, "  ", "    ")
	if err != nil {
		return err
	}

	output = []byte(xml.Header + string(output))

	return ioutil.WriteFile(file, output, 0644)
}

func filter(links []string, KeepFn func(string) bool) []string {
	var ret []string
	for _, link := range links {
		if KeepFn(link) {
			ret = append(ret, link)
		}
	}
	return ret
}

func withPrefix(pfx string) func(string) bool {
	return func(url string) bool {
		return strings.HasPrefix(url, pfx)
	}
}

func buildLink(base string, links []string) []string {
	var ret []string
	for _, link := range links {
		switch {
		case strings.HasPrefix(link, "/"):
			ret = append(ret, base+link)
		case strings.HasPrefix(link, "http"):
			ret = append(ret, link)
		}
	}

	return ret
}
