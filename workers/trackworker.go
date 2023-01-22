package workers

import (
	"api"
	"objects"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

var trackById map[int]struct {
	Id          int    `json:"id"`
	ReleaseDate string `json:"release_date"`
	Title       string `json:"title_short"`
	Preview     string `json:"preview"`
	Album       struct {
		Id    int    `json:"id"`
		Title string `json:"title"`
		Cover string `json:"cover_medium"`
		Genre api.DeezerGenre
	} `json:"album"`
} = map[int]struct {
	Id          int    `json:"id"`
	ReleaseDate string `json:"release_date"`
	Title       string `json:"title_short"`
	Preview     string `json:"preview"`
	Album       struct {
		Id    int    `json:"id"`
		Title string `json:"title"`
		Cover string `json:"cover_medium"`
		Genre api.DeezerGenre
	} `json:"album"`
}{}

var groupByTrackId map[int]objects.Group = map[int]objects.Group{}

func addTrackByGroup(trackRequest api.DeezerTrackRequest, group objects.Group) {
	for _, t := range trackRequest.List {
		trackById[t.Id] = t
		groupByTrackId[t.Id] = group
	}
}

type TrackResearch struct {
	GroupId int
	Track   struct {
		Id          int    `json:"id"`
		ReleaseDate string `json:"release_date"`
		Title       string `json:"title_short"`
		Preview     string `json:"preview"`
		Album       struct {
			Id    int    `json:"id"`
			Title string `json:"title"`
			Cover string `json:"cover_medium"`
			Genre api.DeezerGenre
		} `json:"album"`
	}
}

func FiltreAllTrackByName(filter string) []TrackResearch {
	list := []TrackResearch{}
	for _, t := range maps.Values(trackById) {
		if strings.Contains(strings.ToUpper(t.Title), strings.ToUpper(filter)) {
			list = append(list, TrackResearch{groupByTrackId[t.Id].Id, t})
		}
	}

	slices.SortFunc(list, func(a, b TrackResearch) bool {
		return strings.Index(strings.ToUpper(a.Track.Title), strings.ToUpper(filter)) < strings.Index(strings.ToUpper(b.Track.Title), strings.ToUpper(filter))
	})

	return list
}
