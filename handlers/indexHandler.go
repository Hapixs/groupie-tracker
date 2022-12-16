package handlers

import (
	"api"
	"net/http"
	"text/template"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(indexTmpl))
	data := HtmlData{
		Artist: api.GetAllArtist(),
	}
	tmpl.Execute(w, data)
}
