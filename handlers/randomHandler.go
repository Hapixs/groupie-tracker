package handlers

import (
	"api"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func randomHandler(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().Unix())
	id := rand.Intn(len(api.GroupMap))
	http.Redirect(w, r, "/group/"+strconv.Itoa(id)+"/", http.StatusSeeOther)
}
