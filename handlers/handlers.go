package handlers

import (
	"api"
	"net/http"
	"os"
	"strconv"
	"utils"
)

const (
	homeTemplatePath     = "static/templates/index.html"
	notfoundTempaltePath = "static/templates/notfound.html"
	groupTemplatePath    = "static/templates/group.html"
)

type HtmlData struct {
	InfoMessage  string
	ErrorMessage string
	ErrorCode    int

	PageName string

	Groups []api.Group
	//Test      string
	Group       api.Group
	Fragments   map[string](string)
	ProjectName string
}

func InitHandlers() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/home", homeHandler)
	http.HandleFunc("/test", testHandler)
	http.HandleFunc("/notfound", errorHandler)

	for k := range api.GroupMap {
		link := "/group/" + strconv.Itoa(k) + "/"
		http.HandleFunc(link, groupHandler)
	}

	http.HandleFunc("/random", randomHandler)
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
