package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

var host = getenv("PSQL_HOST", "localhost")
var port = getenv("PSQL_PORT", "5440")
var user = getenv("PSQL_USER", "postgres")
var password = getenv("PSQL_PWDcas", "justin")
var dbname = getenv("PSQL_DB_NAME", "postgres")

func main() {
	r := mux.NewRouter()
	db, err := New()
	if err != nil {
		log.Fatal("error connecting to db")
	}

	r.HandleFunc("/books", db.CreateBook).Methods("POST")
	r.HandleFunc("/books", db.ReadBooks).Methods("GET")
	r.HandleFunc("/books/{id}", db.ReadBook).Methods("GET")
	r.HandleFunc("/books/{id}", db.UpdateBook).Methods("PUT")
	r.HandleFunc("/books/{id}", db.DeleteBook).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", r))

}

func getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

//Book ...
type Book struct {
	ID     int
	Name   string
	Author string
	Date   time.Time
}

type BookResponse struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Author string `json:"author"`
	Date   string `json:"date"`
}

//Database ...
type Database struct {
	db *sql.DB
}

func validateResponse(books []Book) []BookResponse {

	location, _ := time.LoadLocation("UTC")

	resp := make([]BookResponse, len(books))
	for i, v := range books {
		resp[i].ID = v.ID
		resp[i].Name = v.Name
		resp[i].Author = v.Author
		resp[i].Date = v.Date.In(location).Format(time.RFC1123)
	}
	return resp
}

//New constructor that return database
func New() (*Database, error) {
	var err error

	psqlInfo := fmt.Sprintf(`host=%s port=%s user=%s
	password=%s dbname=%s sslmode=disable`, host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	if err = db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("You connected to your database.")

	return &Database{
		db: db,
	}, nil

}

//CreateBook ...
func (d *Database) CreateBook(w http.ResponseWriter, r *http.Request) {
	var bk Book

	json.NewDecoder(r.Body).Decode(&bk)

	now := time.Now()

	_, err := d.db.Exec("insert into books(id, name, author, date) values($1, $2, $3, $4)", bk.ID, bk.Name, bk.Author, now)
	if err != nil {
		http.Error(w, http.StatusText(500)+": could not add new books. Try changing id", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(http.StatusText(200)))
}

//ReadBooks ...
func (d *Database) ReadBooks(w http.ResponseWriter, r *http.Request) {
	rows, err := d.db.Query("select * from books order by id asc")

	if err != nil {
		http.Error(w, http.StatusText(500)+": could not list books", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	bks := []Book{}
	for rows.Next() {
		bk := Book{}
		err := rows.Scan(&bk.ID, &bk.Name, &bk.Author, &bk.Date)

		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		bks = append(bks, bk)
	}

	resp := validateResponse(bks)

	json.NewEncoder(w).Encode(resp)

}

//ReadBook ...
func (d *Database) ReadBook(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	idInt, _ := strconv.Atoi(id)

	if id == "" {
		http.Error(w, http.StatusText(400)+": missing parameter in url", http.StatusBadRequest)
		return
	}

	row := d.db.QueryRow("select * from books where id = $1", idInt)

	bk := Book{}
	err := row.Scan(&bk.ID, &bk.Name, &bk.Author, &bk.Date)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	location, _ := time.LoadLocation("UTC")

	resp := BookResponse{
		ID:     bk.ID,
		Name:   bk.Name,
		Author: bk.Author,
		Date:   bk.Date.In(location).Format(time.RFC1123),
	}

	json.NewEncoder(w).Encode(resp)
}

//UpdateBook ...
func (d *Database) UpdateBook(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	idInt, _ := strconv.Atoi(id)

	if id == "" {
		http.Error(w, http.StatusText(400)+": missing parameter in url", http.StatusBadRequest)
		return
	}

	bk := Book{}

	json.NewDecoder(r.Body).Decode(&bk)

	_, err := d.db.Exec("update books set id = $1, name = $2, author = $3, date = $4 where id = $1", idInt, bk.Name, bk.Author, bk.Date)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(http.StatusText(200)))
}

//DeleteBook ...
func (d *Database) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	idInt, _ := strconv.Atoi(id)

	if id == "" {
		http.Error(w, http.StatusText(400)+": missing parameter in url", http.StatusBadRequest)
		return
	}

	_, err := d.db.Exec("delete from books where id = $1", idInt)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
}
