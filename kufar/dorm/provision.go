package main

import (
	"database/sql"
	"encoding/json"
	"log"
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
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": invalid creds", http.StatusInternalServerError)
			return
		}

		var pr StudProvisions
		if err := json.NewDecoder(r.Body).Decode(&pr); err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusBadRequest)+": eror unmarshalling json", http.StatusBadRequest)
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

		if userID != pr.ID && AdminID != userID && WorkerID != userID {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(http.StatusText(http.StatusUnauthorized) + ": you dont have access to this resource"))
			return
		}

		if _, err := d.db.Exec("insert into provisions(bedsheet, pillow, towel, blanket, curtain) values($1, $2, $3, $4, $5)", pr.Bedsheet, pr.Pillow, pr.Towel, pr.Blanket, pr.Curtain); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not add new provision. Try changing id", http.StatusInternalServerError)
			return
		}
		w.Write([]byte(http.StatusText(http.StatusOK) + ": created provision"))

	}
}

//ReadProvisions ...
func ReadProvisions(d Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := d.db.Query("select id, bedsheet, pillow, towel, blanket, curtain from provisions order by id asc")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not list provisions", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		prs := []StudProvisions{}
		for rows.Next() {

			pr := StudProvisions{}

			if err := rows.Scan(&pr.ID, &pr.Bedsheet, &pr.Pillow, &pr.Towel, &pr.Blanket, &pr.Curtain); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not scan db", http.StatusInternalServerError)
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
			http.Error(w, http.StatusText(http.StatusBadRequest)+": missing parameter in url", http.StatusBadRequest)
			return
		}

		idInt, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not convert to integer", http.StatusInternalServerError)
			return
		}

		pr := StudProvisions{}

		row := d.db.QueryRow("select id, bedsheet, pillow, towel, blanket, curtain from provisions where id = $1", idInt)

		err = row.Scan(&pr.ID, &pr.Bedsheet, &pr.Pillow, &pr.Towel, &pr.Blanket, &pr.Curtain)
		switch {
		case err == sql.ErrNoRows:
			http.NotFound(w, r)
			return
		case err != nil:
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": error scanning db", http.StatusInternalServerError)
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

		if userID != bkUserID && WorkerID != userID && AdminID != userID {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": stop right there criminal scum!", http.StatusInternalServerError)
			return
		}

		prReq := StudProvisionsReqRes{}
		if err := json.NewDecoder(r.Body).Decode(&prReq); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest)+": error unmarshalling json", http.StatusBadRequest)
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
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": cant update student", http.StatusInternalServerError)
			return
		}

		w.Write([]byte(http.StatusText(http.StatusOK) + ": updated student"))
	}
}

//DeleteProvision ...
func DeleteProvision(d Database) func(w http.ResponseWriter, r *http.Request) {
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

		if userID != bkUserID && WorkerID != userID && AdminID != userID {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": stop right there criminal scum!", http.StatusInternalServerError)
			return
		}

		var pr StudProvisions
		if err := json.NewDecoder(r.Body).Decode(&pr); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest)+": Error unmarshalling json", http.StatusBadRequest)
			return
		}

		if _, err = d.db.Exec("delete from provisions where id = $1", idInt); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not delete student", http.StatusInternalServerError)
			return
		}

		w.Write([]byte(http.StatusText(http.StatusOK) + ": deleted student"))

	}
}
