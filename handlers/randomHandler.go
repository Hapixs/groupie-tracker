package handlers

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"
	"workers"

	"golang.org/x/exp/maps"
)

func randomHandler(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())
	id := rand.Intn((len(maps.Values(workers.GroupMap)) - 1) + 1)
	http.Redirect(w, r, "/group/"+strconv.Itoa(id)+"/", http.StatusSeeOther)
}
