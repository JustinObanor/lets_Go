package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

//SignUp ...
func SignUp(d Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		d.SignUpUser(w, r)
	}
}

//SignIn ...
func SignIn(d Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		d.SignInUser(w, r)
	}
}

//Logout ...
func Logout(d Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		d.LogoutUser(w, r)
	}
}

//Create ...
func Create(d Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		d.CreateBook(w, r)
	}
}

//Read ...
func Read(d Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		d.ReadBook(w, r)
	}
}

//ReadAll ...
func ReadAll(d Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		d.ReadBooks(w, r)
	}
}

//Update ...
func Update(d Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//r.Header.Set("WWW-Authenticate", `Basic realm="Restricted"`)

		//username and password not gotten :(
		uname, pword, ok := r.BasicAuth()

		//cant the the gredentials :(
		b64 := r.Header.Get("Authorization")

		s := strings.Split(b64, " ")
		fmt.Println(s)
		sDec, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError)+": could not decode base64 string", http.StatusInternalServerError)
			return
		}

		cred := strings.Split(string(sDec), ":")
		if string(cred[0]) != uname && string(cred[1]) != pword {
			http.Error(w, http.StatusText(http.StatusUnauthorized)+" authorization failed! ", http.StatusUnauthorized)
			return
		}
		if string(cred[0]) == uname && string(cred[1]) == pword && ok != false {
			http.Error(w, http.StatusText(http.StatusOK)+":(correct creds) welcome! "+strings.Title(uname), http.StatusOK)
			return
		}

		d.UpdateBook(w, r)
	}
}

//Delete ...
func Delete(d Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		d.DeleteBook(w, r)
	}
}
