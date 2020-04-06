package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"unicode"
)

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

func scraper(url string) (map[string]int, error) {
	m := make(map[string]int)
	resp, err := http.Get(url)
	if err != nil {
		return m, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return m, err
	}

	accumulateWords(countWords(string(b)))
	return totalWords, nil
}

func main() {
	m, err := scraper("https://tools.ietf.org/rfc/rfc1.txt")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(m)
}
