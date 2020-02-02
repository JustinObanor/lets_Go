package main

import (
	"fmt"
	"strings"
)

func main() {
	x := []string{"emma", "james", "paul"}
	for i := 0; i < len(x); i++ {
		if strings.Compare(x[i], "emma") == 0 {
			fmt.Println("found")
			return
		} 
	}
	fmt.Println("not found")
}
