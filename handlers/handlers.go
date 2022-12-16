package handlers

import (
	"api"
	"net/http"
	"os"
	"utils"
)

const (
	homeTemplatePath     = "static/templates/index.html"
	notfoundTempaltePath = "static/templates/notfound.html"
)

type HtmlData struct {
	Artist    []api.ApiArtist
	Test      string
	Fragments map[string](string)
}

func InitHandlers() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/home", homeHandler)
	http.HandleFunc("/test", testHandler)
}

func PrepareDataWithFragments(data *HtmlData) {
	data.Fragments = map[string](string){}
	fragmentFolder, err := os.ReadDir("static/templates/fragments/")
	if err != nil {
		println("Error when listing the fragement folder.")
		println(err.Error())
		return
	}

	for _, fl := range fragmentFolder {
		data.Fragments[fl.Name()] = utils.LoadFragmentAsString(fl.Name(), data)
	}
}
