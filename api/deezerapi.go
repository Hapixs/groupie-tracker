package api

import (
	"strconv"
	"strings"
	"time"
	"utils"
)

type DeezerSearch struct {
	Data []struct {
		SearchArtist struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"artist"`
	} `json:"data"`
}

type DeezerGroup struct {
	Sharelink string `json:"share"`
}

type DeezerTrackRequest struct {
	List []struct {
		Id          int    `json:"id"`
		ReleaseDate string `json:"release_date"`
		Title       string `json:"title"`
		Preview     string `json:"preview"`
		Album       struct {
			Id    int    `json:"id"`
			Title string `json:"title"`
			Cover string `json:"cover_medium"`
			Genre DeezerGenre
		} `json:"album"`
	} `json:"data"`
}

type DeezerAlbum struct {
	Id       int `json:"id"`
	Genre_Id int `json:"genre_id"`
}

type DeezerGenre struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture_medium"`
}

type DeezerApiCache struct {
	RequestAndResult map[string]string
}

func SearchForDeezerGroupId(artist string, update bool) DeezerSearch {
	var request DeezerSearch
	artist = utils.FormatArtistName(artist)
	url := "https://api.deezer.com/search?q=" + artist
	GetFromApi(url, &request, update, time.Second/10, nil)
	return request
}

func GetDeezerGroup(id int, update bool) DeezerGroup {
	var request DeezerGroup
	url := "https://api.deezer.com/artist/" + strconv.Itoa(id)
	GetFromApi(url, &request, update, time.Second/10, nil)
	return request
}

func GetDeezerTopTrack(groupId, amount int, update bool) DeezerTrackRequest {
	var request DeezerTrackRequest
	url := "https://api.deezer.com/artist/" + strconv.Itoa(groupId) + "/top?limit=" + strconv.Itoa(amount)
	GetFromApi(url, &request, update, time.Second/10, nil)
	GetTracksReleaseDate(&request, update)
	return request
}

func GetDeezerAlbumInformation(albumId int, update bool) DeezerAlbum {
	var request DeezerAlbum
	url := "https://api.deezer.com/album/" + strconv.Itoa(albumId)
	GetFromApi(url, &request, update, time.Second/10, nil)
	return request
}

func GetGenreById(id int, update bool) DeezerGenre {
	var request DeezerGenre
	url := "https://api.deezer.com/genre/" + strconv.Itoa(id)
	GetFromApi(url, &request, update, time.Second/10, nil)
	return request
}

func GetTracksReleaseDate(req *DeezerTrackRequest, update bool) {
	for _, k := range req.List {
		url := "https://api.deezer.com/track/" + strconv.Itoa(k.Id)
		GetFromApi(url, &k, update, time.Second/10, nil)
		println(k.ReleaseDate)
	}
}

type DeezerInformations struct {
	Group     DeezerGroup
	TrackList DeezerTrackRequest
}

func GetDeezerInformationsFromName(name string, update bool) DeezerInformations {
	infos := DeezerInformations{}
	s := SearchForDeezerGroupId(name, update)
	if len(s.Data) <= 0 {
		println("No groupe found for " + name)
		return DeezerInformations{}
	}
	groupId := -1
	for _, data := range s.Data {
		if strings.Contains(data.SearchArtist.Name, name) {
			groupId = data.SearchArtist.Id
			break
		}
	}
	if groupId < 0 {
		groupId = s.Data[0].SearchArtist.Id
	}
	infos.Group = GetDeezerGroup(groupId, update)
	trackRequest := GetDeezerTopTrack(groupId, 10, update)
	UpdateGenreForTracksAlbum(&trackRequest)
	infos.TrackList = trackRequest
	return infos
}

func UpdateGenreForTracksAlbum(request *DeezerTrackRequest) {
	for id, track := range request.List {
		album := GetDeezerAlbumInformation(track.Album.Id, false)
		genre := GetGenreById(album.Genre_Id, false)
		request.List[id].Album.Genre = genre
	}
}
