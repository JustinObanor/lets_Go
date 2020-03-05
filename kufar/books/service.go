package main

import (
	"net/http"
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

var secrets map[Credentials]BookResponse

//Update ...
func Update(d Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// r.SetBasicAuth(cred.User, cred.Password)

		// realm := "Access to the user books"
		// if _, ok := secrets[cred]; !ok {
		// 	w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`," charset="UTF-8"`)
		// 	w.WriteHeader(http.StatusUnauthorized)
		// 	w.Write([]byte(http.StatusText(http.StatusUnauthorized) + ": check credentials"))
		// 	return
		// }

		// cred.User, cred.Password, _ = r.BasicAuth()

		d.UpdateBook(w, r)
	}
}

//Delete ...
func Delete(d Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		d.DeleteBook(w, r)
	}
}
