package main

import "fmt"

func main() {
	s := "abc"
	fmt.Println(cipher(s))
}

func cipher(text string) string {
	runes := []rune(text)
	for index, value := range runes {
		value = value + rune(1)
		runes[index] = value
	}
	return string(runes)
}
