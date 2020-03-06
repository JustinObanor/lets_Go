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

const cost = 12

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
	r.HandleFunc("/signin", SignIn(*db)).Methods("POST")
	r.HandleFunc("/logout", Logout(*db)).Methods("POST")
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
	UUID   int
}

//BookResponse struct
type BookResponse struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Author string `json:"author"`
	Date   string `json:"date"`
	UUID   int    `json:"uuid"`
}

//CredentialsRequest ...
type CredentialsRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
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
		UUID:   books.UUID,
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
	credReq := CredentialsRequest{}
	if err := json.NewDecoder(r.Body).Decode(&credReq); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+": error unmarshalling json", http.StatusInternalServerError)
		return
	}

	pword, err := bcrypt.GenerateFromPassword([]byte(credReq.Password), cost)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not add password", http.StatusInternalServerError)
		return
	}

	cred := Credentials{Username: credReq.Username}
	if _, err = d.db.Exec("insert into credentials(username, password) values($1, $2)", cred.Username, string(pword)); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not log in user. Username already exists", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(http.StatusText(http.StatusOK)))
}

//SignInUser ...
func (d Database) SignInUser(w http.ResponseWriter, r *http.Request) {
	credReq := CredentialsRequest{}
	if err := json.NewDecoder(r.Body).Decode(&credReq); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+"Error unmarshalling json", http.StatusInternalServerError)
		return
	}

	row := d.db.QueryRow("select password from credentials where username = $1", credReq.Username)

	cred := Credentials{}
	err := row.Scan(&cred.Password)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not scan db", http.StatusInternalServerError)
		return
	}

	session, err := store.Get(r, "my-cookie")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not get cookie", http.StatusInternalServerError)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(cred.Password), []byte(credReq.Password)); err != nil {
		session.AddFlash("Incorrect credentials")
		if err = session.Save(r, w); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not check password", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/signup", http.StatusSeeOther)
		return
	}

	//cred.Username = credReq.Username

	session.Values["user"] = &cred.UUID
	if err = session.Save(r, w); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not assign cookie", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(http.StatusText(http.StatusOK)))
}

//LogoutUser ...
func (d Database) LogoutUser(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "my-cookie")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not get cookie", http.StatusInternalServerError)
		return
	}

	session.Values["user"] = Credentials{}
	session.Options.MaxAge = -1

	if err = session.Save(r, w); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not save cookie", http.StatusInternalServerError)
		return
	}
}

//CreateBook ...
func (d Database) CreateBook(w http.ResponseWriter, r *http.Request) {
	var bk Book
	if err := json.NewDecoder(r.Body).Decode(&bk); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+": Error unmarshalling json", http.StatusInternalServerError)
		return
	}

	now := time.Now().UTC()

	if _, err := d.db.Exec("insert into books(id, name, author, date, uuid) values($1, $2, $3, $4, $5)", bk.ID, bk.Name, bk.Author, now, bk.UUID); err != nil {
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
		if err := rows.Scan(&bk.ID, &bk.Name, &bk.Author, &bk.Date, &bk.UUID); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not scan db", http.StatusInternalServerError)
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
		http.Error(w, http.StatusText(http.StatusInternalServerError)+": error marshelling json", http.StatusInternalServerError)
		return
	}
}

//ReadBook ...
func (d Database) ReadBook(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not convert to integer", http.StatusInternalServerError)
		return
	}

	if id == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest)+": missing parameter in url", http.StatusBadRequest)
		return
	}

	row := d.db.QueryRow("select * from books where id = $1", idInt)

	bk := Book{}
	err = row.Scan(&bk.ID, &bk.Name, &bk.Author, &bk.Date, &bk.UUID)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(http.StatusInternalServerError)+": error scanning db", http.StatusInternalServerError)
		return
	}

	resp := convertToResponse(bk)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+": Error marshalling json", http.StatusInternalServerError)
		return
	}
}

//UpdateBook ...
func (d Database) UpdateBook(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not convert to integer", http.StatusInternalServerError)
		return
	}

	if id == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest)+": missing parameter in url", http.StatusBadRequest)
		return
	}

	bk := Book{}

	if err := json.NewDecoder(r.Body).Decode(&bk); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+": Error unmarshalling json", http.StatusInternalServerError)
		return
	}

	if _, err = d.db.Exec("update books set id = $1, name = $2, author = $3, date = $4, uuid = $5 where id = $1", idInt, bk.Name, bk.Author, bk.Date, bk.UUID); err != nil {
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
		http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not convert to integer", http.StatusInternalServerError)
		return
	}

	if id == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest)+": missing parameter in url", http.StatusBadRequest)
		return
	}

	if _, err = d.db.Exec("delete from books where id = $1", idInt); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not delete book", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
