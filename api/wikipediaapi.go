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
	Query WikiQuery `json:"query"`
}

type WikiQuery struct {
	Page map[int](WikiData) `json:"pages"`
}

type WikiData struct {
	Thumbnail WikiThumbnail `json:"thumbnail"`
}

type WikiThumbnail struct {
	Source string `json:"source"`
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
