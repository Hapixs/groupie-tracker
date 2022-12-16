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
		Test:   "yep, fragments are working ! (and with variables ! :D) ",
	}

	PrepareDataWithFragments(&data)
	tmpl.Execute(w, data)
}
