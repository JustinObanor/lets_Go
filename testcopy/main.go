package main

import "fmt"

func main() {
	s := "yoyo mastery"
	reverse(s)
}

func reverse(s string) {
	if s == "" || len(s) <= 1 {
		fmt.Println(s)
	} else {
		fmt.Print(s[len(s)-1 : len(s)-1+1])
		reverse(s[:len(s)-1])
	}
}
