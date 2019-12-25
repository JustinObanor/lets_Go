package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"os"
	"strings"

	_ "github.com/lib/pq"
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

//DataDB struct containing id and info to be saved in db
type DataDB struct {
	ID  int
	Msg string
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

	//http.HandleFunc("/", msgIndex)

	var name = "justin"
	var msg = "abc"

	justin := DataDB{
		ID:  1,
		Msg: fmt.Sprintf("%s %s", name, msg),
	}
	//s := []DataDB{justin}

	msgs := strings.Split(justin.Msg, " ")
	fmt.Println(msgs[0])
	fmt.Println(msgs[len(msgs)-1])

	tpl, err := template.New("datas").Parse(`{{range .}} Hello {{. msgs[0]}}, name is {{. msgs[len(msgs)-1]}}`)
	if err != nil {
		panic(err)
	}
	// for _, v := range datas {
	// 	s := []string{v.Data.Msg, v.Data.Name}
	// 	_, err = db.Exec("INSERT INTO datadb (id, data) VALUES($1, $2)", v.ID, pq.Array(s))
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	tpl.Execute(os.Stdout, msgs)
}

// func createMsg(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != "POST" {
// 		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
// 		return
// 	}
// 	d := DataDB{}
// 	out, err := json.Marshal(d)
// 	if err != nil {
// 		panic(err)
// 	}

// 	i := r.FormValue("id")
// 	d.Data = r.FormValue("out")
// }

// func msgIndex(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != "GET" {
// 		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
// 		return
// 	}

// 	rows, err := db.Query("SELECT * FROM datadb")
// 	if err != nil {
// 		http.Error(w, http.StatusText(500), 500)
// 		return
// 	}
// 	defer rows.Close()

// 	ds := make([]DataDB, 0)
// 	for rows.Next() {
// 		d := DataDB{}
// 		err := rows.Scan(&d.ID, &d.Data)
// 		if err != nil {
// 			http.Error(w, http.StatusText(500), 500)
// 			return
// 		}
// 		ds = append(ds, d)
// 	}
// 	if err = rows.Err(); err != nil {
// 		http.Error(w, http.StatusText(500), 500)
// 		return
// 	}
// }
