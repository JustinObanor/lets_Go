package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
)

const (
	urlPrefix = "https://tools.ietf.org/rfc/rfc"
	urlSuffix = ".txt"
	lenPrefix = len(urlPrefix)
	lenSuffix = len(urlSuffix)
	nLow      = 1
	nHigh     = 1000
	workers   = 10
)

var jwg sync.WaitGroup
var awg sync.WaitGroup
var totalWords = make(map[string]int)

func countWords(text string) map[string]int {
	wordCounts := make(map[string]int, nHigh)

	texts := strings.FieldsFunc(text, func(c rune) bool {
		return !unicode.IsLetter(c) || unicode.IsNumber(c)
	})

	for _, word := range texts {
		if len(word) > 12 {
			wordCounts[word]++
		}
	}
	return wordCounts
}

func accumulateWords(wordCounts map[string]int) {
	for key, value := range wordCounts {
		totalWords[key] += value
	}
}

func scraper(url string) (map[string]int, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return countWords(string(b)), nil
}

func main() {
	jobs := make(chan string, nHigh)
	results := make(chan map[string]int, nHigh)

	awg.Add(1)
	go func() {
		for r := range results {
			accumulateWords(r)
			for k, v := range totalWords {
				fmt.Printf("%s : %d\n", k, v)
			}
		}
		awg.Done()
		close(results)
	}()

	jwg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			for j := range jobs {
				data, err := scraper(j)
				if err != nil {
					fmt.Println(err)
				}
				results <- data
			}
			jwg.Done()
		}()
	}

	go func() {
		var b strings.Builder

		for i := nLow; i <= nHigh; i++ {
			b.Grow(lenPrefix + 3 + lenSuffix)

			b.WriteString(urlPrefix)
			b.WriteString(strconv.Itoa(i))
			b.WriteString(urlSuffix)

			jobs <- b.String()

			b.Reset()
		}
		close(jobs)
	}()

	jwg.Wait()
	close(results)
awg/W
	
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
