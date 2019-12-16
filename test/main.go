package main

import "fmt"

func main() {
	var length int

	fmt.Println("Input length :")
	fmt.Scan(&length)

	var width int

	fmt.Println("Input width :")
	fmt.Scan(&width)

	var choice int

	fmt.Println("Input choice :")
	fmt.Scan(&choice)

	calc(choice, length, width)
}

func calc(choice, l, w int) {
	switch choice {
	case 1:
		fmt.Println(l * w)
	case 2:
		fmt.Println(2 * (l + w))
	}
}
