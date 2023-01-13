package handlers

import (
	"api"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

func groupHandler(w http.ResponseWriter, r *http.Request) {
	splitedUrl := strings.Split(r.URL.String(), "/")
	if splitedUrl[1] != "group" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	groupId, err := strconv.Atoi(splitedUrl[2])

	if err != nil {
		http.Redirect(w, r, "/home?error=Une erreur est survenue", http.StatusSeeOther)
		return
	}

	group := api.GetGroupFromId(groupId)
	tmpl := template.Must(template.ParseFiles(groupTemplatePath))
	data := HtmlData{
		//Groups:      api.GetCachedGroups(),
		Group:       group,
		ProjectName: "Chazam",
	}

	println(len(group.GroupAlternatives))

	PrepareDataWithFragments(&data)
	tmpl.Execute(w, data)
}
