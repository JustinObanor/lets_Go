package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	cyoa "github.com/lets_Go/practice/cyoa/cyoaweb"
)

func main() {
	mux := http.NewServeMux()
	fileName := flag.String("file", "doc/gopher.json", "file that contains the story")
	port := flag.Int("port", 3000, "port to run server on")

	flag.Parse()

	file, err := os.Open(*fileName)
	if err != nil {
		panic(err)
	}

	story, err := cyoa.ParseJSON(file)
	if err != nil {
		panic(err)
	}

	handler := cyoa.NewHandler(story, cyoa.WithPathFunc(pathFn))
	mux.Handle("/story/", handler)

	mux.Handle("/", cyoa.NewHandler(story))

	log.Printf("server starting on port :%d", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
}

func pathFn(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/story" || path == "/story/" {
		path = "/story/intro"
	}
	return path[len("/story/"):]
}
