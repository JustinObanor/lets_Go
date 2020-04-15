package main

import (
	"fmt"
	"sort"
	"strings"
)

type WordCount struct {
	word  rune
	count int
}

func main() {
	s := "hello mr big cucumber"
	ss := calculate(s)

	res := sorter(ss)
	for i := 0; i < len(res); i++ {
		fmt.Printf("%c : %d\n", res[i].word, res[i].count)
	}
}

func calculate(s string) map[rune]int {
	m := make(map[rune]int)
	res := strings.ReplaceAll(s, " ", "")
	for _, letter := range res {
		if _, ok := m[letter]; ok {
			m[letter]++
		} else {
			m[letter] = 1
		}
	}
	return m
}

func sorter(m map[rune]int) []WordCount {
	wordCounts := make([]WordCount, 0, len(m))
	for key, val := range m {
		wordCounts = append(wordCounts, WordCount{word: key, count: val})
	}

	sort.Slice(wordCounts, func(i, j int) bool {
		return wordCounts[i].count > wordCounts[j].count
	})
	return wordCounts
}
