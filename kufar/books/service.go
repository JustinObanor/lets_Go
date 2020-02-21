package main

import "net/http"

func Create(d Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		d.CreateBook(w, r)
	}
}

func Read(d Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		d.ReadBook(w, r)
	}
}

func ReadAll(d Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		d.ReadBooks(w, r)
	}
}

func Update(d Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		d.UpdateBook(w, r)
	}
}

func Delete(d Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		d.DeleteBook(w, r)
	}
}
