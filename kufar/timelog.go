package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "justin"
	password = "1999"
	dbname   = "timestamp"
)

var db *sql.DB

func init() {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
}

//DataLog is the json attributes to be parsed
type DataLog struct {
	ID              int `json:"$id"`
	CurrentFileTime int `json:"currentFileTime"`
}

type insert int

type all int

//ServeHTTP gets the body
func (i insert) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	

	url := "http://worldclockapi.com/api/json/utc/now"
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	data := DataLog{}

	b, err := ioutil.ReadAll(res.Body)

	s := make([]DataLog, 0)

	json.Unmarshal(b, &data)
	s = append(s, data)

	for _, v := range s {
		t := v.CurrentFileTime

		_, err = db.Exec("INSERT INTO timetable (time) VALUES ($1)", t)
		if err != nil {
			panic(err)
		}

		fmt.Fprintf(w,"Current Info : [ID : %d\t  Time : %v]\n", v.ID, v.CurrentFileTime)
	}
}

//SelectAll func selects all rows from db
func (s all) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM timetable")
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	stamps := make([]DataLog, 0)
	for rows.Next() {
		stamp := DataLog{}
		err = rows.Scan(&stamp.ID, &stamp.CurrentFileTime)
		if err != nil {
			panic(err)
		}
		stamps = append(stamps, stamp)
	}
	if err = rows.Err(); err != nil {
		panic(err)
	}
	for _, v := range stamps {
		fmt.Fprintf(w,"Time %d : %d\n", v.ID, v.CurrentFileTime)
	}
}

func main() {
	var i insert
	var s all

	mux := http.NewServeMux()
	mux.Handle("/insert", i)
	mux.Handle("/all", s)

	tout := http.TimeoutHandler(mux, time.Second, "Timeout!")

	log.Fatal(http.ListenAndServe(":8080", tout))
}
