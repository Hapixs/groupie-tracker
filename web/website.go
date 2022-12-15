package tracker

import (
	"net/http"
)

const (
	port = ":8080"
)

func Run() {
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(port, nil)
}
