package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		op2(r.Context())
	})
	log.Println(http.ListenAndServe(":8080", nil))
}

func op2(ctx context.Context) {
	for {
		select {
		case <-time.After(time.Second):
			fmt.Println("pinged")
		case <-ctx.Done():
			fmt.Println("request cancelled")
			return
		}
	}
}

func op1() error {
	return errors.New("my error")
}
