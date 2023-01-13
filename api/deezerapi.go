package api

import (
	"encoding/json"
	"io"
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

type DeezerGroup struct {
	Sharelink string `json:"share"`
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

type DeezerApiCache struct {
	RequestAndResult map[string]struct {
		TrackRequest DeezerTrackRequest
		Search       DeezerSearch
		Group        DeezerGroup
	}
}

var ApiData = DeezerApiCache{}

func CallDeezerApi[T any](url string, structure *T) {
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
	time.Sleep(time.Second / 10)
}

func SearchForDeezerGroupId(artist string, update bool) DeezerSearch {
	var request DeezerSearch
	artist = RemoveAccents(artist)
	artist = strings.Join(strings.Split(artist, " "), "%20")
	url := "https://api.deezer.com/search?q=" + artist
	val, ok := ApiData.RequestAndResult[url]
	if ok && !update {
		return val.Search
	}
	CallDeezerApi(url, &request)
	data := ApiData.RequestAndResult[url]
	data.Search = request
	ApiData.RequestAndResult[url] = data
	time.Sleep(time.Second / 10)
	return request
}

func GetDeezerGroup(id int, update bool) DeezerGroup {
	var request DeezerGroup
	url := "https://api.deezer.com/artist/" + strconv.Itoa(id)
	val, ok := ApiData.RequestAndResult[url]
	if ok && !update {
		return val.Group
	}
	CallDeezerApi(url, &request)
	data := ApiData.RequestAndResult[url]
	data.Group = request
	ApiData.RequestAndResult[url] = data
	time.Sleep(time.Second / 10)
	return request
}

func GetDeezerTopTrack(groupId, amount int, update bool) DeezerTrackRequest {
	var request DeezerTrackRequest
	url := "https://api.deezer.com/artist/" + strconv.Itoa(groupId) + "/top?limit=" + strconv.Itoa(amount)
	val, ok := ApiData.RequestAndResult[url]
	if ok && !update {
		return val.TrackRequest
	}
	CallDeezerApi(url, &request)
	data := ApiData.RequestAndResult[url]
	data.TrackRequest = request
	ApiData.RequestAndResult[url] = data
	time.Sleep(time.Second / 10)
	return request
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

	groupId := s.Data[0].SearchArtist.Id
	infos.Group = GetDeezerGroup(groupId, update)
	infos.TrackList = GetDeezerTopTrack(groupId, 10, update)

	return infos
}

func LoadAllDeezerInformations() {
	println("Asyncron, loading all deezer informations for groups")
	ApiData = DeezerApiCache{
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
	UpdateAllDeezerInformations(false)
	SaveDeezerApiCache()
	println("All deezer informations are loaded !")
	go DeezerApiUpdateroutine()
}

func UpdateAllDeezerInformations(forceUpdate bool) {
	for k, v := range GroupMap {
		v.DZInformations = GetDeezerInformationsFromName(v.Name, forceUpdate)
		mutex.Lock()
		GroupMap[k] = v
		mutex.Unlock()
	}
}

func DeezerApiUpdateroutine() {
	for {
		time.Sleep(time.Minute * 5)
		println("Updating all deezer informations...")
		UpdateAllDeezerInformations(true)
		SaveDeezerApiCache()
		println("Update done!")
	}
}

func SaveDeezerApiCache() {
	save, err := json.Marshal(ApiData)
	if err != nil {
		println("Error: JSON error")
		return
	}
	_, err = os.OpenFile("deezerdata.json", int(os.ModePerm), os.ModePerm)
	if err != nil {
		println("Seams like the first start of this web server.")
		println("Some operations may take several minutes.")
	}

	file, err := os.Create("deezerdata.json")
	if err != nil {
		println("Error 2 with deezerdata")
		return
	}

	file.Write(save)
	file.Close()
}
