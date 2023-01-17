package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
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

func CallWikipediaApi[T any](url string, structure *T) {
	response, err := http.Get(url)
	if err != nil {
		println("Error when calling " + url)
		return
	}
	defer response.Body.Close()
	content, err := io.ReadAll(response.Body)
	if err != nil {
		println("Error when reading body of " + url)
		return
	}
	err = json.Unmarshal(content, structure)
	if err != nil {
		println("Error when pasing body of " + url)
		return
	}
}

func GetWikipediaImage(artist string) WikiRequest {
	var request WikiRequest
	artist = RemoveAccents(artist)
	artist = strings.Join(strings.Split(artist, " "), "%20")
	url := "https://en.wikipedia.org/w/api.php?action=query&titles=" + artist + "&prop=pageimages&format=json&pithumbsize=100"
	CallWikipediaApi(url, &request)
	return request
}

func GetWikipediaPageLink(artist string) string {
	var request WikiRequest
	artist = RemoveAccents(artist)
	artist = strings.Join(strings.Split(artist, " "), "%20")
	url := "https://en.wikipedia.org/w/api.php?action=query&titles=" + artist + "&prop=info&inprop=url&format=json"
	CallWikipediaApi(url, &request)
	for _, page := range request.Query.Page {
		return page.WikiUrl
	}
	return ""
}
