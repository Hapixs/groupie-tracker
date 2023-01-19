package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"workers"
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

	group := workers.GroupMap[groupId]
	group.DefineMostValuableGenreForGroup()

	tmpl := template.Must(template.ParseFiles(groupTemplatePath))
	data := HtmlData{
		Group:       group,
		ProjectName: "Chazam",
	}

	PrepareDataWithFragments(&data)
	tmpl.Execute(w, data)
}
