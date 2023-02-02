package handlers

import (
	"encoding/json"
	"net/http"
	"objects"
	"workers"
)

type ApiRequest struct {
	Groups  []objects.Group  `json:"groups"`
	Artists []objects.Artist `json:"artists"`
	Tracks  []objects.Track  `json:"tracks"`
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	data := ApiRequest{}

	data.Groups = workers.FilterGroupsByName(query)
	data.Artists = workers.FiltreAllArtistByName(query)
	data.Tracks = workers.FiltreAllTrackByName(query)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
}
