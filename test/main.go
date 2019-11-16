package main

import (
	"fmt"
	"sort"
)

func main() {
	x := []int{1,2,3,4,5,6}
	sort.Sort(sort.Reverse(sort.IntSlice(x)))
	fmt.Println(x)
}
