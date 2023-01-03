package handlers

import (
	"net/http"
	"text/template"
)

func errorHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(notfoundTempaltePath))
	data := HtmlData{}
	PrepareDataWithFragments(&data)
	tmpl.Execute(w, data)
}
