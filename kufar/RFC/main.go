package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
	"unicode"
)

const (
	urlBase = "https://tools.ietf.org/rfc/rfc%d.txt"
	nLow    = 1
	nHigh   = 100
	workers = 10
)

var wg sync.WaitGroup
var totalWords = make(map[string]int)

func countWords(text string) map[string]int {
	wordCounts := make(map[string]int)

	texts := strings.FieldsFunc(text, func(c rune) bool {
		return !unicode.IsLetter(c) || unicode.IsNumber(c)
	})

	for _, word := range texts {
		if len(word) > 4 {
			wordCounts[word]++
		}
	}
	return wordCounts
}

func accumulateWords(wordCounts map[string]int) {
	for key, value := range wordCounts {
		totalWords[key] = totalWords[key] + value
	}
}

func scraper(url string) (map[string]int, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(urlBase)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	accumulateWords(countWords(string(b)))

	return totalWords, nil
}

func main() {
	jobs := make(chan string)

	wg.Add(1)
	go func() {
		var b strings.Builder
		for i := nLow; i <= nHigh; i++ {
			fmt.Fprintf(&b, urlBase, i)
		}
		jobs <- b.String()

		close(jobs)
		wg.Done()
	}()

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			for j := range jobs {
				data, err := scraper(j)
				if err != nil {
					fmt.Println(err)
				}

				for k, v := range data {
					fmt.Printf("%s : %d\n", k, v)
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()
}

// type elem struct {
// 	word  string
// 	count int
// }

// type elemHeap []elem

// func (h elemHeap) Len() int           { return len(h) }
// func (h elemHeap) Less(i, j int) bool { return h[i].count < h[j].count }
// func (h elemHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
// func (h elemHeap) Push(x interface{}) { /* not used */ }
// func (h elemHeap) Pop() interface{}   { /* not used */ return nil }

// func frequentWords(m map[string]int, nbrWords int) []elem {
// 	h := elemHeap(make([]elem, nbrWords))
// 	for word, count := range m {
// 		if count > h[0].count {
// 			h[0] = elem{word: word, count: count}
// 			heap.Fix(h, 0)
// 		}
// 	}
// 	sort.Slice(h, func(i, j int) bool { return h[i].count > h[j].count })
// 	return h
// }

// func write(filename string, data []elem) error {
// 	f, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
// 	if err != nil {
// 		return fmt.Errorf("couldnt open file %s", outfile)
// 	}

// 	w := csv.NewWriter(f)
// 	w.Flush()
// 	s := make([]string, len(data))
// 	for _, v := range data
// 	//not complete
// 	return w.Write(data)
// }
