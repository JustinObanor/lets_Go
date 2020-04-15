package main

import (
	"fmt"
	"unicode"
)

func main() {
	s := "Here is the code to count words with more than 3 letters. It’s more efficient to provide the code than to explain the algorithm. Note that it will consider a number as a word"
	m := countWordsIn(s)
	for k, v := range m {
		fmt.Printf("%d : %s\n", v, k)
	}
}

func countWordsIn(text string) map[string]int {
	//wordBegPos = index of first letter in all words
	var wordBegPos, runeCount int
	wordCounts := make(map[string]int)
	for i, c := range text {
		if unicode.IsLetter(c) {
			if runeCount == 0 {
				wordBegPos = i
			}
			runeCount++
			fmt.Println(wordBegPos)
			continue
		}
		if runeCount > 3 {
			word := text[wordBegPos:i]
			count := wordCounts[word]
			count++
			wordCounts[word] = count
		}
		runeCount = 0
	}
	return wordCounts
}
