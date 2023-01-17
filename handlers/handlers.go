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

	LastResearchInput    string
	LastResearchCategory string

	Groups      []api.Group
	Group       api.Group
	Fragments   map[string](string)
	ProjectName string

	DeezerGenres  []api.DeezerGenre
	GroupByGenres map[api.DeezerGenre][]api.Group
}

func InitHandlers() {
	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

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

func CheckForMessageQuery(r *http.Request, data *HtmlData) {
	info := r.URL.Query().Get("info")
	err := r.URL.Query().Get("error")

	data.InfoMessage = info
	data.ErrorMessage = err

}
