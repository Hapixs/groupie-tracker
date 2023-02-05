package handlers

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"
	"workers"
)

func randomHandler(w http.ResponseWriter, r *http.Request) {
	rand.Seed(time.Now().UnixNano())
	id := rand.Intn(len(workers.GroupList)) + 1
	http.Redirect(w, r, "/group/"+strconv.Itoa(id)+"/", http.StatusSeeOther)
}
