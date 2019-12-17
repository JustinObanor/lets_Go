package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

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
}
