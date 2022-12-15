package tracker

import (
	"html/template"
	"net/http"
)

const (
	port = ":8080"
)

func Run () {
	http.HandleFunc("/", index)
	http.ListenAndServe(port, nil)
}