package main

import "fmt"

func main() {
	value1 := "cat"
	reversed1 := reverse(value1)
	fmt.Println(value1)
	fmt.Println(reversed1)
}

func reverse(s string) string {
	data := []rune(s)
	res := make([]rune, len(s))
	for i := len(s) - 1; i >= 0; i-- {
		res = append(res, data[i])
	}
	return string(res)
}
