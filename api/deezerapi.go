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
	RequestAndResult map[string]struct {
		TrackRequest DeezerTrackRequest
		Search       DeezerSearch
		Group        DeezerGroup
		Album        DeezerAlbum
		Genre        DeezerGenre
	}
}

var ApiData = DeezerApiCache{}
var GroupByGenreMap = map[DeezerGenre](*Group){}

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

func GetDeezerAlbumInformation(albumId int, update bool) DeezerAlbum {
	var request DeezerAlbum
	url := "https://api.deezer.com/album/" + strconv.Itoa(albumId)
	val, ok := ApiData.RequestAndResult[url]
	if ok && !update {
		return val.Album
	}
	CallDeezerApi(url, &request)
	data := ApiData.RequestAndResult[url]
	data.Album = request
	ApiData.RequestAndResult[url] = data
	time.Sleep(time.Second / 10)
	return request
}

func GetGenreById(id int, update bool) DeezerGenre {
	var request DeezerGenre
	url := "https://api.deezer.com/genre/" + strconv.Itoa(id)
	val, ok := ApiData.RequestAndResult[url]
	if ok && !update {
		return val.Genre
	}
	CallDeezerApi(url, &request)
	data := ApiData.RequestAndResult[url]
	data.Genre = request
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

func LoadAllDeezerInformations() {
	println("Asyncron, loading all deezer informations for groups")
	ApiData = DeezerApiCache{
		RequestAndResult: make(map[string]struct {
			TrackRequest DeezerTrackRequest
			Search       DeezerSearch
			Group        DeezerGroup
			Album        DeezerAlbum
			Genre        DeezerGenre
		}),
	}

	content, err := os.ReadFile("deezerdata.json")
	if err == nil {
		json.Unmarshal(content, &ApiData)
	} else {
		println("Seams like the first start of this web server.")
		println("Some operations may take several minutes.")
	}
	UpdateAllDeezerInformations(false)
	SaveDeezerApiCache()
	println("All deezer informations are loaded !")
	go DeezerApiUpdateroutine()
}

func UpdateAllDeezerInformations(forceUpdate bool) {
	for k, v := range GroupMap {
		v.DZInformations = GetDeezerInformationsFromName(v.Name, forceUpdate)
		DefineMostValuableGenreForGroup(&v)
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

	file, err := os.Create("deezerdata.json")
	if err != nil {
		println("Error 2 with deezerdata")
		return
	}

	file.Write(save)
	file.Close()
}

func DefineMostValuableGenreForGroup(group *Group) {
	table := map[DeezerGenre](int){}
	for _, track := range group.DZInformations.TrackList.List {
		i := 1
		val, ok := table[track.Album.Genre]
		if ok {
			i += val
		}
		table[track.Album.Genre] = i
	}
	var top DeezerGenre = DeezerGenre{}
	for k, v := range table {
		if v > table[top] {
			top = k
		}
	}
	group.MostValuableGenre = top
	GroupByGenreMap[top] = group
}

func GetDeezerGenreList() []DeezerGenre {
	list := []DeezerGenre{}
	for k := range GroupByGenreMap {
		if k.Name == "" {
			continue
		}
		list = append(list, k)
	}
	return list
}
