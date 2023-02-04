package handlers

import (
	"net/http"
	"text/template"
)

func advancedHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(advancedTemplatePath))

	data := HtmlData{
		ProjectName: "Chazam",
		PageName:    "advanced",
	}

	if r.Method == "POST" {
		// Todo: something
		println("Advanced handler POST")
	}

	PrepareDataWithFragments(&data)
	tmpl.Execute(w, data)
}
