package handlers

import (
	"api"
	"net/http"
	"os"
	"utils"
)

const (
	indexTmpl = "static/templates/index.html"
)

type HtmlData struct {
	Artist    []api.ApiArtist
	Test      string
	Fragments map[string](string)
}

func InitHandlers() {
	http.HandleFunc("/", indexHandler)
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
