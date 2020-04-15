package main

import "fmt"

func main() {
	s := "cat"
	fmt.Println(reverser(s))
}

func reverser(s string) string {
	var ss []byte
	for i := len(s) - 1; i >= 0; i-- {
		ss = append(ss, s[i])
	}
	return string(ss)
}
