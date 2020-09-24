package main

import (
	"fmt"
	"strings"
)

func caesar(r rune, shift int) rune {
	// Shift character by specified number of places.
	// ... If beyond range, shift backward or forward.
	s := int(r) + shift
	if s > 'z' {
		return rune(s - 26)
	} else if s < 'a' {
		return rune(s + 26)
	}
	return rune(s)
}

func main() {
	value := "TEST"
	fmt.Println(value)

	value2 := strings.Map(func(r rune) rune {
		return caesar(r, 5)
	}, value)

	value3 := strings.Map(func(r rune) rune {
		return caesar(r, -5)
	}, value2)
	fmt.Println(value2, value3)

	value4 := strings.Map(func(r rune) rune {
		return caesar(r, 1)
	}, value)

	value5 := strings.Map(func(r rune) rune {
		return caesar(r, -1)
	}, value4)
	fmt.Println(value4, value5)

	value = "exxegoexsrgi"
	result := strings.Map(func(r rune) rune {
		return caesar(r, -4)
	}, value)
	fmt.Println(value, result)
}
