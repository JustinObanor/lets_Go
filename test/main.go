package main

import (
	"fmt"
	"strings"
)

func main() {
	keyLower := "abcdefghijklmnopqrstuvwxyz"
	keyUpper := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	s := "middle-Outz"
	delta := 2
	var res string

	for _, r := range s {
		switch {
		case strings.IndexRune(keyLower, r) > 0:
			res += string(rotate(r, delta, []rune(keyLower)))
		case strings.IndexRune(keyUpper, r) > 0:
			res += string(rotate(r, delta, []rune(keyUpper)))
		default:
			res += string(r)
		}
	}

	fmt.Println(res)
}

func rotate(r rune, delta int, key []rune) rune {
	keyStr := string(key)
	idx := strings.IndexRune(keyStr, r)
	if idx < 0 {
		panic("char not found in key")
	}

	idx = (idx + delta) % len(key)

	kRune := []rune(key)

	return kRune[idx]
}
