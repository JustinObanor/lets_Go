package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

var db *sql.DB
var tpl *template.Template

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://justin:1999@localhost/bookstore_go?sslmode=disable")
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("You connected to your database.")

	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

// export fields to templates
// fields changed to uppercase
type Book struct {
	Studentid    string
	Firstname    string
	Lastname     string
	Classcode    float32
	Roomnumber   string
	FeesToBePaid string
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/books", booksIndex)
	http.HandleFunc("/books/show", booksShow)
	http.HandleFunc("/books/create", booksCreateForm)
	http.HandleFunc("/books/create/process", booksCreateProcess)
	http.HandleFunc("/books/update", booksUpdateForm)
	http.HandleFunc("/books/update/process", booksUpdateProcess)
	http.HandleFunc("/books/delete/process", booksDeleteProcess)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/books", http.StatusSeeOther)
}

func booksIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer rows.Close()

	bks := make([]Book, 0)
	for rows.Next() {
		bk := Book{}
		err := rows.Scan(&bk.Studentid, &bk.Firstname, &bk.Lastname, &bk.Classcode, &bk.Roomnumber, &bk.FeesToBePaid) // order matters
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		bks = append(bks, bk)
	}
	if err = rows.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	tpl.ExecuteTemplate(w, "books.gohtml", bks)
}

func booksShow(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	studentid := r.FormValue("studentid")
	if studentid == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	row := db.QueryRow("SELECT * FROM books WHERE studentid = $1", studentid)

	bk := Book{}
	err := row.Scan(&bk.Studentid, &bk.Firstname, &bk.Lastname, &bk.Classcode, &bk.Roomnumber, &bk.FeesToBePaid)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	tpl.ExecuteTemplate(w, "show.gohtml", bk)
}

func booksCreateForm(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "create.gohtml", nil)
}

func booksCreateProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	// get form values
	bk := Book{}
	bk.Studentid = r.FormValue("studentid")
	bk.Firstname = r.FormValue("firstname")
	bk.Lastname = r.FormValue("lastname")
	p := r.FormValue("classcode")
	bk.Roomnumber = r.FormValue("roomnumber")
	bk.FeesToBePaid = r.FormValue("feestobepaid")

	// validate form values
	if bk.Studentid == "" || bk.Firstname == "" || bk.Lastname == "" || p == "" || bk.Roomnumber == "" || bk.FeesToBePaid == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	// convert form values
	f64, err := strconv.ParseFloat(p, 32)
	if err != nil {
		http.Error(w, http.StatusText(406)+"Please hit back and enter a number for the classcode", http.StatusNotAcceptable)
		return
	}
	bk.Classcode = float32(f64)

	// insert values
	_, err = db.Exec("INSERT INTO books (studentid, firstname, lastname, classcode, roomnumber, feestobepaid) VALUES ($1, $2, $3, $4, $5, $6)", bk.Studentid, bk.Firstname, bk.Lastname, bk.Classcode, bk.Roomnumber, bk.FeesToBePaid)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	// confirm insertion
	tpl.ExecuteTemplate(w, "created.gohtml", bk)
}

func booksUpdateForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	studentid := r.FormValue("studentid")
	if studentid == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	row := db.QueryRow("SELECT * FROM books WHERE studentid = $1", studentid)

	bk := Book{}
	err := row.Scan(&bk.Studentid, &bk.Firstname, &bk.Lastname, &bk.Classcode, &bk.Roomnumber, &bk.FeesToBePaid)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	tpl.ExecuteTemplate(w, "update.gohtml", bk)
}

func booksUpdateProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	// get form values
	bk := Book{}
	bk.Studentid = r.FormValue("studentid")
	bk.Firstname = r.FormValue("firstname")
	bk.Lastname = r.FormValue("lastname")
	p := r.FormValue("classcode")
	bk.Roomnumber = r.FormValue("roomnumber")
	bk.FeesToBePaid = r.FormValue("feestobepaid")

	// validate form values
	if bk.Studentid == "" || bk.Firstname == "" || bk.Lastname == "" || p == "" || bk.Roomnumber == "" || bk.FeesToBePaid == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	// convert form values
	f64, err := strconv.ParseFloat(p, 32)
	if err != nil {
		http.Error(w, http.StatusText(406)+"Please hit back and enter a number for the classcode", http.StatusNotAcceptable)
		return
	}
	bk.Classcode = float32(f64)

	// insert values
	_, err = db.Exec("UPDATE books SET studentid = $1, firstname=$2, lastname=$3, classcode=$4, roomnumber=$5, feestobepaid=$6 WHERE studentid=$1;", bk.Studentid, bk.Firstname, bk.Lastname, bk.Classcode, bk.Roomnumber, bk.FeesToBePaid)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	// confirm insertion
	tpl.ExecuteTemplate(w, "updated.gohtml", bk)
}

func booksDeleteProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	studentid := r.FormValue("studentid")
	if studentid == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	// delete book
	_, err := db.Exec("DELETE FROM books WHERE studentid=$1;", studentid)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/books", http.StatusSeeOther)
}
