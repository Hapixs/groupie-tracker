package handlers

import (
	"net/http"
	"text/template"
	"workers"

	"golang.org/x/exp/maps"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	data := HtmlData{
		Group: maps.Values(workers.GroupMap)[0],
	}
	println(len(data.Group.DateLocations))
	data.InfoMessage = "Test info message"
	data.ErrorMessage = "Test error message"
	PrepareDataWithFragments(&data)
	tmpl := template.Must(template.ParseFiles("static/templates/test.html"))
	tmpl.Execute(w, data)
}
