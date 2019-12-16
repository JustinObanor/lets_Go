package main

import "fmt"

func main() {
	var i = 0
	fmt.Println("In main, i is", i)
	foo(i) //only the copy is modified
	fmt.Println("After foo call, in main i is", i)
	bar(&i) //the value if modified
	fmt.Println("After bar call, in main i is", i)
}

func foo(i int) {
	i = 1
}

func bar(i *int) {
	*i = 1
}
