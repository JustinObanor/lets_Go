package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	f, _ := os.Open("a.txt")
	defer f.Close()

	s := bufio.NewScanner(f)
	s.Split(bufio.ScanLines)

	if !s.Scan() {
		err := s.Err()
		if err == nil {
			log.Println("EOF")
		} else {
			fmt.Println(err.Error())
		}
	}
	// fmt.Println(s.Text())

	for s.Scan() {
		fmt.Println(s.Text())
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
}
