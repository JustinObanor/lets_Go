package main

import (
	"fmt"
	"time"
)

func main() {
	age := 20
	bday(age)        //pass by value(a copy) and add age by 1
	fmt.Println(age) //21
	//age = bday(age)  assigning the new modified value back to age
	birthyear := calc(age) //pass by value(a copy) but not the prevous copy
	fmt.Println(birthyear)
}

func calc(age int) int {
	return time.Now().Year() - age
}

func bday(age int) int {
	return age + 1
}
