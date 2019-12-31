package main

import "fmt"

type Marks struct {
	fname, lname string
}

func main() {
	marks := make([]Marks, 0)
	marks = append(marks, Marks{
		fname: "Justin",
		lname: "Obanor",
	})
	c := make(chan []Marks)
	go func() {
		c <- marks
	}()
	for i := 0; i < len(marks); i++ {
		v := <-c
		fmt.Println(v[i].fname, v[i].lname)
	}
}
