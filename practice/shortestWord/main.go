package main

import (
	"fmt"
	"strings"
)

func main() {
	s := "Let see if can find the legnth of a the shortest word"
	fmt.Println(shortestWord(s))
}

func shortestWord(s string) string {
	word, length := "", len(s)
	for _, v := range strings.Split(s, " ") {
		if len(v) < length {
			word, length = v, len(v)
		}
	}
	return word
}
