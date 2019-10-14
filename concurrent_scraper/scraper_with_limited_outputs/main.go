package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var start time.Time

func init() {
	start = time.Now()
}

type site struct {
	url string
	id  int
}

func main() {
	urlBase := flag.String("url", "https://tools.ietf.org/rfc/rfc%d.txt", "The URL you wish to scrape, containing \"%d\" where the id should be substituted")
	idLow := flag.Int("from", 1, "The first ID that should be searched in the URL")
	idHigh := flag.Int("to", 1000, "The last ID that should be searched in the URL")
	concurrency := flag.Int("concurrency", 1000, "How many scrapers to run in parallel. (More scrapers are faster, but more prone to rate limiting or bandwith issues)")
	outfile := flag.String("output", "output.csv", "Filename to export the CSV results")
	body := flag.String("body > pre", "body", "JQuery-style query for the body element")
	flag.Parse()

	queries := []string{*body}
	headers := []string{"BODY"}
	headers = append([]string{"URL HEADERS"}, headers...)

	task := make(chan site)

	results := make(chan map[string]int)

	go func() {
		for i := *idLow; i < *idHigh; i++ {
			url := fmt.Sprintf(*urlBase, i)
			task <- site{url: url, id: i}
		}
		close(task)
	}()

	var wg sync.WaitGroup

	wg.Add(*concurrency)

	go func() {
		wg.Wait()
		close(results)
	}()

	for i := 0; i < *concurrency; i++ {
		go func() {
			for v := range task {
				site, err := fetch(v.url, v.id, queries)
				if err != nil {
					fmt.Printf("error fetching data from site %v", err)
				}
				wordCount, err := format(site)
				if err != nil {
					fmt.Printf("error fetching data from site %v", err)
				}
				results <- wordCount
			}
			wg.Done()
		}()
	}

	err := writeSites(results, *outfile, headers)
	if err != nil {
		log.Printf("could not dump CSV: %v", err)
	}
	wg.Wait()
	log.Printf("Time to get %d number of RFC files - %v", *concurrency, time.Since(start))
}

func fetch(url string, id int, queries []string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("couldnt get site %s : %v", url, err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("couldnt parse site %s : %v", url, err)
	}

	s := []string{url, strconv.Itoa(id)}

	for _, v := range queries {
		s = append(s, strings.TrimSpace(doc.Find(v).Text()))
	}
	return s, nil
}

func format(queries []string) (map[string]int, error) {
	words := strings.Split(strings.Join(queries, ""), " ")
	m := make(map[string]int)
	for _, word := range words {
		_, ok := m[word]
		if ok {
			m[word]++
		} else {
			m[word] = 1
		}
	}
	return m, nil
}

func writeSites(results chan map[string]int, outfile string, headers []string) error {
	file, err := os.Create(outfile)
	if err != nil {
		return fmt.Errorf("error creating csv file %s : %v", outfile, err)
	}
	defer file.Close()

	w := csv.NewWriter(file)

	defer w.Flush()

	if err := w.Write(headers); err != nil {
		return fmt.Errorf("error writing records to csv %v", err)
	}
	for v := range results {
		for key, value := range v {
			r := make([]string, len(strconv.Itoa(value)))
			r = append(r, key)
			r = append(r, strconv.Itoa(value))
			if err := w.Write(r); err != nil {
				return fmt.Errorf("error writing record to csv %v", err)
			}
		}
	}
	if err != w.Error() {
		return fmt.Errorf("couldnt write to csv file %s : %v", outfile, err)
	}
	return nil
}
