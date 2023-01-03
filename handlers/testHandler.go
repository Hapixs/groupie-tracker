package handlers

import (
	"net/http"
	"text/template"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	data := HtmlData{}
	data.InfoMessage = "Test info message"
	data.ErrorMessage = "Test error message"
	PrepareDataWithFragments(&data)
	tmpl := template.Must(template.ParseFiles(homeTemplatePath))
	tmpl.Execute(w, data)
}
