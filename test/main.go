package main

import (
	"fmt"
	"sort"
)

func main() {
	s := []string{"f", "d", "b", "a", "b"}
	sort.Sort(sort.StringSlice(s))
	fmt.Println(s)
}