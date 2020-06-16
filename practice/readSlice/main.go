package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	input := ""
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		input = scanner.Text()
	}

	fmt.Println("Hello, World.")
	fmt.Println(input)
}
