package tracker

import (
	"html/template"
	"net/http"
)

const (
	indexTmpl = "index.html"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(indexTmpl))
	tmpl.Execute(w, nil)
}
