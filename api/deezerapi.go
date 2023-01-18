package api

import (
	"strconv"
	"strings"
	"time"
	"utils"
)

type DeezerSearch struct {
	Data []struct {
		SearchArtist struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
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
	RequestAndResult map[string]string
}

var ApiData = DeezerApiCache{}
var GroupByGenreMap = map[DeezerGenre]([]Group){}

func SearchForDeezerGroupId(artist string, update bool) DeezerSearch {
	var request DeezerSearch
	artist = utils.FormatArtistName(artist)
	url := "https://api.deezer.com/search?q=" + artist
	GetFromApi(url, &request, update, time.Second/10)
	return request
}

func GetDeezerGroup(id int, update bool) DeezerGroup {
	var request DeezerGroup
	url := "https://api.deezer.com/artist/" + strconv.Itoa(id)
	GetFromApi(url, &request, update, time.Second/10)
	return request
}

func GetDeezerTopTrack(groupId, amount int, update bool) DeezerTrackRequest {
	var request DeezerTrackRequest
	url := "https://api.deezer.com/artist/" + strconv.Itoa(groupId) + "/top?limit=" + strconv.Itoa(amount)
	GetFromApi(url, &request, update, time.Second/10)
	return request
}

func GetDeezerAlbumInformation(albumId int, update bool) DeezerAlbum {
	var request DeezerAlbum
	url := "https://api.deezer.com/album/" + strconv.Itoa(albumId)
	GetFromApi(url, &request, update, time.Second/10)
	return request
}

func GetGenreById(id int, update bool) DeezerGenre {
	var request DeezerGenre
	url := "https://api.deezer.com/genre/" + strconv.Itoa(id)
	GetFromApi(url, &request, update, time.Second/10)
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
	groupId := -1
	for _, data := range s.Data {
		if strings.Contains(data.SearchArtist.Name, name) {
			groupId = data.SearchArtist.Id
			break
		}
	}
	if groupId < 0 {
		groupId = s.Data[0].SearchArtist.Id
	}
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
	UpdateAllDeezerInformations(false)
	SaveApiCache()
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
		println("Updating all deezer informations in background...")
		UpdateAllDeezerInformations(true)
		SaveApiCache()
	}
}

func DefineMostValuableGenreForGroup(group *Group) {
	var top DeezerGenre = DeezerGenre{Name: ""}
	table := map[DeezerGenre](int){top: 0}
	for _, track := range group.DZInformations.TrackList.List {
		i := 1
		val, ok := table[track.Album.Genre]
		if ok {
			i += val
		}
		table[track.Album.Genre] = i
	}
	for k, v := range table {
		if v > table[top] {
			top = k
		}
	}
	group.MostValuableGenre = top
	groups := []Group{*group}
	val, ok := GroupByGenreMap[top]
	if ok {
		for _, v := range val {
			if group.Id == v.Id {
				return
			}
		}
		groups = append(groups, val...)
	}
	GroupByGenreMap[top] = groups
}

func UpdateAlternativeGroupsForGroup(group *Group) {
	group.GroupAlternatives = GroupByGenreMap[group.MostValuableGenre]
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
