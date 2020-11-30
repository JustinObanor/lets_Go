package main

import (
	"fmt"
	"net/url"
)

func main() {
	s := ""

	u , _:= url.Parse(s)

	fmt.Println(u.Path[1:])
	fmt.Println("test")
}