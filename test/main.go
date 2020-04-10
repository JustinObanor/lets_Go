package main

import "fmt"

func main() {
	commits := map[int]int{
		1: 3711,
		2: 2138,
		3: 1908,
		4: 912,
	}
	for k, v := range commits {
		fmt.Println(k, v)
	}
}
