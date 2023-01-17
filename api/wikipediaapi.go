package api

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type WikiRequest struct {
	Query struct {
		Page map[int](struct {
			Thumbnail struct {
				Source string `json:"source"`
			} `json:"thumbnail"`
		}) `json:"pages"`
	} `json:"query"`
}

func GetWikipediaImage(artist string) WikiRequest {
	artist = RemoveAccents(artist)
	artist = strings.Join(strings.Split(artist, " "), "%20")
	rand.Seed(time.Now().UnixMilli())
	time.Sleep(time.Duration(rand.Intn(150)))
	url := "https://en.wikipedia.org/w/api.php?action=query&titles=" + artist + "&prop=pageimages&format=json&pithumbsize=100"
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var request WikiRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		println("Error when parsing wikipedia api response for " + artist)
		println(err.Error())
		return WikiRequest{}
	}
	return request
}

func GetWikipediaPageLink(artist string) string {
	artist = RemoveAccents(artist)
	artist = strings.Join(strings.Split(artist, " "), "%20")
	rand.Seed(time.Now().UnixMilli())
	time.Sleep(time.Duration(rand.Intn(150)))
	url := "https://en.wikipedia.org/w/api.php?action=query&titles=" + artist + "&prop=info&inprop=url&format=json"
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var request WikiRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		println("Error when parsing wikipedia api response for " + artist)
		println(err.Error())
		return ""
	}
	for _, page := range request.Query.Page {
		return page.Thumbnail.Source
	}
	return ""
}