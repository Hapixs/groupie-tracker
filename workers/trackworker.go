package workers

import (
	"api"
	"objects"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

func addTrackByGroup(trackRequest api.DeezerTrackList, group objects.Group) {
	for _, t := range trackRequest.List {
		track := objects.Track{
			GroupId:     group.Id,
			GroupName:   group.Name,
			Id:          t.Id,
			ReleaseDate: t.ReleaseDate,
			Title:       t.Title,
			Preview:     t.Preview,
			Album: struct {
				Id    int                `json:"id"`
				Title string             `json:"title"`
				Cover string             `json:"cover_medium"`
				Genre objects.MusicGenre `json:"genre"`
			}{
				Id:    t.Album.Id,
				Title: t.Album.Title,
				Cover: t.Album.Cover,
				Genre: objects.MusicGenre{
					Id:        t.Album.Genre.Id,
					Name:      t.Album.Genre.Name,
					Picture:   t.Album.Genre.Picture,
					PictureXl: t.Album.Genre.PictureXl,
					FontName:  "aria",
				},
			},
		}
		mutex.Lock()
		trackById[t.Id] = track
		mutex.Unlock()
	}
}

func FiltreAllTrackByName(filter string) []objects.Track {
	list := []objects.Track{}
	for _, t := range maps.Values(trackById) {
		if strings.Contains(strings.ToUpper(t.Title), strings.ToUpper(filter)) {
			list = append(list, t)
		}
	}

	slices.SortFunc(list, func(a, b objects.Track) bool {
		return strings.Index(strings.ToUpper(a.Title), strings.ToUpper(filter)) < strings.Index(strings.ToUpper(b.Title), strings.ToUpper(filter))
	})

	return list
}
