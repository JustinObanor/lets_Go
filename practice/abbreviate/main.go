package main

import (
	"fmt"
	"strings"
	"unicode"
)

func main() {
fmt.Println(Abbreviate("World Health Organization"))
}

func Abbreviate(s string) (abv string) {
	words := strings.FieldsFunc(strings.Title(strings.ReplaceAll(s, "'", "")), func(r rune) bool {
		return !unicode.IsLetter(r)
	})
	for _, word := range words {
		abv += word[0:1]
	}
	return
}
