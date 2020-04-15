package main

import (
	"fmt"
	"math/rand"
	"time"
)

//find number in the middle

/*
or this
rand.Seed(time.Now().UnixNano())
rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
*/

func main() {
	x := []int{1, 2, 5, 3, 6, 8, 9, 8, 54, 3, 5, 7, 89, 4, 3}
	rand.Seed(time.Now().UnixNano())
	for i := len(x) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		x[i], x[j] = x[j], x[i]
	}
	fmt.Println(x)
}

func maxNum(x []int) {

}
