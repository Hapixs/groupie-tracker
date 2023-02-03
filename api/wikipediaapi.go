package api

import (
	"strings"
	"time"
	"utils"
)

type WikiRequest struct {
	Query struct {
		Page map[int](struct {
			Thumbnail struct {
				Source string `json:"source"`
			} `json:"thumbnail"`
			WikiUrl string `json:"fullurl"`
		}) `json:"pages"`
	} `json:"query"`
}

func GetWikipediaImage(artist string) WikiRequest {
	var request WikiRequest
	artist = utils.FormatArtistName(artist)
	url := "https://fr.wikipedia.org/w/api.php?action=query&titles=" + artist + "&prop=pageimages&format=json&pithumbsize=100"
	GetFromApi(url, &request, false, time.Millisecond, nil)
	return request
}

func GetWikipediaPageLink(artist string) string {
	if strings.EqualFold(artist, "Gary Maurice Lucas Jr.") {
		return "https://fr.wikipedia.org/wiki/Joyner_Lucas"
	}
	var request WikiRequest
	artist = utils.FormatArtistName(artist)
	url := "https://fr.wikipedia.org/w/api.php?action=query&titles=" + artist + "&prop=info&inprop=url&format=json"
	GetFromApi(url, &request, false, time.Millisecond, nil)
	for _, page := range request.Query.Page {
		return page.WikiUrl
	}
	return "/"
}
