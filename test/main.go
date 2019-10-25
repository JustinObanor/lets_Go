package main

import "fmt"

func main() {
	var word string
	fmt.Println("Input word to cypher")
	fmt.Scan(&word)

	var key rune
	fmt.Print("key: ")
	fmt.Scan(&key)

	var action int
	fmt.Println("Input action: [1]Enconde [2]Decode")
	fmt.Scan(&action)

	switch action {
	case 1:
		for _, v := range word {
			v = v + key
			fmt.Printf("%c", v)
		}

	case 2:
		for _, v := range word {
			v = v - key
			fmt.Printf("%c", v)
		}
	default:
		fmt.Println("Wrong option")
	}
}
