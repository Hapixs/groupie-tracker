package tracker

import (
	"html/template"
	"net/http"
	"tracker/api-user"
)

const (
	indexTmpl = "web/templates/index.html"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(indexTmpl))
	data := Index{
		Artist: tracker.GetAllArtist(),
	}
	tmpl.Execute(w, data)
}
