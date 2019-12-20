package main

import "fmt"

func main() {
	err := hey(1, 2, "yo")
	fmt.Println(err)
}

func hey(v ...interface{}) error {

}
