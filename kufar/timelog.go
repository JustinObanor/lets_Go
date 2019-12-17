package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "justin"
	password = "1999"
	dbname   = "timestamp"
)

/*
//TimeTable table from db
type TimeTable struct {
	ID   int
	Time string
}

var db *sql.DB

func init() {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
}

//Insert func inserts row to db
func Insert() {
	t := time.Now().Format("3:04")

	_, err := db.Exec("INSERT INTO timetable (time) VALUES ($1)", t)
	if err != nil {
		panic(err)
	}

	fmt.Printf("The current time is : %s\n", t)
}

//SelectAll func selects all rows from db
func SelectAll() {
	rows, err := db.Query("SELECT * FROM timetable")
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	stamps := make([]TimeTable, 0)
	for rows.Next() {
		stamp := TimeTable{}
		err = rows.Scan(&stamp.ID, &stamp.Time)
		if err != nil {
			panic(err)
		}
		stamps = append(stamps, stamp)
	}
	if err = rows.Err(); err != nil {
		panic(err)
	}
	for _, v := range stamps {
		fmt.Printf("Time  %d : %s\n", v.ID, v.Time)
	}
}
*/

//TimeLog is the data being parsed from site
type TimeLog struct {
	TimeLog []DataLog `json:"timelog"`
}

//DataLog is the json attributes to be parsed
type DataLog struct {
	ID              int    `json:"$id"`
	CurrentDateTime string `json:"currentDateTime"`
	CurrentFileTime int64  `json:"currentFileTime"`
}

//ServeHTTP gets the body
func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := "http://worldclockapi.com/api/json/utc/now"
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	data := TimeLog{}
	s := make([]TimeLog, 0)

	b, err := ioutil.ReadAll(res.Body)

	json.Unmarshal(b, &data)

	s = append(s, data)

	for i := 0; i < len(data.TimeLog); i++ {
		fmt.Printf("ID : %d\t Date : %s\t Time : %v", data.TimeLog[i].ID, data.TimeLog[i].CurrentDateTime, data.TimeLog[i].CurrentFileTime)
	}
}

func main() {
	http.HandleFunc("/", ServeHTTP)
	log.Fatal(http.ListenAndServe(":8080", nil))

	/*
		Insert()
		SelectAll()
	*/
}
