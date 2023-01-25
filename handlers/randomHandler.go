package handlers

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"
	"workers"
)

func randomHandler(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().Unix())
	id := rand.Intn(len(workers.GroupMap) - 1)
	http.Redirect(w, r, "/group/"+strconv.Itoa(id)+"/", http.StatusSeeOther)
}
