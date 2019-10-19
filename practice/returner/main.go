package main

import (
	"fmt"
	"sort"
)

type myString []rune

func main() {
	bigArr := []int{1, 23, 4, 6, 3, 2, 6, 8, 9, 9, 5, 43, 7, 8, 5, 3, 2, 46, 0}
	// sort.Ints(bigArr)
	// fmt.Println(bigArr)
	sort.Sort(sort.Reverse(sort.IntSlice(bigArr)))
	fmt.Println(bigArr)
	fmt.Println(bigArr[:3])

	s := "this is my is is my wow so is is my my so hey no yes"
	ss := sortString(s)
	fmt.Println(ss)

}
func (s myString) Less(i, j int) bool { return s[i] < s[j] }
func (s myString) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s myString) Len() int           { return len(s) }

func sortString(s string) string {
	r := []rune(s)
	sort.Sort(myString(r))
	return string(r)
}
