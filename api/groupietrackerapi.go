package api

import (
	"strconv"
	"time"
)

type ApiArtist struct {
	Id             int                   `json:"id"`
	Image          string                `json:"image"`
	Name           string                `json:"name"`
	Members        []string              `json:"members"`
	CreationDate   int                   `json:"creationDate"`
	FirstAlbum     string                `json:"firstAlbum"`
	DatesLocations map[string]([]string) `json:"datesLocations"`
}

func GetGroupieArtistList() []ApiArtist {
	var request []ApiArtist
	url := "https://groupietrackers.herokuapp.com/api/artists"
	GetFromApi(url, &request, false, time.Millisecond)
	return request
}

func UpdateGroupRelation(artist ApiArtist) {
	request := &artist
	url := "https://groupietrackers.herokuapp.com/api/relation/" + strconv.Itoa(artist.Id)
	GetFromApi(url, &request, false, time.Millisecond)
}
