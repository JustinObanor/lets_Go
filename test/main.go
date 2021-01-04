package main

import "fmt"

type blah struct {
	a, b string
	c, d int
}

var pss *blah
var pss2 *blah = nil

func main() {
	fmt.Println(pss == nil)
	fmt.Println(pss2 == nil)
}
