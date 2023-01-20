package handlers

import (
	"net/http"
	"text/template"
	"workers"

	"golang.org/x/exp/maps"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(homeTemplatePath))
	groups := workers.GroupByGenreMap

	data := HtmlData{
		ProjectName:  "Chazam",
		PageName:     "home",
		DeezerGenres: workers.GetDeezerGenreList(),
	}

	CheckForMessageQuery(r, &data)

	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			data.ErrorMessage = "Une erreur est survenue lors de l'application des filtres.."
		} else {
			str := r.Form.Get("search_input")
			category := r.Form.Get("search_category")
			groups = workers.FilterGroups(str)
			switch category {
			case "all":
				//groups = api.GetGroupListFiltredByAll(str)
			case "name":
				//groups = api.GetGroupListFiltredByName(str)
			case "date":
				//groups = api.GetGroupListFiltredByDate(str) // Todo: translate to french
			case "places":
				// groups = api.GetGroupListFiltredByLocation(str)
			}
			data.LastResearchCategory = category
			data.LastResearchInput = str
		}
	}

	data.GroupByGenres = groups
	for _, k := range maps.Values(groups) {
		data.Groups = append(data.Groups, k...)
	}

	PrepareDataWithFragments(&data)
	tmpl.Execute(w, data)
}
