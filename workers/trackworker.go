package workers

import (
	"api"
	"objects"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

func buildTrackListFromGroupId(deezerGroupId, localGroupId int, localGroupName string) []*objects.Track {
	tracks := make([]*objects.Track, 0)

	type deezerTrackList struct {
		List []struct {
			Id          int    `json:"id"`
			ReleaseDate string `json:"release_date"`
			Title       string `json:"title_short"`
			Preview     string `json:"preview"`
			Album       struct {
				Id int `json:"id"`
			} `json:"album"`
		} `json:"data"`
	}

	type deezerAlbum struct {
		Id       int    `json:"id"`
		Title    string `json:"title"`
		Cover    string `json:"cover_medium"`
		Genre_Id int    `json:"genre_id"`
	}

	deezerTrackListRequest := deezerTrackList{}
	api.GetFromApi(
		"https://api.deezer.com/artist/"+strconv.Itoa(deezerGroupId)+"/top?limit=10",
		&deezerTrackListRequest,
		false,
		time.Second/9,
		nil)

	for _, deezerTrack := range deezerTrackListRequest.List {
		deezerAlbumRequest := deezerAlbum{}

		api.GetFromApi(
			"https://api.deezer.com/album/"+strconv.Itoa(deezerTrack.Album.Id),
			&deezerAlbumRequest,
			false,
			time.Second/9,
			nil)

		track := new(objects.Track)
		track.Genre = *buildGenderFromDeezerId(deezerAlbumRequest.Id)
		track.GroupId = localGroupId
		track.GroupName = localGroupName
		track.Id = deezerTrack.Id
		track.Preview = deezerTrack.Preview
		track.Title = deezerTrack.Title
		track.ReleaseDate = deezerTrack.ReleaseDate
		tracks = append(tracks, track)

		mutex.Lock()
		TrackList = append(TrackList, track)
		TrackById[track.Id] = track
		mutex.Unlock()
	}
	return tracks
}

func FiltreAllTrackByName(filter string) []*objects.Track {
	list := []*objects.Track{}
	for _, t := range maps.Values(TrackById) {
		if strings.Contains(strings.ToUpper(t.Title), strings.ToUpper(filter)) {
			list = append(list, t)
		}
	}

	slices.SortFunc(list, func(a, b *objects.Track) bool {
		return strings.Index(strings.ToUpper(a.Title), strings.ToUpper(filter)) < strings.Index(strings.ToUpper(b.Title), strings.ToUpper(filter))
	})

	return list
}
