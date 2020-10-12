package main

import (
	"fmt"
	"strconv"
)

type wallet struct {
	cash int
}

func (w *wallet) String() string {
	return strconv.Itoa(w.cash)
}

func main() {
	w := &wallet{cash: 100}
	fmt.Printf("%#v\n", w)
	fmt.Printf("%s\n", w)
}
