package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/lets_Go/quiz/internal"
)

func newDB() (*internal.Database, error) {
	var err error

	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		return nil, err
	}

	log.Println("boltdb connected")

	return &internal.Database{
		DB: db,
	}, nil
}

func main() {
	db, err := newDB()
	if err != nil {
		log.Println(err)
	}

	defer db.DB.Close()

	mux := defaultMux()
	ymlFile := flag.String("f", "data.yml", "yaml file where containing the URLS")
	flag.Parse()

	db.DB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(internal.DBBucket))
		if err != nil {
			log.Printf("create bucket: %s\n", err)
		}

		b.Put([]byte("/yaml-godoc"),
			[]byte("https://godoc.org/gopkg.in/yaml.v2"))

		b.Put([]byte("/urlshort-godoc"),
			[]byte("https://godoc.org/github.com/gophercises/urlshort"))

		return nil
	})

	mapHandler := db.MapHandler(mux)

	yaml, err := internal.ReadYAML(*ymlFile)
	if err != nil {
		log.Println(err)
	}

	yamlHandler, err := db.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	log.Println(http.ListenAndServe(":8080", yamlHandler))
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
