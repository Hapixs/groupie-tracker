package handlers

import (
	"net/http"
	"text/template"
)

func errorHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(notfoundTempaltePath))
	data := HtmlData{
		ErrorMessage: "La page demandé n'a pas été trouvé",
	}
	PrepareDataWithFragments(&data)
	tmpl.Execute(w, data)
}
