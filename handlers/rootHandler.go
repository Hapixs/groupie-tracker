package handlers

import (
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.String() == "/" {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	println(r.URL.String())
	http.Redirect(w, r, "/notfound", http.StatusSeeOther)
}
