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

	"github.com/go-chi/chi"
	"github.com/go-redis/redis"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

const (
	cost = 12
)

var host = getenv("PSQL_HOST", "localhost")
var port = getenv("PSQL_PORT", "5432")
var user = getenv("PSQL_USER", "postgres")
var password = getenv("PSQL_PWDcas", "justin1999")
var dbname = getenv("PSQL_DB_NAME", "dorm")

var location, _ = time.LoadLocation("Europe/Minsk")

//Student ...
type Student struct {
	ID                  int
	Firstname, Lastname string
	Date                time.Time
	UUID                int
	StudRoom            StudRoom
	StudFloor           StudFloor
}

//StudRoom ...
type StudRoom struct {
	ID, Room int
}

//StudFloor ...
type StudFloor struct {
	ID, Floor int
}

//Credentials ...
type Credentials struct {
	UUID     int
	Username string
	Password string
}

//StudentRequest ...
type StudentRequest struct {
	ID        int       `json:"id"`
	FName     string    `json:"firstname"`
	LName     string    `json:"lastname"`
	StudRoom  StudRoom  `json:"studroom"`
	StudFloor StudFloor `json:"studfloor"`
}

//StudentResponse ...
type StudentResponse struct {
	ID        int       `json:"id"`
	FName     string    `json:"firstname"`
	LName     string    `json:"lastname"`
	Date      string    `json:"date"`
	UUID      int       `json:"uuid"`
	StudRoom  StudRoom  `json:"studroom"`
	StudFloor StudFloor `json:"studfloor"`
}

//CredentialsRequest ...
type CredentialsRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//FloorCodeResReq ...
type FloorCodeResReq struct {
	Floor int `json:"floor"`
	Code  int `json:"code"`
}

//Cache interface
type Cache interface {
	Get(string) (Student, error)
	Set(string, *Student) error
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

	r.Route("/student", func(r chi.Router) {
		r.Post("/", CreateStudent(*db))
		r.Get("/", ReadStudents(*db))

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", ReadStudent(*db, c))
			r.Put("/", UpdateStudent(*db, c))
			r.Delete("/", DeleteStudent(*db, c))
		})
	})

	r.Route("/provisions", func(r chi.Router) {
		r.Post("/", CreateProvision(*db))
		r.Get("/", ReadProvisions(*db))

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", ReadProvision(*db))
			r.Put("/", UpdateProvision(*db))
			r.Delete("/", DeleteProvision(*db))
		})
	})

	r.Route("/room", func(r chi.Router) {
		r.Post("/", CreateRoom(*db))
		r.Get("/", ReadRooms(*db))

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", ReadRoom(*db))
			r.Put("/", UpdateRoom(*db))
			r.Delete("/", DeleteRoom(*db))
		})
	})

	r.Route("/worker", func(r chi.Router) {
		r.Post("/", CreateWorker(*db))
		r.Get("/", ReadWorkers(*db))

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", ReadWorker(*db))
			r.Put("/", UpdateWorker(*db))
			r.Delete("/", DeleteWorker(*db))
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

func convertToResponse(s Student) StudentResponse {
	return StudentResponse{
		ID:        s.ID,
		FName:     s.Firstname,
		LName:     s.Lastname,
		Date:      s.Date.In(location).Format(time.RFC1123),
		UUID:      s.UUID,
		StudRoom:  StudRoom{s.StudRoom.ID, s.StudRoom.Room},
		StudFloor: StudFloor{s.StudFloor.ID, s.StudFloor.Floor},
	}
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

//CreateStudent ...
func CreateStudent(d Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, valid := d.CheckAuth(&r.Header)
		if !valid {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": invalid creds", http.StatusInternalServerError)
			return
		}

		var st Student
		if err := json.NewDecoder(r.Body).Decode(&st); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest)+": eror unmarshalling json", http.StatusBadRequest)
			return
		}

		now := time.Now().UTC()

		if userID != st.UUID && userID != 0 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(http.StatusText(http.StatusUnauthorized) + ": you dont have access to this resource"))
			return
		}

		sroom, err := json.Marshal(st.StudRoom)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not add new student. Try changing id", http.StatusInternalServerError)
			return
		}

		sfloor, err := json.Marshal(st.StudFloor)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not add new student. Try changing id", http.StatusInternalServerError)
			return
		}

		if _, err := d.db.Exec("insert into student(id, firstname, lastname, date, uuid, studroom, studfloor) values($1, $2, $3, $4, $5, $6, $7)", st.ID, st.Firstname, st.Lastname, now, st.UUID, string(sroom), string(sfloor)); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not add new student. Try changing id", http.StatusInternalServerError)
			return
		}
		w.Write([]byte(http.StatusText(http.StatusOK) + ": created student"))
	}
}

//ReadStudents ...
func ReadStudents(d Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := d.db.Query("select id, firstname, lastname, date, uuid ,studroom, studfloor from student order by id asc")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not list students", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		sroom, sfloor := "", ""

		stds := []Student{}
		for rows.Next() {

			std := Student{}

			if err := rows.Scan(&std.ID, &std.Firstname, &std.Lastname, &std.Date, &std.UUID, &sroom, &sfloor); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not scan db", http.StatusInternalServerError)
				return
			}

			if err = json.Unmarshal([]byte(sroom), &std.StudRoom); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not scan db", http.StatusInternalServerError)
				return
			}

			if err = json.Unmarshal([]byte(sfloor), &std.StudFloor); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not scan db", http.StatusInternalServerError)
				return
			}

			stds = append(stds, std)
		}
		var resps []StudentResponse
		for _, s := range stds {
			resp := convertToResponse(s)
			resps = append(resps, resp)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resps); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest)+": error marshalling json", http.StatusBadRequest)
			return
		}
	}
}

func (r *rediscache) Get(id string) (Student, error) {
	var s Student
	val, err := r.redis.Get(id).Result()
	if err == redis.Nil || err != nil {
		return s, err
	}

	if err := json.Unmarshal([]byte(val), &s); err != nil {
		return s, err
	}

	return s, nil
}

func (r *rediscache) Set(id string, s *Student) error {
	b, err := json.Marshal(&s)
	if err != nil {
		return err
	}

	return r.redis.Set(id, string(b), time.Hour).Err()
}

func (r *rediscache) Remove(id string) error {
	return r.redis.Del(id).Err()
}

//ReadStudent ...
func ReadStudent(d Database, c Cache) func(w http.ResponseWriter, r *http.Request) {
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

		st, err := c.Get(id)

		sroom, sfloor := "", ""

		if err != nil {
			row := d.db.QueryRow("select id, firstname, lastname, date, uuid ,studroom, studfloor from student where id = $1", idInt)

			err = row.Scan(&st.ID, &st.Firstname, &st.Lastname, &st.Date, &st.UUID, &sroom, &sfloor)
			switch {
			case err == sql.ErrNoRows:
				http.NotFound(w, r)
				return
			case err != nil:
				http.Error(w, http.StatusText(http.StatusInternalServerError)+": error scanning db", http.StatusInternalServerError)
				return
			}

			if err = json.Unmarshal([]byte(sroom), &st.StudRoom); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not scan db", http.StatusInternalServerError)
				return
			}

			if err = json.Unmarshal([]byte(sfloor), &st.StudFloor); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not scan db", http.StatusInternalServerError)
				return
			}
			c.Set(id, &st)

		}

		resp := convertToResponse(st)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest)+": Error marshalling json", http.StatusBadRequest)
			return
		}

	}
}

//UpdateStudent ...
func UpdateStudent(d Database, c Cache) func(w http.ResponseWriter, r *http.Request) {
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

		if userID != bkUserID && userID != 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": stop right there criminal scum!", http.StatusInternalServerError)
			return
		}

		var stReq StudentRequest
		if err := json.NewDecoder(r.Body).Decode(&stReq); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest)+": error unmarshalling json", http.StatusBadRequest)
			return
		}

		st := Student{
			ID:        stReq.ID,
			Firstname: stReq.FName,
			Lastname:  stReq.LName,
			StudRoom:  stReq.StudRoom,
			StudFloor: stReq.StudFloor,
		}

		sroom, err := json.Marshal(st.StudRoom)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not add new student. Try changing id", http.StatusInternalServerError)
			return
		}

		sfloor, err := json.Marshal(st.StudFloor)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not add new student. Try changing id", http.StatusInternalServerError)
			return
		}

		if _, err = d.db.Exec("update student set id = $1, firstname = $2, lastname = $3, studroom = $4, studfloor = $5 where id = $1", idInt, st.Firstname, st.Lastname, string(sroom), string(sfloor)); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": cant update student", http.StatusInternalServerError)
			return
		}

		if err := c.Remove(id); err != nil {
			return
		}

		w.Write([]byte(http.StatusText(http.StatusOK) + ": updated student"))
	}
}

//DeleteStudent ...
func DeleteStudent(d Database, c Cache) func(w http.ResponseWriter, r *http.Request) {
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

		if userID != bkUserID && userID != 0 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(http.StatusText(http.StatusUnauthorized) + ": you dont have access to this resource"))
			return
		}

		var st Student
		if err := json.NewDecoder(r.Body).Decode(&st); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest)+": Error unmarshalling json", http.StatusBadRequest)
			return
		}

		if _, err = d.db.Exec("delete from student where id = $1", idInt); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not delete student", http.StatusInternalServerError)
			return
		}

		if err := c.Remove(id); err != nil {
			return
		}

		w.Write([]byte(http.StatusText(http.StatusOK) + ": deleted student"))
	}
}
