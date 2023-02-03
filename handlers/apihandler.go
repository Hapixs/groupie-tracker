package handlers

import (
	"encoding/json"
	"net/http"
	"objects"
	"strconv"
	"workers"
)

type ApiRequest struct {
	Groups  []objects.Group  `json:"groups"`
	Artists []objects.Artist `json:"artists"`
	Tracks  []objects.Track  `json:"tracks"`
}

func searchApiHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	data := ApiRequest{}

	data.Groups = workers.FilterGroupsByName(query)
	data.Artists = workers.FiltreAllArtistByName(query)
	data.Tracks = workers.FiltreAllTrackByName(query)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
}

func groupApiHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("id")

	id, err := strconv.Atoi(query)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode("Error")
	}

	g := workers.GroupMap[id]

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(g)
}
