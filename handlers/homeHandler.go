package handlers

import (
	"api"
	"net/http"
	"strings"
	"text/template"

	"golang.org/x/exp/slices"
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
			str := strings.ToUpper(r.Form.Get("search_input"))
			switch r.Form.Get("search_category") {
			case "name":
				slices.SortStableFunc(groups, func(a, b api.Group) bool {
					return strings.Index(strings.ToUpper(a.Name), str) > strings.Index(strings.ToUpper(b.Name), str)
				})
				sortedGroups := []api.Group{}
				for _, k := range groups {
					if strings.Contains(strings.ToUpper(k.Name), str) {
						sortedGroups = append(sortedGroups, k)
					}
				}
				groups = sortedGroups
			case "date":
			case "places":
			}
		}
	}

	data.Groups = groups

	PrepareDataWithFragments(&data)
	tmpl.Execute(w, data)
}
