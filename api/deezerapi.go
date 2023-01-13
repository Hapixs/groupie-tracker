package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type DeezerSearch struct {
	Data []struct {
		SearchArtist struct {
			Id int `json:"id"`
		} `json:"artist"`
	} `json:"data"`
}

func SearchForDeezzeGroupId(artist string) DeezerSearch {
	artist = RemoveAccents(artist)
	artist = strings.Join(strings.Split(artist, " "), "%20")

	url := "https://api.deezer.com/search?q=" + artist

	val, ok := ApiData.RequestAndResult[url]
	if ok {
		return val.Search
	}

	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	content, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var request DeezerSearch
	err = json.Unmarshal(content, &request)
	if err != nil {
		return DeezerSearch{}
	}
	if len(request.Data) <= 0 {
		println("-----")
		println(string(content))
		println("-----")
	}
	data := ApiData.RequestAndResult[url]
	data.Search = request
	ApiData.RequestAndResult[url] = data
	time.Sleep(time.Second / 10)
	return request
}

type DeezerGroup struct {
	Sharelink string `json:"share"`
}

func GetDeezerGroup(id int) DeezerGroup {
	url := "https://api.deezer.com/artist/" + strconv.Itoa(id)
	val, ok := ApiData.RequestAndResult[url]
	if ok {
		return val.Group
	}

	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	content, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var request DeezerGroup
	err = json.Unmarshal(content, &request)
	if err != nil {
		return DeezerGroup{}
	}
	data := ApiData.RequestAndResult[url]
	data.Group = request
	ApiData.RequestAndResult[url] = data
	time.Sleep(time.Second / 10)
	return request
}

type DeezerTrackRequest struct {
	List []struct {
		Title   string `json:"title"`
		Preview string `json:"preview"`
		Album   struct {
			Id        int    `json:"id"`
			Title     string `json:"title"`
			Cover     string `json:"cover_medium"`
			GenreList struct {
				List []struct {
					Name string `json:"name"`
				} `json:"data"`
			} `json:"genres"`
		} `json:"album"`
	} `json:"data"`
}

func GetDeezerTopTrack(groupId, amount int) DeezerTrackRequest {
	url := "https://api.deezer.com/artist/" + strconv.Itoa(groupId) + "/top?limit=" + strconv.Itoa(amount)
	val, ok := ApiData.RequestAndResult[url]

	if ok {
		return val.TrackRequest
	}
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	content, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var request DeezerTrackRequest
	err = json.Unmarshal(content, &request)
	if err != nil {
		return DeezerTrackRequest{}
	}
	data := ApiData.RequestAndResult[url]
	data.TrackRequest = request
	ApiData.RequestAndResult[url] = data
	time.Sleep(time.Second / 10)
	return request
}

// PROVIDER

type DeezerInformations struct {
	Group     DeezerGroup
	TrackList DeezerTrackRequest
}

func GetDeezerInformationsFromName(name string) DeezerInformations {
	infos := DeezerInformations{}

	s := SearchForDeezzeGroupId(name)

	if len(s.Data) <= 0 {
		println("Group not found .. for " + name)
		return DeezerInformations{}
	}

	groupId := s.Data[0].SearchArtist.Id

	infos.Group = GetDeezerGroup(groupId)
	infos.TrackList = GetDeezerTopTrack(groupId, 10)

	return infos
}

func LoadAllDeezerInformations() {
	println("Asyncron, loading all deezer informations for groups")
	ApiData = SaveAPIInformation{
		RequestAndResult: make(map[string]struct {
			TrackRequest DeezerTrackRequest
			Search       DeezerSearch
			Group        DeezerGroup
		}),
	}
	content, err := os.ReadFile("deezerdata.json")
	if err == nil {
		json.Unmarshal(content, &ApiData)
	}
	LoadAllDeezerInformationsFromApi()

	save, err := json.Marshal(ApiData)
	if err != nil {
		println("Error: JSON error")
		return
	}

	file, err := os.Create("deezerdata.json")
	if err != nil {
		println("Error 2 with deezerdata")
		return
	}

	file.Write(save)
	file.Close()

	println("All deezer informations are loaded !")
}

func LoadAllDeezerInformationsFromApi() {
	for k, v := range GroupMap {
		v.DZInformations = GetDeezerInformationsFromName(v.Name)
		mutex.Lock()
		GroupMap[k] = v
		mutex.Unlock()
	}
}

type SaveAPIInformation struct {
	RequestAndResult map[string]struct {
		TrackRequest DeezerTrackRequest
		Search       DeezerSearch
		Group        DeezerGroup
	}
}

var ApiData = SaveAPIInformation{}
