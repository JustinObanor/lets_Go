package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"log"

	"github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "justin"
	password = "1999"
	dbname   = "messages"
)

var db *sql.DB
var tpl *template.Template

//Data struct containing id and info to be saved in datadb
type Data struct {
	Msg  string
	Name string
}

//DataDB struct containing id and info to be saved in db
type DataDB struct {
	ID int
	Data
}

func init() {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

}

func main() {
	//http.HandleFunc("/", msgIndex)x

	justin := DataDB{ID: 1, Data: Data{Msg: "Justin", Name: "ABC"}}
	ruben := DataDB{ID: 2, Data: Data{Msg: "Ruben", Name: "tratata"}}
	petyaT := DataDB{ID: 1, Data: Data{Msg: "Petya Tereodor Pidgallo", Name: "ololol"}}

	datas := []DataDB{justin, ruben, petyaT}

	tpl, err := template.New("msgs").Parse(` {{range .}}	
	 {{.ID}} Hello {{.Data.Msg}}, my name is {{.Data.Name}}
	{{end}}
	`)
	if err != nil {
		panic(err)
	}
	
	var s bytes.Buffer

	err = tpl.Execute(&s, datas)
	if err != nil {
		log.Fatalln(err)
	}

	result := s.String()
	fmt.Println(result)

	for _, v := range datas {
		s := []string{v.Data.Msg, v.Data.Name}
		_, err = db.Exec("INSERT INTO datadb (id, data) VALUES($1, $2)", v.ID, pq.Array(s))
		if err != nil {
			panic(err)
		}
	}
}

/*
func createMsg(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}
	d := DataDB{}

	i := r.FormValue("id")
	s := fmt.Sprintf("%s %s", d.Data.Msg, d.Data.Name)
	s = r.FormValue("data")

	_, err := db.Exec("INSERT INTO datadb (id, data) VALUES($1, $2)", d.ID, pq.Array(d.Data))
	if err != nil {
		panic(err)
	}

}

func msgIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query("SELECT * FROM datadb")
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer rows.Close()

	ds := make([]DataDB, 0)
	for rows.Next() {
		d := DataDB{}
		err := rows.Scan(&d.ID, &d.Data)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		ds = append(ds, d)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
}
*/
