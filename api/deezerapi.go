package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type DeezerSearch struct {
	Data []DeezerSearchTrackdata `json:"data"`
}

type DeezerSearchTrackdata struct {
	SearchArtist DeezerSearchArtist `json:"artist"`
}

type DeezerSearchArtist struct {
	Id int `json:"id"`
}

func SearchForDeezzeGroupId(artist string) DeezerSearch {
	artist = RemoveAccents(artist)
	artist = strings.Join(strings.Split(artist, " "), "%20")
	rand.Seed(time.Now().UnixMilli())
	time.Sleep(time.Duration(rand.Intn(500)))
	url := "https://api.deezer.com/search?q=" + artist
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var request DeezerSearch
	err = json.Unmarshal(body, &request)
	if err != nil {
		println("Error when parsing deezer api response for" + artist)
		println(err.Error())
		return DeezerSearch{}
	}
	if len(request.Data) <= 0 {
		println("-----")
		println(string(body))
		println("-----")
	}
	return request
}

type DeezerGroup struct {
	Sharelink string `json:"share"`
}

func GetDeezerGroup(id int) DeezerGroup {
	rand.Seed(time.Now().UnixMilli())
	time.Sleep(time.Duration(rand.Intn(500)))
	url := "https://api.deezer.com/artist/" + strconv.Itoa(id)
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var request DeezerGroup
	err = json.Unmarshal(body, &request)
	if err != nil {
		println("Error when parsing wikipedia api response for" + strconv.Itoa(id))
		return DeezerGroup{}
	}
	return request
}

type DeezerTrackRequest struct {
	List []DeezerTrack `json:"data"`
}

type DeezerTrack struct {
	Title   string      `json:"title"`
	Preview string      `json:"preview"`
	Album   DeezerAlbum `json:"album"`
}

type DeezerAlbum struct {
	Id        int                  `json:"id"`
	Title     string               `json:"title"`
	Cover     string               `json:"cover_medium"`
	GenreList DeezerAlbumGenreList `json:"genres"`
}

type DeezerAlbumGenreList struct {
	List []DeezerAlbumGenre `json:"data"`
}

type DeezerAlbumGenre struct {
	Name string `json:"name"`
}

func GetDeezerTopTrack(groupId, amount int) DeezerTrackRequest {
	rand.Seed(time.Now().UnixMilli())
	time.Sleep(time.Duration(rand.Intn(500)))
	url := "https://api.deezer.com/artist/" + strconv.Itoa(groupId) + "/top?limit=" + strconv.Itoa(amount)
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var request DeezerTrackRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		println("Error when parsing wikipedia api response for" + strconv.Itoa(groupId))
		return DeezerTrackRequest{}
	}
	for _, v := range request.List {
		v.Album = GetAlbumInformation(v.Album.Id)
		time.Sleep(time.Second / 20)
	}
	return request
}

func GetAlbumInformation(id int) DeezerAlbum {
	url := "https://api.deezer.com/album/" + strconv.Itoa(id)
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	var request DeezerAlbum
	err = json.Unmarshal(body, &request)
	if err != nil {
		println("oupsi")
		return DeezerAlbum{}
	}
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
	for k, v := range GroupMap {
		v.DZInformations = GetDeezerInformationsFromName(v.Name)
		mutex.Lock()
		GroupMap[k] = v
		mutex.Unlock()
		time.Sleep(time.Second)
	}
	LoadGenreMap()
	LoadAlternativesForGroups()
	println("All deezer informations are loaded !")
}

var genreMap = map[string]([]Group){}

func LoadGenreMap() {
	for _, group := range GroupMap {
		for _, track := range group.DZInformations.TrackList.List {
			for _, genre := range track.Album.GenreList.List {
				println("yes")
				val, ok := genreMap[genre.Name]
				if !ok {
					val = []Group{}
				}
				val = append(val, group)
				genreMap[genre.Name] = val
			}
		}
	}
}

func LoadAlternativesForGroups() {
	for id, group := range GroupMap {
		group.GroupAlternatives = []Group{}
		for _, track := range group.DZInformations.TrackList.List {
			for _, genre := range track.Album.GenreList.List {
				group.GroupAlternatives = append(group.GroupAlternatives, genreMap[genre.Name]...)
			}
		}
		mutex.Lock()
		GroupMap[id] = group
		mutex.Unlock()
	}
}
