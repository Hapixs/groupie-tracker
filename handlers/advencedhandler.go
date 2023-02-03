package handlers

import (
	"net/http"
	"text/template"
)

func advencedHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(homeTemplatePath))

	data := HtmlData{
		ProjectName: "Chazam",
		PageName:    "advenced",
	}

	if r.Method == "POST" {
		// Todo: something
	}

	PrepareDataWithFragments(&data)
	tmpl.Execute(w, data)
}
