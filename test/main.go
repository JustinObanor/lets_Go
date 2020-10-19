package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	f, err := ioutil.ReadDir("project")
	if err != nil{
		fmt.Println(err)
	}
	
	fmt.Println(f)
}