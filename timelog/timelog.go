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

var (
	host     = getenv("PSQL_HOST", "db_time")
	port     = getenv("PSQL_PORT", "5432")
	user     = getenv("PSQL_USER", "postgres")
	password = getenv("PSQL_PWDcas", "postgres")
	dbname   = getenv("PSQL_DB_NAME", "timelog")
)

var url = "http://worldclockapi.com/api/json/utc/now"

//DataLog ...
type DataLog struct {
	ID   int    `json:"id"`
	Time string `json:"currentDateTime"`
}

//Database ...
type Database struct {
	*sql.DB
}

type client struct {
	*http.Client
}

func newDB() (*Database, error) {
	psqlInfo := fmt.Sprintf(`host=%s port=%s user=%s
	password=%s dbname=%s sslmode=disable`, host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	log.Println("postgres connected")
	return &Database{
		db,
	}, nil
}

//Close closes the conn
func (db Database) Close() error {
	return db.Close()
}

func newClient() *client {
	return &client{
		&http.Client{
			Timeout: time.Second * 5,
		},
	}
}

func getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func timeLog(url string) ([]DataLog, error) {
	client := newClient()

	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	s := make([]DataLog, 0)
	data := DataLog{}

	err = json.Unmarshal(b, &data)
	if err != nil {
		return nil, err
	}

	s = append(s, data)

	return s, nil
}

func main() {
	c := make(chan DataLog)
	var wg sync.WaitGroup

	db, err := newDB()
	if err != nil {
		log.Println(err)
	}

	go func() {
		wg.Add(1)
		for value := range c {
			_, err := db.Exec("insert into datalog (time) values ($1)", value.Time)
			if err != nil {
				log.Println(err)
			}
		}
	}()

	ticker := time.NewTicker(time.Second)

	for range ticker.C {
		wg.Add(1)
		d, err := timeLog(url)
		if err != nil {
			fmt.Println(err)
		}

		go func() {
			for _, value := range d {
				c <- DataLog{value.ID, value.Time}
			}
		}()
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("select * from datalog order by id desc limit 1")
		if err != nil {
			log.Println(err)
		}
		defer rows.Close()

		v := DataLog{}

		if err := rows.Scan(&v.ID, &v.Time); err != nil {
			log.Println(err)
		}

		fmt.Fprintf(w, "[ID:%d\tTime:%v]\n", v.ID, v.Time)
	})

	http.HandleFunc("/all", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("select * from datalog")
		if err != nil {
			log.Println(err)
		}

		defer rows.Close()

		stamps := make([]DataLog, 0)
		data := DataLog{}
		for rows.Next() {
			if err := rows.Scan(&data.ID, &data.Time); err != nil {
				log.Println(err)
			}
			stamps = append(stamps, data)
		}

		for _, v := range stamps {
			fmt.Fprintf(w, "[ID:%d\tTime:%v]\n", v.ID, v.Time)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
	wg.Wait()
}
