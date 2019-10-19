package main

import (
	"fmt"
	"sort"
	"strings"
)

type sorted []rune

func (s sorted) Len() int           { return len(s) }
func (s sorted) Less(i, j int) bool { return s[i] < s[j] }
func (s sorted) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func sorter(s string) string {
	r := []rune(s)
	sort.Sort(sorted(r))
	return string(r)
}

func main() {
	s := "did i do a this right"
	fmt.Println(longest(s))
}

func longest(s string) string {
	res := sorter(s)
	words := strings.Split(s, " ")
	for _, word := range words {
		if len(word) < len(res) {
			res = word
		}
	}
	return res
}
