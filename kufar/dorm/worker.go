package main

import (
	"database/sql"
	"encoding/json"
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
			w.WriteHeader(http.StatusUnauthorized)
			res := Response{
				Status:  http.StatusUnauthorized,
				Message: "invalid creds",
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest)+": error marshalling json", http.StatusBadRequest)
				return
			}
			return
		}

		var wk Worker
		if err := json.NewDecoder(r.Body).Decode(&wk); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			res := Response{
				Status:  http.StatusBadRequest,
				Message: "error marshalling json",
				Error:   err,
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest)+": error marshalling json", http.StatusBadRequest)
				return
			}
			return
		}

		AdminID, err := d.getCredUUID("admin")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res := Response{
				Status:  http.StatusInternalServerError,
				Message: "cant get admin id",
				Error:   err,
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest)+": error marshalling json", http.StatusBadRequest)
				return
			}
			return
		}

		WorkerID, err := d.getCredUUID("worker")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res := Response{
				Status:  http.StatusInternalServerError,
				Message: "could not get worker id",
				Error:   err,
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError)+": error marshalling json", http.StatusBadRequest)
				return
			}
			return
		}

		if AdminID != userID && WorkerID != userID {
			w.WriteHeader(http.StatusUnauthorized)
			res := Response{
				Status:  http.StatusUnauthorized,
				Message: "stop right there criminal scum!",
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest)+": error marshalling json", http.StatusBadRequest)
				return
			}
			return
		}

		workfloor, err := json.Marshal(wk.WorkFloor)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res := Response{
				Status:  http.StatusInternalServerError,
				Message: "could not add new worker",
				Error:   err,
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest)+": error marshalling json", http.StatusBadRequest)
				return
			}
			return
		}

		if _, err := d.db.Exec("insert into worker(firstname, lastname, workfloor, workdays) values($1, $2, $3, $4)", wk.Firstname, wk.Lastname, string(workfloor), wk.WorkDays); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res := Response{
				Status:  http.StatusInternalServerError,
				Message: "could not add new worker",
				Error:   err,
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest)+": error marshalling json", http.StatusBadRequest)
				return
			}
			return
		}

		res := Response{
			Status:  http.StatusOK,
			Message: "created worker",
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest)+": error marshalling json", http.StatusBadRequest)
			return
		}
	}
}

//ReadWorkers ...
func ReadWorkers(d Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := d.db.Query("select id, firstname, lastname, workfloor, workdays from worker order by id asc")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res := Response{
				Status:  http.StatusInternalServerError,
				Message: "could not list workers",
				Error:   err,
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest)+": error marshalling json", http.StatusBadRequest)
				return
			}
			return
		}
		defer rows.Close()

		workfloor := ""

		wks := []Worker{}
		for rows.Next() {

			wk := Worker{}

			if err := rows.Scan(&wk.ID, &wk.Firstname, &wk.Lastname, &workfloor, &wk.WorkDays); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				res := Response{
					Status:  http.StatusInternalServerError,
					Message: "could not scan db",
					Error:   err,
				}

				w.Header().Set("Content-Type", "application/json")
				if err := json.NewEncoder(w).Encode(res); err != nil {
					http.Error(w, http.StatusText(http.StatusBadRequest)+": error marshalling json", http.StatusBadRequest)
					return
				}
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
			w.WriteHeader(http.StatusBadRequest)
			res := Response{
				Status:  http.StatusBadRequest,
				Message: "missing parameter in url",
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest)+": error marshalling json", http.StatusBadRequest)
				return
			}
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res := Response{
				Status:  http.StatusInternalServerError,
				Message: "could not convert to integer",
				Error:   err,
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError)+": error marshalling json", http.StatusBadRequest)
				return
			}
			return
		}

		wk := Worker{}

		workfloor := ""

		row := d.db.QueryRow("select id, firstname, lastname, workfloor, workdays from worker where id = $1", idInt)

		err = row.Scan(&wk.ID, &wk.Firstname, &wk.Lastname, &workfloor, &wk.WorkDays)
		switch {
		case err == sql.ErrNoRows:
			w.WriteHeader(http.StatusNotFound)
			res := Response{
				Status:  http.StatusNotFound,
				Message: "no such worker",
				Error:   err,
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError)+": error marshalling json", http.StatusBadRequest)
				return
			}
			return

		case err != nil:
			w.WriteHeader(http.StatusBadRequest)
			res := Response{
				Status:  http.StatusBadRequest,
				Message: "could not scan db",
				Error:   err,
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError)+": error marshalling json", http.StatusBadRequest)
				return
			}
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
			w.WriteHeader(http.StatusUnauthorized)
			res := Response{
				Status:  http.StatusUnauthorized,
				Message: "invalid creds",
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest)+": error marshalling json", http.StatusBadRequest)
				return
			}
			return
		}

		id := chi.URLParam(r, "id")

		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			res := Response{
				Status:  http.StatusBadRequest,
				Message: "missing parameter in url",
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest)+": error marshalling json", http.StatusBadRequest)
				return
			}
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res := Response{
				Status:  http.StatusInternalServerError,
				Message: "could not convert to integer",
				Error:   err,
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError)+": error marshalling json", http.StatusBadRequest)
				return
			}
			return
		}

		AdminID, err := d.getCredUUID("admin")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res := Response{
				Status:  http.StatusInternalServerError,
				Message: "could not get admin id",
				Error:   err,
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError)+": error marshalling json", http.StatusBadRequest)
				return
			}
			return
		}

		WorkerID, err := d.getCredUUID("worker")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res := Response{
				Status:  http.StatusInternalServerError,
				Message: "could not get worker id",
				Error:   err,
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError)+": error marshalling json", http.StatusBadRequest)
				return
			}
			return
		}

		if userID != AdminID && userID != WorkerID {
			w.WriteHeader(http.StatusUnauthorized)
			res := Response{
				Status:  http.StatusUnauthorized,
				Message: "stop right there criminal scum!",
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest)+": error marshalling json", http.StatusBadRequest)
				return
			}
			return
		}

		wkReq := WorkerResReq{}
		if err := json.NewDecoder(r.Body).Decode(&wkReq); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			res := Response{
				Status:  http.StatusBadRequest,
				Message: "error unmarshalling json",
				Error:   err,
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError)+": error marshalling json", http.StatusBadRequest)
				return
			}
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
			w.WriteHeader(http.StatusInternalServerError)
			res := Response{
				Status:  http.StatusInternalServerError,
				Message: "cant update worker",
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest)+": error marshalling json", http.StatusBadRequest)
				return
			}
			return
		}

		res := Response{
			Status:  http.StatusOK,
			Message: "updated worker",
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest)+": error marshalling json", http.StatusBadRequest)
			return
		}
	}
}

//DeleteWorker ...
func DeleteWorker(d Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, valid := d.CheckAuth(&r.Header)
		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			res := Response{
				Status:  http.StatusUnauthorized,
				Message: "invalid creds",
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest)+": error marshalling json", http.StatusBadRequest)
				return
			}
			return
		}

		id := chi.URLParam(r, "id")

		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			res := Response{
				Status:  http.StatusBadRequest,
				Message: "missing parameter in url",
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest)+": error marshalling json", http.StatusBadRequest)
				return
			}
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res := Response{
				Status:  http.StatusInternalServerError,
				Message: "could not convert to integer",
				Error:   err,
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError)+": error marshalling json", http.StatusBadRequest)
				return
			}
			return
		}

		AdminID, err := d.getCredUUID("admin")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res := Response{
				Status:  http.StatusInternalServerError,
				Message: "could not get admin id",
				Error:   err,
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError)+": error marshalling json", http.StatusBadRequest)
				return
			}
			return
		}

		WorkerID, err := d.getCredUUID("worker")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res := Response{
				Status:  http.StatusInternalServerError,
				Message: "could not get worker id",
				Error:   err,
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError)+": error marshalling json", http.StatusBadRequest)
				return
			}
			return
		}

		if userID != AdminID && userID != WorkerID {
			w.WriteHeader(http.StatusUnauthorized)
			res := Response{
				Status:  http.StatusUnauthorized,
				Message: "stop right there criminal scum!",
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest)+": error marshalling json", http.StatusBadRequest)
				return
			}
			return
		}

		if _, err = d.db.Exec("delete from worker where id = $1", idInt); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res := Response{
				Status:  http.StatusInternalServerError,
				Message: "cant delete worker",
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest)+": error marshalling json", http.StatusBadRequest)
				return
			}
			return
		}

		res := Response{
			Status:  http.StatusOK,
			Message: "deleted worker of id " + id,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest)+": error marshalling json", http.StatusBadRequest)
			return
		}
	}
}
