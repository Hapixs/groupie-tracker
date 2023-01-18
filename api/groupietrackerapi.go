package api

import (
	"strconv"
	"time"
)

type ApiArtist struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type ApiRelation struct {
	Id             int                   `json:"id"`
	DatesLocations map[string]([]string) `json:"datesLocations"`
}

func getAllArtist() []ApiArtist {
	var request []ApiArtist
	url := "https://groupietrackers.herokuapp.com/api/artists"
	GetFromApi(url, &request, false, time.Millisecond)
	return request
}

func GetRelationInfo(id int) ApiRelation {
	var request ApiRelation
	url := "https://groupietrackers.herokuapp.com/api/relation/" + strconv.Itoa(id)
	GetFromApi(url, &request, false, time.Millisecond)
	return request
}
