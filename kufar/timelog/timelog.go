package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

var host = getenv("PSQL_HOST", "localhost")
var port = getenv("PSQL_PORT", "5432")
var user = getenv("PSQL_USER", "justin")
var password = getenv("PSQL_PWDcas", "1999")
var dbname = getenv("PSQL_DB_NAME", "timestamp")

var wg sync.WaitGroup
var db *sql.DB
var w http.ResponseWriter
var data DataLog

func init() {
	var err error
	psqlInfo := fmt.Sprintf(`host=%s port=%s user=%s
	password=%s dbname=%s sslmode=disable`, host, port, user, password, dbname)

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
}

func getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

type DataLog struct {
	ID   int    `json:$id`
	Time string `json:"currentDateTime"`
}

func getter(url string) []DataLog {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	data = DataLog{}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	s := make([]DataLog, 0)

	err = json.Unmarshal(b, &data)
	if err != nil {
		log.Fatal(err)
	}
	s = append(s, data)

	return s
}

func worker(url string, c chan DataLog, d chan bool) {
	s := getter(url)

	for _, v := range s {
		c <- DataLog{v.ID, v.Time}
	}

	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-ticker.C:
			s := getter(url)
			for _, v := range s {
				c <- DataLog{v.ID, v.Time}
			}

		case <-d:
			fmt.Println("Worker finished work")
			wg.Done()
			return
		}
	}
}

func puller(c chan DataLog, d chan bool) {
	insQuery := "INSERT INTO timetable (time) VALUES ($1)"
	for {
		select {
		case v := <-c:
			_, err := db.Exec(insQuery, v.Time)
			if err != nil {
				panic(err)
			}
		case <-d:
			fmt.Println("Puller finished pulling")
			wg.Done()
			return
		}
	}
}

func recentTime(w http.ResponseWriter, _ *http.Request) {
	selQuery := "select * from timetable order by id desc limit 1"
	rows, err := db.Query(selQuery)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	s := make([]DataLog, 0)
	for rows.Next() {
		err = rows.Scan(&data.ID, &data.Time)
		if err != nil {
			panic(err)
		}
		s = append(s, data)
	}
	if err = rows.Err(); err != nil {
		panic(err)
	}

	for _, v := range s {
		fmt.Fprintf(w, "Current Time : [ID : %d\t  Time : %v]\n", v.ID, v.Time)
	}

}

func allTime(w http.ResponseWriter, _ *http.Request) {
	selQuery := "SELECT * FROM timetable"
	rows, err := db.Query(selQuery)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	stamps := make([]DataLog, 0)
	for rows.Next() {
		err = rows.Scan(&data.ID, &data.Time)
		if err != nil {
			panic(err)
		}
		stamps = append(stamps, data)
	}
	if err = rows.Err(); err != nil {
		panic(err)
	}

	for _, v := range stamps {
		fmt.Fprintf(w, "Time %d : %v\n", v.ID, v.Time)
	}
}

func main() {
	url := "http://worldclockapi.com/api/json/utc/now"

	c := make(chan DataLog)
	d := make(chan bool)

	wg.Add(2)
	go worker(url, c, d)
	go puller(c, d)

	http.HandleFunc("/", recentTime)
	http.HandleFunc("/all", allTime)
	log.Fatal(http.ListenAndServe(":8080", nil))
	wg.Wait()
}
