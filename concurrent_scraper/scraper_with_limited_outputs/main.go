package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/PuerkitoBio/goquery"
)

type Task struct {
	url string
	id  int
}

type MyMap struct {
	key   string
	value int
}

func main() {
	var wg sync.WaitGroup

	urlBase := flag.String("urlbase", "https://tools.ietf.org/rfc/rfc%d.txt", "the url to search")
	idLow := flag.Int("idlow", 1, "minumum url id to go to")
	idHigh := flag.Int("idhigh", 2, "maximum url id to go to")
	concurrency := flag.Int("concurrency", 1, "amount of parallel tasks to run")
	body := flag.String("body > pre", "body", "jquery thingy")
	outfile := flag.String("outfile", "file.csv", "csv file to parse to")
	flag.Parse()

	query := []string{*body}
	headers := []string{"URL", "HEADERS", "BODY"}

	tasks := make(chan Task)

	results := make(chan []MyMap)

	go func() {
		for i := *idLow; i < *idHigh; i++ {
			url := fmt.Sprintf(*urlBase, i)
			tasks <- Task{url: url, id: i}
		}
	}()

	wg.Add(*concurrency)

	go func() {
		wg.Wait()
		close(results)
	}()

	for i := 0; i < *concurrency; i++ {
		go func() {
			for value := range tasks {
				site, err := scraper(value.url, value.id, query)
				if err != nil {
					fmt.Printf("error parsing data %v", err)
				}
				for i := 0; i < 20; i++ {
					fmt.Printf("%s\t : %d\n", site[i].key, site[i].value)
				}
				results <- site
			}
			wg.Done()
		}()
	}

	err := writeSlice(results, *outfile, headers)
	if err != nil {
		fmt.Printf("couldnt write data to csv file %v", err)
	}

	wg.Wait()
}

func scraper(url string, id int, query []string) ([]MyMap, error) {
	client := &http.Client{
		Timeout: time.Minute,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error making request %v", err)
	}
	req.Header.Set("User-Agent", "Not Firefox")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request %v", err)
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error getting doc %v", err)
	}

	data := []string{url, strconv.Itoa(id)}

	for _, value := range query {
		data = append(data, strings.TrimSpace(doc.Find(value).Text()))
	}

	m := make(map[string]int)

	f := func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	}

	data = strings.FieldsFunc(strings.Join(data, " "), f)

	for _, value := range data {
		if _, ok := m[value]; ok {
			m[value]++
		} else {
			m[value] = 1
		}
	}
	mapSlice := make([]MyMap, 0, len(m))

	for key, value := range m {
		mapSlice = append(mapSlice, MyMap{key, value})
	}
	sort.Slice(mapSlice, func(i, j int) bool { return mapSlice[i].value > mapSlice[j].value })

	return mapSlice, err
}

func writeSlice(results chan []MyMap, outfile string, headers []string) error {
	f, err := os.OpenFile(outfile, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("couldnt open file %s", outfile)
	}
	defer f.Close()

	w := csv.NewWriter(f)

	w.Flush()

	if err = w.Write(headers); err != nil {
		return fmt.Errorf("couldnt write headers to file %s", outfile)
	}
	data := make([]string, 0, len(results))
	for v := range results {
		for i := 0; i < 20; i++ {
			data = append(data, v[i].key)
			data = append(data, strconv.Itoa(v[i].value))
			if err = w.Write(data); err != nil {
				return fmt.Errorf("couldnt write headers to file %s", outfile)
			}
		}
	}
	if err = w.Error(); err != nil {
		return fmt.Errorf("couldnt write headers to file %s", outfile)
	}
	return nil
}
