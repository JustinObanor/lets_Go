package main

import "fmt"

func main() {
	s := []byte("emma")
	fmt.Printf("%p\n", s)
	for i := range s {
		fmt.Printf("%p\n", &s[i])
	}
}
