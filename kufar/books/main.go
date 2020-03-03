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
	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

var host = getenv("PSQL_HOST", "db")
var port = getenv("PSQL_PORT", "5432")
var user = getenv("PSQL_USER", "postgres")
var password = getenv("PSQL_PWDcas", "postgres")
var dbname = getenv("PSQL_DB_NAME", "books")

func main() {
	r := mux.NewRouter()
	db, err := New()
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}

	defer db.Close()

	r.HandleFunc("/signup", SignUp(*db)).Methods("POST")
	r.HandleFunc("/signin/{uuid}", SignIn(*db)).Methods("GET")
	r.HandleFunc("/books", Create(*db)).Methods("POST")
	r.HandleFunc("/books", ReadAll(*db)).Methods("GET")
	r.HandleFunc("/books/{id}", Read(*db)).Methods("GET")
	r.HandleFunc("/books/{id}", Update(*db)).Methods("PUT")
	r.HandleFunc("/books/{id}", Delete(*db)).Methods("DELETE")
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

//BookResponse struct
type BookResponse struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Author string `json:"author"`
	Date   string `json:"date"`
}

//Credentials ...
type Credentials struct {
	UUID     int    
	Username string 
	Password string 
}

//Database ...
type Database struct {
	db *sql.DB
}

var location, _ = time.LoadLocation("Europe/Minsk")

func convertToResponse(books Book) BookResponse {

	return BookResponse{
		ID:     books.ID,
		Name:   books.Name,
		Author: books.Author,
		Date:   books.Date.In(location).Format(time.RFC1123),
	}
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
	log.Println("You connected to your database.")

	return &Database{
		db: db,
	}, nil

}

//Close closes the conn
func (d Database) Close() error {
	return d.db.Close()
}

//SignUpUser ...
func (d Database) SignUpUser(w http.ResponseWriter, r *http.Request) {
	cred := Credentials{}
	if err := json.NewDecoder(r.Body).Decode(&cred); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+"Error unmarshalling json", http.StatusInternalServerError)
		return
	}

	pword, err := bcrypt.GenerateFromPassword([]byte(cred.Password), 8)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if _, err = d.db.Exec("insert into credentials(uuid, username, password) values($1, $2, $3)", cred.UUID, cred.Username, string(pword)); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not log in user. Try different uuid", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(http.StatusText(http.StatusOK)))

}

//SignInUser ...
func (d Database) SignInUser(w http.ResponseWriter, r *http.Request) {
	uuid := mux.Vars(r)["uuid"]
	uuidInt, err := strconv.Atoi(uuid)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if uuid == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest)+": missing parameter in url", http.StatusBadRequest)
		return
	}

	row := d.db.QueryRow("select uuid, password from credentials where uuid=$1", uuidInt)

	creds := Credentials{}

	dbCreds := Credentials{}

	err = row.Scan(&creds.UUID, &creds.Password)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	log.Println("1",dbCreds.Password, "2",creds.Password)

	if err = bcrypt.CompareHashAndPassword([]byte(dbCreds.Password), []byte(creds.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}


	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(creds); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

//CreateBook ...
func (d Database) CreateBook(w http.ResponseWriter, r *http.Request) {
	var bk Book
	if err := json.NewDecoder(r.Body).Decode(&bk); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+"Error unmarshalling json", http.StatusInternalServerError)
		return
	}

	now := time.Now().UTC()

	_, err := d.db.Exec("insert into books(id, name, author, date) values($1, $2, $3, $4)", bk.ID, bk.Name, bk.Author, now)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not add new books. Try changing id", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

//ReadBooks ...
func (d Database) ReadBooks(w http.ResponseWriter, r *http.Request) {
	rows, err := d.db.Query("select * from books order by id asc")

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not list books", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	bks := []Book{}
	for rows.Next() {
		bk := Book{}
		err := rows.Scan(&bk.ID, &bk.Name, &bk.Author, &bk.Date)

		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		bks = append(bks, bk)
	}
	var resps []BookResponse
	for _, book := range bks {
		resp := convertToResponse(book)
		resps = append(resps, resp)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resps); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

//ReadBook ...
func (d Database) ReadBook(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if id == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest)+": missing parameter in url", http.StatusBadRequest)
		return
	}

	row := d.db.QueryRow("select * from books where id = $1", idInt)

	bk := Book{}
	err = row.Scan(&bk.ID, &bk.Name, &bk.Author, &bk.Date)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := convertToResponse(bk)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

//UpdateBook ...
func (d Database) UpdateBook(w http.ResponseWriter, r *http.Request) {
	log.Println("chu blyat")
	id := mux.Vars(r)["id"]
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if id == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest)+": missing parameter in url", http.StatusBadRequest)
		return
	}

	bk := Book{}

	if err := json.NewDecoder(r.Body).Decode(&bk); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	_, err = d.db.Exec("update books set id = $1, name = $2, author = $3, date = $4 where id = $1", idInt, bk.Name, bk.Author, bk.Date)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

//DeleteBook ...
func (d Database) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if id == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest)+": missing parameter in url", http.StatusBadRequest)
		return
	}

	_, err = d.db.Exec("delete from books where id = $1", idInt)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
