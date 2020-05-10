package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

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

//CredentialsResponse ...
type CredentialsResponse struct {
	Message  int    `json:"message"`
	UUID     int    `json:"uuid"`
	Username string `json:"username"`
	Password string `json:"password"`
}

//Login ...
type Login struct {
	Login   string `json:"login"`
	Message int    `json:"message"`
	Token   string `json:"token"`
	Role    string `json:"role"`
	Rights  bool   `json:"rights"`
	UUID    int    `json:"uuid"`
}

// StupidMiddleware ...
func StupidMiddleware(id string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if id != "0" {
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (d Database) getBookUserID(bookID int) (int, error) {
	row := d.db.QueryRow("select uuid from student where id = $1", bookID)

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

func (d Database) getStudentID(fname,lname string) (int, error) {
	row := d.db.QueryRow("select id from student where firstname = $1 and lastname = $2", fname, lname)

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

func (d Database) getCredUUID(uname string) (int, error) {
	row := d.db.QueryRow("select uuid from credentials where username = $1", uname)

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

		if credReq.Username == "" && credReq.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			res := Response{
				Status:  http.StatusBadRequest,
				Message: "missing username and/or password",
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest)+": error marshalling json", http.StatusBadRequest)
				return
			}
			return
		}

		pword, err := bcrypt.GenerateFromPassword([]byte(credReq.Password), cost)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res := Response{
				Status:  http.StatusInternalServerError,
				Message: "could not add password",
				Error:   err,
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest)+": error marshalling json", http.StatusBadRequest)
				return
			}
			return
		}

		cred := Credentials{Username: credReq.Username}
		if _, err = d.db.Exec("insert into credentials(username, password) values($1, $2)", cred.Username, string(pword)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res := Response{
				Status:  http.StatusInternalServerError,
				Message: "could not log in user " + cred.Username,
				Error:   err,
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest)+": error marshalling json", http.StatusBadRequest)
				return
			}
			return
		}

		id, err := d.getCredUUID(cred.Username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res := Response{
				Status:  http.StatusInternalServerError,
				Message: "could not get user id",
				Error:   err,
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError)+": error marshalling json", http.StatusBadRequest)
				return
			}
			return
		}

		credRes := CredentialsResponse{
			Message:  http.StatusOK,
			UUID:     id,
			Username: credReq.Username,
			Password: credReq.Password,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(credRes); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest)+": error marshalling json", http.StatusBadRequest)
			return
		}
	}
}

//LogIn ...
func LogIn(d Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		credReq := CredentialsRequest{}
		if err := json.NewDecoder(r.Body).Decode(&credReq); err != nil {
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

		if credReq.Username == "" && credReq.Password == "" {
			w.WriteHeader(http.StatusBadRequest)
			res := Response{
				Status:  http.StatusBadRequest,
				Message: "missing username and/or password",
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest)+": error marshalling json", http.StatusBadRequest)
				return
			}
			return
		}

		row := d.db.QueryRow("select username, password from credentials where username = $1", credReq.Username)

		dbCred := Credentials{}
		err := row.Scan(&dbCred.Username, &dbCred.Password)

		switch {
		case err == sql.ErrNoRows:
			w.WriteHeader(http.StatusNotFound)
			res := Response{
				Status:  http.StatusNotFound,
				Message: "no such user",
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

		if err = bcrypt.CompareHashAndPassword([]byte(dbCred.Password), []byte(credReq.Password)); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			res := Response{
				Status:  http.StatusUnauthorized,
				Message: "wrong creds",
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest)+": error marshalling json", http.StatusBadRequest)
				return
			}
			return
		}

		var b strings.Builder
		b.WriteString(credReq.Username)
		b.WriteString(":")
		b.WriteString(credReq.Password)

		token := base64.StdEncoding.EncodeToString([]byte(b.String()))

		b.Reset()
		b.WriteString("Basic ")
		b.WriteString(token)

		var login Login

		id, err := d.getCredUUID(dbCred.Username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			res := Response{
				Status:  http.StatusInternalServerError,
				Message: "could not get user id",
				Error:   err,
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError)+": error marshalling json", http.StatusBadRequest)
				return
			}
			return
		}

		login.Login = dbCred.Username
		login.Message = http.StatusOK
		login.Token = b.String()

		if dbCred.Username == "admin" {
			login.Role += "admin"
			login.Rights = true
		} else if dbCred.Username == "worker" {
			login.Role += "worker"
			login.Rights = true
		} else {
			login.Role += "student"
			login.Rights = false
		}
		login.UUID = id

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(login); err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest)+": error marshalling json", http.StatusBadRequest)
			return
		}
	}
}
