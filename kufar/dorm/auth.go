package main

import (
	"database/sql"
	"encoding/base64"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

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
