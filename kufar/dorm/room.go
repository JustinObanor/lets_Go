package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

//RoomItems ...
type RoomItems struct {
	ID      int
	Room    int
	Chairs  int
	Tables  int
	Shelves int
}

//RoomItemsReqRes ...
type RoomItemsReqRes struct {
	ID      int `json:"ID"`
	Room    int `json:"Room"`
	Chairs  int `json:"Chair"`
	Tables  int `json:"Table"`
	Shelves int `json:"Shelve"`
}

//CreateRoom ...
func CreateRoom(d Database) func(w http.ResponseWriter, r *http.Request) {
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

		var rm RoomItems
		if err := json.NewDecoder(r.Body).Decode(&rm); err != nil {
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
				Message: "could not admin id",
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
				Message: "could not worker id",
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

		if _, err := d.db.Exec("insert into room(room, chairs, tables, shelves) values($1, $2, $3, $4)", rm.Room, rm.Chairs, rm.Tables, rm.Shelves); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res := Response{
				Status:  http.StatusInternalServerError,
				Message: "could not add new room",
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
			Message: "created room",
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest)+": error marshalling json", http.StatusBadRequest)
			return
		}
	}
}

//ReadRooms ...
func ReadRooms(d Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := d.db.Query("select id, room, chairs, tables, shelves from room order by id asc")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res := Response{
				Status:  http.StatusInternalServerError,
				Message: "could not list rooms",
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

		rms := []RoomItems{}
		for rows.Next() {

			rm := RoomItems{}

			if err := rows.Scan(&rm.ID, &rm.Room, &rm.Chairs, &rm.Tables, &rm.Shelves); err != nil {
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

			rms = append(rms, rm)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(&rms); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest)+": error marshalling json", http.StatusBadRequest)
			return
		}
	}
}

//ReadRoom ...
func ReadRoom(d Database) func(w http.ResponseWriter, r *http.Request) {
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

		rm := RoomItems{}

		row := d.db.QueryRow("select id, room, chairs, tables, shelves from room where id = $1", idInt)

		err = row.Scan(&rm.ID, &rm.Room, &rm.Chairs, &rm.Tables, &rm.Shelves)
		switch {
		case err == sql.ErrNoRows:
			w.WriteHeader(http.StatusNotFound)
			res := Response{
				Status:  http.StatusNotFound,
				Message: "no such room",
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
		if err := json.NewEncoder(w).Encode(rm); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest)+": Error marshalling json", http.StatusBadRequest)
			return
		}
	}
}

//UpdateRoom ...
func UpdateRoom(d Database) func(w http.ResponseWriter, r *http.Request) {
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

		rmReq := RoomItemsReqRes{}
		if err := json.NewDecoder(r.Body).Decode(&rmReq); err != nil {
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

		rm := RoomItems{
			Room:    rmReq.Room,
			Chairs:  rmReq.Chairs,
			Tables:  rmReq.Tables,
			Shelves: rmReq.Shelves,
		}

		if _, err = d.db.Exec("update room set id = $1, room = $2, chairs = $3, tables = $4, shelves = $5 where id = $1", idInt, rm.Room, rm.Chairs, rm.Tables, rm.Shelves); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res := Response{
				Status:  http.StatusInternalServerError,
				Message: "cant update room",
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
			Message: "updated room",
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest)+": error marshalling json", http.StatusBadRequest)
			return
		}
	}
}

//DeleteRoom ...
func DeleteRoom(d Database) func(w http.ResponseWriter, r *http.Request) {
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

		if _, err = d.db.Exec("delete from room where id = $1", idInt); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res := Response{
				Status:  http.StatusInternalServerError,
				Message: "cant delete room",
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
			Message: "deleted room of id " + id,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest)+": error marshalling json", http.StatusBadRequest)
			return
		}
	}
}
