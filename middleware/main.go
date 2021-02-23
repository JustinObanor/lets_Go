package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
)

var (
	isDev = strings.EqualFold("dev", os.Getenv("ENV"))
)

type responseWriter struct {
	http.ResponseWriter
	status int
	writes [][]byte
}

// type ResponseWriter interface {
// 	Header() http.Header
// 	Write([]byte) (int, error)
// 	WriteHeader(statusCode int)
// }

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/panic/", panicDemo)
	mux.HandleFunc("/panic-after/", panicAfterDemo)
	mux.HandleFunc("/", hello)

	log.Println("Server listeneing on port :3000")
	log.Fatal(http.ListenAndServe(":3000", recoverMw(mux)))
}

func recoverMw(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic: %v\nStack trace: %s", r, debug.Stack())

				if !isDev {
					http.Error(w, "something went wrong :(", http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "<h1>Panic: %v</h1><pre>%s</pre>", r, debug.Stack())
			}
		}()

		nw := &responseWriter{ResponseWriter: w}
		next.ServeHTTP(nw, r)
		nw.flush()
	}
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
}

func (rw *responseWriter) Write(bs []byte) (int, error) {
	rw.writes = append(rw.writes, bs)
	return len(bs), nil
}

func (rw *responseWriter) flush() error {
	if rw.status != 0 {
		rw.ResponseWriter.WriteHeader(rw.status)
	}

	for _, bs := range rw.writes {
		_, err := rw.ResponseWriter.Write(bs)
		if err != nil {
			return err
		}
	}

	return nil
}

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := rw.ResponseWriter.(http.Hijacker)
	if ok {
		return hijacker.Hijack()
	}

	return nil, nil, errors.New("hijacking not supported")
}

func panicDemo(w http.ResponseWriter, r *http.Request) {
	funcThatPanics()
}

func panicAfterDemo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Hello!</h1>")
	funcThatPanics()
}

func funcThatPanics() {
	panic("Oh no!")
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello!</h1>")
}
