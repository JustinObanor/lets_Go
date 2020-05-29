package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	file, err := os.OpenFile("info.txt", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal("open fail", err)
	}
	defer file.Close()

	b := bufio.NewWriter(file)

	if _, err := b.Write([]byte("The Go gopher is an iconic mascot!")); err != nil {
		log.Println(err)
	}

	fmt.Println(b.Buffered())
	fmt.Println(b.Available())

	if err := b.Flush(); err != nil {
		log.Fatal("write fail ", err)
	}

	bss, _ := ioutil.ReadFile("info.txt")
	fmt.Printf("%s", bss)

}
