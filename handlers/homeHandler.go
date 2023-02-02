package handlers

import (
	"api"
	"net/http"
	"objects"
	"text/template"
	"workers"

	"golang.org/x/exp/maps"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(homeTemplatePath))
	groupsbygenre := workers.GroupByGenreMap
	groups := []objects.Group{}

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
			groups = workers.FilterGroupsByName(str)
			data.TrackResearch = workers.FiltreAllTrackByName(str)
			data.ArtistResearch = workers.FiltreAllArtistByName(str)
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
	data.GroupByGenres = groupsbygenre
	if len(groups) <= 0 && len(data.TrackResearch) <= 0 {
		for _, k := range maps.Values(groupsbygenre) {
			data.Groups = append(data.Groups, k...)
		}
	} else {
		data.GroupByGenres = map[api.DeezerGenre][]objects.Group{}
		data.Groups = groups
	}

	PrepareDataWithFragments(&data)
	tmpl.Execute(w, data)
}
