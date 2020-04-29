package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-redis/redis"

	_ "github.com/lib/pq"
)

const (
	cost = 12
)

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
			http.Error(w, http.StatusText(http.StatusBadRequest)+": error unmarshalling json", http.StatusBadRequest)
			return
		}

		now := time.Now().UTC()

		AdminID, err := d.getCredUUID("admin")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": cant get id", http.StatusInternalServerError)
			return
		}

		WorkerID, err := d.getCredUUID("worker")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": cant get id", http.StatusInternalServerError)
			return
		}

		if userID != st.UUID && userID != AdminID && userID != WorkerID {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": stop right there criminal scum!", http.StatusInternalServerError)
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

		if _, err := d.db.Exec("insert into student(firstname, lastname, date, uuid, studroom, studfloor) values($1, $2, $3, $4, $5, $6)", st.Firstname, st.Lastname, now, st.UUID, string(sroom), string(sfloor)); err != nil {
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

		AdminID, err := d.getCredUUID("admin")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": cant get id", http.StatusInternalServerError)
			return
		}

		WorkerID, err := d.getCredUUID("worker")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": cant get id", http.StatusInternalServerError)
			return
		}

		if userID != bkUserID && userID != AdminID && userID != WorkerID {
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
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": cant verify authority", http.StatusInternalServerError)
			return
		}

		AdminID, err := d.getCredUUID("admin")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": cant get id", http.StatusInternalServerError)
			return
		}

		WorkerID, err := d.getCredUUID("worker")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": cant get id", http.StatusInternalServerError)
			return
		}

		if userID != bkUserID && userID != AdminID && userID != WorkerID {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": stop right there criminal scum!", http.StatusInternalServerError)
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
