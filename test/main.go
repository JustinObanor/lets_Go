package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return
	}

	fmt.Println("Hello, World.")
	fmt.Println(string(input))
}
