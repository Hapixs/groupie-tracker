package handlers

import (
	"api"
	"net/http"
)

const (
	indexTmpl = "static/templates/index.html"
)

type HtmlData struct {
	Artist []api.Artist
}

func InitHandlers() {
	http.HandleFunc("/", indexHandler)
}
