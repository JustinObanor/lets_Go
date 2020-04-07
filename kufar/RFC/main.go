package main

import (
	"container/heap"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"unicode"
)

type elem struct {
	word  string
	count int
}

type elemHeap []elem

var totalWords = make(map[string]int)

func countWords(text string) map[string]int {
	var wordBegPos, runeCount int
	wordCounts := make(map[string]int)
	for i, c := range text {
		if unicode.IsLetter(c) {
			if runeCount == 0 {
				wordBegPos = i
			}
			runeCount++
			continue
		}

		if runeCount > 4 {
			word := text[wordBegPos:i]
			count := wordCounts[word] // return 0 if word is not in wordCounts
			count++
			wordCounts[word] = count
		}
		runeCount = 0
	}
	return wordCounts
}

//for accumulating maps from different RFC sites
func accumulateWords(wordCounts map[string]int) {
	for key, value := range wordCounts {
		totalWords[key] = totalWords[key] + value
	}
}

func (h elemHeap) Len() int           { return len(h) }
func (h elemHeap) Less(i, j int) bool { return h[i].count < h[j].count }
func (h elemHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h elemHeap) Push(x interface{}) { /* not used */ }
func (h elemHeap) Pop() interface{}   { /* not used */ return nil }

func frequentWords(m map[string]int, nbrWords int) []elem {
	h := elemHeap(make([]elem, nbrWords))
	for word, count := range m {
		if count > h[0].count {
			h[0] = elem{word: word, count: count}
			heap.Fix(h, 0)
		}
	}
	sort.Slice(h, func(i, j int) bool { return h[i].count > h[j].count })
	return h
}

func write(filename string, data []elem) error {
	f, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("couldnt open file %s", outfile)
	}

	w := csv.NewWriter(f)
	w.Flush()

	s := make([]string, len(data))

	for _, v := range data 
	//not complete

	return w.Write(data)
}

func scraper(url string) ([]elem, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	accumulateWords(countWords(string(b)))

	return frequentWords(totalWords, 20), nil
}

func main() {
	elem, err := scraper("https://tools.ietf.org/rfc/rfc1.txt")
	if err != nil {
		fmt.Println(err)
	}
	if err := write("rfc.txt", elem); err != nil {
		fmt.Println(err)
	}
}
