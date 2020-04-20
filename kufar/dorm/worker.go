package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

//Worker ...
type Worker struct {
	ID        int
	Firstname string
	Lastname  string
	WorkFloor WorkFloor
	WorkDays  string
}

//WorkFloor ...
type WorkFloor struct {
	ID    int
	Floor Floor
}

//Floor ...
type Floor struct {
	Floor, Code int
}

// WorkerResReq ...
type WorkerResReq struct {
	ID        int       `json:"id"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	WorkFloor WorkFloor `json:"workfloor"`
	WorkDays  string    `json:"workdays"`
}

//CreateWorker ...
func CreateWorker(d Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, valid := d.CheckAuth(&r.Header)
		if !valid {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": invalid creds", http.StatusInternalServerError)
			return
		}

		var wk Worker
		if err := json.NewDecoder(r.Body).Decode(&wk); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest)+": eror unmarshalling json", http.StatusBadRequest)
			return
		}

		if userID != wk.ID && userID != 0 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(http.StatusText(http.StatusUnauthorized) + ": you dont have access to this resource"))
			return
		}

		workfloor, err := json.Marshal(wk.WorkFloor)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not add new worker. Try changing id", http.StatusInternalServerError)
			return
		}

		if _, err := d.db.Exec("insert into worker(id, firstname, lastname, workfloor, workdays) values($1, $2, $3, $4, $5)", wk.ID, wk.Firstname, wk.Lastname, string(workfloor), wk.WorkDays); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not add new provision. Try changing id", http.StatusInternalServerError)
			return
		}
		w.Write([]byte(http.StatusText(http.StatusOK) + ": created worker"))

	}
}

//ReadWorkers ...
func ReadWorkers(d Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := d.db.Query("select id, firstname, lastname, workfloor, workdays from worker order by id asc")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not list students", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		workfloor := ""

		wks := []Worker{}
		for rows.Next() {

			wk := Worker{}

			if err := rows.Scan(&wk.ID, &wk.Firstname, &wk.Lastname, &workfloor, &wk.WorkDays); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not scan db", http.StatusInternalServerError)
				return
			}

			if err = json.Unmarshal([]byte(workfloor), &wk.WorkFloor); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not scan db", http.StatusInternalServerError)
				return
			}

			wks = append(wks, wk)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(&wks); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest)+": error marshalling json", http.StatusBadRequest)
			return
		}
	}
}

//ReadWorker ...
func ReadWorker(d Database) func(w http.ResponseWriter, r *http.Request) {
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

		wk := Worker{}

		workfloor := ""

		row := d.db.QueryRow("select id, firstname, lastname, workfloor, workdays from worker where id = $1", idInt)

		err = row.Scan(&wk.ID, &wk.Firstname, &wk.Lastname, &workfloor, &wk.WorkDays)
		switch {
		case err == sql.ErrNoRows:
			http.NotFound(w, r)
			return
		case err != nil:
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": error scanning db", http.StatusInternalServerError)
			return
		}

		if err = json.Unmarshal([]byte(workfloor), &wk.WorkFloor); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not scan db", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(wk); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest)+": Error marshalling json", http.StatusBadRequest)
			return
		}

	}
}

//UpdateWorker ...
func UpdateWorker(d Database) func(w http.ResponseWriter, r *http.Request) {
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

		wkReq := WorkerResReq{}
		if err := json.NewDecoder(r.Body).Decode(&wkReq); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest)+": error unmarshalling json", http.StatusBadRequest)
			return
		}

		wk := Worker{
			ID:        wkReq.ID,
			Firstname: wkReq.Firstname,
			Lastname:  wkReq.Lastname,
			WorkFloor: WorkFloor{
				ID: wkReq.WorkFloor.ID,
				Floor: Floor{
					Floor: wkReq.WorkFloor.Floor.Floor,
					Code:  wkReq.WorkFloor.Floor.Code,
				},
			},
			WorkDays: wkReq.WorkDays,
		}

		workfloor, err := json.Marshal(wk.WorkFloor)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not add new student. Try changing id", http.StatusInternalServerError)
			return
		}

		if _, err = d.db.Exec("update worker set id = $1, firstname = $2, lastname = $3, workfloor = $4, workdays = $5 where id = $1", idInt, wk.Firstname, wk.Lastname, string(workfloor), wk.WorkDays); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": cant update student", http.StatusInternalServerError)
			return
		}

		w.Write([]byte(http.StatusText(http.StatusOK) + ": updated worker"))
	}
}

//DeleteWorker ...
func DeleteWorker(d Database) func(w http.ResponseWriter, r *http.Request) {
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

		var wk Worker
		if err := json.NewDecoder(r.Body).Decode(&wk); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest)+": Error unmarshalling json", http.StatusBadRequest)
			return
		}

		if _, err = d.db.Exec("delete from worker where id = $1", idInt); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not delete student", http.StatusInternalServerError)
			return
		}

		w.Write([]byte(http.StatusText(http.StatusOK) + ": deleted worker"))

	}
}
