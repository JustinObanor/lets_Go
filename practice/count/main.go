package main

import (
	"fmt"
	"sort"
	"strings"
)

// WordCount holds word and count pair
type WordCount struct {
	word  string
	count int
}

func main() {
	s := []string{"this is my is is my wow so is is my my so hey no yes"}
	words := strings.Split(strings.Join(s, ""), " ")

	// count same words in s
	m := make(map[string]int)
	for _, word := range words {
		if _, ok := m[word]; ok {
			m[word]++
		} else {
			m[word] = 1
		}
	}

	// create and fill slice of word-count pairs for sorting by count
	wordCounts := make([]WordCount, 0, len(m))
	for key, val := range m {
		wordCounts = append(wordCounts, WordCount{word: key, count: val})
	}

	// sort wordCount slice by decreasing count number
	sort.Slice(wordCounts, func(i, j int) bool {
		return wordCounts[i].count > wordCounts[j].count
	})

	// display the three most frequent words
	for i := 0; i < len(wordCounts) && i < 3; i++ {
		fmt.Println(wordCounts[i].word, ":", wordCounts[i].count)
	}
}
