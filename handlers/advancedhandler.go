package handlers

import (
	"net/http"
	"text/template"
	"workers"
)

func advancedHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(advancedTemplatePath))

	MembersNumber := []int{}

	for i := 0; i < 7; i++ {
		MembersNumber = append(MembersNumber, i)
	}

	data := HtmlData{
		ProjectName: "Chazam",
		PageName:    "advanced",
		Groups:      workers.FilterGroupsByName(""),

		ASearch_MembersNumber: MembersNumber,
		ASearch_Locations:     workers.GetListOfLocation(),
	}

	if r.Method == "POST" {
		println("Advanced handler POST")
	}

	PrepareDataWithFragments(&data)
	tmpl.Execute(w, data)
}
