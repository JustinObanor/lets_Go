package cyoa

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

var tpl *template.Template

type Story map[string]Chapter

type HandlerOpts func(h *handle)

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

type handle struct {
	story    Story
	tmpl     *template.Template
	pathFunc func(r *http.Request) string
}

func init() {
	tpl = template.Must(template.ParseFiles("doc/tpl.gohtml"))
}

func NewHandler(s Story, opts ...HandlerOpts) http.Handler {
	h := handle{
		story:    s,
		tmpl:     tpl,
		pathFunc: defaultPathFunc,
	}

	for _, opt := range opts {
		opt(&h)
	}

	return h
}

func WithTemplate(t *template.Template) HandlerOpts {
	return func(h *handle) {
		h.tmpl = t
	}
}

func WithPathFunc(f func(r *http.Request) string) HandlerOpts {
	return func(h *handle) {
		h.pathFunc = f
	}
}

func defaultPathFunc(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)

	if path == "" || path == "/" {
		path = "/intro"
	}

	return path[1:]
}

func ParseJSON(r io.Reader) (story Story, err error) {
	if err = json.NewDecoder(r).Decode(&story); err != nil {
		return nil, err
	}
	return
}

func (h handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.pathFunc(r)

	if chapter, ok := h.story[path]; ok {
		if err := h.tmpl.Execute(w, chapter); err != nil {
			log.Println(err)
			http.Error(w, "Something went wrong...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found.", http.StatusNotFound)
}
