package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"os"

	_ "github.com/lib/pq"
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

//Data struct containing message and name
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
	justin := DataDB{ID:1,	Data: Data {Msg: "Justin", Name: "ABC"}}
	ruben :=  DataDB{ID:2,	Data: Data {Msg: "Ruben", Name: "tratata"}}
	petyaT := DataDB{ID:1,	Data: Data {Msg: "Petya Tereodor Pidgallo", Name: "ololol"}}

	datas := []DataDB{justin, ruben, petyaT}

	tpl, err := template.New("msgs").Parse(` {{range .}}
	 Hello {{.Data.Msg}}, my name is {{.Data.Name}}
	 {{end}}
	 `)
	if err != nil {
		panic(err)
	}

	err = tpl.Execute(os.Stdout, datas)

	if err != nil {
		log.Fatalln(err)
	}

	for _, v := range datas {
		s := []string{v.Data.Msg,v.Data.Name}
		_, err = db.Exec("INSERT INTO datadb (id, data) VALUES($1, $2)", v.ID, pq.Array(s))
		if err != nil {
			panic(err)
		}
	}
}
