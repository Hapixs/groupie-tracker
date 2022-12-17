package handlers

import (
	"api"
	"net/http"
	"text/template"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(homeTemplatePath))
	data := HtmlData{
		Groups: api.GetCachedGroups(),
		Test:   "yep, fragments are working ! (and with variables ! :D) ",
	}
	PrepareDataWithFragments(&data)
	tmpl.Execute(w, data)
}
