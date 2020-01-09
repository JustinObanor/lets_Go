package main

import "text/template"

import "log"

import "net/http"

import "strconv"

type ToDo struct{
	Name string
	Done bool
}

func main() {
	//making a template, and passing a function into it
	tpl, err := template.New("template.html").Funcs(template.FuncMap{"IsNotDone":IsNotDone}).ParseFiles("template.html")
	if err != nil{
		log.Fatal("Cant expand template", err)
		return
        default:
            fmt.Println("no value received")
        }
		}
}