package main

import (
	"net/http"
)

//SignUp ...
func SignUp(d Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		SignUpUser(d)
	}
}

//Create ...
func Create(d Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		CreateBook(d)
	}
}

//Read ...
func Read(d Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ReadBook(d)
	}
}

//ReadAll ...
func ReadAll(d Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ReadBooks(d)
	}
}

//Update ...
func Update(d Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		UpdateBook(d)
	}
}

//Delete ...
func Delete(d Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		DeleteBook(d)
	}
}
