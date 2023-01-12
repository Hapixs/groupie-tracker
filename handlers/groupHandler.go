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
	println(r.URL.String())
	if splitedUrl[1] != "group" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	groupId, err := strconv.Atoi(splitedUrl[2])

	if err != nil {
		// TODO: redirect with error message
	}

	group := api.GetGroupFromId(groupId)
	tmpl := template.Must(template.ParseFiles(groupTemplatePath))
	data := HtmlData{
		//Groups:      api.GetCachedGroups(),
		Group:       group,
		ProjectName: "Chazam",
	}
	PrepareDataWithFragments(&data)
	tmpl.Execute(w, data)
}
