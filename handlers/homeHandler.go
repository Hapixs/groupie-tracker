package handlers

import (
	"api"
	"net/http"
	"text/template"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(homeTemplatePath))
	groups := api.GetCachedGroups()
	data := HtmlData{
		ProjectName: "Chazam",
		PageName:    "home",
	}

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			data.ErrorMessage = "Une erreur est survenue lors de l'application des filtres.."
		} else {
			str := r.Form.Get("search_input")
			category := r.Form.Get("search_category")
			switch category {
			case "name":
				groups = api.GetGroupListFiltredByName(str)
			case "date":
				groups = api.GetGroupListFiltredByDate(str) // Todo: translate to french
			case "places":
				groups = api.GetGroupListFiltredByLocation(str)
			}
			data.LastResearchCategory = category
			data.LastResearchInput = str
		}
	}

	data.Groups = groups

	PrepareDataWithFragments(&data)
	tmpl.Execute(w, data)
}
