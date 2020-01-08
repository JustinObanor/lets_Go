package main

import "text/template"

import "log"

import "net/http"

import "strconv"

type ToDo struct{
	Name string
	Done bool
}

func IsNotDone(todo ToDo)bool{
	return !todo.Done
}

func main() {
	//making a template, and passing a function into it
	tpl, err := template.New("template.html").Funcs(template.FuncMap{"IsNotDone":IsNotDone}).ParseFiles("template.html")
	if err != nil{
		log.Fatal("Cant expand template", err)
		return
	}
	todos := []ToDo{
		{"Study Go", false},
		{"Attend lectures per web dev", false},
		{"...", false},
		{"Profit", false},
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		if r.Method == http.MethodPost{
			param := r.FormValue("id")
			index, _ := strconv.ParseInt(param, 10,0)
			todos[index].Done = true
		}
		err := tpl.Execute(w,todos)
		if err != nil{
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	log.Fatal(http.ListenAndServe(":8080", nil)) 
}