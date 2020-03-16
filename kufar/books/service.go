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
		d.UpdateBook(w, r)
	}
}

//Delete ...
func Delete(d Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		d.DeleteBook(w, r)
	}
}
