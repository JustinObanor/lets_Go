package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	// open the main.go source file - likely this one
	file, err := os.Open("testdata")
	if err != nil {
		log.Fatal(err)
	}

	fInfo, err := file.Readdir(-1)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range fInfo {
		fmt.Println(v.Name())
	}
}
