package main

import (
	"fmt"
	"net/http"
)

type hotdog int

func main() {
	var h hotdog
	http.ListenAndServe(":8080", h)
}

func (h hotdog) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Sup fag")
}
