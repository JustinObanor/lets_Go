package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
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

var wg sync.WaitGroup
var db *sql.DB
var w http.ResponseWriter
var data DataLog

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

func getter(url string) []DataLog {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	data = DataLog{}

	b, err := ioutil.ReadAll(res.Body)

	s := make([]DataLog, 0)

	json.Unmarshal(b, &data)

	return s
}

//Worker func loads data to chan
func Worker(url string, c chan DataLog, d chan bool) {
	s := getter(url)

	s = append(s, data)

	for _, v := range s {
		c <- DataLog{v.ID, v.CurrentFileTime}
	}

	ticker := time.NewTicker(10 * time.Second)

	for {
		select {
		case <-ticker.C:
			s := getter(url)

			for _, v := range s {
				c <- DataLog{v.ID, v.CurrentFileTime}
			}

		case <-d:
			fmt.Fprint(w, "Worker finished work")
			wg.Done()
			return
		}
	}
}

//Puller retrieves data from chan
func Puller(c chan DataLog, d chan bool) {
	for {
		select {
		case v := <-c:
			_, err := db.Exec("INSERT INTO timetable (time) VALUES ($1)", v.CurrentFileTime)
			if err != nil {
				panic(err)
			}
		case <-d:
			fmt.Fprint(w, "Puller finished pulling")
			wg.Done()
			return
		}
	}
}

type row int

func (o row) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("select * from timetable order by id desc limit 1")
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	s := make([]DataLog, 0)
	for rows.Next() {
		ctime := DataLog{}
		err = rows.Scan(&ctime.ID, &ctime.CurrentFileTime)
		if err != nil {
			panic(err)
		}
		s = append(s, ctime)
	}
	if err = rows.Err(); err != nil {
		panic(err)
	}

	for _, v := range s {
		fmt.Fprintf(w, "Current Time : [ID : %d\t  Time : %v]\n", v.ID, v.CurrentFileTime)
	}

}

type all int

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
		fmt.Fprintf(w, "Time %d : %d\n", v.ID, v.CurrentFileTime)
	}
}

func main() {
	var s all
	var o row
	url := "http://worldclockapi.com/api/json/utc/now"

	c := make(chan DataLog)
	d := make(chan bool)

	wg.Add(2)
	go Worker(url, c, d)
	go Puller(c, d)

	mux := http.NewServeMux()
	mux.Handle("/", o)
	mux.Handle("/all", s)
	log.Fatal(http.ListenAndServe(":8080", mux))
	wg.Wait()
}
