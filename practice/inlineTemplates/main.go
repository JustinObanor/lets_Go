package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)



func main() {
	tpl, _ := template.New("main").Parse(`<div style="display: inline-block; border: 1px solid #aaa;
	border-radius: 3px; padding:30px; margin:20px;">
		<pre>{{.}}</pre>
		</div>
		`)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		c := http.Client{}
		res, err := c.Get("http://artii.herokuapp.com/make?text=" + path)
		if err != nil {
			log.Println(err)
		}
		defer res.Body.Close()

		body, _ := ioutil.ReadAll(res.Body)

		w.WriteHeader(http.StatusOK)
		tpl.Execute(w, string(body))
		//log.Println(string(body))
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
