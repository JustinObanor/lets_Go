package main

import "fmt"

func main() {
	draw(4)
}

func draw(height int) {
	for i := 0; i < height; i++ {
		for j := 0; j < i; j++ {
			fmt.Print("#")
		}
		fmt.Printf("\n")
	}
}
