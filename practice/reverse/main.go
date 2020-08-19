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

func reverse(s string) string {
	ss := []byte(s)
	for i := 0; i < len(s)/2; i++ {
		ss[i], ss[len(s)-i-1] = ss[len(s)-i-1], ss[i]
	}
	return string(ss)
}

func reverse2(s string) {
	for i := len(s) - 1; i >= 0; i-- {
		fmt.Print(s[i : i+1])
	}
}

func reverse3(s string) {
	if s == "" || len(s) <= 1 {
		fmt.Println(s)
	} else {
		fmt.Print(s[len(s)-1 : len(s)-1+1])
		reverse(s[:len(s)-1])
	}
}
