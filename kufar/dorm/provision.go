package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

//StudProvisions ...
type StudProvisions struct {
	ID       int
	Bedsheet int
	Pillow   int
	Towel    int
	Blanket  int
	Curtain  int
}

//StudProvisionsReqRes ...
type StudProvisionsReqRes struct {
	ID       int `json:"id"`
	Bedsheet int `json:"bedsheet"`
	Pillow   int `json:"pillow"`
	Towel    int `json:"towel"`
	Blanket  int `json:"blanket"`
	Curtain  int `json:"curtain"`
}

//CreateProvision ...
func CreateProvision(d Database) func(w http.ResponseWriter, r *http.Request) {
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

		var pr StudProvisions
		if err := json.NewDecoder(r.Body).Decode(&pr); err != nil {
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

		if _, err := d.db.Exec("insert into provisions(id, bedsheet, pillow, towel, blanket, curtain) values($1, $2, $3, $4, $5, $6)",pr.ID, pr.Bedsheet, pr.Pillow, pr.Towel, pr.Blanket, pr.Curtain); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not add new provision. Try changing id", http.StatusInternalServerError)
			return
		}

		res := Response{
			Status:  http.StatusOK,
			Message: "created provisions",
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest)+": error marshalling json", http.StatusBadRequest)
			return
		}
	}
}

//ReadProvisions ...
func ReadProvisions(d Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := d.db.Query("select id, bedsheet, pillow, towel, blanket, curtain from provisions order by id asc")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res := Response{
				Status:  http.StatusInternalServerError,
				Message: "could not list provisions",
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

		prs := []StudProvisions{}
		for rows.Next() {

			pr := StudProvisions{}

			if err := rows.Scan(&pr.ID, &pr.Bedsheet, &pr.Pillow, &pr.Towel, &pr.Blanket, &pr.Curtain); err != nil {
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

			prs = append(prs, pr)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(&prs); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest)+": error marshalling json", http.StatusBadRequest)
			return
		}
	}
}

//ReadProvision ...
func ReadProvision(d Database) func(w http.ResponseWriter, r *http.Request) {
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

		pr := StudProvisions{}

		row := d.db.QueryRow("select id, bedsheet, pillow, towel, blanket, curtain from provisions where id = $1", idInt)

		err = row.Scan(&pr.ID, &pr.Bedsheet, &pr.Pillow, &pr.Towel, &pr.Blanket, &pr.Curtain)
		switch {
		case err == sql.ErrNoRows:
			w.WriteHeader(http.StatusNotFound)
			res := Response{
				Status:  http.StatusNotFound,
				Message: "no such provision",
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

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(pr); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest)+": Error marshalling json", http.StatusBadRequest)
			return
		}

	}
}

//UpdateProvision ...
func UpdateProvision(d Database) func(w http.ResponseWriter, r *http.Request) {
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

		if WorkerID != userID && AdminID != userID {
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

		prReq := StudProvisionsReqRes{}
		if err := json.NewDecoder(r.Body).Decode(&prReq); err != nil {
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

		st := StudProvisions{
			ID:       prReq.ID,
			Bedsheet: prReq.Bedsheet,
			Pillow:   prReq.Pillow,
			Towel:    prReq.Towel,
			Blanket:  prReq.Blanket,
			Curtain:  prReq.Curtain,
		}

		if _, err = d.db.Exec("update provisions set id = $1, bedsheet = $2, pillow = $3, towel = $4, blanket = $5, curtain = $6 where id = $1", idInt, st.Bedsheet, st.Pillow, st.Towel, st.Blanket, st.Curtain); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res := Response{
				Status:  http.StatusInternalServerError,
				Message: "cant update provision",
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
			Message: "updated provoision",
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest)+": error marshalling json", http.StatusBadRequest)
			return
		}
	}
}

//DeleteProvision ...
func DeleteProvision(d Database) func(w http.ResponseWriter, r *http.Request) {
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

		if WorkerID != userID && AdminID != userID {
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

		if _, err = d.db.Exec("delete from provisions where id = $1", idInt); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not delete student", http.StatusInternalServerError)
			return
		}

		w.Write([]byte(http.StatusText(http.StatusOK) + ": deleted student"))

	}
}
