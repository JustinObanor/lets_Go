package main

import (
	"fmt"
	"strings"
)

func main() {
	s := "CivIc"
	a := "Was it a car or a cat I saw"
	fmt.Println(isPalindrone(s))
	fmt.Println(isPalindrone(a))
}

func isPalindrone(s string) bool {
	for i := 0; i < len(s)/2; i++ {
		s = strings.ToLower(s)
		s = strings.ReplaceAll(s, " ", "")
		if s[i] != s[len(s)-i-1] {
			return false
		}
	}
	return true
}
