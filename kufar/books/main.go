package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-redis/redis"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

const (
	cost = 12
)

var host = getenv("PSQL_HOST", "db")
var port = getenv("PSQL_PORT", "5432")
var user = getenv("PSQL_USER", "postgres")
var password = getenv("PSQL_PWDcas", "postgres")
var dbname = getenv("PSQL_DB_NAME", "book")

var location, _ = time.LoadLocation("Europe/Minsk")

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

//Book ...
type Book struct {
	ID     int
	Name   string
	Author string
	Date   time.Time
	UserID int
}

//BookRequest ...
type BookRequest struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Author string `json:"author"`
}

//BookResponse ...
type BookResponse struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Author string `json:"author"`
	Date   string `json:"date"`
	UserID int    `json:"userid"`
}

//Credentials ...
type Credentials struct {
	UUID     int
	Username string
	Password string
}

//CredentialsRequest ...
type CredentialsRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//Database ...
type Database struct {
	db *sql.DB
}

type rediscache struct {
	redis *redis.Client
}

//Cache interface
type Cache interface {
	Get(string) (Book, error)
	Set(string, *Book) error
	Remove(string) error
}

func main() {
	db, err := newDB()
	if err != nil {
		log.Fatalf("error connecting to db: %v", err)
	}
	defer db.Close()

	c, err := newRedisCacheClient()
	if err != nil {
		log.Fatalf("error connecting to redis: %v", err)
	}

	r := chi.NewRouter()
	r.Post("/signup", SignUpUser(*db))

	r.Middlewares()

	r.Route("/books", func(r chi.Router) {
		r.Use(StupidMiddleware("1"))
		r.Post("/", CreateBook(*db))
		r.Get("/", ReadBooks(*db))

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", ReadBook(*db, c))
			r.Put("/", UpdateBook(*db, c))
			r.Delete("/", DeleteBook(*db, c))
		})
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}

func getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

//StupidMiddleware ...
func StupidMiddleware(id string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if id == "8" {
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func convertToResponse(books Book) BookResponse {
	return BookResponse{
		ID:     books.ID,
		Name:   books.Name,
		Author: books.Author,
		Date:   books.Date.In(location).Format(time.RFC1123),
		UserID: books.UserID,
	}
}

// newRedisCacheClient ...
func newRedisCacheClient() (*rediscache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "db_redis:6379",
		Password: os.Getenv("Password"),
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	log.Println("redis connected")

	return &rediscache{
		redis: client,
	}, nil
}

func newDB() (*Database, error) {
	var err error

	psqlInfo := fmt.Sprintf(`host=%s port=%s user=%s
	password=%s dbname=%s sslmode=disable`, host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	if err = db.Ping(); err != nil {
		return nil, err
	}
	log.Println("postgres connected")

	return &Database{
		db: db,
	}, nil
}

//Close closes the conn
func (d Database) Close() error {
	return d.db.Close()
}

func (d Database) getBookUserID(bookID int) (int, error) {
	row := d.db.QueryRow("select userid from books where id = $1", bookID)

	var userID int
	err := row.Scan(&userID)

	switch {
	case err == sql.ErrNoRows:
		return 0, err
	case err != nil:
		return 0, err
	}
	return userID, nil
}

//CheckAuth ...
func (d Database) CheckAuth(header *http.Header) (int, bool) {
	cred := header.Get("Authorization")
	if cred == "" {
		return 0, false
	}

	s := strings.Split(cred, " ")
	if len(s) != 2 || s[0] != "Basic" || s[1] == "" {
		return 0, false
	}

	b, err := base64.StdEncoding.DecodeString(s[1])
	if err != nil {
		return 0, false
	}

	bStr := string(b)
	if !strings.ContainsRune(bStr, ':') {
		return 0, false
	}

	creds := strings.Split(bStr, ":")
	if len(creds) != 2 || s[0] == "" || s[1] == "" {
		return 0, false
	}

	row := d.db.QueryRow("select uuid, username, password from credentials where username = $1", creds[0])

	dbCred := Credentials{}
	err = row.Scan(&dbCred.UUID, &dbCred.Username, &dbCred.Password)

	switch {
	case err == sql.ErrNoRows:
		return 0, false
	case err != nil:
		return 0, false
	}

	if err = bcrypt.CompareHashAndPassword([]byte(dbCred.Password), []byte(creds[1])); err != nil {
		return 0, false
	}
	return dbCred.UUID, true
}

//SignUpUser ...
func SignUpUser(d Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		credReq := CredentialsRequest{}
		if err := json.NewDecoder(r.Body).Decode(&credReq); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": error unmarshalling json", http.StatusInternalServerError)
			return
		}

		if credReq.Username == "" && credReq.Password == "" {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": missing username or password", http.StatusInternalServerError)
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

		w.Write([]byte(http.StatusText(http.StatusOK) + ": signup passed"))
	}
}

//CreateBook ...
func CreateBook(d Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, valid := d.CheckAuth(&r.Header)
		if !valid {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": invalid creds", http.StatusInternalServerError)
			return
		}

		var bk Book
		if err := json.NewDecoder(r.Body).Decode(&bk); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest)+": eror unmarshalling json", http.StatusBadRequest)
			return
		}

		now := time.Now().UTC()

		if userID != bk.UserID {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(http.StatusText(http.StatusUnauthorized) + ": you dont have access to this resource"))
			return
		}

		if _, err := d.db.Exec("insert into books(id, name, author, date, userid) values($1, $2, $3, $4, $5)", bk.ID, bk.Name, bk.Author, now, bk.UserID); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not add new books. Try changing id", http.StatusInternalServerError)
			return
		}
		w.Write([]byte(http.StatusText(http.StatusOK) + ": created book"))
	}
}

//ReadBooks ...
func ReadBooks(d Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := d.db.Query("select id, name, author, date, userid from books order by id asc")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not list books", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		bks := []Book{}
		for rows.Next() {
			bk := Book{}
			if err := rows.Scan(&bk.ID, &bk.Name, &bk.Author, &bk.Date, &bk.UserID); err != nil {
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
			http.Error(w, http.StatusText(http.StatusBadRequest)+": error marshalling json", http.StatusBadRequest)
			return
		}
	}
}

func (r *rediscache) Get(id string) (Book, error) {
	var bk Book
	val, err := r.redis.Get(id).Result()
	if err == redis.Nil || err != nil {
		return bk, err
	}

	if err := json.Unmarshal([]byte(val), &bk); err != nil {
		return bk, err
	}

	return bk, nil
}

func (r *rediscache) Set(id string, bk *Book) error {
	b, err := json.Marshal(&bk)
	if err != nil {
		return err
	}

	return r.redis.Set(id, string(b), time.Hour).Err()
}

func (r *rediscache) Remove(id string) error {
	return r.redis.Del(id).Err()
}

//ReadBook ...
func ReadBook(d Database, c Cache) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		if id == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest)+": missing parameter in url", http.StatusBadRequest)
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not convert to integer", http.StatusInternalServerError)
			return
		}

		bk, err := c.Get(id)

		if err != nil {
			row := d.db.QueryRow("select id, name, author, date, userid from books where id = $1", idInt)

			err = row.Scan(&bk.ID, &bk.Name, &bk.Author, &bk.Date, &bk.UserID)
			switch {
			case err == sql.ErrNoRows:
				http.NotFound(w, r)
				return
			case err != nil:
				http.Error(w, http.StatusText(http.StatusInternalServerError)+": error scanning db", http.StatusInternalServerError)
				return
			}
			c.Set(id, &bk)
		}

		resp := convertToResponse(bk)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest)+": Error marshalling json", http.StatusBadRequest)
			return
		}

	}
}

//UpdateBook ...
func UpdateBook(d Database, c Cache) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, valid := d.CheckAuth(&r.Header)
		if !valid {
			return
		}

		id := chi.URLParam(r, "id")

		if id == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest)+": missing parameter in url", http.StatusBadRequest)
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not convert to integer", http.StatusInternalServerError)
			return
		}

		bkUserID, err := d.getBookUserID(idInt)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not convert to integer", http.StatusInternalServerError)
			return
		}

		if userID != bkUserID {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": stop right there criminal scum!", http.StatusInternalServerError)
			return
		}

		var bkReq BookRequest
		if err := json.NewDecoder(r.Body).Decode(&bkReq); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest)+": error unmarshalling json", http.StatusBadRequest)
			return
		}

		bk := Book{ID: bkReq.ID, Name: bkReq.Name, Author: bkReq.Author}
		if _, err = d.db.Exec("update books set id = $1, name = $2, author = $3 where id = $1", idInt, bk.Name, bk.Author); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if err := c.Remove(id); err != nil {
			return
		}

		w.Write([]byte(http.StatusText(http.StatusOK) + ": updated book"))
	}
}

//DeleteBook ...
func DeleteBook(d Database, c Cache) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, valid := d.CheckAuth(&r.Header)
		if !valid {
			return
		}

		id := chi.URLParam(r, "id")

		if id == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest)+": missing parameter in url", http.StatusBadRequest)
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not convert to integer", http.StatusInternalServerError)
			return
		}

		bkUserID, err := d.getBookUserID(idInt)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not convert to integer", http.StatusInternalServerError)
			return
		}

		if userID != bkUserID {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(http.StatusText(http.StatusUnauthorized) + ": you dont have access to this resource"))
			return
		}

		var bk BookRequest
		if err := json.NewDecoder(r.Body).Decode(&bk); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest)+": Error unmarshalling json", http.StatusBadRequest)
			return
		}

		if _, err = d.db.Exec("delete from books where id = $1", idInt); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not delete book", http.StatusInternalServerError)
			return
		}

		if err := c.Remove(id); err != nil {
			return
		}

		w.Write([]byte(http.StatusText(http.StatusOK) + ": deleted book"))
	}
}
